// Code generated by protoc-gen-go. DO NOT EDIT.
// source: app/table/so_blockProducerVote.proto

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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type SoBlockProducerVote struct {
	VoterId              *prototype.BpVoterId         `protobuf:"bytes,1,opt,name=voter_id,json=voterId,proto3" json:"voter_id,omitempty"`
	BlockProducerId      *prototype.BpBlockProducerId `protobuf:"bytes,2,opt,name=block_producer_id,json=blockProducerId,proto3" json:"block_producer_id,omitempty"`
	VoteTime             *prototype.TimePointSec      `protobuf:"bytes,3,opt,name=vote_time,json=voteTime,proto3" json:"vote_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *SoBlockProducerVote) Reset()         { *m = SoBlockProducerVote{} }
func (m *SoBlockProducerVote) String() string { return proto.CompactTextString(m) }
func (*SoBlockProducerVote) ProtoMessage()    {}
func (*SoBlockProducerVote) Descriptor() ([]byte, []int) {
	return fileDescriptor_5d279fa3103478f1, []int{0}
}

func (m *SoBlockProducerVote) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoBlockProducerVote.Unmarshal(m, b)
}
func (m *SoBlockProducerVote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoBlockProducerVote.Marshal(b, m, deterministic)
}
func (m *SoBlockProducerVote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoBlockProducerVote.Merge(m, src)
}
func (m *SoBlockProducerVote) XXX_Size() int {
	return xxx_messageInfo_SoBlockProducerVote.Size(m)
}
func (m *SoBlockProducerVote) XXX_DiscardUnknown() {
	xxx_messageInfo_SoBlockProducerVote.DiscardUnknown(m)
}

var xxx_messageInfo_SoBlockProducerVote proto.InternalMessageInfo

func (m *SoBlockProducerVote) GetVoterId() *prototype.BpVoterId {
	if m != nil {
		return m.VoterId
	}
	return nil
}

func (m *SoBlockProducerVote) GetBlockProducerId() *prototype.BpBlockProducerId {
	if m != nil {
		return m.BlockProducerId
	}
	return nil
}

func (m *SoBlockProducerVote) GetVoteTime() *prototype.TimePointSec {
	if m != nil {
		return m.VoteTime
	}
	return nil
}

type SoMemBlockProducerVoteByVoterId struct {
	VoterId              *prototype.BpVoterId `protobuf:"bytes,1,opt,name=voter_id,json=voterId,proto3" json:"voter_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *SoMemBlockProducerVoteByVoterId) Reset()         { *m = SoMemBlockProducerVoteByVoterId{} }
func (m *SoMemBlockProducerVoteByVoterId) String() string { return proto.CompactTextString(m) }
func (*SoMemBlockProducerVoteByVoterId) ProtoMessage()    {}
func (*SoMemBlockProducerVoteByVoterId) Descriptor() ([]byte, []int) {
	return fileDescriptor_5d279fa3103478f1, []int{1}
}

func (m *SoMemBlockProducerVoteByVoterId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoMemBlockProducerVoteByVoterId.Unmarshal(m, b)
}
func (m *SoMemBlockProducerVoteByVoterId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoMemBlockProducerVoteByVoterId.Marshal(b, m, deterministic)
}
func (m *SoMemBlockProducerVoteByVoterId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoMemBlockProducerVoteByVoterId.Merge(m, src)
}
func (m *SoMemBlockProducerVoteByVoterId) XXX_Size() int {
	return xxx_messageInfo_SoMemBlockProducerVoteByVoterId.Size(m)
}
func (m *SoMemBlockProducerVoteByVoterId) XXX_DiscardUnknown() {
	xxx_messageInfo_SoMemBlockProducerVoteByVoterId.DiscardUnknown(m)
}

var xxx_messageInfo_SoMemBlockProducerVoteByVoterId proto.InternalMessageInfo

func (m *SoMemBlockProducerVoteByVoterId) GetVoterId() *prototype.BpVoterId {
	if m != nil {
		return m.VoterId
	}
	return nil
}

type SoMemBlockProducerVoteByBlockProducerId struct {
	BlockProducerId      *prototype.BpBlockProducerId `protobuf:"bytes,1,opt,name=block_producer_id,json=blockProducerId,proto3" json:"block_producer_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *SoMemBlockProducerVoteByBlockProducerId) Reset() {
	*m = SoMemBlockProducerVoteByBlockProducerId{}
}
func (m *SoMemBlockProducerVoteByBlockProducerId) String() string { return proto.CompactTextString(m) }
func (*SoMemBlockProducerVoteByBlockProducerId) ProtoMessage()    {}
func (*SoMemBlockProducerVoteByBlockProducerId) Descriptor() ([]byte, []int) {
	return fileDescriptor_5d279fa3103478f1, []int{2}
}

func (m *SoMemBlockProducerVoteByBlockProducerId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoMemBlockProducerVoteByBlockProducerId.Unmarshal(m, b)
}
func (m *SoMemBlockProducerVoteByBlockProducerId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoMemBlockProducerVoteByBlockProducerId.Marshal(b, m, deterministic)
}
func (m *SoMemBlockProducerVoteByBlockProducerId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoMemBlockProducerVoteByBlockProducerId.Merge(m, src)
}
func (m *SoMemBlockProducerVoteByBlockProducerId) XXX_Size() int {
	return xxx_messageInfo_SoMemBlockProducerVoteByBlockProducerId.Size(m)
}
func (m *SoMemBlockProducerVoteByBlockProducerId) XXX_DiscardUnknown() {
	xxx_messageInfo_SoMemBlockProducerVoteByBlockProducerId.DiscardUnknown(m)
}

var xxx_messageInfo_SoMemBlockProducerVoteByBlockProducerId proto.InternalMessageInfo

func (m *SoMemBlockProducerVoteByBlockProducerId) GetBlockProducerId() *prototype.BpBlockProducerId {
	if m != nil {
		return m.BlockProducerId
	}
	return nil
}

type SoMemBlockProducerVoteByVoteTime struct {
	VoteTime             *prototype.TimePointSec `protobuf:"bytes,1,opt,name=vote_time,json=voteTime,proto3" json:"vote_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *SoMemBlockProducerVoteByVoteTime) Reset()         { *m = SoMemBlockProducerVoteByVoteTime{} }
func (m *SoMemBlockProducerVoteByVoteTime) String() string { return proto.CompactTextString(m) }
func (*SoMemBlockProducerVoteByVoteTime) ProtoMessage()    {}
func (*SoMemBlockProducerVoteByVoteTime) Descriptor() ([]byte, []int) {
	return fileDescriptor_5d279fa3103478f1, []int{3}
}

func (m *SoMemBlockProducerVoteByVoteTime) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoMemBlockProducerVoteByVoteTime.Unmarshal(m, b)
}
func (m *SoMemBlockProducerVoteByVoteTime) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoMemBlockProducerVoteByVoteTime.Marshal(b, m, deterministic)
}
func (m *SoMemBlockProducerVoteByVoteTime) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoMemBlockProducerVoteByVoteTime.Merge(m, src)
}
func (m *SoMemBlockProducerVoteByVoteTime) XXX_Size() int {
	return xxx_messageInfo_SoMemBlockProducerVoteByVoteTime.Size(m)
}
func (m *SoMemBlockProducerVoteByVoteTime) XXX_DiscardUnknown() {
	xxx_messageInfo_SoMemBlockProducerVoteByVoteTime.DiscardUnknown(m)
}

var xxx_messageInfo_SoMemBlockProducerVoteByVoteTime proto.InternalMessageInfo

func (m *SoMemBlockProducerVoteByVoteTime) GetVoteTime() *prototype.TimePointSec {
	if m != nil {
		return m.VoteTime
	}
	return nil
}

type SoListBlockProducerVoteByVoterId struct {
	VoterId              *prototype.BpVoterId `protobuf:"bytes,1,opt,name=voter_id,json=voterId,proto3" json:"voter_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *SoListBlockProducerVoteByVoterId) Reset()         { *m = SoListBlockProducerVoteByVoterId{} }
func (m *SoListBlockProducerVoteByVoterId) String() string { return proto.CompactTextString(m) }
func (*SoListBlockProducerVoteByVoterId) ProtoMessage()    {}
func (*SoListBlockProducerVoteByVoterId) Descriptor() ([]byte, []int) {
	return fileDescriptor_5d279fa3103478f1, []int{4}
}

func (m *SoListBlockProducerVoteByVoterId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoListBlockProducerVoteByVoterId.Unmarshal(m, b)
}
func (m *SoListBlockProducerVoteByVoterId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoListBlockProducerVoteByVoterId.Marshal(b, m, deterministic)
}
func (m *SoListBlockProducerVoteByVoterId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoListBlockProducerVoteByVoterId.Merge(m, src)
}
func (m *SoListBlockProducerVoteByVoterId) XXX_Size() int {
	return xxx_messageInfo_SoListBlockProducerVoteByVoterId.Size(m)
}
func (m *SoListBlockProducerVoteByVoterId) XXX_DiscardUnknown() {
	xxx_messageInfo_SoListBlockProducerVoteByVoterId.DiscardUnknown(m)
}

var xxx_messageInfo_SoListBlockProducerVoteByVoterId proto.InternalMessageInfo

func (m *SoListBlockProducerVoteByVoterId) GetVoterId() *prototype.BpVoterId {
	if m != nil {
		return m.VoterId
	}
	return nil
}

type SoUniqueBlockProducerVoteByVoterId struct {
	VoterId              *prototype.BpVoterId `protobuf:"bytes,1,opt,name=voter_id,json=voterId,proto3" json:"voter_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *SoUniqueBlockProducerVoteByVoterId) Reset()         { *m = SoUniqueBlockProducerVoteByVoterId{} }
func (m *SoUniqueBlockProducerVoteByVoterId) String() string { return proto.CompactTextString(m) }
func (*SoUniqueBlockProducerVoteByVoterId) ProtoMessage()    {}
func (*SoUniqueBlockProducerVoteByVoterId) Descriptor() ([]byte, []int) {
	return fileDescriptor_5d279fa3103478f1, []int{5}
}

func (m *SoUniqueBlockProducerVoteByVoterId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoUniqueBlockProducerVoteByVoterId.Unmarshal(m, b)
}
func (m *SoUniqueBlockProducerVoteByVoterId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoUniqueBlockProducerVoteByVoterId.Marshal(b, m, deterministic)
}
func (m *SoUniqueBlockProducerVoteByVoterId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoUniqueBlockProducerVoteByVoterId.Merge(m, src)
}
func (m *SoUniqueBlockProducerVoteByVoterId) XXX_Size() int {
	return xxx_messageInfo_SoUniqueBlockProducerVoteByVoterId.Size(m)
}
func (m *SoUniqueBlockProducerVoteByVoterId) XXX_DiscardUnknown() {
	xxx_messageInfo_SoUniqueBlockProducerVoteByVoterId.DiscardUnknown(m)
}

var xxx_messageInfo_SoUniqueBlockProducerVoteByVoterId proto.InternalMessageInfo

func (m *SoUniqueBlockProducerVoteByVoterId) GetVoterId() *prototype.BpVoterId {
	if m != nil {
		return m.VoterId
	}
	return nil
}

func init() {
	proto.RegisterType((*SoBlockProducerVote)(nil), "table.so_blockProducerVote")
	proto.RegisterType((*SoMemBlockProducerVoteByVoterId)(nil), "table.so_mem_blockProducerVote_by_voter_id")
	proto.RegisterType((*SoMemBlockProducerVoteByBlockProducerId)(nil), "table.so_mem_blockProducerVote_by_block_producer_id")
	proto.RegisterType((*SoMemBlockProducerVoteByVoteTime)(nil), "table.so_mem_blockProducerVote_by_vote_time")
	proto.RegisterType((*SoListBlockProducerVoteByVoterId)(nil), "table.so_list_blockProducerVote_by_voter_id")
	proto.RegisterType((*SoUniqueBlockProducerVoteByVoterId)(nil), "table.so_unique_blockProducerVote_by_voter_id")
}

func init() {
	proto.RegisterFile("app/table/so_blockProducerVote.proto", fileDescriptor_5d279fa3103478f1)
}

var fileDescriptor_5d279fa3103478f1 = []byte{
	// 311 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x93, 0x5d, 0x4b, 0xfb, 0x30,
	0x14, 0xc6, 0xe9, 0xff, 0x8f, 0x6f, 0xf1, 0x42, 0x2c, 0x43, 0xe6, 0x6e, 0x94, 0x32, 0x51, 0xc4,
	0x35, 0xa8, 0xe0, 0x07, 0xf0, 0x6e, 0x78, 0x23, 0x43, 0x04, 0x87, 0x70, 0x68, 0xd2, 0xb0, 0x05,
	0x9b, 0x9e, 0xd8, 0x9c, 0x0a, 0x03, 0x3f, 0xa3, 0x9f, 0x49, 0x9a, 0x69, 0xb5, 0x74, 0xa8, 0x8c,
	0xdd, 0x94, 0xe6, 0xbc, 0xfc, 0xf2, 0x3c, 0x0f, 0x84, 0xf5, 0x13, 0x6b, 0x39, 0x25, 0x22, 0x53,
	0xdc, 0x21, 0x88, 0x0c, 0xe5, 0xd3, 0x6d, 0x81, 0x69, 0x29, 0x55, 0x71, 0x8f, 0xa4, 0x62, 0x5b,
	0x20, 0x61, 0xb8, 0xe6, 0x27, 0x7a, 0x5d, 0x7f, 0xa2, 0x99, 0x55, 0xdc, 0x94, 0x19, 0x69, 0xd0,
	0xe9, 0x7c, 0xa0, 0xd7, 0xf9, 0xea, 0x54, 0x9f, 0x79, 0x35, 0x7a, 0x0b, 0x58, 0x67, 0x11, 0x35,
	0x3c, 0x67, 0x9b, 0x2f, 0x48, 0xaa, 0x00, 0x9d, 0x76, 0x83, 0xc3, 0xe0, 0x64, 0xfb, 0x62, 0x2f,
	0xae, 0x09, 0xb1, 0xb0, 0xf0, 0xd9, 0x1d, 0x6d, 0xf8, 0xbf, 0x61, 0x1a, 0xde, 0xb0, 0x5d, 0xcf,
	0x01, 0xfb, 0x01, 0xaa, 0x76, 0xff, 0xf9, 0xdd, 0x83, 0xe6, 0x6e, 0x6b, 0x6c, 0xb4, 0xd3, 0x50,
	0x30, 0x4c, 0xc3, 0x2b, 0xb6, 0x55, 0x71, 0x81, 0xb4, 0x51, 0xdd, 0xff, 0x1e, 0xb2, 0xff, 0x0d,
	0x52, 0x95, 0xc1, 0xa2, 0xce, 0x09, 0x9c, 0x92, 0x23, 0xaf, 0xf5, 0x4e, 0x1b, 0x15, 0x3d, 0xb0,
	0xbe, 0x43, 0x30, 0xca, 0xb4, 0x3d, 0x81, 0x98, 0xd5, 0xaa, 0x97, 0xf0, 0x17, 0xbd, 0xb2, 0xc1,
	0x4f, 0xe8, 0x96, 0xa9, 0xc5, 0x81, 0x04, 0xcb, 0x05, 0x12, 0x01, 0x3b, 0xfa, 0xcd, 0x98, 0x0f,
	0xab, 0x99, 0x5c, 0xf0, 0xf7, 0xe4, 0xc6, 0xfe, 0x82, 0x4c, 0x3b, 0x5a, 0x7d, 0x74, 0x8f, 0xec,
	0xd8, 0x21, 0x94, 0xb9, 0x7e, 0x2e, 0xd5, 0xca, 0xe9, 0xd7, 0x67, 0xe3, 0xd3, 0x89, 0xa6, 0x69,
	0x29, 0x62, 0x89, 0x86, 0x4b, 0x74, 0x72, 0x9a, 0xe8, 0x9c, 0x4b, 0xcc, 0x49, 0xe5, 0x84, 0x6e,
	0x30, 0x41, 0x5e, 0x3f, 0x22, 0xb1, 0xee, 0x69, 0x97, 0xef, 0x01, 0x00, 0x00, 0xff, 0xff, 0x56,
	0x9b, 0x7c, 0x09, 0x58, 0x03, 0x00, 0x00,
}
