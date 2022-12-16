package conf

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port            int16    `json:"port"`
	AnnounceList    []string `json:"announce-list"`
	DownloadPath    string   `json:"downloadPath"`
	TorrentFilePath string   `json:"torrentFilePath"`
}

var Conf *Config

func ReadConf() error {
	bytes, err := os.ReadFile("conf.json")
	if err != nil {
		return err
	}

	Conf = &Config{}
	err = json.Unmarshal(bytes, Conf)
	if err != nil {
		return err
	}

	return nil
}
