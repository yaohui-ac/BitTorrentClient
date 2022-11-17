package p2p

import "BitTorrentClient/peers"

// MaxBlockSize is the largest number of bytes a request can ask for
// 一个块 通常为16KB

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
