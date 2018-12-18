// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rcv.proto

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

type UdpAck struct {
	SerialNo             uint64   `protobuf:"varint,1,opt,name=serialNo,proto3" json:"serialNo,omitempty"`
	DataType             uint32   `protobuf:"varint,2,opt,name=dataType,proto3" json:"dataType,omitempty"`
	Ack                  uint32   `protobuf:"varint,3,opt,name=ack,proto3" json:"ack,omitempty"`
	Resend               []uint32 `protobuf:"varint,4,rep,packed,name=resend,proto3" json:"resend,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UdpAck) Reset()         { *m = UdpAck{} }
func (m *UdpAck) String() string { return proto.CompactTextString(m) }
func (*UdpAck) ProtoMessage()    {}
func (*UdpAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f416513267849f0, []int{0}
}

func (m *UdpAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UdpAck.Unmarshal(m, b)
}
func (m *UdpAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UdpAck.Marshal(b, m, deterministic)
}
func (m *UdpAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UdpAck.Merge(m, src)
}
func (m *UdpAck) XXX_Size() int {
	return xxx_messageInfo_UdpAck.Size(m)
}
func (m *UdpAck) XXX_DiscardUnknown() {
	xxx_messageInfo_UdpAck.DiscardUnknown(m)
}

var xxx_messageInfo_UdpAck proto.InternalMessageInfo

func (m *UdpAck) GetSerialNo() uint64 {
	if m != nil {
		return m.SerialNo
	}
	return 0
}

func (m *UdpAck) GetDataType() uint32 {
	if m != nil {
		return m.DataType
	}
	return 0
}

func (m *UdpAck) GetAck() uint32 {
	if m != nil {
		return m.Ack
	}
	return 0
}

func (m *UdpAck) GetResend() []uint32 {
	if m != nil {
		return m.Resend
	}
	return nil
}

func init() {
	proto.RegisterType((*UdpAck)(nil), "packet.UdpAck")
}

func init() { proto.RegisterFile("rcv.proto", fileDescriptor_8f416513267849f0) }

var fileDescriptor_8f416513267849f0 = []byte{
	// 129 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2c, 0x4a, 0x2e, 0xd3,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2b, 0x48, 0x4c, 0xce, 0x4e, 0x2d, 0x51, 0xca, 0xe2,
	0x62, 0x0b, 0x4d, 0x29, 0x70, 0x4c, 0xce, 0x16, 0x92, 0xe2, 0xe2, 0x28, 0x4e, 0x2d, 0xca, 0x4c,
	0xcc, 0xf1, 0xcb, 0x97, 0x60, 0x54, 0x60, 0xd4, 0x60, 0x09, 0x82, 0xf3, 0x41, 0x72, 0x29, 0x89,
	0x25, 0x89, 0x21, 0x95, 0x05, 0xa9, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0xbc, 0x41, 0x70, 0xbe, 0x90,
	0x00, 0x17, 0x73, 0x62, 0x72, 0xb6, 0x04, 0x33, 0x58, 0x18, 0xc4, 0x14, 0x12, 0xe3, 0x62, 0x2b,
	0x4a, 0x2d, 0x4e, 0xcd, 0x4b, 0x91, 0x60, 0x51, 0x60, 0xd6, 0xe0, 0x0d, 0x82, 0xf2, 0x92, 0xd8,
	0xc0, 0x56, 0x1b, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x7a, 0xc0, 0x71, 0xfc, 0x87, 0x00, 0x00,
	0x00,
}
