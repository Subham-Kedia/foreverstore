package message

type Message struct {
	From    string
	Payload any
}

type DataMessage struct {
	Data []byte
	Key  string
}
