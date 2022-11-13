package torrent

import (
	"BitTorrentClient/ben_code"
	"io"
)

func SingleTorrentParser(reader io.Reader) (s *SingleTorrent, err error) {
	s = new(SingleTorrent)
	err = ben_code.Unmarshal(reader, s)
	return s, err
}
