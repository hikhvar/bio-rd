package packet

import (
	"bytes"

	bnet "github.com/bio-routing/bio-rd/net"
	"github.com/taktv6/tflow2/convert"
)

// MultiProtocolUnreachNLRI represents network layer withdraw information for one prefix of an IP address family (rfc4760)
type MultiProtocolUnreachNLRI struct {
	AFI      uint16
	SAFI     uint8
	Prefixes []bnet.Prefix
}

func (n *MultiProtocolUnreachNLRI) serialize(buf *bytes.Buffer) uint8 {
	tempBuf := bytes.NewBuffer(nil)
	tempBuf.Write(convert.Uint16Byte(n.AFI))
	tempBuf.WriteByte(n.SAFI)
	for _, pfx := range n.Prefixes {
		tempBuf.Write(serializePrefix(pfx))
	}

	buf.Write(tempBuf.Bytes())

	return uint8(tempBuf.Len())
}

func deserializeMultiProtocolUnreachNLRI(b []byte) (MultiProtocolUnreachNLRI, error) {
	n := MultiProtocolUnreachNLRI{}
	prefix := make([]byte, len(b)-3)

	fields := []interface{}{
		&n.AFI,
		&n.SAFI,
		&prefix,
	}
	err := decode(bytes.NewBuffer(b), fields)
	if err != nil {
		return MultiProtocolUnreachNLRI{}, err
	}

	idx := uint8(0)
	n.Prefixes = make([]bnet.Prefix, 0)
	for idx < uint8(len(prefix)) {
		l := numberOfBytesForPrefixLength(prefix[idx])

		pfx, err := deserializePrefix(prefix[idx+1:idx+1+l], prefix[idx], n.AFI)
		if err != nil {
			return MultiProtocolUnreachNLRI{}, err
		}
		n.Prefixes = append(n.Prefixes, pfx)

		idx = idx + l + 1
	}

	return n, nil
}
