package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	
	"k8s.io/klog"
)

// PuttyInitializePatch contains JSON patches to apply if a field does not exist.
type PuttyInitializePatch struct {
	Path string `json:"path"`
	PathPieces []string `json:"-"`
	Patch []interface{} `json:"patch"`
}

// PuttyPatch contains JSON patches to apply.
type PuttyPatch struct {
	InitializePatches []PuttyInitializePatch `json:"initializePatches"`
	Patch []interface{} `json:"patch"`
}

// Config contains server configuration.
type Config struct {
	CertFile string
	KeyFile  string
	TLS bool
	Port string
	PatchFile string
	Patch PuttyPatch
}

func (c *Config) addFlags() {
	klog.InitFlags(nil)
	flag.StringVar(&c.CertFile, "tls-cert-file", "cert.pem", ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")
	flag.StringVar(&c.KeyFile, "tls-private-key-file", "key.pem", ""+
		"File containing the default x509 private key matching --tls-cert-file.")
	flag.BoolVar(&c.TLS, "tls", true, "Enable TLS")

	flag.StringVar(&c.Port, "port", "443", "Port")

	flag.StringVar(&c.PatchFile, "patch-file", "", ""+
		"File containing Putty patch definition.")
}

func configTLS(config Config) *tls.Config {
	sCert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		klog.Fatal(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{sCert},
	}
}

func loadPatch(config *Config) {
	var data []byte

	if config.PatchFile != "" {
		var err error
		data, err = ioutil.ReadFile(config.PatchFile)
		if err != nil {
			klog.Fatal(err)
		}
	} else if os.Getenv("PUTTY_PATCH") != "" {
		data = []byte(os.Getenv("PUTTY_PATCH"))
	} else {
		klog.Fatal("Missing -patch-file or PUTTY_PATCH")
	}
	klog.Info(fmt.Sprintf("Configured with JSON patch: %s", data))

	var puttyPatch PuttyPatch
	json.Unmarshal(data, &puttyPatch)	

	for index, initializePatch := range puttyPatch.InitializePatches {
		var path []string
		for _, piece := range strings.Split(initializePatch.Path, "/")[1:] {
			path = append(path, strings.Title(piece))
		}
		puttyPatch.InitializePatches[index].PathPieces = path
	}

	config.Patch = puttyPatch
}
