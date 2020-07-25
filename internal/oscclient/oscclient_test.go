package oscclient

import (
	"bytes"
	"testing"
)

func TestPad(t *testing.T) {

	var m OscClient

	cases := []struct {
		Req string
		Res []byte
	}{
		{
			Req: "/time/hour/0",
			Res: []byte{47, 116, 105, 109, 101, 47, 104, 111, 117, 114, 47, 48, 0, 0, 0, 0},
		},
		{
			Req: "/time/hour/12",
			Res: []byte{47, 116, 105, 109, 101, 47, 104, 111, 117, 114, 47, 49, 50, 0, 0, 0},
		},
	}

	for _, c := range cases {
		reqB := []byte(c.Req)
		resp := m.pad(reqB)

		if !bytes.Equal(resp, c.Res) {
			t.Errorf("pad is failed, %+v : %+v", resp, c.Res)
		}
	}
}
