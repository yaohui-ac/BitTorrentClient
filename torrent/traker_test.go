package torrent

import (
	"BitTorrentClient/p2p"
	"fmt"
	"testing"
)

func TestBuildTrackerURL(t *testing.T) {
	to := &TorrentFile{
		Announce: "http://bttracker.debian.org:6969/announce",
		// all
		//InfoHash: [20]byte{85, 57, 15, 31, 20, 98, 244, 202, 73, 202, 4, 111, 119, 35, 156, 68, 240, 78, 88, 43},
		Torrent: p2p.Torrent{
			InfoHash: [20]byte{159, 41, 44, 147, 235, 13, 189, 215, 255, 122, 74, 165, 81, 170, 161, 234, 124, 175, 224, 4},
			Length:   353370112,
		},
	}

	peerID := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	const port uint16 = 6882

	url, err := to.buildTrackerUrl(peerID, port)
	fmt.Println("url: ", url, err)
	//expected := "http://bttracker.debian.org:6969/announce?compact=1&downloaded=0&info_hash=%D8%F79%CE%C3%28%95l%CC%5B%BF%1F%86%D9%FD%CF%DB%A8%CE%B6&left=351272960&peer_id=%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14&port=6881&uploaded=0"
}
