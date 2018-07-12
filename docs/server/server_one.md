# SERVER_ONE

The server_one expose REST API to start one or more test and then retrieve the result.

## Prerequisites

We assume that the user of the test are using kubernetes. So if the server is deployed it must have the right
permissions to create, delete jobs, pods. If the server is run on the local machine one must provide a valid
kubeconfig.ymal file to the server.

## Compiling the server_one

 `make build-server-one` or `make build`

 This will make static binary under /bin/server_one/ directory

## Starting the server

The server needs tree parameters

- `--address`  The IP address on which the server will listen for requests.

- `--port`  No need for explanation

- `--kubeconfig` The path to a valid kubeconfig file. Defaults to `./conf/kubeconfig.yaml`

> NOTE: If the server is run on local machine the pods will receive for an address to send their data

>       `--address:--port`, so they may not have connection with their master.

## Start a Test

To start a test one will need the following parameters:

- `test_name` This is the name of the test. There can not be more than one test with the same name.The name later will be used for retrieving the results.

- `pods`      The number of pods which are going to dump messages.

- `elastic_api`  The entry point of the elastic search cluster (e.g. myelastic-search.com:9200).

- `msgcount`  The number of full text(tales) which are going to be dumped. This content of this tale is it /resources/tex.tx. If one wants to change the text just have to replace the content. The text is separated by `'.'` `'!'` and `'?'`.Because elastic search has beeb limited to 10000 hits on normal search the maximum number must be no more than 454.The default text has 22 sentences. 454 x 22 = 9988 lines which will be stored as same number records in the elastic.

- `logtime_ms`  The time duration in ms in which the given number of messages will be dumped.

              For example, if you set the value to 30000 this will dump all the messages denoted by msgcount for 30 seconds.

- `time_to_wait_after_logging_sec`  The time which the worker will is going to sleep after the logging is done. This sleep duration is needed because at the time when the logging is all of some logs may not be yet in the elastic search cluster due to slow work of the fluent(d/bit)

- `namespace`  The namespace of the pods. If the namespace does not exist one must create it first.


All of the arguments must be passed as JSON.
Example:
`'{"test_name":"test001","pods":10,"elastic_api":"elasticsearch-logging.garden:9200","msgcount":10,"logtime_ms":3000,"time_to_wait_after_logging_sec":120,"namespace":"garden"}'`

This data must be sent as POST request body to:

`http://host:port/test/start`

Example:

`curl -X POST -d '{"test_name":"test001","pods":10,"elastic_api":"elasticsearch-logging.garden:9200","msgcount":10,"logtime_ms":3000,"time_to_wait_after_logging_sec":120,"namespace":"garden"}' -H "Content-type: application/json" http://host:port/test/start`

> Do not forget to specify the header "Content-type: application/json"

On succeed one has to see the following message:

```
The test is running!

Job test001 created!
```

Now server_one will spawn 10 flooding pods. They will gather they're logs and will analyze them. 
After that the pods will send the result to server_one. Then the server_one will delete the pods.

## Getting the result

All the results are returned as JSON format.

### The result for all tests
GET request to `http://host:port/status/all`

### The result for a single test

#### Long version
This will return the test status plus the status of all the workers(pods)
GET request to `http://host:port/status/long?test=`
#### Short version
This will return the base status of the test but not the reason of failure
GET request to `http://host:port/status/short?test=test001`
#### Normal version
This will return the base status plus the status only of those workers which 
have found an error during the analyzing.
GET request to `http://host:port/status/normal?test=test001`

## Destroy test
The test status will remain on the server until it is destroyed.
To destroy the test one must send POST request to:
`http://host:port/test/destroy?test=<test name>`