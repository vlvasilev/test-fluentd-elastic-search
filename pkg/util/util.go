package util

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/I330716/test-fluentd-elastic-search/pkg/types"
)

func ReadFile(filePath string) string {
	data, err := ioutil.ReadFile(filePath)
	check(err)
	return string(data)
}

func WriteToFile(data, filePath string) {
	file, err := os.Create(filePath)
	defer file.Close()
	check(err)
	file.WriteString(data)
	file.Sync()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func separeteSentences(text string) []string {
	return strings.FieldsFunc(text, func(r rune) bool {
		if r == '.' || r == '!' || r == '?' {
			return true
		}
		return false
	})
}

func GetPodLastHash(pod string) string {
	tokens := strings.FieldsFunc(pod, func(r rune) bool {
		return r == '-'
	})
	if len(tokens) > 1 {
		return tokens[len(tokens)-1]
	}
	return pod
}

func trimStrings(unescapedStrings []string) {
	for index, s := range unescapedStrings {
		unescapedStrings[index] = strings.TrimSpace(s)
	}
}

func replaceNewLineInStrings(strs []string) {
	for index, s := range strs {
		strs[index] = strings.Replace(s, "\n", " ", -1)
	}
}

func MakeMessage(filePath string, number uint64) *types.Message {
	text := ReadFile(filePath)
	sentences := separeteSentences(text)
	trimStrings(sentences)
	replaceNewLineInStrings(sentences)

	msg := &types.Message{
		Number:    number,
		Sentences: make([]types.Sentence, len(sentences)),
	}

	for index, sentence := range sentences {
		msg.Sentences[index] = types.Sentence{
			Number: (index + 1),
			Text:   sentence,
		}
	}

	return msg
}

func GetMessageSize(worker string, msg *types.Message) uint64 {
	var size uint64
	records := msg.ToRecords(worker)
	for _, record := range records {
		size += uint64(record.JsonSize())
	}
	return size
}

func Exe_cmd(cmd string) ([]byte, error) {
	//fmt.Println("command is ", cmd)
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		return []byte{}, err
	}
	//wg.Done() // Need to signal to waitgroup that this goroutine is done
	return out, nil
}

func GetLogstashIndex() string {
	current_time := time.Now().Local()
	return string(current_time.AppendFormat([]byte(`logstash-`), "2006-01-02"))
}
