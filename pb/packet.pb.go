// Code generated by protoc-gen-go. DO NOT EDIT.
// source: packet.proto

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

type UDPPacketData struct {
	SerialNo             uint64   `protobuf:"varint,1,opt,name=serialNo,proto3" json:"serialNo,omitempty"`
	TotalCnt             uint32   `protobuf:"varint,16,opt,name=totalCnt,proto3" json:"totalCnt,omitempty"`
	PosNum               uint32   `protobuf:"varint,2,opt,name=posNum,proto3" json:"posNum,omitempty"`
	DataType             uint32   `protobuf:"varint,3,opt,name=dataType,proto3" json:"dataType,omitempty"`
	TryCnt               uint32   `protobuf:"varint,17,opt,name=tryCnt,proto3" json:"tryCnt,omitempty"`
	Data                 []byte   `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UDPPacketData) Reset()         { *m = UDPPacketData{} }
func (m *UDPPacketData) String() string { return proto.CompactTextString(m) }
func (*UDPPacketData) ProtoMessage()    {}
func (*UDPPacketData) Descriptor() ([]byte, []int) {
	return fileDescriptor_e9ef1a6541f9f9e7, []int{0}
}

func (m *UDPPacketData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UDPPacketData.Unmarshal(m, b)
}
func (m *UDPPacketData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UDPPacketData.Marshal(b, m, deterministic)
}
func (m *UDPPacketData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UDPPacketData.Merge(m, src)
}
func (m *UDPPacketData) XXX_Size() int {
	return xxx_messageInfo_UDPPacketData.Size(m)
}
func (m *UDPPacketData) XXX_DiscardUnknown() {
	xxx_messageInfo_UDPPacketData.DiscardUnknown(m)
}

var xxx_messageInfo_UDPPacketData proto.InternalMessageInfo

func (m *UDPPacketData) GetSerialNo() uint64 {
	if m != nil {
		return m.SerialNo
	}
	return 0
}

func (m *UDPPacketData) GetTotalCnt() uint32 {
	if m != nil {
		return m.TotalCnt
	}
	return 0
}

func (m *UDPPacketData) GetPosNum() uint32 {
	if m != nil {
		return m.PosNum
	}
	return 0
}

func (m *UDPPacketData) GetDataType() uint32 {
	if m != nil {
		return m.DataType
	}
	return 0
}

func (m *UDPPacketData) GetTryCnt() uint32 {
	if m != nil {
		return m.TryCnt
	}
	return 0
}

func (m *UDPPacketData) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*UDPPacketData)(nil), "packet.UDPPacketData")
}

func init() { proto.RegisterFile("packet.proto", fileDescriptor_e9ef1a6541f9f9e7) }

var fileDescriptor_e9ef1a6541f9f9e7 = []byte{
	// 159 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x48, 0x4c, 0xce,
	0x4e, 0x2d, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x83, 0xf0, 0x94, 0x96, 0x33, 0x72,
	0xf1, 0x86, 0xba, 0x04, 0x04, 0x80, 0x79, 0x2e, 0x89, 0x25, 0x89, 0x42, 0x52, 0x5c, 0x1c, 0xc5,
	0xa9, 0x45, 0x99, 0x89, 0x39, 0x7e, 0xf9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x2c, 0x41, 0x70, 0x3e,
	0x48, 0xae, 0x24, 0xbf, 0x24, 0x31, 0xc7, 0x39, 0xaf, 0x44, 0x42, 0x40, 0x81, 0x51, 0x83, 0x37,
	0x08, 0xce, 0x17, 0x12, 0xe3, 0x62, 0x2b, 0xc8, 0x2f, 0xf6, 0x2b, 0xcd, 0x95, 0x60, 0x02, 0xcb,
	0x40, 0x79, 0x20, 0x3d, 0x29, 0x89, 0x25, 0x89, 0x21, 0x95, 0x05, 0xa9, 0x12, 0xcc, 0x10, 0x3d,
	0x30, 0x3e, 0x48, 0x4f, 0x49, 0x51, 0x25, 0xc8, 0x34, 0x41, 0x88, 0x1e, 0x08, 0x4f, 0x48, 0x88,
	0x8b, 0x05, 0xa4, 0x46, 0x82, 0x45, 0x81, 0x51, 0x83, 0x27, 0x08, 0xcc, 0x4e, 0x62, 0x03, 0x3b,
	0xdc, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xe4, 0x9a, 0xb2, 0xe9, 0xc8, 0x00, 0x00, 0x00,
}
