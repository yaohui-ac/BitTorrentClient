package torrent

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
