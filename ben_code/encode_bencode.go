package ben_code

import (
	"BitTorrentClient/consts"
	"bytes"
	"io"
	"sort"
	"strconv"
)

func (b *BenCode) EncodeToBytes() []byte {
	//统一接口解析
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
func EncodeStr(w io.Writer, str string) {
	_, _ = w.Write(EncodeStrToByte(str))
}
func EncodeIntToByte(val int64) []byte {
	strVal := strconv.FormatInt(val, 10)
	message := "i" + strVal + "e"
	return []byte(message)
}
func EncodeInt(w io.Writer, val int64) {
	_, _ = w.Write(EncodeIntToByte(val))
}
func EncodeListToByte(benCodes []*BenCode) []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteRune('l')
	for _, item := range benCodes {
		buffer.Write(item.EncodeToBytes())
	}
	buffer.WriteRune('e')
	return buffer.Bytes()
}
func EncodeDictToByte(dict map[string]*BenCode) []byte {

	sortedList := make([]*benCodeKV, 0)
	for k, v := range dict {
		sortedList = append(sortedList, &benCodeKV{
			key: k,
			val: v,
		})
	}
	sort.Slice(sortedList, func(i, j int) bool {
		return sortedList[i].key < sortedList[j].key
	})
	buffer := new(bytes.Buffer)
	buffer.WriteRune('d')
	for _, item := range sortedList {
		strKeyLen := strconv.FormatInt(int64(len(item.key)), 10)
		buffer.WriteRune(':')
		buffer.WriteString(strKeyLen)
		buffer.Write(item.val.EncodeToBytes())
	}
	buffer.WriteRune('e')
	return buffer.Bytes()

}
func (b *BenCode) EncodeToWriter(w io.Writer) error {
	_, err := w.Write(b.EncodeToBytes())
	return err
}
