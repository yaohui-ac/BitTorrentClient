package p2p

import "BitTorrentClient/peers"

// MaxBlockSize is the largest number of bytes a request can ask for
// 一个块 通常为16KB
const MaxBlockSize = 16384

// MaxBacklog is the number of unfulfilled requests a client can have in its pipeline
const MaxBacklog = 5

// Torrent holds data required to download a torrent from a list of peers
type Torrent struct {
	Peers       []peers.Peer
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
	File        []FileInfo
}

type FileInfo struct {
	Path   string
	Length int
}

type pieceWork struct {
	index  int
	hash   [20]byte
	length int
}

type pieceResult struct {
	index int
	buf   []byte
}

//type pieceProgress struct {
//	index  int
//	client *client.Client
//	// 当前分片的 大小
//	buf        []byte
//	downloaded int
//	requested  int
//	backlog    int
//}
