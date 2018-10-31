// Code generated by protoc-gen-go. DO NOT EDIT.
// source: table.proto

package table

import (
	fmt "fmt"
	type_proto "github.com/coschain/contentos-go/proto/type-proto"
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

type SoAccount struct {
	Name                 *type_proto.AccountName   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	CreatedTime          *type_proto.TimePointSec  `protobuf:"bytes,2,opt,name=created_time,json=createdTime,proto3" json:"created_time,omitempty"`
	Creator              *type_proto.AccountName   `protobuf:"bytes,3,opt,name=creator,proto3" json:"creator,omitempty"`
	PubKey               *type_proto.PublicKeyType `protobuf:"bytes,4,opt,name=pub_key,json=pubKey,proto3" json:"pub_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *SoAccount) Reset()         { *m = SoAccount{} }
func (m *SoAccount) String() string { return proto.CompactTextString(m) }
func (*SoAccount) ProtoMessage()    {}
func (*SoAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_448a2743262f7a00, []int{0}
}

func (m *SoAccount) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoAccount.Unmarshal(m, b)
}
func (m *SoAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoAccount.Marshal(b, m, deterministic)
}
func (m *SoAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoAccount.Merge(m, src)
}
func (m *SoAccount) XXX_Size() int {
	return xxx_messageInfo_SoAccount.Size(m)
}
func (m *SoAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_SoAccount.DiscardUnknown(m)
}

var xxx_messageInfo_SoAccount proto.InternalMessageInfo

func (m *SoAccount) GetName() *type_proto.AccountName {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *SoAccount) GetCreatedTime() *type_proto.TimePointSec {
	if m != nil {
		return m.CreatedTime
	}
	return nil
}

func (m *SoAccount) GetCreator() *type_proto.AccountName {
	if m != nil {
		return m.Creator
	}
	return nil
}

func (m *SoAccount) GetPubKey() *type_proto.PublicKeyType {
	if m != nil {
		return m.PubKey
	}
	return nil
}

type SKeyAccountByCreatedTime struct {
	CreatedTime          *type_proto.TimePointSec `protobuf:"bytes,1,opt,name=created_time,json=createdTime,proto3" json:"created_time,omitempty"`
	Name                 *type_proto.AccountName  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *SKeyAccountByCreatedTime) Reset()         { *m = SKeyAccountByCreatedTime{} }
func (m *SKeyAccountByCreatedTime) String() string { return proto.CompactTextString(m) }
func (*SKeyAccountByCreatedTime) ProtoMessage()    {}
func (*SKeyAccountByCreatedTime) Descriptor() ([]byte, []int) {
	return fileDescriptor_448a2743262f7a00, []int{1}
}

func (m *SKeyAccountByCreatedTime) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SKeyAccountByCreatedTime.Unmarshal(m, b)
}
func (m *SKeyAccountByCreatedTime) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SKeyAccountByCreatedTime.Marshal(b, m, deterministic)
}
func (m *SKeyAccountByCreatedTime) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SKeyAccountByCreatedTime.Merge(m, src)
}
func (m *SKeyAccountByCreatedTime) XXX_Size() int {
	return xxx_messageInfo_SKeyAccountByCreatedTime.Size(m)
}
func (m *SKeyAccountByCreatedTime) XXX_DiscardUnknown() {
	xxx_messageInfo_SKeyAccountByCreatedTime.DiscardUnknown(m)
}

var xxx_messageInfo_SKeyAccountByCreatedTime proto.InternalMessageInfo

func (m *SKeyAccountByCreatedTime) GetCreatedTime() *type_proto.TimePointSec {
	if m != nil {
		return m.CreatedTime
	}
	return nil
}

func (m *SKeyAccountByCreatedTime) GetName() *type_proto.AccountName {
	if m != nil {
		return m.Name
	}
	return nil
}

type SoKeyAccount struct {
	PubKey               *type_proto.PublicKeyType `protobuf:"bytes,1,opt,name=pub_key,json=pubKey,proto3" json:"pub_key,omitempty"`
	Type                 int32                     `protobuf:"varint,2,opt,name=type,proto3" json:"type,omitempty"`
	Name                 *type_proto.AccountName   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *SoKeyAccount) Reset()         { *m = SoKeyAccount{} }
func (m *SoKeyAccount) String() string { return proto.CompactTextString(m) }
func (*SoKeyAccount) ProtoMessage()    {}
func (*SoKeyAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_448a2743262f7a00, []int{2}
}

func (m *SoKeyAccount) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoKeyAccount.Unmarshal(m, b)
}
func (m *SoKeyAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoKeyAccount.Marshal(b, m, deterministic)
}
func (m *SoKeyAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoKeyAccount.Merge(m, src)
}
func (m *SoKeyAccount) XXX_Size() int {
	return xxx_messageInfo_SoKeyAccount.Size(m)
}
func (m *SoKeyAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_SoKeyAccount.DiscardUnknown(m)
}

var xxx_messageInfo_SoKeyAccount proto.InternalMessageInfo

func (m *SoKeyAccount) GetPubKey() *type_proto.PublicKeyType {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *SoKeyAccount) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *SoKeyAccount) GetName() *type_proto.AccountName {
	if m != nil {
		return m.Name
	}
	return nil
}

type SoPost struct {
	Idx                  uint32                   `protobuf:"varint,1,opt,name=idx,proto3" json:"idx,omitempty"`
	Name                 *type_proto.AccountName  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	PostTime             *type_proto.TimePointSec `protobuf:"bytes,3,opt,name=post_time,json=postTime,proto3" json:"post_time,omitempty"`
	Content              string                   `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	LikeCount            uint32                   `protobuf:"varint,5,opt,name=like_count,json=likeCount,proto3" json:"like_count,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *SoPost) Reset()         { *m = SoPost{} }
func (m *SoPost) String() string { return proto.CompactTextString(m) }
func (*SoPost) ProtoMessage()    {}
func (*SoPost) Descriptor() ([]byte, []int) {
	return fileDescriptor_448a2743262f7a00, []int{3}
}

func (m *SoPost) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoPost.Unmarshal(m, b)
}
func (m *SoPost) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoPost.Marshal(b, m, deterministic)
}
func (m *SoPost) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoPost.Merge(m, src)
}
func (m *SoPost) XXX_Size() int {
	return xxx_messageInfo_SoPost.Size(m)
}
func (m *SoPost) XXX_DiscardUnknown() {
	xxx_messageInfo_SoPost.DiscardUnknown(m)
}

var xxx_messageInfo_SoPost proto.InternalMessageInfo

func (m *SoPost) GetIdx() uint32 {
	if m != nil {
		return m.Idx
	}
	return 0
}

func (m *SoPost) GetName() *type_proto.AccountName {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *SoPost) GetPostTime() *type_proto.TimePointSec {
	if m != nil {
		return m.PostTime
	}
	return nil
}

func (m *SoPost) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *SoPost) GetLikeCount() uint32 {
	if m != nil {
		return m.LikeCount
	}
	return 0
}

func init() {
	proto.RegisterType((*SoAccount)(nil), "table.so_account")
	proto.RegisterType((*SKeyAccountByCreatedTime)(nil), "table.s_key_account_by_created_time")
	proto.RegisterType((*SoKeyAccount)(nil), "table.so_key_account")
	proto.RegisterType((*SoPost)(nil), "table.so_post")
}

func init() { proto.RegisterFile("table.proto", fileDescriptor_448a2743262f7a00) }

var fileDescriptor_448a2743262f7a00 = []byte{
	// 345 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0x41, 0x4b, 0xf3, 0x40,
	0x10, 0x65, 0x9b, 0xb6, 0xf9, 0x3a, 0xfd, 0x14, 0xd9, 0x8b, 0xb1, 0x50, 0x90, 0x9e, 0x04, 0x69,
	0x83, 0x16, 0x3c, 0xf5, 0xe6, 0xd1, 0x5b, 0xf0, 0xbe, 0x64, 0xb7, 0x4b, 0xbb, 0xb4, 0xc9, 0x2c,
	0xdd, 0x0d, 0x98, 0x3f, 0xe0, 0xc1, 0x7f, 0xe4, 0xaf, 0xf1, 0xaf, 0xc8, 0x4e, 0x1b, 0x0d, 0x08,
	0xda, 0xde, 0xe6, 0xcd, 0xbe, 0x99, 0x79, 0xef, 0x25, 0x30, 0xf4, 0xb9, 0xdc, 0xea, 0x99, 0xdd,
	0xa1, 0x47, 0xde, 0x23, 0x30, 0x5a, 0xac, 0x8c, 0x5f, 0x57, 0x72, 0xa6, 0xb0, 0x48, 0x15, 0x3a,
	0xb5, 0xce, 0x4d, 0x99, 0x2a, 0x2c, 0xbd, 0x2e, 0x3d, 0xba, 0xe9, 0x0a, 0x53, 0xa2, 0xa7, 0xbe,
	0xb6, 0x7a, 0xfa, 0x5d, 0xee, 0x97, 0x4c, 0x3e, 0x18, 0x80, 0x43, 0x91, 0x2b, 0x85, 0x55, 0xe9,
	0xf9, 0x2d, 0x74, 0xcb, 0xbc, 0xd0, 0x09, 0xbb, 0x66, 0x37, 0xc3, 0xfb, 0xcb, 0x3d, 0x89, 0xe8,
	0x07, 0x86, 0x08, 0xcf, 0x19, 0x91, 0xf8, 0x02, 0xfe, 0xab, 0x9d, 0xce, 0xbd, 0x5e, 0x0a, 0x6f,
	0x0a, 0x9d, 0x74, 0x68, 0xe8, 0xaa, 0x35, 0x14, 0xda, 0xc2, 0xa2, 0x29, 0xbd, 0x70, 0x5a, 0x65,
	0xc3, 0x03, 0xfd, 0xd9, 0x14, 0x9a, 0xdf, 0x41, 0x4c, 0x10, 0x77, 0x49, 0xf4, 0xfb, 0xb5, 0x86,
	0xc7, 0xe7, 0x10, 0xdb, 0x4a, 0x8a, 0x8d, 0xae, 0x93, 0x2e, 0x8d, 0x8c, 0x5a, 0x23, 0xb6, 0x92,
	0x5b, 0xa3, 0xc2, 0xa3, 0x08, 0x38, 0xeb, 0xdb, 0x4a, 0x3e, 0xe9, 0x7a, 0xf2, 0xc6, 0x60, 0xec,
	0xa8, 0xdd, 0x2c, 0x95, 0xb5, 0x68, 0xeb, 0xfe, 0xe1, 0x83, 0x9d, 0xe4, 0xa3, 0x89, 0xac, 0x73,
	0x44, 0x64, 0x93, 0x57, 0x06, 0xe7, 0x0e, 0xdb, 0x6a, 0xda, 0xa6, 0xd8, 0xb1, 0xa6, 0x38, 0x87,
	0x6e, 0xc0, 0x74, 0xb4, 0x97, 0x51, 0xfd, 0x25, 0x24, 0x3a, 0x46, 0xc8, 0x3b, 0x83, 0xd8, 0xa1,
	0xb0, 0xe8, 0x3c, 0xbf, 0x80, 0xc8, 0x2c, 0x5f, 0xe8, 0xfa, 0x59, 0x16, 0xca, 0x93, 0x3c, 0xf1,
	0x07, 0x18, 0x84, 0x35, 0xfb, 0xec, 0xa2, 0xbf, 0xb2, 0xfb, 0x17, 0xb8, 0x14, 0x5c, 0x02, 0xf1,
	0xe1, 0x37, 0xa5, 0xaf, 0x39, 0xc8, 0x1a, 0xc8, 0xc7, 0x00, 0x5b, 0xb3, 0xd1, 0x82, 0x4e, 0x25,
	0x3d, 0xd2, 0x35, 0x08, 0x9d, 0xc7, 0xd0, 0x90, 0x7d, 0x5a, 0x3e, 0xff, 0x0c, 0x00, 0x00, 0xff,
	0xff, 0x2b, 0x82, 0x47, 0x2a, 0x0e, 0x03, 0x00, 0x00,
}
