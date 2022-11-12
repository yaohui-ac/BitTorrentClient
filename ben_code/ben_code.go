package ben_code

import "BitTorrentClient/consts"

type BenCode struct {
	BenType  consts.BenObjType
	BenValue consts.BenValue
}

func NewBenCode() *BenCode {
	return &BenCode{}
}
func (b *BenCode) Str() (string, error) {
	if b.BenType != consts.STRING {
		return "", consts.BenObjTypeErr
	}
	return b.BenValue.(string), nil
}
func (b *BenCode) Integer() (int64, error) {
	if b.BenType != consts.INTEGER {
		return 0, consts.BenObjTypeErr
	}
	return b.BenValue.(int64), nil
}

func (b *BenCode) List() ([]*BenCode, error) {
	if b.BenType != consts.LIST {
		return nil, consts.BenObjTypeErr
	}
	return b.BenValue.([]*BenCode), nil
}
func (b *BenCode) Dict() (map[string]*BenCode, error) {
	if b.BenType != consts.DICT {
		return nil, consts.BenObjTypeErr
	}
	return b.BenValue.(map[string]*BenCode), nil
}
