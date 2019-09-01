#!/bin/sh

set -e

echo "[Uninstalling]"
kubectl delete ValidatingWebhookConfiguration/putty 2> /dev/null || /bin/true
kubectl delete svc/putty 2> /dev/null || /bin/true
kubectl delete deploy/putty 2> /dev/null || /bin/true

echo "[Installing]"
kubectl apply -f example/putty-deploy.yaml
kubectl apply -f example/putty-svc.yaml
sleep 5
kubectl apply -f example/validating-webhook-configuration.yaml
echo "[Done]"
