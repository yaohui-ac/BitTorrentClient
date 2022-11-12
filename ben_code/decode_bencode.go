package ben_code

import "io"

func BytesToBenCode(bytes []byte) (*BenCode, error) {
	// todo
	return nil, nil
}
func DecodeString(w *io.Reader)(str string, err error) {
	return "", nil
}
func DecodeInteger(w *io.Reader)(val int64, err error) {
	return 0, err
}
func DecodeList(w *io.Reader)([]*BenCode, error) {
	return nil, nil
}
func DecodeDict(w *io.Reader) (map[string]*BenCode, error) {
	return nil, nil
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
