package p2p

import (
	"BitTorrentClient/peers"
	"BitTorrentClient/torrentClient"
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"runtime"
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
func (t *Torrent) calculateBoundsForPiece(index int) (begin int, end int) {
	begin = index * t.PieceLength
	end = begin + t.PieceLength
	if end > t.Length {
		end = t.Length
	}
	return begin, end
}
func (t *Torrent) caculatePieceSize(index int) int {
	begin, end := t.calculateBoundsForPiece(index)
	return end - begin
}
func (t *Torrent) startDownloadWorker(peer peers.Peer, WorkQueue chan *pieceWork, result chan *pieceResult) {
	c, err := torrentClient.NewTorrentClient(peer, t.PeerID, t.InfoHash)
	if err != nil {
		log.Printf("Could not handshake with %s. Disconnecting\n", peer.IP)
		return
	}
	defer func() {
		_ = c.Conn.Close()
	}()
	log.Printf("Completed handshake with %s\n", peer.IP)
	_ = c.SendUnchoke()
	_ = c.SendInterested()

	for pw := range WorkQueue {
		if !c.Bitfield.HasPiece(pw.index) {
			//下载完成 重新写入
			WorkQueue <- pw //
			continue
		}
		buf, err := attemptDownloadPiece(c, pw)
		if err != nil {
			log.Println("Exiting", err)
			// 出错了，把任务放回队列
			WorkQueue <- pw // Put piece back on the queue
			return
		}
		// 校验
		err = checkIntegrity(pw, buf)
		if err != nil {
			log.Printf("Piece #%d failed integrity check\n", pw.index)
			WorkQueue <- pw // Put piece back on the queue
			continue
		}
		//告知对端
		_ = c.SendHave(pw.index)
		result <- &pieceResult{pw.index, buf}
	}
}
func (t *Torrent) Download() ([]byte, error) {
	log.Println("Starting download: Name ", t.Name)
	workQueue := make(chan *pieceWork, len(t.PieceHashes))
	results := make(chan *pieceResult)
	for i, hash := range t.PieceHashes {
		length := t.caculatePieceSize(i)
		//生产写入
		workQueue <- &pieceWork{i, hash, length}
	}
	for _, peer := range t.Peers {
		go t.startDownloadWorker(peer, workQueue, results)
	}

	buf := make([]byte, t.Length)
	donePieces := 0
	for donePieces < len(t.PieceHashes) {
		res := <-results
		begin, end := t.calculateBoundsForPiece(res.index)
		copy(buf[begin:end], res.buf)
		donePieces++

		percent := float64(donePieces) / float64(len(t.PieceHashes)) * 100
		numWorkers := runtime.NumGoroutine() - 1
		log.Printf("(%0.2f%%) Downloaded piece #%d from %d peers\n", percent, res.index, numWorkers)

	}
	close(workQueue)
	return buf, nil

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
