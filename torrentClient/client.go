package torrentClient

import (
	"BitTorrentClient/peers"
	"BitTorrentClient/util"
	"bytes"
	"fmt"
	"net"
	"time"
)

type Client struct {
	//对一个peer的连接
	Conn     net.Conn
	Choked   bool
	Bitfield util.Bitfield
	peer     peers.Peer
	//
	infoHash [20]byte
	peerID   [20]byte
}

func FinishHandshake(conn net.Conn, infoHash, peerID [20]byte) error {
	_ = conn.SetDeadline(time.Now().Add(time.Second * 3))
	defer func() {
		_ = conn.SetDeadline(time.Time{}) //取消
	}()
	req := peers.NewHandshake(infoHash, peerID)
	_, err := conn.Write(req.Serialize())
	if err != nil {
		return err
	}
	resp, err := peers.HandshakeRead(conn)
	if err != nil {
		return err
	}
	if !bytes.Equal(resp.InfoHash[:], infoHash[:]) {
		return fmt.Errorf("Expected infohash %x but got %x", resp.InfoHash, infoHash)
	}
	return nil
}
func GetBitField(conn net.Conn) (util.Bitfield, error) {
	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer func() {
		_ = conn.SetDeadline(time.Time{}) //取消
	}()

	msg, err := peers.ReadMessage(conn)

	if err != nil {
		return nil, err
	}
	if msg == nil || msg.ID != peers.MsgBitfield {
		err := fmt.Errorf("Expected bitfield but got %s", msg.String())
		return nil, err
	}
	return msg.Payload, nil
}
func NewTorrentClient(peer peers.Peer, peerID, infoHash [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), time.Second*4)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	err = FinishHandshake(conn, infoHash, peerID)
	if err != nil {
		_ = conn.Close() //不能defer
		return nil, err
	}
	bitMap, err := GetBitField(conn)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return &Client{
		Conn:     conn,
		Choked:   true,
		Bitfield: bitMap,
		peer:     peer,
		infoHash: infoHash,
		peerID:   peerID,
	}, nil
}

//ReadMessage 从连接读取并使用一条消息
func (c *Client) ReadMessage() (*peers.Message, error) {
	msg, err := peers.ReadMessage(c.Conn)
	return msg, err
}
func (c *Client) SendRequest(index, begin, length int) error {
	req := peers.FormatRequest(index, begin, length)
	_, err := c.Conn.Write(req.Serialize())
	return err
}
func (c *Client) SendInterested() error {
	msg := peers.Message{
		ID: peers.MsgInterested,
	}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendNotInterested sends a NotInterested message to the peer
func (c *Client) SendNotInterested() error {
	msg := peers.Message{ID: peers.MsgNotInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendUnchoke sends an Unchoke message to the peer
// 向对方发送取消锁定消息
func (c *Client) SendUnchoke() error {
	msg := peers.Message{ID: peers.MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendHave sends a Have message to the peer
func (c *Client) SendHave(index int) error {
	msg := peers.FormatHave(index)
	_, err := c.Conn.Write(msg.Serialize())
	return err
}