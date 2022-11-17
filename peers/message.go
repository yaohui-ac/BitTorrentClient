package peers

import (
	"encoding/binary"
	"fmt"
	"io"
)

type messageID uint8

const (
	// MsgChoke chokes the receiver
	MsgChoke messageID = 0
	// MsgUnchoke unchokes the receiver
	MsgUnchoke messageID = 1
	// MsgInterested expresses interest in receiving data
	MsgInterested messageID = 2
	// MsgNotInterested expresses disinterest in receiving data
	MsgNotInterested messageID = 3
	// MsgHave alerts the receiver that the sender has downloaded a piece
	MsgHave messageID = 4
	// MsgBitfield encodes which pieces that the sender has downloaded
	MsgBitfield messageID = 5
	// MsgRequest requests a block of data from the receiver
	MsgRequest messageID = 6
	// MsgPiece delivers a block of data to fulfill a request
	MsgPiece messageID = 7
	// MsgCancel cancels a request
	MsgCancel messageID = 8
)

// Message stores ID and payload of a message
type Message struct {
	ID      messageID
	Payload []byte
}

func (m *Message) String() string {
	if m == nil {
		return "KeepAlive"
	}
	switch m.ID {
	case MsgChoke:
		return "Choke"
	case MsgUnchoke:
		return "Unchoke"
	case MsgInterested:
		return "Interested"
	case MsgNotInterested:
		return "NotInterested"
	case MsgHave:
		return "Have"
	case MsgBitfield:
		return "Bitfield"
	case MsgRequest:
		return "Request"
	case MsgPiece:
		return "Piece"
	case MsgCancel:
		return "Cancel"
	default:
		return fmt.Sprintf("Unknown#%d", m.ID)
	}
}
func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}
	length := uint32(len(m.Payload) + 1)
	buf := make([]byte, length+4) //消息体长度 + 消息体内容
	binary.BigEndian.PutUint32(buf[0:4], length)
	// 消息内容
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	return buf
}

// FormatRequest 创建request报文格式
func FormatRequest(index, begin, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return &Message{
		ID:      MsgRequest,
		Payload: payload,
	}
}

// FormatHave 创建have报文格式
func FormatHave(index int) *Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	return &Message{
		ID:      MsgHave,
		Payload: payload,
	}
}

//ParsePiece 解析piece报文消息并将其有效负载复制到缓冲区buf buf需要预先开好
func ParsePiece(index int, buf []byte, msg *Message) (int, error) {
	if msg.ID != MsgPiece {
		return 0, fmt.Errorf("Expected PIECE (ID %d), got ID %d", MsgPiece, msg.ID)
	}
	// 有 8byte 的长度 [index(4), begin(4)]
	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("Payload too short. %d < 8", len(msg.Payload))
	}
	parsedIndex := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
	if parsedIndex != index {
		return 0, fmt.Errorf("Expected index %d, got %d", index, parsedIndex)
	}
	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	if begin >= len(buf) {
		return 0, fmt.Errorf("Begin offset too high. %d >= %d", begin, len(buf))
	}
	// 数据体
	data := msg.Payload[8:]
	// 溢出
	if begin+len(data) > len(buf) {
		return 0, fmt.Errorf("Data too long [%d] for offset %d with length %d", len(data), begin, len(buf))
	}
	// 字节 复制 填充
	copy(buf[begin:], data)
	return len(data), nil

}

// ParseHave 解析have报文
func ParseHave(msg *Message) (int, error) {
	if msg.ID != MsgHave {
		return 0, fmt.Errorf("Expected HAVE (ID %d), got ID %d", MsgHave, msg.ID)
	}

	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("Expected payload length 4, got length %d", len(msg.Payload))
	}

	index := int(binary.BigEndian.Uint32(msg.Payload))
	return index, nil
}

func ReadMessage(r io.Reader) (*Message, error) {
	// 消息长度 4 byte
	// 消息ID 1 byte
	// 消息 payload [5:]
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint16(lengthBuf)
	if length == 0 {
		//keepalive
		return nil, nil
	}
	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, err
	}
	return &Message{
		ID:      messageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}, nil
}
