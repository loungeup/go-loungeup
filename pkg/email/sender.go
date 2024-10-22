package email

type Sender interface {
	Send(email *Email) error
}
