package p2p

import (
	"BitTorrentClient/peers"
	"BitTorrentClient/torrentClient"
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"time"
)

const MaxBlockSize = 16 * (1 << 10)
const MaxBacklog = 5 //并发度控制
type pieceWork struct {
	//下载任务
	index  int
	hash   [20]byte
	length int
}

type pieceResult struct {
	index int
	buf   []byte
}

type pieceProgress struct {
	//下载过程描述
	index  int
	client *torrentClient.Client
	// 当前分片的 大小
	buf        []byte
	downloaded int
	requested  int
	backlog    int
}

func (p *pieceProgress) readMessage() error {
	msg, err := p.client.ReadMessage()
	if err != nil {
		return err
	}
	if msg == nil { // keep-alive
		return nil
	}
	switch msg.ID {
	case peers.MsgUnchoke:
		p.client.Choked = false
	case peers.MsgChoke:
		p.client.Choked = true
	case peers.MsgHave:
		index, err := peers.ParseHave(msg)
		if err != nil {
			p.client.Bitfield.SetPiece(index)
		}
	case peers.MsgPiece:
		n, err := peers.ParsePiece(p.index, p.buf, msg)
		if err != nil {
			return nil
		}
		p.downloaded += n
		p.backlog--
	}
	return nil
}

func (t *Torrent) Download([]byte, error) {
	log.Println("Starting download: Name ", t.Name)
	//workQueue := make(chan *pieceWork, len(t.PieceHashes))
	//results := make(chan *pieceResult)
	//for i, hash := range t.PieceHashes {
	//	length := t.
	//}

}
func attemptDownloadPiece(c *torrentClient.Client, pw *pieceWork) ([]byte, error) {
	state := pieceProgress{
		index:  pw.index,
		client: c,
		buf:    make([]byte, pw.length),
	}
	_ = c.Conn.SetDeadline(time.Now().Add(30 * time.Second))
	defer func() {
		_ = c.Conn.SetDeadline(time.Time{}) // Disable the deadline
	}()
	for state.downloaded < pw.length {
		fmt.Println(state.index, state.downloaded)
		if !state.client.Choked {
			//
			for state.backlog < MaxBacklog && state.requested < pw.length {
				blockSize := MaxBlockSize
				if pw.length-state.requested < blockSize {
					blockSize = pw.length - state.requested
				}
				err := c.SendRequest(pw.index, state.requested, blockSize)
				if err != nil {
					return nil, err
				}
				state.backlog++
				state.requested += blockSize
			}
		}

		err := state.readMessage()
		if err != nil {
			return nil, err
		}
	}
	return state.buf, nil
}

//checkIntegrity 检查hash
func checkIntegrity(pw *pieceWork, buf []byte) error {
	hash := sha1.Sum(buf)
	if !bytes.Equal(hash[:], pw.hash[:]) {
		return fmt.Errorf("Index %d failed integrity check", pw.index)
	}
	return nil
}
