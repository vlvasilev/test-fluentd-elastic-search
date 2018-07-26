#!/bin/bash
EL_ES_TO_TURNON=${1:0}

NAMESPACES+=$(kubectl get ns -l garden.sapcloud.io/role=shoot | awk '{if(NR>1)print $1}')

# if [ "$EL_ES_TO_TURNON" -le "0" -o "$EL_ES_TO_TURNON" -gt "${#NAMESPACES[@]}" ]; then
#     EL_ES_TO_TURNON=${#NAMESPACES[@]}
# fi

declare -i TURNED_ON=0
for namespace in ${NAMESPACES}; do
    if [ "$TURNED_ON" -lt "$EL_ES_TO_TURNON" ]; then
        kubectl -n $namespace scale statefulset elasticsearch-logging --replicas=1
        TURNED_ON+=1
    else
        break
    fi
done