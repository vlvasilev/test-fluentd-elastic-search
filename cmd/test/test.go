package main

import (
	"encoding/json"
	"fmt"

	"github.com/I330716/test-fluentd-elastic-search/pkg/operations/requests"
)

type EnrolData struct {
	WorkerName                string `json:"worker"`
	TimeToWaitAfterLoggingSec uint   `json:"time_waited_after_logging_sec"`
}

func printHTTPResponse(status string, body []byte, err error) {
	fmt.Println("Status: " + status)
	fmt.Println("Body: ", string(body))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	enrolData := EnrolData{"Snail", 60}
	json, _ := json.Marshal(&enrolData)
	status, body, err := requests.MakeJsonHTTPPOSTRequest("localhost:8000/enrol?worker=Snail", json)
	printHTTPResponse(status, body, err)
	logData := []byte("{\"number_of_dumped_messages\":5,\"number_of_messages_to_dump\":5,\"dumped_size_in_kilo_bytes\":19,\"elapsed_time\":1.061350535}")
	status, body, err = requests.MakeJsonHTTPPOSTRequest("localhost:8000/logging_status?worker=Snail", logData)
	printHTTPResponse(status, body, err)
	analyseData := []byte("{\"number_of_sentences\":22,\"number_of_messages\":5,\"number_of_read_messages\":455,\"starting_message\":187741,\"ending_message\":188195,\"error\":true,\"report\":\"Different number of read messages. Read:  455  Need:  5\\nMessage:  187741  has different count of sentences!\\nThe count is  14  and must be  22 .\\nStarts from  9  and ends at  22 .\\nMessage:  188195  has different count of sentences!\\nThe count is  20  and must be  22 .\\nStarts from  1  and ends at  20 .\\n\"}")
	status, body, err = requests.MakeJsonHTTPPOSTRequest("localhost:8000/analyse_status?worker=Snail", analyseData)
	printHTTPResponse(status, body, err)
	//bytes, _ := requests.GetCurrentDateFluentdRecords("a76a3480b761611e89fe8727c31f3fa0-977657185.eu-west-1.elb.amazonaws.com:9200", "worker", "frrds")
	//fmt.Println(string(bytes))
	// number, err := requests.GetRecordNumbers("a76a3480b761611e89fe8727c31f3fa0-977657185.eu-west-1.elb.amazonaws.com:9200", "logstash-2018.07.04", "fluentd", "worker", "frrds")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(number)
	// }

	// status, _, err := requests.MakeGetRequest("http://a76a3480b761611e89fe8727c31f3fa0-977657185.eu-west-1.elb.amazonaws.com:9200")
	// if err != nil || status != "200 OK" {
	// 	fmt.Println("False " + status)
	// } else {
	// 	fmt.Println("True")
	// }

}
