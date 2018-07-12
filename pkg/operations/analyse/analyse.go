package analise

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/I330716/test-fluentd-elastic-search/pkg/types"
	"github.com/I330716/test-fluentd-elastic-search/pkg/util"
)

func getLogs(filePath string) []byte {
	data, err := ioutil.ReadFile(filePath)
	check(err)
	return data
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetRecords(logFile string) types.Records {
	rawLogs := getLogs(logFile)
	records := types.Records{}
	err := json.Unmarshal(rawLogs, &records)
	check(err)
	return records
}

func GetRecordsFromJSON(jsonData []byte) (*types.Records, error) {
	records := new(types.Records)
	if len(jsonData) < 1 {
		return nil, errors.New("empty data to unmarshal in GetRecordsFromJSON")
	}
	err := json.Unmarshal(jsonData, records)
	if err != nil {
		log.Printf(string(jsonData))
		return nil, errors.New(err.Error() + " : " + string(jsonData))
		//check(err)
	}

	return records, nil
}

func doesMessageLenghtIsCorrenct(msg *types.Message, maxSentences int) (bool, string) {
	if len(msg.Sentences) != maxSentences {
		var buffer bytes.Buffer
		currentMessageSize := len(msg.Sentences)
		buffer.WriteString(fmt.Sprintln("Message: ", msg.Number, " has different count of sentences!"))
		buffer.WriteString(fmt.Sprintln("The count is ", currentMessageSize, " and must be ", maxSentences, "."))
		if len(msg.Sentences) > 0 {
			buffer.WriteString(fmt.Sprintln("Starts from ", msg.Sentences[0].Number, " and ends at ", msg.Sentences[currentMessageSize-1].Number, "."))
			sentenceNumbers := make([]byte, maxSentences)
			for _, sentence := range msg.Sentences {
				if sentence.Number > 0 && sentence.Number <= maxSentences {
					sentenceNumbers[sentence.Number-1] = 1
				}
			}
			buffer.WriteString("\nMissing: ")
			var separator string
			for index, value := range sentenceNumbers {
				if value == 0 {
					buffer.WriteString(separator + strconv.Itoa(index+1))
					separator = ","
				}
			}
		}
		return false, buffer.String()
	}
	return true, ""
}

func doMessagesHaveTheSameContent(origin, other *types.Message) (bool, string) {
	var buffer bytes.Buffer
	var isMistmatch bool
	buffer.WriteString(fmt.Sprintln("Message: ", other.Number, " has different content:"))
	for index, originSentence := range origin.Sentences {
		if originSentence != other.Sentences[index] {
			isMistmatch = true
			buffer.WriteString(fmt.Sprintln(index, " > ", originSentence))
			buffer.WriteString(fmt.Sprintln(index, " < ", other.Sentences[index]))
		}
	}

	if isMistmatch {
		return false, buffer.String()
	}
	return true, ""
}

func analyzeMessages(msg1, msg2 *types.Message) string {
	result, difference := doesMessageLenghtIsCorrenct(msg2, len(msg1.Sentences))
	if !result {
		return difference
	}

	result, difference = doMessagesHaveTheSameContent(msg1, msg2)
	if !result {
		return difference
	}

	return "No difference"
}

func Analyze(originMessage *types.Message, messages *types.Messages, msgCount uint) *types.AnalyseState {
	state := &types.AnalyseState{}
	var buffer bytes.Buffer
	state.NumberOfSentences = uint(len(originMessage.Sentences))
	state.NumberOfMessages = uint(msgCount)
	state.NumberOfReadMessages = uint(len(*messages))
	state.StartingMessage = uint((*messages)[0].Number)
	state.EndingMessage = uint((*messages)[state.NumberOfReadMessages-1].Number)

	if state.NumberOfMessages != state.NumberOfReadMessages {
		state.Error = true
		buffer.WriteString(fmt.Sprintln("Different number of read messages. Read: ", state.NumberOfReadMessages, " Need: ", state.NumberOfMessages))
	}

	for _, currMsg := range *messages {
		if !originMessage.IsEqual(currMsg) {
			state.Error = true
			difference := analyzeMessages(originMessage, &currMsg)
			buffer.WriteString(difference)
		}
	}

	state.Report = buffer.String()
	return state

}

func Analyse(logFile *string, msgFile *string, msgCount uint) *types.AnalyseState {
	records := GetRecords(*logFile)
	messages := records.ToMessages()
	originMessage := util.MakeMessage(*msgFile, 1)
	return Analyze(originMessage, messages, msgCount)
}
