package pidgypost

type Metadata struct {
	Sender string
}

type Message struct {
	contents string
	metadata Metadata
}

func NewMessage(contents string) *Message {
	return &Message{contents: contents, metadata: Metadata{""}}
}

func (message *Message) SetSender(sender string) {
	message.metadata.Sender = sender
}

func (message *Message) Contents() (contents string) {
	return message.contents
}

func (message *Message) Sender() (sender string) {
	return message.metadata.Sender
}
