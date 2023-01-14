package mi

import (
	"bufio"
	"net"
)

//Client is used to connect to OpenVPN Management Interface
type Client struct {
	MINetwork string
	MIAddress string
}

//NewClient initializes Management Interface client structure
func NewClient(network, address string) *Client {
	c := &Client{
		MINetwork: network, //Management Interface Network
		MIAddress: address, //Management Interface Address
	}

	return c
}

//GetPid returns process id of OpenVPN server
func (c *Client) GetPid() (int64, error) {
	str, err := c.Execute("pid")
	if err != nil {
		return -1, err
	}
	return ParsePid(str)
}

//GetVersion returns version of OpenVPN server
func (c *Client) GetVersion() (*Version, error) {
	str, err := c.Execute("version")
	if err != nil {
		return nil, err
	}
	return ParseVersion(str)
}

//GetStatus returns list of connected clients and routing table
func (c *Client) GetStatus() (*Status, error) {
	str, err := c.Execute("status 2")
	if err != nil {
		return nil, err
	}
	return ParseStatus(str)
}

//GetLoadStats returns number of connected clients and total number of network traffic
func (c *Client) GetLoadStats() (*LoadStats, error) {
	str, err := c.Execute("load-stats")
	if err != nil {
		return nil, err
	}
	return ParseStats(str)
}

//KillSession kills OpenVPN connection
func (c *Client) KillSession(cname string) (string, error) {
	str, err := c.Execute("kill " + cname)
	if err != nil {
		return "", err
	}
	return ParseKillSession(str)
}

//Signal sends signal to daemon
func (c *Client) Signal(signal string) error {
	str, err := c.Execute("signal " + signal)
	if err != nil {
		return err
	}
	return ParseSignal(str)
}

//Execute connects to the OpenVPN server, sends command and reads response
func (c *Client) Execute(cmd string) (string, error) {
	conn, err := net.Dial(c.MINetwork, c.MIAddress)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	buf := bufio.NewReader(conn)
	buf.ReadString('\n') //read welcome message
	err = SendCommand(conn, cmd)
	if err != nil {
		return "", err
	}

	return ReadResponse(buf)
}
