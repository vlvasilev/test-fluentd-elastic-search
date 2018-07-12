#!/bin/bash
#start a test
curl -X POST -d '{"test_name":"test001","pods":10,"elastic_api":"elasticsearch-logging.garden:9200","msgcount":10,"logtime_ms":3000,"time_to_wait_after_logging_sec":120,"namespace":"garden"}' -H "Content-type: application/json" http://localhost:8000/test/start
#enrol a worker
curl -X POST -d '{"worker":"bot666","test_name":"test001","time_waited_after_logging_sec":2}' -H "Content-type: application/json" http://localhost:8000/enrol
#send logging status
curl -X POST -d"{\"number_of_dumped_messages\":5,\"number_of_messages_to_dump\":5,\"dumped_size_in_kilo_bytes\":19,\"elapsed_time\":1.061350535}" -H "Content-type: application/json" http://localhost:8000/logging_status?worker=bot666&&test=test001
#send analysing status
curl -X POST -d '{"number_of_sentences":22,"number_of_messages":5,"number_of_read_messages":455,"starting_message":187741,"ending_message":188195,"error":true,"report":"Different number of read messages. Read:  455  Need:  5\nMessage:  187741  has different count of sentences!\nThe count is  14  and must be  22 .\nStarts from  9  and ends at  22 .\nMessage:  188195  has different count of sentences!\nThe count is  20  and must be  22 .\nStarts from  1  and ends at  20 .\n"}' -H "Content-type: application/json" http://localhost:8000/analyse_status?worker=bot666\&test=test001

#retrieve record from elastic search cluster 
curl -X POST -H "Content-type: application/json" -d '{"query":{"term":{"worker":"test001"}},"_source":["worker","message","sentence","text"],"from":0,"size":10000}' http://localhost:9200/logstash-2018.07.12/_search?pretty=true 