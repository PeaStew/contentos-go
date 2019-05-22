// Code generated by protoc-gen-go. DO NOT EDIT.
// source: app/table/so_global.proto

package table

import (
	fmt "fmt"
	prototype "github.com/coschain/contentos-go/prototype"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type SoGlobal struct {
	Id                   int32                        `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Props                *prototype.DynamicProperties `protobuf:"bytes,2,opt,name=props,proto3" json:"props,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *SoGlobal) Reset()         { *m = SoGlobal{} }
func (m *SoGlobal) String() string { return proto.CompactTextString(m) }
func (*SoGlobal) ProtoMessage()    {}
func (*SoGlobal) Descriptor() ([]byte, []int) {
	return fileDescriptor_8b138d292b065e72, []int{0}
}

func (m *SoGlobal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoGlobal.Unmarshal(m, b)
}
func (m *SoGlobal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoGlobal.Marshal(b, m, deterministic)
}
func (m *SoGlobal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoGlobal.Merge(m, src)
}
func (m *SoGlobal) XXX_Size() int {
	return xxx_messageInfo_SoGlobal.Size(m)
}
func (m *SoGlobal) XXX_DiscardUnknown() {
	xxx_messageInfo_SoGlobal.DiscardUnknown(m)
}

var xxx_messageInfo_SoGlobal proto.InternalMessageInfo

func (m *SoGlobal) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *SoGlobal) GetProps() *prototype.DynamicProperties {
	if m != nil {
		return m.Props
	}
	return nil
}

type SoMemGlobalById struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SoMemGlobalById) Reset()         { *m = SoMemGlobalById{} }
func (m *SoMemGlobalById) String() string { return proto.CompactTextString(m) }
func (*SoMemGlobalById) ProtoMessage()    {}
func (*SoMemGlobalById) Descriptor() ([]byte, []int) {
	return fileDescriptor_8b138d292b065e72, []int{1}
}

func (m *SoMemGlobalById) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoMemGlobalById.Unmarshal(m, b)
}
func (m *SoMemGlobalById) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoMemGlobalById.Marshal(b, m, deterministic)
}
func (m *SoMemGlobalById) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoMemGlobalById.Merge(m, src)
}
func (m *SoMemGlobalById) XXX_Size() int {
	return xxx_messageInfo_SoMemGlobalById.Size(m)
}
func (m *SoMemGlobalById) XXX_DiscardUnknown() {
	xxx_messageInfo_SoMemGlobalById.DiscardUnknown(m)
}

var xxx_messageInfo_SoMemGlobalById proto.InternalMessageInfo

func (m *SoMemGlobalById) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

type SoMemGlobalByProps struct {
	Props                *prototype.DynamicProperties `protobuf:"bytes,1,opt,name=props,proto3" json:"props,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *SoMemGlobalByProps) Reset()         { *m = SoMemGlobalByProps{} }
func (m *SoMemGlobalByProps) String() string { return proto.CompactTextString(m) }
func (*SoMemGlobalByProps) ProtoMessage()    {}
func (*SoMemGlobalByProps) Descriptor() ([]byte, []int) {
	return fileDescriptor_8b138d292b065e72, []int{2}
}

func (m *SoMemGlobalByProps) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoMemGlobalByProps.Unmarshal(m, b)
}
func (m *SoMemGlobalByProps) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoMemGlobalByProps.Marshal(b, m, deterministic)
}
func (m *SoMemGlobalByProps) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoMemGlobalByProps.Merge(m, src)
}
func (m *SoMemGlobalByProps) XXX_Size() int {
	return xxx_messageInfo_SoMemGlobalByProps.Size(m)
}
func (m *SoMemGlobalByProps) XXX_DiscardUnknown() {
	xxx_messageInfo_SoMemGlobalByProps.DiscardUnknown(m)
}

var xxx_messageInfo_SoMemGlobalByProps proto.InternalMessageInfo

func (m *SoMemGlobalByProps) GetProps() *prototype.DynamicProperties {
	if m != nil {
		return m.Props
	}
	return nil
}

type SoUniqueGlobalById struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SoUniqueGlobalById) Reset()         { *m = SoUniqueGlobalById{} }
func (m *SoUniqueGlobalById) String() string { return proto.CompactTextString(m) }
func (*SoUniqueGlobalById) ProtoMessage()    {}
func (*SoUniqueGlobalById) Descriptor() ([]byte, []int) {
	return fileDescriptor_8b138d292b065e72, []int{3}
}

func (m *SoUniqueGlobalById) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoUniqueGlobalById.Unmarshal(m, b)
}
func (m *SoUniqueGlobalById) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoUniqueGlobalById.Marshal(b, m, deterministic)
}
func (m *SoUniqueGlobalById) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoUniqueGlobalById.Merge(m, src)
}
func (m *SoUniqueGlobalById) XXX_Size() int {
	return xxx_messageInfo_SoUniqueGlobalById.Size(m)
}
func (m *SoUniqueGlobalById) XXX_DiscardUnknown() {
	xxx_messageInfo_SoUniqueGlobalById.DiscardUnknown(m)
}

var xxx_messageInfo_SoUniqueGlobalById proto.InternalMessageInfo

func (m *SoUniqueGlobalById) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func init() {
	proto.RegisterType((*SoGlobal)(nil), "table.so_global")
	proto.RegisterType((*SoMemGlobalById)(nil), "table.so_mem_global_by_id")
	proto.RegisterType((*SoMemGlobalByProps)(nil), "table.so_mem_global_by_props")
	proto.RegisterType((*SoUniqueGlobalById)(nil), "table.so_unique_global_by_id")
}

func init() { proto.RegisterFile("app/table/so_global.proto", fileDescriptor_8b138d292b065e72) }

var fileDescriptor_8b138d292b065e72 = []byte{
	// 216 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x90, 0xcd, 0x4a, 0xc5, 0x30,
	0x10, 0x85, 0x49, 0xa1, 0x82, 0x11, 0x5c, 0x54, 0x91, 0xab, 0x20, 0x5c, 0x0a, 0x42, 0x11, 0x4d,
	0xc0, 0xbe, 0x81, 0x7b, 0x41, 0xba, 0x74, 0x13, 0xf2, 0x47, 0x3b, 0xd0, 0x64, 0x62, 0x93, 0x2e,
	0xfa, 0xf6, 0xd2, 0x54, 0x2a, 0x28, 0x08, 0x77, 0x33, 0x30, 0x73, 0x0e, 0xdf, 0x07, 0x43, 0x6f,
	0x65, 0x08, 0x3c, 0x49, 0x35, 0x5a, 0x1e, 0x51, 0xf4, 0x23, 0x2a, 0x39, 0xb2, 0x30, 0x61, 0xc2,
	0xaa, 0xcc, 0xe7, 0xbb, 0xeb, 0xbc, 0xa5, 0x25, 0x58, 0xbe, 0x8e, 0x2d, 0xac, 0xdf, 0xe9, 0xf9,
	0xde, 0xaf, 0x2e, 0x69, 0x01, 0xe6, 0x40, 0x8e, 0xa4, 0x29, 0xbb, 0x02, 0x4c, 0xd5, 0xd2, 0x32,
	0x4c, 0x18, 0xe2, 0xa1, 0x38, 0x92, 0xe6, 0xe2, 0xe5, 0x9e, 0xed, 0x08, 0x66, 0x16, 0x2f, 0x1d,
	0x68, 0xb1, 0xe6, 0x76, 0x4a, 0x60, 0x63, 0xb7, 0x75, 0xeb, 0x07, 0x7a, 0x15, 0x51, 0x38, 0xeb,
	0xbe, 0xa9, 0x42, 0x2d, 0x02, 0xcc, 0x6f, 0x76, 0xfd, 0x46, 0x6f, 0xfe, 0xd4, 0x32, 0xe0, 0xc7,
	0x4a, 0x4e, 0xb0, 0x36, 0x19, 0x37, 0x7b, 0xf8, 0x9c, 0xed, 0xbf, 0xe2, 0xd7, 0xa7, 0x8f, 0xc7,
	0x1e, 0xd2, 0x30, 0x2b, 0xa6, 0xd1, 0x71, 0x8d, 0x51, 0x0f, 0x12, 0x3c, 0xd7, 0xe8, 0x93, 0xf5,
	0x09, 0xe3, 0x73, 0x8f, 0x7c, 0x7f, 0xa6, 0x3a, 0xcb, 0xf2, 0xf6, 0x2b, 0x00, 0x00, 0xff, 0xff,
	0x25, 0x65, 0x59, 0x1a, 0x60, 0x01, 0x00, 0x00,
}
