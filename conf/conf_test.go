package conf

import (
	"fmt"
	"testing"
)

func TestReadConf(t *testing.T) {
	err := ReadConf()
	fmt.Println(err, Conf.Port)
}
