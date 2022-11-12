package consts

import "errors"

type BenObjType int
type BenValue interface{}

var (
	BenObjTypeErr error
)

func init() {
	BenObjTypeErr = errors.New("BenCode Object Type Error")
}

const (
	STRING  BenObjType = 1
	INTEGER BenObjType = 2
	LIST    BenObjType = 3
	DICT    BenObjType = 4
)
