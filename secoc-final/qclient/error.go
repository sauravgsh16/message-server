package qclient

import (
	"github.com/sauravgsh16/secoc-third/secoc-final/proto"
)

var (
	ErrUnexpectedFrame = proto.NewHardError(505, "Unexpected Frame", 0, 0)
	ErrClosed          = proto.NewHardError(504, "Communication attempt on close Channel/Connection", 0, 0)
)
