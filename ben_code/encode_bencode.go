package ben_code

import (
	"BitTorrentClient/consts"
	"bytes"
	"strconv"
)

func (b *BenCode) EncodeToBytes() []byte {
	switch b.BenType {

	case consts.STRING:
		str, _ := b.Str()
		return EncodeStrToByte(str)
	case consts.INTEGER:
		val, _ := b.Integer()
		return EncodeIntToByte(val)
	case consts.LIST:
		lists, _ := b.List()
		return EncodeListToByte(lists)
	case consts.DICT:
		dict, _ := b.Dict()
		return EncodeDictToByte(dict)
	}

	return nil
}
func EncodeStrToByte(str string) []byte {
	Len := len(str)
	strLen := strconv.FormatInt(int64(Len), 10)
	message := strLen + ":" + str
	return []byte(message)
}
func EncodeIntToByte(val int64) []byte {
	strVal := strconv.FormatInt(val, 10)
	message := "i" + strVal + "e"
	return []byte(message)
}
func EncodeListToByte(benCodes []*BenCode) []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteRune('l')
	for i := range benCodes {
		buffer.Write(benCodes[i].EncodeToBytes())
	}
	buffer.WriteRune('e')
	return buffer.Bytes()
}
func EncodeDictToByte(dict map[string]*BenCode) []byte {

	return nil
}
