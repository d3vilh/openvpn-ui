package mi

// Version .
type Version struct {
	OpenVPN    string
	Management string
}

// LoadStats .
type LoadStats struct {
	NClients int64
	BytesIn  int64
	BytesOut int64
}

// Status .
type Status struct {
	Title        string
	Time         string
	TimeT        string
	ClientList   []*OVClient
	RoutingTable []*RoutingPath
}

// OVClient .
type OVClient struct {
	CommonName      string
	RealAddress     string
	VirtualAddress  string
	VirtualIPv6     string
	BytesReceived   uint64
	BytesSent       uint64
	ConnectedSince  string
	ConnectedSinceT string
	Username        string
	ClientID        string
	PeerID          string
	// DataCipher      string
}

// RoutingPath .
type RoutingPath struct {
	VirtualAddress string
	CommonName     string
	RealAddress    string
	LastRef        string
	LastRefT       string
}
