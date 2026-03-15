package app

import (
	"github.com/haochend413/Munina/internal/app/events"
	"github.com/haochend413/Munina/internal/app/iobuf"
)

type Model struct {
	IO       iobuf.Model
	EventMgr events.EventMgr
}
