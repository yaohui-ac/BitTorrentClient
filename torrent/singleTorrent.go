package torrent

import (
	"BitTorrentClient/ben_code"
	"BitTorrentClient/p2p"
	"io"
)

//SingleInfo  单文件Torrent文件信息
type SingleInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
	// 文件发布者
	Publisher string `bencode:"publisher,omitempty"`
	// 文件发布者的网址
	PublisherUrl     string `bencode:"publisher-url,omitempty"`
	NameUtf8         string `bencode:"name.utf-8,omitempty"`
	PublisherUtf8    string `bencode:"publisher.utf-8,omitempty"`
	PublisherUrlUtf8 string `bencode:"publisher-url.utf-8,omitempty"`
	MD5Sum           string `bencode:"md5sum,omitempty"`
	Private          bool   `bencode:"private,omitempty"`
}
type SingleTorrent struct {
	// `bencode:""`
	// tracker服务器的URL 字符串
	Announce string `bencode:"announce"`
	// 备用tracker服务器列表 列表
	// 发现 announce-list 后面跟了两个l(ll) announce-listll
	AnnounceList [][]string `bencode:"announce-list"`
	// 种子的创建时间 整数
	CreatDate int64 `bencode:"creation date"`
	// 备注 字符串
	Comment string `bencode:"comment"`
	// 创建者 字符串
	CreatedBy string     `bencode:"created by"`
	Info      SingleInfo `bencode:"info"`
	// 包含一系列ip和相应端口的列表，是用于连接DHT初始node
	Nodes [][]interface{} `bencode:"nodes"`
	// 文件的默认编码
	Encoding string `bencode:"encoding"`
	// 备注的utf-8编码
	CommentUtf8 string `bencode:"comment.utf-8"`
}

func SingleTorrentParser(reader io.Reader) (s *SingleTorrent, err error) {
	s = new(SingleTorrent)
	err = ben_code.Unmarshal(reader, s)
	return s, err
}

func (s *SingleTorrent) toTorrentFile() (TorrentFile, error) {
	infoHash, err := hash(s.Info)
	if err != nil {
		return TorrentFile{}, err
	}
	pieceHashs, err := splitPieceHashes(s.Info.Pieces)
	torrentFile := TorrentFile{
		Announce: s.Announce,
		Torrent: p2p.Torrent{
			InfoHash:    infoHash,
			PieceHashes: pieceHashs,
			PieceLength: s.Info.PieceLength,
			Length:      s.Info.Length,
			Name:        s.Info.Name,
		},
	}

	torrentFile.AnnounceList = []string{}
	for _, v := range s.AnnounceList {
		torrentFile.AnnounceList = append(torrentFile.AnnounceList, v[0])
	}
	return torrentFile, nil

}
func (s *SingleTorrent) writeFile() {
	//none
}
