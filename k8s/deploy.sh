#!/bin/bash -e

kubectl apply -f rk-configmap.yml
kubectl apply -f rk-pv.yml
kubectl apply -f rk-pvc.yml
kubectl apply -f rk-service.yml
kubectl apply -f rk-statefulset.yml