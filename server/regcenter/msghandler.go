package regcenter

import (
	"io"
)

type msgHandle struct {
	ws io.WriteSeeker
	f func(data interface{}, writer io.Writer) error
}


type MsgHandler interface {
	SetWriteSeeker(ws io.WriteSeeker)
	GetWriteSeeker() io.WriteSeeker
	SetHandler(f func(data interface{},writer io.Writer) error)
	GetHandler() func(data interface{},writer io.Writer) error
}

func NewMsgHandler() MsgHandler {
	mh := &msgHandle{}

	return mh
}

func (mh *msgHandle)SetWriteSeeker(ws io.WriteSeeker)  {
	mh.ws = ws
}

func (mh *msgHandle)GetWriteSeeker() io.WriteSeeker {
	return mh.ws
}

func (mh *msgHandle)SetHandler(f func(data interface{},writer io.Writer) error)  {
	mh.f = f
}

func (mh *msgHandle)GetHandler() func(data interface{},writer io.Writer) error  {
	return mh.f
}