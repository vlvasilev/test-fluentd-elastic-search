package logging

import (
	"fmt"
	"time"

	"github.com/I330716/test-fluentd-elastic-search/pkg/types"
	"github.com/I330716/test-fluentd-elastic-search/pkg/util"
)

func GenerateLogs(filePath string, workerName string, msgCount uint64, dataSize uint64, sleepTimeMiliSeconds int) *types.LoggingState {
	start := time.Now()
	msg := util.MakeMessage(filePath, uint64(1))
	dumpedMsg, dumpedSize := generateLogs(msg, workerName, msgCount, dataSize, sleepTimeMiliSeconds)
	t := time.Now()
	elapsed := t.Sub(start)

	state := new(types.LoggingState)
	state.DumpedKBytes = dumpedSize / uint64(1000)
	state.NumberOfDumpedMessages = dumpedMsg
	state.NumberOfMessagesToDump = msgCount
	state.ElapsedTime = time.Duration.Seconds(elapsed)
	return state
}

func generateLogs(msg *types.Message, workerName string, msgCount uint64, dataSize uint64, sleepTimeMiliSeconds int) (uint64, uint64) {
	currentDataSize := uint64(0)
	for index := uint64(1); index <= msgCount; index++ {
		msg.Number = index
		records := msg.ToRecords(workerName)
		if currentDataSize < dataSize {
			for _, record := range records {
				fmt.Println(record.ToJSON())
				sleepTime := sleepTimeMiliSeconds * int(time.Millisecond)
				time.Sleep(time.Duration(sleepTime))
				currentDataSize += uint64(len(record.ToJSON()))
			}
		} else {
			break
		}
	}
	return msg.Number, currentDataSize
}

func GenerateLogsWithMsg(msg *types.Message, workerName string, msgCount uint64, dataSize uint64, sleepTimeMiliSeconds int) *types.LoggingState {
	start := time.Now()
	dumpedMsg, dumpedSize := generateLogs(msg, workerName, msgCount, dataSize, sleepTimeMiliSeconds)
	t := time.Now()
	elapsed := t.Sub(start)

	state := new(types.LoggingState)
	state.DumpedKBytes = dumpedSize / uint64(1000)
	state.NumberOfDumpedMessages = dumpedMsg
	state.NumberOfMessagesToDump = msgCount
	state.ElapsedTime = time.Duration.Seconds(elapsed)
	return state
}
