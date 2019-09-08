package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

var config Config

// toAdmissionResponse is a helper function to create an AdmissionResponse
// with an embedded error
func toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

func mutate(ar v1beta1.AdmissionReview, puttyPatch PuttyPatch) *v1beta1.AdmissionResponse {
	klog.V(2).Info("mutating pods")
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if ar.Request.Resource != podResource {
		klog.Errorf("expect resource to be %s", podResource)
		return nil
	}

	raw := ar.Request.Object.Raw
	pod := corev1.Pod{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true

	var patch []interface{}
	for _, initializePatch := range puttyPatch.InitializePatches {
		if !Exists(pod, initializePatch.PathPieces) {
			patch = append(patch, initializePatch.Patch ...)
		}
	}
	patch = append(patch, puttyPatch.Patch...)

	data, err := json.Marshal(patch)
	if err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}

	reviewResponse.Patch = data
	pt := v1beta1.PatchTypeJSONPatch
	reviewResponse.PatchType = &pt

	return &reviewResponse
}

func handler(w http.ResponseWriter, r *http.Request) {
	klog.V(2).Info(fmt.Sprintf("handling request: %s %s %s",
		r.Method,
		r.URL,
		r.Header.Get("Content-Type")))

	if r.Body == nil {
		klog.Warning("Empty Request.Body")
		w.WriteHeader(400)
		return
	}

	var body []byte
	if data, err := ioutil.ReadAll(r.Body); err != nil {
		body = data
	} else {
		klog.Error(err)
		w.WriteHeader(400)
		return
	}

	requestedAdmissionReview := v1beta1.AdmissionReview{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &requestedAdmissionReview); err != nil {
		klog.Error(err)
		w.WriteHeader(400)
		return
	}

	responseAdmissionReview := v1beta1.AdmissionReview{}
	responseAdmissionReview.Response = mutate(requestedAdmissionReview, config.Patch)
	// https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#response
	responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID

	klog.V(2).Info(fmt.Sprintf("sending response: %v", responseAdmissionReview.Response))
	respBytes, err := json.Marshal(responseAdmissionReview)
	if err != nil {
		klog.Error(err)
		w.WriteHeader(500)
		return
	}
	if _, err := w.Write(respBytes); err != nil {
		klog.Error(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func main() {
	config.addFlags()
	flag.Parse()
	loadPatch(&config)

	http.HandleFunc("/", handler)
	if config.TLS {
		klog.Info(fmt.Sprintf("Listening on :%s with TLS", config.Port))
		server := &http.Server{
			Addr: fmt.Sprintf(":%s", config.Port),
			TLSConfig: loadTLS(config),
		}
		klog.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		klog.Info(fmt.Sprintf("Listening on :%s without TLS", config.Port))
		klog.Info("TLS is REQUIRED for Kubernetes Webhooks")
		server := &http.Server{
			Addr: fmt.Sprintf(":%s", config.Port),
		}
		klog.Fatal(server.ListenAndServe())
	}
}
