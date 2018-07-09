# test-fluentd-elastic-search
Stupid test which generates logs in the stout and then looks for them in the elastic-search cluster


This Beta test is used to test the correctness of the gathering and persisting of logs dumped on the stout in one or many pods in Kubernetes cluster.
The test assumes that there is a fluentd running on kubernetes the cluster which gathers the logs from the Nodes and sends them to elastic search cluster.



First, each worker pod dumps in the stout a text separated sentence by sentence using the following pattern:

 {"worker":"flood-and-analyse-x6j8f","message":10,"sentence":1,"text":"Once, on the African plains, there lived a moody rhino who was very easily angered"}

 {"worker":"flood-and-analyse-x6j8f","message":10,"sentence":2,"text":" One\nday, a giant turtle entered the rhino's territory unaware"}
 .
 .
 {"worker":"flood-and-analyse-x6j8f","message":10,"sentence":22,"text":" and while the\nmonkeys were putting on their sticking plasters, their chief realised it was about time they\nfound a better way to amuse themselves than making fun of others"}

The worker name is the name of the pod. Message is the number of the consecutive tale (the full text). Sentence in the consecutive sentence of the given tale. Text is the sentence content.

Fluentd gathers this logs from /var/log/container in each Node or get them from fluent-bit, thus making the following index pattern in elastic search cluster.

index: logstash-2018.07.9
_type: fluentd
_source:{
    worker: flood-and-analyse-x6j8f
    message: 10
    sentence: 2
    text: " One\nday, a giant turtle entered the rhino's territory unaware"
    log : {"worker":"flood-and-analyse-x6j8f","message":10,"sentence":2,"text":" One\nday, a giant turtle entered the rhino's territory unaware"}
}
.
.

Based on the index, _type and _source we are going to examine the stored logs.
We are going to search records in the elastic search cluster base on the worker name because it is unique.
Because of the reverse index mechanism in elastic we cannot search by full pod name because flood-and-analyse-x6j8f will collapse to tables of 
four tokens(flood, and, analyse, x6j8f). The sole name will give zero records when searching. We are going to use the hash part of the name(x6j8f).


How it works

There is central pod with the name server-one which runs a REST API.
The test is started by POST request to "http://host:port/test/start" with body consisting of parameters in JSON.
This params are:

- pods -> the number of pods which are gong to dump messages.

- elastic_api -> the entry point of the elastic search cluster (e.g. elastic-search.com:9200).

- msgcount -> the number of full text(tales) which are going to be dumped.

- logtime_ms -> the time duration in ms in which the given number of messages will be dumped.

- time_to_wait_after_logging_sec -> the time which the worker will is going to sleep after the logging is done.

- namespace -> the namespace of the pods



example:

curl -X POST -d '{"pods":10,"elastic_api":"elasticsearch-logging.garden:9200","msgcount":10,"logtime_ms":60000,"time_to_wait_after_logging_sec":120,"namespace":"garden"}' -H "Content-type: application/json" http://localhost:8000/test/start



When this request succeed the server-one pod will deploy Job in namespace "garden" with 10 pods which are going to dump 10 messages in one minute after that each pod is going to sleep for 2 minutes.
Before sleep each pod will send the data to the server-one of how much logs it have generated.
After that they will try to extract that logs from elasticsearch-logging.garden:9200/logstash-(YYYY.MM.DD)/fluentd.
Note that we are using logstash so each day the indices must named logstash-(current date).


After gathering all of the information each pod is going to analyze the logs, find if some sentences are missing.
Then it will send the analyze status to the server-one pod.


We can obtain the result by GET request to http://server-one:8000/status/all