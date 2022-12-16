package torrent

import (
	"BitTorrentClient/ben_code"
	"BitTorrentClient/p2p"
	"bytes"
)

type MultipleInfo struct {
	// 每个块的20个字节的SHA1 Hash的值(二进制格式)
	Pieces string `bencode:"pieces"`
	// 每个块的大小，单位字节 整数
	PieceLength int `bencode:"piece length"`
	// 文件长度 整数
	Length int `bencode:"length,omitempty"`
	// 目录名 字符串
	Name  string `bencode:"name"`
	Files []struct {
		// 文件长度 单位字节 整数
		Length int `bencode:"length"`
		// 文件的路径和名字 列表
		Path []string `bencode:"path"`
		// path.utf-8：文件名的UTF-8编码
		PathUtf8 string `bencode:"path.utf-8,omitempty"`
	} `bencode:"files"`
	NameUtf8 string `bencode:"name.utf-8,omitempty"`
}

//MultipleTorrent 多文件
type MultipleTorrent struct {
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
	CreatedBy string       `bencode:"created by"`
	Info      MultipleInfo `bencode:"info"`
	// 包含一系列ip和相应端口的列表，是用于连接DHT初始node
	Nodes [][]interface{} `bencode:"nodes"`
	// 文件的默认编码
	Encoding string `bencode:"encoding"`
	// 备注的utf-8编码
	CommentUtf8 string `bencode:"comment.utf-8"`
}

// multipleParser 多文件解析
func (bto *MultipleTorrent) fileParser(file []byte) error {
	err := ben_code.Unmarshal(bytes.NewReader(file), &bto)
	return err
}

func (bto *MultipleTorrent) toTorrentFile() (TorrentFile, error) {
	infoHash, err := hash(bto.Info)
	if err != nil {
		return TorrentFile{}, err
	}

	// 每个分片的 SHA-1 hash 长度是20 把他们从Pieces中切出来
	pieceHashes, err := splitPieceHashes(bto.Info.Pieces)
	if err != nil {
		return TorrentFile{}, err
	}

	tf := TorrentFile{
		Announce: bto.Announce,
		Torrent: p2p.Torrent{
			InfoHash:    infoHash,
			PieceHashes: pieceHashes,
			PieceLength: bto.Info.PieceLength,
			Length:      bto.Info.Length,
			Name:        bto.Info.Name,
		},
	}

	// 添加 备用节点
	tf.AnnounceList = []string{}
	for _, v := range bto.AnnounceList {
		tf.AnnounceList = append(tf.AnnounceList, v[0])
	}

	// 构建 fileInfo 列表
	var fileInfo []p2p.FileInfo

	for _, v := range bto.Info.Files {
		path := ""
		for _, p := range v.Path {
			path += "/" + p
		}
		fileInfo = append(fileInfo, p2p.FileInfo{
			Path:   path,
			Length: v.Length,
		})
	}
	tf.File = fileInfo

	return tf, nil
}

func (bto *MultipleTorrent) writerFile() {
	//none
}
