package types

type Record struct {
	Worker   string `json:"worker"`
	Message  uint64 `json:"message"`
	Sentence int    `json:"sentence"`
	Text     string `json:"text"`
}

type Sentence struct {
	Number int
	Text   string
}

type Message struct {
	Number    uint64
	Sentences []Sentence
}

type Records []Record

type Messages []Message

type LoggingState struct {
	NumberOfDumpedMessages uint64  `json:"number_of_dumped_messages"`
	NumberOfMessagesToDump uint64  `json:"number_of_messages_to_dump"`
	DumpedKBytes           uint64  `json:"dumped_size_in_kilo_bytes"`
	ElapsedTime            float64 `json:"elapsed_time"`
}

type AnalyseState struct {
	NumberOfSentences    uint   `json:"number_of_sentences"`
	NumberOfMessages     uint   `json:"number_of_messages"`
	NumberOfReadMessages uint   `json:"number_of_read_messages"`
	StartingMessage      uint   `json:"starting_message"`
	EndingMessage        uint   `json:"ending_message"`
	Error                bool   `json:"error"`
	Report               string `json:"report"`
}
