package torrent

import (
	"fmt"
	"os"
	"testing"
)

func TestSingleTorrentParser(t *testing.T) {
	fd, _ := os.Open("../btfiles/debian-iso.torrent")
	b, err := SingleTorrentParser(fd)
	fmt.Println("ann:", b.Announce, "annlist:", b.AnnounceList, "commment:", b.CreatDate, b.Comment, "\nerr:", err)
}
