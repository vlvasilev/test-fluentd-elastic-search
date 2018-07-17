#!bin/bash

SCRIPT_PATH=$(dirname "$0")
if [ -z "$SCRIPT_PATH" ]; then
  SCRIPT_PATH="."
fi

kubectl apply -f ${SCRIPT_PATH}/namespace/namespace.yaml
kubectl apply -f ${SCRIPT_PATH}/elastic_search/
kubectl apply -f ${SCRIPT_PATH}/fluent/fluentd/
kubectl apply -f ${SCRIPT_PATH}/fluent/fluent_bit/

