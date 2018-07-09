package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/I330716/test-fluentd-elastic-search/pkg/operations/analise"
)

func main() {
	runtime.GOMAXPROCS(1)
	var logFile = flag.String("logsfile", "logsfile.json", "The file in which the records from elastic_serach are stored in JSON format")
	var msgFile = flag.String("examplemsg", "text.txt", "The file in which the sent message is stored which is for comparsion")
	var msgCount = flag.Uint("msgcount", 0, "The number of messages which are going to be recieved.")

	flag.Parse()
	state := analise.Analyse(logFile, msgFile, *msgCount)
	fmt.Println(state.ToJSON())
}
