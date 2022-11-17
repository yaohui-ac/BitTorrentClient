package torrent

import (
	"BitTorrentClient/ben_code"
	"BitTorrentClient/conf"
	"BitTorrentClient/consts"
	"BitTorrentClient/peers"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type bencodeTrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (t *TorrentFile) buildTrackerUrl(peerID [20]byte, port uint16) (*url.URL, error) {
	baseUrl, err := url.Parse(t.Announce)
	if err != nil {
		return &url.URL{}, err
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(conf.Port))},
		"uploaded":   []string{"0"}, //默认值
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	baseUrl.RawQuery = params.Encode()
	return baseUrl, nil
}

func (t *TorrentFile) requestPeers(peerID [20]byte, port uint16) error {
	base, err := t.buildTrackerUrl(peerID, port)
	if err != nil {
		return err
	}
	// fmt.Println(base)
	switch base.Scheme {
	case "http", "https":
		c := &http.Client{
			Timeout: time.Second * 30,
		}
		resp, err := c.Get(base.String())
		if err != nil {
			return err
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		trackerResp := bencodeTrackerResp{}
		all, err := ioutil.ReadAll(resp.Body)
		err = ben_code.Unmarshal(bytes.NewReader(all), &trackerResp)
		if err != nil {
			return err
		}
		t.Peers, err = peers.Unmarshal([]byte(trackerResp.Peers))
		return err
	case "udp":
		_, err := net.Dial("udp", base.Host)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return err

	default:
		return consts.ProtocolErr
	}
}
