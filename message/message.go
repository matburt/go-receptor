package message

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v3"
	"github.com/spf13/viper"
)

// FrameType is the type of frame being transmitted as part of a message
type FrameType int

// Frame constants identifying the type of message following the Frame
const (
	Header FrameType = iota
	Payload
	Command
)

//CommandMessage represents command sent across the mesh
type CommandMessage struct {
	Cmd        string            `json:"cmd"`
	NodeID     string            `json:"id"`
	ExpireTime int64             `json:"expire_time"`
	Meta       map[string]string `json:"meta"` // placeholder for a more complex type
}

// Frame represents message metadata for storage and transmission
type Frame struct {
	Type      byte
	Version   byte
	Length    uint32
	Ident     uint32
	MessageID uuid.UUID
}

//FramedMessage represents a message
type FramedMessage struct {
	MessageID uuid.UUID
	Header    []byte
	payload   []byte
}

//MakeHiMessage generates a hi command message
func MakeHiMessage() (*FramedMessage, error) {
	header := CommandMessage{"hi", viper.GetString("node_id"), time.Now().Unix(), make(map[string]string)}
	bHeader, err := json.Marshal(header)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate Hi message %v", err)
	}
	return NewFramedMessageB(bHeader), nil
}

//NewFramedMessage returns a new FramedMessage
func NewFramedMessage() *FramedMessage {
	newUUID, _ := uuid.NewV4()
	return &FramedMessage{newUUID, []byte{}, nil}
}

//NewFramedMessageB creates a new Framed Message from header bytes
func NewFramedMessageB(b []byte) *FramedMessage {
	newUUID, _ := uuid.NewV4()
	return &FramedMessage{newUUID, b, nil}
}

// Serialize a FramedMessage for network transmission
func (m *FramedMessage) Serialize(b *bytes.Buffer) {
	fmt.Println("Serializing", m.MessageID)
	var messageType FrameType
	if len(m.payload) == 0 {
		messageType = Command
	} else {
		messageType = Header
	}
	headerFrame := Frame{
		byte(messageType),
		1,
		uint32(len(m.Header)),
		1,
		m.MessageID}
	headerFrame.Serialize(b)
	binary.Write(b, binary.BigEndian, m.Header)
	// TODO: Do something with payload
}

// DeSerializeFramedMessage from bytes into a FramedMessage
func DeSerializeFramedMessage(b *bytes.Buffer, frame *Frame) *FramedMessage {
	newFramedMessage := new(FramedMessage)
	headerPayload := make([]byte, frame.Length)
	binary.Read(b, binary.BigEndian, headerPayload)
	newFramedMessage.MessageID = frame.MessageID
	newFramedMessage.Header = headerPayload[:]
	newFramedMessage.payload = nil
	return newFramedMessage
}

// Serialize a Frame for network transmission
func (f *Frame) Serialize(b *bytes.Buffer) {
	b.Write([]byte{f.Type, f.Version})
	binary.Write(b, binary.BigEndian, f.Ident)
	binary.Write(b, binary.BigEndian, f.Length)
	binary.Write(b, binary.BigEndian, f.MessageID.Bytes()[:8])
	binary.Write(b, binary.BigEndian, f.MessageID.Bytes()[8:])
}

// DeSerializeFrame from bytes into a Frame refrence
// NOTE: Replace struct uuid with [16]byte and we can binary.Read(b....&Frame)
func DeSerializeFrame(b *bytes.Buffer) *Frame {
	newFrame := new(Frame)
	binary.Read(b, binary.BigEndian, newFrame.Type)
	binary.Read(b, binary.BigEndian, newFrame.Version)
	binary.Read(b, binary.BigEndian, newFrame.Ident)
	binary.Read(b, binary.BigEndian, newFrame.Length)
	var h, l [8]byte
	binary.Read(b, binary.BigEndian, h)
	binary.Read(b, binary.BigEndian, l)
	message := h[:8]
	newFrame.MessageID, _ = uuid.FromBytes(append(message, l[:8]...))
	return newFrame
}
