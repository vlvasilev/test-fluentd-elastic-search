#!/bin/bash

while [[ 1==1 ]]; 
do
	sleep 300
	make run-test PODS=15 MESSAGES=400 LOGGING_TIME=1200000 TTWAL=240 MASTER=a66c1f69591af11e8a01612362c3d869-381879666.eu-west-1.elb.amazonaws.com:8000
	sleep 4000
	curl -XPOST http://a66c1f69591af11e8a01612362c3d869-381879666.eu-west-1.elb.amazonaws.com:8000/test/destroy_all
done
