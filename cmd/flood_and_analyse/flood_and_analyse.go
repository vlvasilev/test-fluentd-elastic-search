package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/I330716/test-fluentd-elastic-search/pkg/elastic"
	"github.com/I330716/test-fluentd-elastic-search/pkg/master"
	"github.com/I330716/test-fluentd-elastic-search/pkg/operations/analise"
	"github.com/I330716/test-fluentd-elastic-search/pkg/operations/logging"
	"github.com/I330716/test-fluentd-elastic-search/pkg/types"
	"github.com/I330716/test-fluentd-elastic-search/pkg/util"
)

func generateLogs(workerName *string, msg *types.Message, msgCount uint64, logSize uint64, loggingTime uint64) *types.LoggingState {
	var sleepTimeMs uint64
	if msgCount == 0 {
		msgCount = logSize / util.GetMessageSize(*workerName, msg)
	}
	sleepTimeMs = loggingTime / (msgCount * uint64(len(msg.Sentences)))
	loggingState := logging.GenerateLogsWithMsg(msg, *workerName, msgCount, logSize, int(sleepTimeMs))
	return loggingState
}

func getExtraxtedMessages(extactedDataJSON *[]byte) (*types.Messages, error) {
	records := analise.GetRecordsFromJSON(extactedDataJSON)
	messages := records.ToMessages()
	return messages, nil
}

func getErrorSetAnalyseStatus(err error) *types.AnalyseState {
	status := new(types.AnalyseState)
	status = new(types.AnalyseState)
	status.NumberOfMessages = 0
	status.NumberOfReadMessages = 0
	status.NumberOfSentences = 0
	status.StartingMessage = uint(0)
	status.EndingMessage = 0
	status.Report = err.Error()
	status.Error = true
	return status
}

func main() {

	var textFile = flag.String("textfile", "text.txt", "The path where is the text which is going to be dumped")
	var loggingTime = flag.Uint64("logtime", 0, "The time duration(ms) in which the log are going to be dumped")
	var msgCount = flag.Uint64("msgcount", 0, "The number of full texts which are going to be dumped in stdout")
	var logSize = flag.Uint64("logsize", 18446744073709551615, "The size of logs in bytes which will be printed in the stdout")
	var workerName = flag.String("workername", "unknown", "The name of the worker which will be posted in the log message.")
	var timeToWaitAfterLogging = flag.Uint("time_to_wait_after_logging_sec", 120, "The time which the program is going to sleep after the logging is done")
	var elasticSearchAPI = flag.String("elastic_end_point", "localhost:9200", "The address and port of the elastic search API")
	var masterEndPoint = flag.String("master", "localhost:33661", "The master api where the worker will send data")
	flag.Parse()

	var loggingState *types.LoggingState
	var analyseState *types.AnalyseState
	var elasticServer elastic.ElasticSearchClient
	var masterClient master.MasterClient
	var errorString string

	elasticServer.Init(*elasticSearchAPI)
	masterClient.Init(*masterEndPoint, *workerName)

	msg := util.MakeMessage(*textFile, 0)
	//enrol to the master server
	err := masterClient.Enrol(*timeToWaitAfterLogging)
	if err != nil {
		errorString += err.Error() + "\n"
	}
	//generate some logs
	loggingState = generateLogs(workerName, msg, *msgCount, *logSize, *loggingTime)
	//sent the loging status to the master server
	err = masterClient.SendLoggingStatus(loggingState)
	if err != nil {
		errorString += err.Error() + "\n"
	}
	//wait some time until the fluentd gathers all of the logs and send in to elasticsearch cluster
	time.Sleep(time.Duration(*timeToWaitAfterLogging * uint(1000) * uint(time.Millisecond)))
	//get the log records from the elasticsearch cluster
	data, err := elasticServer.GetCurrenLogstashRecords("fluentd", "worker", util.GetPodLastHash(*workerName))
	if err != nil {
		analyseState = getErrorSetAnalyseStatus(err)
	} else {
		messages, err := getExtraxtedMessages(&data)
		data = nil
		if err != nil {
			analyseState = getErrorSetAnalyseStatus(err)
		} else {
			analyseState = analise.Analyze(msg, messages, uint(loggingState.NumberOfDumpedMessages))
		}
	}
	//now send the analyse satatus to the server
	masterClient.SendAnalyseStatus(analyseState)
	if err != nil {
		errorString += err.Error() + "\n"
	}

	fmt.Println(loggingState.ToJSON())
	fmt.Println(analyseState.ToJSON())
	fmt.Println(errorString)
}
