package main

import (
	"BitTorrentClient/conf"
	"fmt"
)

func main() {
	err := conf.ReadConf()
	if err != nil {
		fmt.Println(err)
		return
	}

}
