package peers

import (
	"fmt"
	"io"
)

//Handshake peers 通信握手
type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

func NewHandshake(info, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: info,
		PeerID:   peerID,
	}
}
func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+20+20+8+1)
	buf[0] = byte(len(h.Pstr))
	curr := 1
	// 协议名
	curr += copy(buf[curr:], h.Pstr)
	// 8 个 0
	curr += copy(buf[curr:], make([]byte, 8)) // 8 reserved bytes
	// info hash
	curr += copy(buf[curr:], h.InfoHash[:])
	// peerId
	curr += copy(buf[curr:], h.PeerID[:])
	return buf
}
func HandshakeRead(r io.Reader) (*Handshake, error) {
	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}

	PstrLen := int(lengthBuf[0])
	if PstrLen == 0 {
		return nil, fmt.Errorf("pstrlen cannot be 0")
	}

	// 根据协议 应该要有 48+19 个字节
	HandshakeBuf := make([]byte, PstrLen+48)
	_, err = io.ReadFull(r, HandshakeBuf)
	if err != nil {
		return nil, err
	}

	var infoHash, peerID [20]byte

	copy(infoHash[:], HandshakeBuf[PstrLen+8:PstrLen+8+20])
	copy(peerID[:], HandshakeBuf[PstrLen+8+20:])

	h := Handshake{
		Pstr:     string(HandshakeBuf[0:PstrLen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	return &h, nil
}
