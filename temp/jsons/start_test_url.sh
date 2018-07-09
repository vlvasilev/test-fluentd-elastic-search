#!/bin/bash
curl -X POST -d '{"pods":10,"elastic_api":"elasticsearch-logging.garden:9200","msgcount":10,"logtime_ms":3000,"time_to_wait_after_logging_sec":120,"namespace":"garden"}' -H "Content-type: application/json" http://localhost:8000/test/start
