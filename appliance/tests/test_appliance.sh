#!/bin/bash
PODS_TO_SPAWN=${1:-10}
MESSAGES_TO_PRINT=${2:-10}
LOGGING_DURATION=${3:-30000}
TIME_TO_WAIT_AFTER_LOGGING=${4:-60}
MASTER=${5:-"localhost:8000"}



NAMESPACES+=$(kubectl get ns -l garden.sapcloud.io/role=shoot | awk '{if(NR>1)print $1}')

for namespace in $NAMESPACES; do
    TEST_REQUEST="{\"test_name\":\"load-test-${namespace}\",\"pods\":${PODS_TO_SPAWN},\"elastic_api\":\"elasticsearch-logging.${namespace}:9200\",\"msgcount\":${MESSAGES_TO_PRINT},\"logtime_ms\":${LOGGING_DURATION},\"time_to_wait_after_logging_sec\":${TIME_TO_WAIT_AFTER_LOGGING},\"namespace\":\"${namespace}\"}"
    curl -X POST -d $TEST_REQUEST -H "Content-type: application/json" http://${MASTER}/test/start
done

TEST_REQUEST="{\"test_name\":\"load-test-garden\",\"pods\":${PODS_TO_SPAWN},\"elastic_api\":\"elasticsearch-logging.garden:9200\",\"msgcount\":${MESSAGES_TO_PRINT},\"logtime_ms\":${LOGGING_DURATION},\"time_to_wait_after_logging_sec\":${TIME_TO_WAIT_AFTER_LOGGING},\"namespace\":\"garden\"}"
curl -X POST -d $TEST_REQUEST -H "Content-type: application/json" http://${MASTER}/test/start
