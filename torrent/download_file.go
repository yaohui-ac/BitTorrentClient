package torrent

import (
	"BitTorrentClient/conf"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func NewTorrent() (f TorrentFile, err error) {
	file, err := os.Open(conf.Conf.TorrentFilePath)
	if err != nil {
		return f, err
	}
	defer func() {
		_ = file.Close()
	}()
	all, err := ioutil.ReadAll(file)
	if err != nil {
		return f, err
	}
	var t torrent

	if strings.Contains(string(all), "5:files") {
		// 多个文件
		fmt.Println("检测到多个文件")
		// t = &MultipleTorrent{}
	} else {
		// 单个文件
		fmt.Println("检测到单个文件")
		t = &SingleTorrent{}
	}

	err = t.fileParser(all)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	f, err = t.toTorrentFile()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	fmt.Println(conf.Conf.TorrentFilePath, "info hash:", hex.EncodeToString(f.InfoHash[:]))

	return f, nil

}
func (t *TorrentFile) DownloadFile(path string) error {
	var peerID [20]byte
	_, err := rand.Read(peerID[:]) //随机生成peerId
	if err != nil {
		return err
	}
	err = t.requestPeers(peerID, uint16(conf.Conf.Port))
	if err != nil {
		return err
	}

	tor := t.Torrent

	buf, err := tor.Download()

	if err != nil {
		return err
	}
	outFile, err := os.Create(path)
	defer func() {
		_ = outFile.Close()
	}()
	if err != nil {
		return err
	}
	_, err = outFile.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
