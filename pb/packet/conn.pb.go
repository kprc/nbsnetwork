// Code generated by protoc-gen-go. DO NOT EDIT.
// source: packet/conn.proto

package packet

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

type UdpConnMsg struct {
	Typ                  uint32   `protobuf:"varint,1,opt,name=typ,proto3" json:"typ,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UdpConnMsg) Reset()         { *m = UdpConnMsg{} }
func (m *UdpConnMsg) String() string { return proto.CompactTextString(m) }
func (*UdpConnMsg) ProtoMessage()    {}
func (*UdpConnMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_79a9896a3408d039, []int{0}
}

func (m *UdpConnMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UdpConnMsg.Unmarshal(m, b)
}
func (m *UdpConnMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UdpConnMsg.Marshal(b, m, deterministic)
}
func (m *UdpConnMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UdpConnMsg.Merge(m, src)
}
func (m *UdpConnMsg) XXX_Size() int {
	return xxx_messageInfo_UdpConnMsg.Size(m)
}
func (m *UdpConnMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_UdpConnMsg.DiscardUnknown(m)
}

var xxx_messageInfo_UdpConnMsg proto.InternalMessageInfo

func (m *UdpConnMsg) GetTyp() uint32 {
	if m != nil {
		return m.Typ
	}
	return 0
}

func (m *UdpConnMsg) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*UdpConnMsg)(nil), "packet.UdpConnMsg")
}

func init() { proto.RegisterFile("packet/conn.proto", fileDescriptor_79a9896a3408d039) }

var fileDescriptor_79a9896a3408d039 = []byte{
	// 99 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2c, 0x48, 0x4c, 0xce,
	0x4e, 0x2d, 0xd1, 0x4f, 0xce, 0xcf, 0xcb, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x83,
	0x08, 0x29, 0x19, 0x71, 0x71, 0x85, 0xa6, 0x14, 0x38, 0xe7, 0xe7, 0xe5, 0xf9, 0x16, 0xa7, 0x0b,
	0x09, 0x70, 0x31, 0x97, 0x54, 0x16, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0xf0, 0x06, 0x81, 0x98, 0x42,
	0x42, 0x5c, 0x2c, 0x29, 0x89, 0x25, 0x89, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x3c, 0x41, 0x60, 0x76,
	0x12, 0x1b, 0xd8, 0x08, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x51, 0x95, 0x2a, 0x42, 0x57,
	0x00, 0x00, 0x00,
}
