package oscclient

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

// OscIface ... interface
type OscIface interface {
	Sendhour(hour int, button int) error
}

// OscClient ... Main Client
type OscClient struct {
	IP       string
	Port     string
	Timepath string
}

// New ... return OscClient
func New() OscIface {
	return &OscClient{
		IP:       "127.0.0.1",
		Port:     "9000",
		Timepath: "/time/hour",
	}
}

// osc padding
func (client *OscClient) pad(b []byte) []byte {
	if len(b)%4 == 0 {
		byte4 := []byte{0, 0, 0, 0}
		b = append(b, byte4...)
		return b
	}

	for i := len(b); (i % 4) != 0; i++ {
		b = append(b, 0)
	}
	return b
}

// make osc message
func (client *OscClient) makeMessage(addr []byte, types []byte, args []byte) []byte {
	var data []byte

	padaddr := client.pad(addr)
	data = append(data, padaddr...)

	padtypes := client.pad(types)
	data = append(data, padtypes...)

	data = append(data, args...)
	return data
}

// Sendhour ... send hour to vrchat
func (client *OscClient) Sendhour(hour int, button int) error {

	// parse argument
	args := new(bytes.Buffer)
	if err := binary.Write(args, binary.BigEndian, int32(button)); err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%d", client.Timepath, hour)

	data := client.makeMessage([]byte(path), []byte(",i"), args.Bytes())

	return client.Send(data)
}

// Send ... send binary message to vrchat
func (client *OscClient) Send(data []byte) error {

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", client.IP, client.Port))
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err = conn.Write(data); err != nil {
		return err
	}

	return nil
}
