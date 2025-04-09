package pidgypost

type Metadata struct {
	SentFromClient bool
}

type Message struct {
	contents string
	metadata Metadata
}

func NewSentMessage(contents string) *Message {
	return &Message{contents: contents, metadata: Metadata{true}}
}

func NewReceivedMessage(contents string) *Message {
	return &Message{contents: contents, metadata: Metadata{false}}
}

func (message *Message) Contents() (contents string) {
	return message.contents
}
