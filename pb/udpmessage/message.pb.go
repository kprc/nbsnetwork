// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package udpmessage

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Udpmsg struct {
	Sn                   uint64   `protobuf:"varint,1,opt,name=sn,proto3" json:"sn,omitempty"`
	Pos                  uint64   `protobuf:"varint,2,opt,name=pos,proto3" json:"pos,omitempty"`
	Data                 []byte   `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Udpmsg) Reset()         { *m = Udpmsg{} }
func (m *Udpmsg) String() string { return proto.CompactTextString(m) }
func (*Udpmsg) ProtoMessage()    {}
func (*Udpmsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}

func (m *Udpmsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Udpmsg.Unmarshal(m, b)
}
func (m *Udpmsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Udpmsg.Marshal(b, m, deterministic)
}
func (m *Udpmsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Udpmsg.Merge(m, src)
}
func (m *Udpmsg) XXX_Size() int {
	return xxx_messageInfo_Udpmsg.Size(m)
}
func (m *Udpmsg) XXX_DiscardUnknown() {
	xxx_messageInfo_Udpmsg.DiscardUnknown(m)
}

var xxx_messageInfo_Udpmsg proto.InternalMessageInfo

func (m *Udpmsg) GetSn() uint64 {
	if m != nil {
		return m.Sn
	}
	return 0
}

func (m *Udpmsg) GetPos() uint64 {
	if m != nil {
		return m.Pos
	}
	return 0
}

func (m *Udpmsg) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type Udpmsgid struct {
	Sn                   uint64   `protobuf:"varint,1,opt,name=sn,proto3" json:"sn,omitempty"`
	Pos                  uint64   `protobuf:"varint,2,opt,name=pos,proto3" json:"pos,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Udpmsgid) Reset()         { *m = Udpmsgid{} }
func (m *Udpmsgid) String() string { return proto.CompactTextString(m) }
func (*Udpmsgid) ProtoMessage()    {}
func (*Udpmsgid) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{1}
}

func (m *Udpmsgid) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Udpmsgid.Unmarshal(m, b)
}
func (m *Udpmsgid) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Udpmsgid.Marshal(b, m, deterministic)
}
func (m *Udpmsgid) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Udpmsgid.Merge(m, src)
}
func (m *Udpmsgid) XXX_Size() int {
	return xxx_messageInfo_Udpmsgid.Size(m)
}
func (m *Udpmsgid) XXX_DiscardUnknown() {
	xxx_messageInfo_Udpmsgid.DiscardUnknown(m)
}

var xxx_messageInfo_Udpmsgid proto.InternalMessageInfo

func (m *Udpmsgid) GetSn() uint64 {
	if m != nil {
		return m.Sn
	}
	return 0
}

func (m *Udpmsgid) GetPos() uint64 {
	if m != nil {
		return m.Pos
	}
	return 0
}

type Udpmsgack struct {
	Uid                  *Udpmsgid `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Arrpos               []uint64  `protobuf:"varint,3,rep,packed,name=arrpos,proto3" json:"arrpos,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Udpmsgack) Reset()         { *m = Udpmsgack{} }
func (m *Udpmsgack) String() string { return proto.CompactTextString(m) }
func (*Udpmsgack) ProtoMessage()    {}
func (*Udpmsgack) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{2}
}

func (m *Udpmsgack) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Udpmsgack.Unmarshal(m, b)
}
func (m *Udpmsgack) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Udpmsgack.Marshal(b, m, deterministic)
}
func (m *Udpmsgack) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Udpmsgack.Merge(m, src)
}
func (m *Udpmsgack) XXX_Size() int {
	return xxx_messageInfo_Udpmsgack.Size(m)
}
func (m *Udpmsgack) XXX_DiscardUnknown() {
	xxx_messageInfo_Udpmsgack.DiscardUnknown(m)
}

var xxx_messageInfo_Udpmsgack proto.InternalMessageInfo

func (m *Udpmsgack) GetUid() *Udpmsgid {
	if m != nil {
		return m.Uid
	}
	return nil
}

func (m *Udpmsgack) GetArrpos() []uint64 {
	if m != nil {
		return m.Arrpos
	}
	return nil
}

func init() {
	proto.RegisterType((*Udpmsg)(nil), "udpmessage.udpmsg")
	proto.RegisterType((*Udpmsgid)(nil), "udpmessage.udpmsgid")
	proto.RegisterType((*Udpmsgack)(nil), "udpmessage.udpmsgack")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 162 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x4d, 0x2d, 0x2e,
	0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2a, 0x4d, 0x29, 0x80, 0x8a,
	0x28, 0xd9, 0x71, 0xb1, 0x81, 0x78, 0xc5, 0xe9, 0x42, 0x7c, 0x5c, 0x4c, 0xc5, 0x79, 0x12, 0x8c,
	0x0a, 0x8c, 0x1a, 0x2c, 0x41, 0x4c, 0xc5, 0x79, 0x42, 0x02, 0x5c, 0xcc, 0x05, 0xf9, 0xc5, 0x12,
	0x4c, 0x60, 0x01, 0x10, 0x53, 0x48, 0x88, 0x8b, 0x25, 0x25, 0xb1, 0x24, 0x51, 0x82, 0x55, 0x81,
	0x51, 0x83, 0x27, 0x08, 0xcc, 0x56, 0xd2, 0xe1, 0xe2, 0x80, 0xe8, 0xcf, 0x4c, 0x21, 0x6c, 0x82,
	0x92, 0x37, 0x17, 0x27, 0x44, 0x75, 0x62, 0x72, 0xb6, 0x90, 0x1a, 0x17, 0x73, 0x69, 0x66, 0x0a,
	0x58, 0x3d, 0xb7, 0x91, 0x88, 0x1e, 0xc2, 0x51, 0x7a, 0x30, 0x13, 0x83, 0x40, 0x0a, 0x84, 0xc4,
	0xb8, 0xd8, 0x12, 0x8b, 0x8a, 0x40, 0x26, 0x31, 0x2b, 0x30, 0x6b, 0xb0, 0x04, 0x41, 0x79, 0x49,
	0x6c, 0x60, 0xdf, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x31, 0x9f, 0x50, 0xec, 0xde, 0x00,
	0x00, 0x00,
}