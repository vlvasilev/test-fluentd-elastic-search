package main

import (
	"flag"
	"runtime"

	"github.com/I330716/test-fluentd-elastic-search/pkg/operations/logging"
	"github.com/I330716/test-fluentd-elastic-search/pkg/util"
)

func main() {
	runtime.GOMAXPROCS(1)
	//home := os.Getenv("HOME")
	//tf := "/go/src/projects/lt/"
	defTextFile := "text.txt"
	// if home != "" {
	// 	defTextFile = home + tf + defTextFile
	// }
	var workerName = flag.String("workername", "unknown", "The name of the worker which will be posted in the log message.")
	var msgCount = flag.Uint64("msgcount", 18446744073709551615, "The number of full texts which are going to be dumped in stdout")
	var logSize = flag.Uint64("logsize", 18446744073709551615, "The size of logs in bytes which will be printed in the stdout")
	var sleepTimeMiliSeconds = flag.Int("sleep", 0, "The time interval between each loged message in ms.")
	var textFile = flag.String("textfile", defTextFile, "The path where is the text which is going to be dumped")
	flag.Parse()

	state := logging.GenerateLogs(*textFile, *workerName, *msgCount, *logSize, *sleepTimeMiliSeconds)

	util.WriteToFile(state.ToJSON(), "logging_output.txt")
}
