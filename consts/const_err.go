package consts

import "errors"

var (
	BenObjTypeErr    error //bencode解析出错
	ProtocolErr      error //traker通信 不支持的协议
	PeerUnmarshalErr error //解析出peers ip port对时出错
)

func init() {
	BenObjTypeErr = errors.New("BenCode Object Type Error")
	ProtocolErr = errors.New("unsupported protocols error")
	PeerUnmarshalErr = errors.New("peers unmarshal err")

}
