package util

type Bitfield []byte

func (bf Bitfield) HasPiece(idx int) bool {
	byteIdx := idx / 8
	offset := idx % 8

	if byteIdx < 0 || byteIdx >= len(bf) {
		return false
	}
	return bf[byteIdx]>>uint(7-offset)&1 != 0
}

func (bf Bitfield) SetPiece(index int) {
	byteIndex := index / 8
	offset := index % 8
	
	if byteIndex < 0 || byteIndex >= len(bf) {
		return
	}
	bf[byteIndex] |= 1 << uint(7-offset)
}
