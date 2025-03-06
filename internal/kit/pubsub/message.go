package pubsub

type Message struct {
	Data       []byte
	Attributes map[string]string
}
