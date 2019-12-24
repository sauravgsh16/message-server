package qclient

import (
	"github.com/sauravgsh16/message-server/proto"
)

// TODO: Find all proto error defs and add here

var (
	ErrUnexpectedFrame = proto.NewHardError(505, "Unexpected Frame", 0, 0)
	ErrClosed          = proto.NewHardError(504, "Communication attempt on close Channel/Connection", 0, 0)
)
