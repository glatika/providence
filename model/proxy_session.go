package model

// TODO: This will implement in later version
// note this just a placeholder

type ProxySession struct {
	Id    int
	Port  int
	Token string
	// if outbound this will proxy to stock
	Outbound bool
	// indicate the proxy just a forwarder
	Pure bool
	Host string
}
