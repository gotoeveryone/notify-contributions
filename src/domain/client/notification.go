package client

type Notification interface {
	Exec(message string) error
}
