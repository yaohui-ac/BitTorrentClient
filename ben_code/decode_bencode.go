package ben_code

import (
	"BitTorrentClient/consts"
	"bufio"
	"errors"
	"io"
)

func BytesToBenCode(bytes []byte) (*BenCode, error) {
	// todo

	return nil, nil
}
func DecodeBencode(r io.Reader) (*BenCode, error) {
	//todo
	bufReader := bufio.NewReader(r)
	preview, err := bufReader.Peek(1)
	benCode := new(BenCode)
	if err != nil {
		return nil, err
	}
	switch {
	case preview[0] >= '0' && preview[0] <= '9':
		//字符串
		str, err := DecodeString(bufReader)
		if err != nil {
			return nil, err
		}
		benCode.BenType = consts.STRING
		benCode.BenValue = str
		return benCode, err
	case preview[0] == 'i':
		num, err := DecodeInteger(bufReader)
		if err != nil {
			return nil, err
		}
		benCode.BenType = consts.INTEGER
		benCode.BenValue = num
		return benCode, err
	case preview[0] == 'l':
		lists, err := DecodeList(bufReader)
		if err != nil {
			return nil, err
		}
		benCode.BenType = consts.LIST
		benCode.BenValue = lists
		return benCode, err
	case preview[0] == 'd':
		dict, err := DecodeDict(bufReader)
		if err != nil {
			return nil, err
		}
		benCode.BenType = consts.DICT
		benCode.BenValue = dict
		return benCode, err
	}

	return nil, err
}

func DecodeString(r io.Reader) (str string, err error) {
	bufReader := bufio.NewReader(r)
	strLen, err := parseInteger(bufReader) //解析长度
	if err != nil {
		return "", errors.New("待定")
	}
	if strLen <= 0 {
		return "", errors.New("xx")
	}

	b, _ := bufReader.ReadByte()
	if b != ':' {
		return "", errors.New("待定")
	}

	str, err = parseString(bufReader, strLen) //解析后面的字符串
	if err != nil {
		return "", errors.New("待定")
	}

	return str, nil
}
func parseInteger(r *bufio.Reader) (int64, error) {
	var num int64 = 0
	sign := 1
	b, _ := r.ReadByte()
	if b == '-' {
		sign = -1
	}
	for {
		bytes, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		if bytes < '0' || bytes > '9' {
			_ = r.UnreadByte()
			break
		}
		num = num*10 + int64(bytes-'0')
	}
	return int64(sign) * num, nil
}
func parseString(r *bufio.Reader, strLen int64) (string, error) {
	bytes := make([]byte, strLen)
	n, err := io.ReadFull(r, bytes)
	if int64(n) != strLen || err != nil {
		return "", errors.New("待定")
	}
	return string(bytes), nil
}

func DecodeInteger(r io.Reader) (int64, error) {
	bufReader := bufio.NewReader(r)
	b, err := bufReader.ReadByte()
	if b != 'i' {
		return 0, errors.New("待定")
	}
	num, err := parseInteger(bufReader)
	b, err = bufReader.ReadByte()
	if b != 'e' {
		return 0, errors.New("待定")
	}

	return num, err
}
func DecodeList(r io.Reader) ([]*BenCode, error) {
	bufReader := bufio.NewReader(r)
	b, err := bufReader.ReadByte()
	if err != nil {
		//todo
	}
	lists := make([]*BenCode, 0)
	if b != 'l' {
		return nil, errors.New("待定")
	}
	for {
		b, _ := bufReader.Peek(1)
		if b[0] == 'e' {
			_, _ = bufReader.ReadByte()
			break //解析终止
		}
		benCode, err := DecodeBencode(r)
		if err != nil {
			//todo
			return nil, err
		}
		lists = append(lists, benCode)
	}
	return lists, nil
}
func DecodeDict(r io.Reader) (map[string]*BenCode, error) {
	bufReader := bufio.NewReader(r)
	mp := make(map[string]*BenCode)
	b, err := bufReader.ReadByte()
	if err != nil {
		//todo
	}
	if b != 'd' {
		return nil, errors.New("待定")
	}

	for {

		b, _ := bufReader.Peek(1)
		if b[0] == 'e' {
			_, _ = bufReader.ReadByte()
			break //解析终止
		}
		key, err := DecodeString(r)
		if err != nil {
			return nil, err
		}
		val, err := DecodeBencode(r)
		if err != nil {
			return nil, err
		}
		mp[key] = val

	}
	return mp, nil
}

func BytesDecodeToStr(bytes []byte) (string, error) {
	return "", nil
}
func BytesDecodeToInt(bytes []byte) (int64, error) {
	return 0, nil
}
func BytesDecodeToList(bytes []byte) ([]*BenCode, error) {
	return nil, nil
}
func BytesDecodeToDict(bytes []byte) (map[string]*BenCode, error) {
	return nil, nil
}
