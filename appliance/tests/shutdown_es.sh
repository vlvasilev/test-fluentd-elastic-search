#!/bin/bash
EL_SRCH_TO_SHUTDOWN=${1:0}

NAMESPACES+=$(kubectl get ns -l garden.sapcloud.io/role=shoot | awk '{if(NR>1)print $1}')

# if [ "$EL_SRCH_TO_SHUTDOWN" -le "0" -o "$EL_SRCH_TO_SHUTDOWN" -gt "${#NAMESPACES[@]}" ]; then
#     EL_SRCH_TO_SHUTDOWN=${#NAMESPACES[@]}
# fi

declare -i SHUTED_DOWN=0
for namespace in $NAMESPACES; do
    if [ "$SHUTED_DOWN" -lt "$EL_SRCH_TO_SHUTDOWN" ]; then
        kubectl -n $namespace scale statefulset elasticsearch-logging --replicas=0
        SHUTED_DOWN+=1
    else
        break    
    fi
done
