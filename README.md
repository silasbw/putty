# Kubernetes External Admission Webhook Test Image

The image tests MutatingAdmissionWebhook and ValidatingAdmissionWebhook. After deploying
it to kubernetes cluster, administrator needs to create a ValidatingWebhookConfiguration
in kubernetes cluster to register remote webhook admission controllers.

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

## Notes

https://github.com/caesarxuchao/example-webhook-admission-controller
