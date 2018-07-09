package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

func (s *Sentence) ToString() string {
	return fmt.Sprintf("sentence:%d %s", s.Number, s.Text)
}

func (s *Sentence) StringSize() int {
	return len(s.ToString())
}

func (s1 *Sentence) IsEqual(s2 Sentence) bool {
	if s1.Number != s2.Number {
		return false
	}
	if s1.Text != s2.Text {
		return false
	}
	return true
}

func (s *Sentence) ToJSON() string {
	json, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(json)
}

func (s *Sentence) JsonSize() int {
	return len(s.ToJSON())
}

func (m *Message) ToString() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintln("message number: %d", m.Number))
	for _, sentence := range m.Sentences {
		buffer.WriteString(fmt.Sprintln("%s", sentence.ToString()))
	}
	return fmt.Sprintf(buffer.String())
}

func (m *Message) StringSize() int {
	return len(m.ToString())
}

func (m *Message) IsEqual(m2 Message) bool {
	if len(m.Sentences) != len(m2.Sentences) {
		return false
	}
	for index := range m.Sentences {
		if m.Sentences[index] != m2.Sentences[index] {
			return false
		}
	}
	return true
}

func (m *Message) ToRecords(workerName string) []Record {
	records := make([]Record, len(m.Sentences))
	for index, sentence := range m.Sentences {
		records[index] = Record{workerName, m.Number, sentence.Number, sentence.Text}
	}
	return records
}

//AddSentence - add sentence to the collection of sentence in the message
func (m *Message) AddSentence(s Sentence) {
	if s.Number > 0 {
		m.Sentences = append(m.Sentences, s)
	}
}

func (r *Record) ToJSON() string {
	json, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(json)
}

func (r *Record) JsonSize() int {
	return len(r.ToJSON())
}

func (s *LoggingState) ToJSON() string {
	json, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(json)
}

//Len - returns the number of elements in the collection of Record objects
func (records Records) Len() int {
	return len(records)
}

//Swap - swap thwo elements defined by indices i and j
func (records Records) Swap(i, j int) {
	records[i], records[j] = records[j], records[i]
}

//Less - return true if Record defined by index i is less than this difined by j
func (records Records) Less(i, j int) bool {
	if records[i].Message != records[j].Message {
		return records[i].Message < records[j].Message
	}
	if records[i].Sentence != records[j].Sentence {
		return records[i].Sentence < records[j].Sentence
	}
	return records[i].Text < records[j].Text
}

//Sort sorts data. It makes one call to data.Len to determine n, and O(n*log(n)) calls to data.Less and data.Swap. The sort is not guaranteed to be stable.
func (records Records) Sort() {
	sort.Sort(records)
}

func isNewMessage(messageNumber *uint64, currentMessageNumber uint64) bool {
	if *messageNumber == 0 || *messageNumber != currentMessageNumber {
		*messageNumber = currentMessageNumber
		return true
	}
	return false
}

//Add add msg to the collection. If msg with no Sentence is ignored
func (messages *Messages) Add(msg Message) {
	if len(msg.Sentences) > 0 {
		*messages = append(*messages, msg)
	}
}

//ToMessages. Return Messages made from the Records
func (records Records) ToMessages() *Messages {
	var messageNumber uint64
	var msg Message
	messages := new(Messages)
	records.Sort()
	for _, record := range records {
		if isNewMessage(&messageNumber, record.Message) {
			messages.Add(msg)
			msg = Message{record.Message, []Sentence{}}
		}
		msg.AddSentence(Sentence{record.Sentence, record.Text})
	}
	messages.Add(msg)
	return messages
}

func (records Records) ToString() string {
	var buffer bytes.Buffer
	for _, record := range records {
		buffer.WriteString(fmt.Sprintln("Message: ", record.Message, " Sentence: ", record.Sentence, " Text: ", record.Text))
	}
	return fmt.Sprintf(buffer.String())
}

func (s *AnalyseState) ToJSON() string {
	json, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(json)
}
