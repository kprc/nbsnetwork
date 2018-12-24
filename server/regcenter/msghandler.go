package regcenter

import (
	"io"
	"sync/atomic"
)

type msgHandle struct {
	fGetWriteSeeker func() io.WriteSeeker
	f func(data interface{}, writer io.Writer) error
	ws io.WriteSeeker
	refcnt int32
}


type MsgHandler interface {
	SetWriteSeeker(ws io.WriteSeeker)
	GetWriteSeeker() io.WriteSeeker
	SetHandler(f func(data interface{},writer io.Writer) error)
	GetHandler() func(data interface{},writer io.Writer) error
	IncRef()
	DecRef()
	GetRef() int32
	SetWSNew(fGetWriteSeeker func() io.WriteSeeker)
	GetWSNew() func() io.WriteSeeker
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

func (mh *msgHandle)IncRef()  {
	atomic.AddInt32(&mh.refcnt,1)
}

func (mh *msgHandle)DecRef(){
	atomic.AddInt32(&mh.refcnt,-1)
}

func (mh *msgHandle)GetRef() int32 {

	return atomic.LoadInt32(&mh.refcnt)
}

func (mh *msgHandle)SetWSNew(fGetWriteSeeker func() io.WriteSeeker) {
	mh.fGetWriteSeeker = fGetWriteSeeker
}

func (mh *msgHandle)GetWSNew() func() io.WriteSeeker  {
	return mh.fGetWriteSeeker
}