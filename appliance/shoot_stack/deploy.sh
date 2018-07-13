#!/bin/bash

SCRIPT_PATH=$(dirname "$0")
if [ -z $SCRIPT_PATH ]; then
  SCRIPT_PATH="."
fi

#delete all shoot namespaces
kubectl delete namespace -l garden.sapcloud.io/role=shoot
sleep 30

#make new shoot namespaces and deploy elastic-search
for i in `seq 1 $1`; do
  sed -e "s/<NAMESPACE>/shoot$i/g" ${SCRIPT_PATH}/namespace-template.yaml > ${SCRIPT_PATH}/namespace.yaml
  kubectl apply -f ${SCRIPT_PATH}/namespace.yaml
  sed -e "s/<NAMESPACE>/shoot$i/g" ${SCRIPT_PATH}/es-statefulset-template.yaml > ${SCRIPT_PATH}/es-statefulset.yaml
  kubectl apply -f ${SCRIPT_PATH}/es-statefulset.yaml
  sed -e "s/<NAMESPACE>/shoot$i/g" ${SCRIPT_PATH}/es-service-template.yaml > ${SCRIPT_PATH}/es-service.yaml
  kubectl apply -f ${SCRIPT_PATH}/es-service.yaml
done

#clean
rm ${SCRIPT_PATH}/namespace.yaml
rm ${SCRIPT_PATH}/es-statefulset.yaml
rm ${SCRIPT_PATH}/es-service.yaml


