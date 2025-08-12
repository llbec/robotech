package tcp

type Conn interface {
	Close() error
	Send(string) error
}
