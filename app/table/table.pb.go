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
	return fileDescriptor_448a2743262f7a00, []int{1}
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
	proto.RegisterType((*SoPost)(nil), "table.so_post")
}

func init() { proto.RegisterFile("table.proto", fileDescriptor_448a2743262f7a00) }

var fileDescriptor_448a2743262f7a00 = []byte{
	// 292 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0x4f, 0x4b, 0x33, 0x31,
	0x10, 0xc6, 0xd9, 0xfe, 0xdb, 0xb7, 0xd3, 0x57, 0x90, 0x5c, 0x8c, 0x05, 0x41, 0x7a, 0x12, 0xa4,
	0x5d, 0xb4, 0xe0, 0xa9, 0x37, 0x8f, 0xde, 0x16, 0xef, 0x61, 0x93, 0x0e, 0x6d, 0x68, 0x37, 0x13,
	0x9a, 0x2c, 0xd8, 0x8f, 0xe6, 0xa7, 0xf1, 0xab, 0x48, 0x66, 0xb7, 0xda, 0x93, 0x7a, 0xcb, 0x33,
	0xf3, 0x7b, 0x66, 0xe6, 0x21, 0x30, 0x89, 0x95, 0xde, 0xe3, 0xc2, 0x1f, 0x28, 0x92, 0x18, 0xb2,
	0x98, 0xae, 0x36, 0x36, 0x6e, 0x1b, 0xbd, 0x30, 0x54, 0x17, 0x86, 0x82, 0xd9, 0x56, 0xd6, 0x15,
	0x86, 0x5c, 0x44, 0x17, 0x29, 0xcc, 0x37, 0x54, 0x30, 0x5e, 0xc4, 0xa3, 0xc7, 0xf9, 0xf7, 0xb3,
	0x1d, 0x32, 0xfb, 0xc8, 0x00, 0x02, 0xa9, 0xca, 0x18, 0x6a, 0x5c, 0x14, 0xf7, 0x30, 0x70, 0x55,
	0x8d, 0x32, 0xbb, 0xcd, 0xee, 0x26, 0x8f, 0x57, 0x2d, 0xc4, 0x78, 0x47, 0xa8, 0xd4, 0x2e, 0x19,
	0x12, 0x2b, 0xf8, 0x6f, 0x0e, 0x58, 0x45, 0x5c, 0xab, 0x68, 0x6b, 0x94, 0x3d, 0x36, 0x5d, 0x9f,
	0x99, 0x52, 0x59, 0x79, 0xb2, 0x2e, 0xaa, 0x80, 0xa6, 0x9c, 0x74, 0xf8, 0xab, 0xad, 0x51, 0x3c,
	0x40, 0xce, 0x92, 0x0e, 0xb2, 0xff, 0xf3, 0xb6, 0x13, 0x27, 0x96, 0x90, 0xfb, 0x46, 0xab, 0x1d,
	0x1e, 0xe5, 0x80, 0x2d, 0xd3, 0x33, 0x8b, 0x6f, 0xf4, 0xde, 0x9a, 0xd4, 0x54, 0x49, 0x97, 0x23,
	0xdf, 0xe8, 0x17, 0x3c, 0xce, 0xde, 0x33, 0xc8, 0x03, 0x29, 0x4f, 0x21, 0x8a, 0x4b, 0xe8, 0xdb,
	0xf5, 0x1b, 0xa7, 0xbb, 0x28, 0xd3, 0xf3, 0x2b, 0x70, 0xef, 0x2f, 0x81, 0x9f, 0x60, 0x9c, 0xc6,
	0xb4, 0x69, 0xfb, 0xbf, 0xa5, 0xfd, 0x97, 0x58, 0x8e, 0x2a, 0x21, 0xef, 0x3e, 0x84, 0xef, 0x1e,
	0x97, 0x27, 0x29, 0x6e, 0x00, 0xf6, 0x76, 0x87, 0x8a, 0x57, 0xc9, 0x21, 0xdf, 0x35, 0x4e, 0x95,
	0xe7, 0x54, 0xd0, 0x23, 0x1e, 0xbe, 0xfc, 0x0c, 0x00, 0x00, 0xff, 0xff, 0x1e, 0x45, 0x7a, 0xc0,
	0xf8, 0x01, 0x00, 0x00,
}
