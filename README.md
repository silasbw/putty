# Putty for Kubernetes

Interpose on Kubernetes operations and apply a [JSON patch](http://jsonpatch.com/) to select resources.

## Examples

### Add [initContainer](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/)

[example/](example/) adds an initContainer to every Pod with a label matching `putty.sbw.io: "true"`.

## Build

```
go build
```

## Run locally

```
./putty -tls=false -port=8080
```

## Run on minikube

```
./init.sh
```

## Publish

```
docker build -t silasbw/putty:latest .
docker push silasbw/putty:latest
```

## Implementation

Putty is a [Dynamic MutatingAdmissionWebhook controller](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/).

## Related

* [Deprecated Admission Webhook example](https://github.com/caesarxuchao/example-webhook-admission-controller)
* [MutatingWebhookConfiguration reference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.15/#mutatingwebhookconfiguration-v1beta1-admissionregistration-k8s-io)
