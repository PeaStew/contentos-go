// Code generated by protoc-gen-go. DO NOT EDIT.
// source: app/table/so_extHourTrx.proto

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

type SoExtHourTrx struct {
	Hour                 *prototype.TimePointSec `protobuf:"bytes,1,opt,name=hour,proto3" json:"hour,omitempty"`
	Count                uint64                  `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *SoExtHourTrx) Reset()         { *m = SoExtHourTrx{} }
func (m *SoExtHourTrx) String() string { return proto.CompactTextString(m) }
func (*SoExtHourTrx) ProtoMessage()    {}
func (*SoExtHourTrx) Descriptor() ([]byte, []int) {
	return fileDescriptor_ccee57d6353b6f3d, []int{0}
}

func (m *SoExtHourTrx) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoExtHourTrx.Unmarshal(m, b)
}
func (m *SoExtHourTrx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoExtHourTrx.Marshal(b, m, deterministic)
}
func (m *SoExtHourTrx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoExtHourTrx.Merge(m, src)
}
func (m *SoExtHourTrx) XXX_Size() int {
	return xxx_messageInfo_SoExtHourTrx.Size(m)
}
func (m *SoExtHourTrx) XXX_DiscardUnknown() {
	xxx_messageInfo_SoExtHourTrx.DiscardUnknown(m)
}

var xxx_messageInfo_SoExtHourTrx proto.InternalMessageInfo

func (m *SoExtHourTrx) GetHour() *prototype.TimePointSec {
	if m != nil {
		return m.Hour
	}
	return nil
}

func (m *SoExtHourTrx) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

type SoMemExtHourTrxByHour struct {
	Hour                 *prototype.TimePointSec `protobuf:"bytes,1,opt,name=hour,proto3" json:"hour,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *SoMemExtHourTrxByHour) Reset()         { *m = SoMemExtHourTrxByHour{} }
func (m *SoMemExtHourTrxByHour) String() string { return proto.CompactTextString(m) }
func (*SoMemExtHourTrxByHour) ProtoMessage()    {}
func (*SoMemExtHourTrxByHour) Descriptor() ([]byte, []int) {
	return fileDescriptor_ccee57d6353b6f3d, []int{1}
}

func (m *SoMemExtHourTrxByHour) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoMemExtHourTrxByHour.Unmarshal(m, b)
}
func (m *SoMemExtHourTrxByHour) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoMemExtHourTrxByHour.Marshal(b, m, deterministic)
}
func (m *SoMemExtHourTrxByHour) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoMemExtHourTrxByHour.Merge(m, src)
}
func (m *SoMemExtHourTrxByHour) XXX_Size() int {
	return xxx_messageInfo_SoMemExtHourTrxByHour.Size(m)
}
func (m *SoMemExtHourTrxByHour) XXX_DiscardUnknown() {
	xxx_messageInfo_SoMemExtHourTrxByHour.DiscardUnknown(m)
}

var xxx_messageInfo_SoMemExtHourTrxByHour proto.InternalMessageInfo

func (m *SoMemExtHourTrxByHour) GetHour() *prototype.TimePointSec {
	if m != nil {
		return m.Hour
	}
	return nil
}

type SoMemExtHourTrxByCount struct {
	Count                uint64   `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SoMemExtHourTrxByCount) Reset()         { *m = SoMemExtHourTrxByCount{} }
func (m *SoMemExtHourTrxByCount) String() string { return proto.CompactTextString(m) }
func (*SoMemExtHourTrxByCount) ProtoMessage()    {}
func (*SoMemExtHourTrxByCount) Descriptor() ([]byte, []int) {
	return fileDescriptor_ccee57d6353b6f3d, []int{2}
}

func (m *SoMemExtHourTrxByCount) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoMemExtHourTrxByCount.Unmarshal(m, b)
}
func (m *SoMemExtHourTrxByCount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoMemExtHourTrxByCount.Marshal(b, m, deterministic)
}
func (m *SoMemExtHourTrxByCount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoMemExtHourTrxByCount.Merge(m, src)
}
func (m *SoMemExtHourTrxByCount) XXX_Size() int {
	return xxx_messageInfo_SoMemExtHourTrxByCount.Size(m)
}
func (m *SoMemExtHourTrxByCount) XXX_DiscardUnknown() {
	xxx_messageInfo_SoMemExtHourTrxByCount.DiscardUnknown(m)
}

var xxx_messageInfo_SoMemExtHourTrxByCount proto.InternalMessageInfo

func (m *SoMemExtHourTrxByCount) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

type SoListExtHourTrxByHour struct {
	Hour                 *prototype.TimePointSec `protobuf:"bytes,1,opt,name=hour,proto3" json:"hour,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *SoListExtHourTrxByHour) Reset()         { *m = SoListExtHourTrxByHour{} }
func (m *SoListExtHourTrxByHour) String() string { return proto.CompactTextString(m) }
func (*SoListExtHourTrxByHour) ProtoMessage()    {}
func (*SoListExtHourTrxByHour) Descriptor() ([]byte, []int) {
	return fileDescriptor_ccee57d6353b6f3d, []int{3}
}

func (m *SoListExtHourTrxByHour) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoListExtHourTrxByHour.Unmarshal(m, b)
}
func (m *SoListExtHourTrxByHour) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoListExtHourTrxByHour.Marshal(b, m, deterministic)
}
func (m *SoListExtHourTrxByHour) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoListExtHourTrxByHour.Merge(m, src)
}
func (m *SoListExtHourTrxByHour) XXX_Size() int {
	return xxx_messageInfo_SoListExtHourTrxByHour.Size(m)
}
func (m *SoListExtHourTrxByHour) XXX_DiscardUnknown() {
	xxx_messageInfo_SoListExtHourTrxByHour.DiscardUnknown(m)
}

var xxx_messageInfo_SoListExtHourTrxByHour proto.InternalMessageInfo

func (m *SoListExtHourTrxByHour) GetHour() *prototype.TimePointSec {
	if m != nil {
		return m.Hour
	}
	return nil
}

type SoListExtHourTrxByCount struct {
	Count                uint64                  `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	Hour                 *prototype.TimePointSec `protobuf:"bytes,2,opt,name=hour,proto3" json:"hour,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *SoListExtHourTrxByCount) Reset()         { *m = SoListExtHourTrxByCount{} }
func (m *SoListExtHourTrxByCount) String() string { return proto.CompactTextString(m) }
func (*SoListExtHourTrxByCount) ProtoMessage()    {}
func (*SoListExtHourTrxByCount) Descriptor() ([]byte, []int) {
	return fileDescriptor_ccee57d6353b6f3d, []int{4}
}

func (m *SoListExtHourTrxByCount) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoListExtHourTrxByCount.Unmarshal(m, b)
}
func (m *SoListExtHourTrxByCount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoListExtHourTrxByCount.Marshal(b, m, deterministic)
}
func (m *SoListExtHourTrxByCount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoListExtHourTrxByCount.Merge(m, src)
}
func (m *SoListExtHourTrxByCount) XXX_Size() int {
	return xxx_messageInfo_SoListExtHourTrxByCount.Size(m)
}
func (m *SoListExtHourTrxByCount) XXX_DiscardUnknown() {
	xxx_messageInfo_SoListExtHourTrxByCount.DiscardUnknown(m)
}

var xxx_messageInfo_SoListExtHourTrxByCount proto.InternalMessageInfo

func (m *SoListExtHourTrxByCount) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *SoListExtHourTrxByCount) GetHour() *prototype.TimePointSec {
	if m != nil {
		return m.Hour
	}
	return nil
}

type SoUniqueExtHourTrxByHour struct {
	Hour                 *prototype.TimePointSec `protobuf:"bytes,1,opt,name=hour,proto3" json:"hour,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *SoUniqueExtHourTrxByHour) Reset()         { *m = SoUniqueExtHourTrxByHour{} }
func (m *SoUniqueExtHourTrxByHour) String() string { return proto.CompactTextString(m) }
func (*SoUniqueExtHourTrxByHour) ProtoMessage()    {}
func (*SoUniqueExtHourTrxByHour) Descriptor() ([]byte, []int) {
	return fileDescriptor_ccee57d6353b6f3d, []int{5}
}

func (m *SoUniqueExtHourTrxByHour) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoUniqueExtHourTrxByHour.Unmarshal(m, b)
}
func (m *SoUniqueExtHourTrxByHour) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoUniqueExtHourTrxByHour.Marshal(b, m, deterministic)
}
func (m *SoUniqueExtHourTrxByHour) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoUniqueExtHourTrxByHour.Merge(m, src)
}
func (m *SoUniqueExtHourTrxByHour) XXX_Size() int {
	return xxx_messageInfo_SoUniqueExtHourTrxByHour.Size(m)
}
func (m *SoUniqueExtHourTrxByHour) XXX_DiscardUnknown() {
	xxx_messageInfo_SoUniqueExtHourTrxByHour.DiscardUnknown(m)
}

var xxx_messageInfo_SoUniqueExtHourTrxByHour proto.InternalMessageInfo

func (m *SoUniqueExtHourTrxByHour) GetHour() *prototype.TimePointSec {
	if m != nil {
		return m.Hour
	}
	return nil
}

func init() {
	proto.RegisterType((*SoExtHourTrx)(nil), "table.so_extHourTrx")
	proto.RegisterType((*SoMemExtHourTrxByHour)(nil), "table.so_mem_extHourTrx_by_hour")
	proto.RegisterType((*SoMemExtHourTrxByCount)(nil), "table.so_mem_extHourTrx_by_count")
	proto.RegisterType((*SoListExtHourTrxByHour)(nil), "table.so_list_extHourTrx_by_hour")
	proto.RegisterType((*SoListExtHourTrxByCount)(nil), "table.so_list_extHourTrx_by_count")
	proto.RegisterType((*SoUniqueExtHourTrxByHour)(nil), "table.so_unique_extHourTrx_by_hour")
}

func init() { proto.RegisterFile("app/table/so_extHourTrx.proto", fileDescriptor_ccee57d6353b6f3d) }

var fileDescriptor_ccee57d6353b6f3d = []byte{
	// 245 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x92, 0xc1, 0x4a, 0x03, 0x31,
	0x10, 0x86, 0xd9, 0xd2, 0x7a, 0x18, 0xf1, 0xb2, 0xf4, 0xd0, 0x56, 0x85, 0xb2, 0xa7, 0x22, 0x76,
	0x03, 0xf5, 0x0d, 0x3c, 0x89, 0xe2, 0xa5, 0xf4, 0xe4, 0x25, 0x6c, 0xc2, 0xd0, 0x0d, 0x34, 0x99,
	0x98, 0x4c, 0xa0, 0x7d, 0x7b, 0x31, 0x85, 0xb5, 0xc2, 0x22, 0x8a, 0xbd, 0x04, 0x86, 0xc9, 0xff,
	0x7d, 0xff, 0x61, 0xe0, 0xb6, 0xf1, 0x5e, 0x70, 0xa3, 0x76, 0x28, 0x22, 0x49, 0xdc, 0xf3, 0x13,
	0xa5, 0xb0, 0x09, 0xfb, 0xda, 0x07, 0x62, 0x2a, 0x47, 0x79, 0x35, 0x1b, 0xe7, 0x89, 0x0f, 0x1e,
	0xc5, 0xe7, 0x73, 0x5c, 0x56, 0x1b, 0xb8, 0xfa, 0x96, 0x29, 0x97, 0x30, 0x6c, 0x29, 0x85, 0x49,
	0x31, 0x2f, 0x16, 0x97, 0xab, 0x69, 0xdd, 0xa5, 0x6a, 0x36, 0x16, 0xa5, 0x27, 0xe3, 0x58, 0x46,
	0xd4, 0xeb, 0xfc, 0xad, 0x1c, 0xc3, 0x48, 0x53, 0x72, 0x3c, 0x19, 0xcc, 0x8b, 0xc5, 0x70, 0x7d,
	0x1c, 0xaa, 0x67, 0x98, 0x46, 0x92, 0x16, 0xed, 0x09, 0x59, 0xaa, 0x83, 0xcc, 0x91, 0xbf, 0x19,
	0xaa, 0x15, 0xcc, 0x7a, 0x59, 0xd9, 0xf4, 0xe5, 0x2f, 0x4e, 0xfd, 0x2f, 0x39, 0xb3, 0x33, 0x91,
	0xcf, 0x50, 0x40, 0xc1, 0x75, 0x3f, 0xec, 0x87, 0x06, 0x9d, 0x63, 0xf0, 0x3b, 0xc7, 0x2b, 0xdc,
	0x44, 0x92, 0xc9, 0x99, 0xf7, 0x84, 0xff, 0xaf, 0xfc, 0x78, 0xff, 0x76, 0xb7, 0x35, 0xdc, 0x26,
	0x55, 0x6b, 0xb2, 0x42, 0x53, 0xd4, 0x6d, 0x63, 0x9c, 0xd0, 0xe4, 0x18, 0x1d, 0x53, 0x5c, 0x6e,
	0x49, 0x74, 0x47, 0xa3, 0x2e, 0x32, 0xed, 0xe1, 0x23, 0x00, 0x00, 0xff, 0xff, 0x8e, 0xda, 0x1b,
	0x24, 0x48, 0x02, 0x00, 0x00,
}
