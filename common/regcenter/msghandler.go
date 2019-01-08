package regcenter

import (
	"io"
	"sync/atomic"
)

type msgHandle struct {
	fGetWriteSeeker func(param interface{}) io.WriteSeeker
	f func(head interface{},data interface{}, snd io.Writer) error
	refcnt int32
}


type MsgHandler interface {
	SetHandler(f func(head interface{},data interface{},snd io.Writer) error)
	GetHandler() func(head interface{},data interface{},snd io.Writer) error
	IncRef()
	DecRef() int32
	GetRef() int32
	SetWSNew(fGetWriteSeeker func(param interface{}) io.WriteSeeker)
	GetWSNew() func(param interface{}) io.WriteSeeker
}

func NewMsgHandler() MsgHandler {
	mh := &msgHandle{}

	return mh
}


func (mh *msgHandle)SetHandler(f func(head interface{},data interface{},snd io.Writer) error)  {
	mh.f = f
}

func (mh *msgHandle)GetHandler() func(head interface{},data interface{},snd io.Writer) error  {
	return mh.f
}

func (mh *msgHandle)IncRef()  {
	atomic.AddInt32(&mh.refcnt,1)
}

func (mh *msgHandle)DecRef() int32{
	return atomic.AddInt32(&mh.refcnt,-1)
}

func (mh *msgHandle)GetRef() int32 {

	return atomic.LoadInt32(&mh.refcnt)
}

func (mh *msgHandle)SetWSNew(fGetWriteSeeker func(param interface{}) io.WriteSeeker) {
	mh.fGetWriteSeeker = fGetWriteSeeker
}

func (mh *msgHandle)GetWSNew() func(param interface{}) io.WriteSeeker  {
	return mh.fGetWriteSeeker
}