package consts

type BenObjType int
type BenValue interface{}

const (
	STRING  BenObjType = 1
	INTEGER BenObjType = 2
	LIST    BenObjType = 3
	DICT    BenObjType = 4
)

const BenCodeTag = "bencode"
