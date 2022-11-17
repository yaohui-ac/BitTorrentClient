package torrent

import (
	"BitTorrentClient/ben_code"
	"BitTorrentClient/p2p"
	"bytes"
	"crypto/sha1"
	"fmt"
)

type TorrentFile struct {
	Announce     string
	AnnounceList []string
	p2p.Torrent  //建立网络连接所需的信息
}

//定义统一接口 统一单文件模式和多文件模式
type torrent interface {
	fileParser(file []byte) error
	toTorrentFile() (TorrentFile, error)
	writeFile()
}

//info 对块进行哈希 -> info hash
func hash(i interface{}) ([20]byte, error) {
	var buf bytes.Buffer
	err := ben_code.Marshal(&buf, i)
	if err != nil {
		return [20]byte{}, err
	}
	hs := sha1.Sum(buf.Bytes())
	return hs, nil
}

//// splitPieceHashes 切割Pieces
func splitPieceHashes(pieces string) ([][20]byte, error) {
	// SHA-1 hash的长度
	hashLen := 20
	buf := []byte(pieces)

	if len(buf)%hashLen != 0 {
		// 片段的长度不正确
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}

	// hash 的总数
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ { // 下标访问切割
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}
