// Code generated by protoc-gen-go. DO NOT EDIT.
// source: packet/packet.proto

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
	RcvSn                uint64   `protobuf:"varint,2,opt,name=rcvSn,proto3" json:"rcvSn,omitempty"`
	TotalCnt             uint32   `protobuf:"varint,16,opt,name=totalCnt,proto3" json:"totalCnt,omitempty"`
	PosNum               uint32   `protobuf:"varint,3,opt,name=posNum,proto3" json:"posNum,omitempty"`
	DataType             uint32   `protobuf:"varint,4,opt,name=dataType,proto3" json:"dataType,omitempty"`
	TryCnt               uint32   `protobuf:"varint,17,opt,name=tryCnt,proto3" json:"tryCnt,omitempty"`
	Len                  int32    `protobuf:"varint,5,opt,name=len,proto3" json:"len,omitempty"`
	TransInfo            []byte   `protobuf:"bytes,18,opt,name=transInfo,proto3" json:"transInfo,omitempty"`
	Finished             bool     `protobuf:"varint,19,opt,name=finished,proto3" json:"finished,omitempty"`
	Data                 []byte   `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UDPPacketData) Reset()         { *m = UDPPacketData{} }
func (m *UDPPacketData) String() string { return proto.CompactTextString(m) }
func (*UDPPacketData) ProtoMessage()    {}
func (*UDPPacketData) Descriptor() ([]byte, []int) {
	return fileDescriptor_57dbb8dc3dbf2351, []int{0}
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

func (m *UDPPacketData) GetRcvSn() uint64 {
	if m != nil {
		return m.RcvSn
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

func (m *UDPPacketData) GetLen() int32 {
	if m != nil {
		return m.Len
	}
	return 0
}

func (m *UDPPacketData) GetTransInfo() []byte {
	if m != nil {
		return m.TransInfo
	}
	return nil
}

func (m *UDPPacketData) GetFinished() bool {
	if m != nil {
		return m.Finished
	}
	return false
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

func init() { proto.RegisterFile("packet/packet.proto", fileDescriptor_57dbb8dc3dbf2351) }

var fileDescriptor_57dbb8dc3dbf2351 = []byte{
	// 224 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0xc1, 0x4a, 0x03, 0x31,
	0x10, 0x86, 0x49, 0xbb, 0xbb, 0xd4, 0xc1, 0x42, 0x9d, 0x8a, 0x0c, 0xe2, 0x21, 0x78, 0xca, 0x49,
	0x0f, 0x3e, 0x82, 0xbd, 0x78, 0x29, 0x25, 0xea, 0x03, 0x8c, 0x6d, 0x8a, 0x8b, 0x6b, 0xb2, 0x64,
	0x47, 0x61, 0x1f, 0xc1, 0xb7, 0x96, 0x24, 0xba, 0x3d, 0xe5, 0xff, 0xf2, 0xf3, 0xe5, 0x87, 0xc0,
	0xba, 0xe7, 0xfd, 0x87, 0x93, 0xfb, 0x72, 0xdc, 0xf5, 0x31, 0x48, 0xc0, 0xa6, 0xd0, 0xed, 0xcf,
	0x0c, 0x96, 0xaf, 0x9b, 0xdd, 0x2e, 0xd3, 0x86, 0x85, 0xf1, 0x1a, 0x16, 0x83, 0x8b, 0x2d, 0x77,
	0xdb, 0x40, 0x4a, 0x2b, 0x53, 0xd9, 0x89, 0xf1, 0x12, 0xea, 0xb8, 0xff, 0x7e, 0xf6, 0x34, 0xcb,
	0x45, 0x81, 0x64, 0x48, 0x10, 0xee, 0x1e, 0xbd, 0xd0, 0x4a, 0x2b, 0xb3, 0xb4, 0x13, 0xe3, 0x15,
	0x34, 0x7d, 0x18, 0xb6, 0x5f, 0x9f, 0x34, 0xcf, 0xcd, 0x1f, 0x25, 0xe7, 0xc0, 0xc2, 0x2f, 0x63,
	0xef, 0xa8, 0x2a, 0xce, 0x3f, 0x27, 0x47, 0xe2, 0x98, 0x5e, 0xbb, 0x28, 0x4e, 0x21, 0x5c, 0xc1,
	0xbc, 0x73, 0x9e, 0x6a, 0xad, 0x4c, 0x6d, 0x53, 0xc4, 0x1b, 0x38, 0x93, 0xc8, 0x7e, 0x78, 0xf2,
	0xc7, 0x40, 0xa8, 0x95, 0x39, 0xb7, 0xa7, 0x8b, 0xb4, 0x71, 0x6c, 0x7d, 0x3b, 0xbc, 0xbb, 0x03,
	0xad, 0xb5, 0x32, 0x0b, 0x3b, 0x31, 0x22, 0x54, 0x69, 0x8f, 0x9a, 0x2c, 0xe5, 0xfc, 0xd6, 0xe4,
	0xaf, 0x79, 0xf8, 0x0d, 0x00, 0x00, 0xff, 0xff, 0xc5, 0xce, 0x56, 0x45, 0x31, 0x01, 0x00, 0x00,
}
