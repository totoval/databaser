package internal

type Connector interface {
	SetConfiger(configer Configer)
	Connector() Connector
	Client() (clientPtr interface{})
	Connect() error
	Disconnect() error
}
