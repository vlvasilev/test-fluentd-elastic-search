#!/bin/sh
sed -e "s/<PODS>/$1/g" -e "s/<EEP>/$2/g" -e "s/<MASTER>/$3/g" -e "s/<MSGCOUNT>/$4/g" -e "s/<LOGTIME>/$5/g" -e "s/<TTWALS>/$6/g" -e "s/<NAMESPACE>/$7/g" jobtemplate.yaml > job.yaml
output=$(kubectl --kubeconfig=/conf/kubeconfig.yaml apply -f job.yaml)
result=$?
echo $output
return result
