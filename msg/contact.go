package pidgypost

type Contact struct {
	name            string
	originated_msgs []Message
	terminated_msgs []Message
}

func NewContact(name string) *Contact {
	return &Contact{name: name, terminated_msgs: make([]Message, 0),
		originated_msgs: make([]Message, 0)}
}

func (contact *Contact) AddOriginated(originated Message) {
	contact.originated_msgs = append(contact.originated_msgs, originated)
}

func (contact *Contact) AddTerminated(terminated Message) {
	contact.terminated_msgs = append(contact.terminated_msgs, terminated)
}
