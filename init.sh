#!/bin/sh

set -e

echo "[Uninstalling]"
kubectl delete MutatingWebhookConfiguration/putty 2> /dev/null || /bin/true
kubectl delete svc/putty 2> /dev/null || /bin/true
kubectl delete deploy/putty 2> /dev/null || /bin/true
kubectl delete cm/putty 2> /dev/null || /bin/true

echo "[Installing]"
kubectl apply -f example/putty-cm.yaml
kubectl apply -f example/putty-deploy.yaml
kubectl apply -f example/putty-svc.yaml
sleep 5
kubectl apply -f example/mutating-webhook-configuration.yaml
echo "[Done]"
