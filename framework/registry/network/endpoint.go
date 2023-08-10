package network

type Endpoint interface {
	Name() string
	Address() string
	MustRegisterNetwork() bool
}
