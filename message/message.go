package message

import "github.com/gofrs/uuid/v3"

// Frame represents message metadata for storage and transmission
type Frame struct {
	Type      byte
	Version   byte
	length    uint32
	Ident     uint32
	MessageID uuid.UUID
}
