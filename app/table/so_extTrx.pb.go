// Code generated by protoc-gen-go. DO NOT EDIT.
// source: app/table/so_extTrx.proto

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

type SoExtTrx struct {
	TrxId                *prototype.Sha256             `protobuf:"bytes,1,opt,name=trx_id,json=trxId,proto3" json:"trx_id,omitempty"`
	BlockHeight          uint64                        `protobuf:"varint,2,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	TrxWrap              *prototype.TransactionWrapper `protobuf:"bytes,3,opt,name=trx_wrap,json=trxWrap,proto3" json:"trx_wrap,omitempty"`
	BlockTime            *prototype.TimePointSec       `protobuf:"bytes,4,opt,name=block_time,json=blockTime,proto3" json:"block_time,omitempty"`
	TrxCreateOrder       *prototype.UserTrxCreateOrder `protobuf:"bytes,5,opt,name=trx_create_order,json=trxCreateOrder,proto3" json:"trx_create_order,omitempty"`
	BlockId              *prototype.Sha256             `protobuf:"bytes,6,opt,name=block_id,json=blockId,proto3" json:"block_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *SoExtTrx) Reset()         { *m = SoExtTrx{} }
func (m *SoExtTrx) String() string { return proto.CompactTextString(m) }
func (*SoExtTrx) ProtoMessage()    {}
func (*SoExtTrx) Descriptor() ([]byte, []int) {
	return fileDescriptor_76957eaae1f8a1bc, []int{0}
}

func (m *SoExtTrx) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoExtTrx.Unmarshal(m, b)
}
func (m *SoExtTrx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoExtTrx.Marshal(b, m, deterministic)
}
func (m *SoExtTrx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoExtTrx.Merge(m, src)
}
func (m *SoExtTrx) XXX_Size() int {
	return xxx_messageInfo_SoExtTrx.Size(m)
}
func (m *SoExtTrx) XXX_DiscardUnknown() {
	xxx_messageInfo_SoExtTrx.DiscardUnknown(m)
}

var xxx_messageInfo_SoExtTrx proto.InternalMessageInfo

func (m *SoExtTrx) GetTrxId() *prototype.Sha256 {
	if m != nil {
		return m.TrxId
	}
	return nil
}

func (m *SoExtTrx) GetBlockHeight() uint64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

func (m *SoExtTrx) GetTrxWrap() *prototype.TransactionWrapper {
	if m != nil {
		return m.TrxWrap
	}
	return nil
}

func (m *SoExtTrx) GetBlockTime() *prototype.TimePointSec {
	if m != nil {
		return m.BlockTime
	}
	return nil
}

func (m *SoExtTrx) GetTrxCreateOrder() *prototype.UserTrxCreateOrder {
	if m != nil {
		return m.TrxCreateOrder
	}
	return nil
}

func (m *SoExtTrx) GetBlockId() *prototype.Sha256 {
	if m != nil {
		return m.BlockId
	}
	return nil
}

type SoMemExtTrxByTrxId struct {
	TrxId                *prototype.Sha256 `protobuf:"bytes,1,opt,name=trx_id,json=trxId,proto3" json:"trx_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SoMemExtTrxByTrxId) Reset()         { *m = SoMemExtTrxByTrxId{} }
func (m *SoMemExtTrxByTrxId) String() string { return proto.CompactTextString(m) }
func (*SoMemExtTrxByTrxId) ProtoMessage()    {}
func (*SoMemExtTrxByTrxId) Descriptor() ([]byte, []int) {
	return fileDescriptor_76957eaae1f8a1bc, []int{1}
}

func (m *SoMemExtTrxByTrxId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoMemExtTrxByTrxId.Unmarshal(m, b)
}
func (m *SoMemExtTrxByTrxId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoMemExtTrxByTrxId.Marshal(b, m, deterministic)
}
func (m *SoMemExtTrxByTrxId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoMemExtTrxByTrxId.Merge(m, src)
}
func (m *SoMemExtTrxByTrxId) XXX_Size() int {
	return xxx_messageInfo_SoMemExtTrxByTrxId.Size(m)
}
func (m *SoMemExtTrxByTrxId) XXX_DiscardUnknown() {
	xxx_messageInfo_SoMemExtTrxByTrxId.DiscardUnknown(m)
}

var xxx_messageInfo_SoMemExtTrxByTrxId proto.InternalMessageInfo

func (m *SoMemExtTrxByTrxId) GetTrxId() *prototype.Sha256 {
	if m != nil {
		return m.TrxId
	}
	return nil
}

type SoMemExtTrxByBlockHeight struct {
	BlockHeight          uint64   `protobuf:"varint,1,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SoMemExtTrxByBlockHeight) Reset()         { *m = SoMemExtTrxByBlockHeight{} }
func (m *SoMemExtTrxByBlockHeight) String() string { return proto.CompactTextString(m) }
func (*SoMemExtTrxByBlockHeight) ProtoMessage()    {}
func (*SoMemExtTrxByBlockHeight) Descriptor() ([]byte, []int) {
	return fileDescriptor_76957eaae1f8a1bc, []int{2}
}

func (m *SoMemExtTrxByBlockHeight) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoMemExtTrxByBlockHeight.Unmarshal(m, b)
}
func (m *SoMemExtTrxByBlockHeight) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoMemExtTrxByBlockHeight.Marshal(b, m, deterministic)
}
func (m *SoMemExtTrxByBlockHeight) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoMemExtTrxByBlockHeight.Merge(m, src)
}
func (m *SoMemExtTrxByBlockHeight) XXX_Size() int {
	return xxx_messageInfo_SoMemExtTrxByBlockHeight.Size(m)
}
func (m *SoMemExtTrxByBlockHeight) XXX_DiscardUnknown() {
	xxx_messageInfo_SoMemExtTrxByBlockHeight.DiscardUnknown(m)
}

var xxx_messageInfo_SoMemExtTrxByBlockHeight proto.InternalMessageInfo

func (m *SoMemExtTrxByBlockHeight) GetBlockHeight() uint64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

type SoMemExtTrxByTrxWrap struct {
	TrxWrap              *prototype.TransactionWrapper `protobuf:"bytes,1,opt,name=trx_wrap,json=trxWrap,proto3" json:"trx_wrap,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *SoMemExtTrxByTrxWrap) Reset()         { *m = SoMemExtTrxByTrxWrap{} }
func (m *SoMemExtTrxByTrxWrap) String() string { return proto.CompactTextString(m) }
func (*SoMemExtTrxByTrxWrap) ProtoMessage()    {}
func (*SoMemExtTrxByTrxWrap) Descriptor() ([]byte, []int) {
	return fileDescriptor_76957eaae1f8a1bc, []int{3}
}

func (m *SoMemExtTrxByTrxWrap) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoMemExtTrxByTrxWrap.Unmarshal(m, b)
}
func (m *SoMemExtTrxByTrxWrap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoMemExtTrxByTrxWrap.Marshal(b, m, deterministic)
}
func (m *SoMemExtTrxByTrxWrap) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoMemExtTrxByTrxWrap.Merge(m, src)
}
func (m *SoMemExtTrxByTrxWrap) XXX_Size() int {
	return xxx_messageInfo_SoMemExtTrxByTrxWrap.Size(m)
}
func (m *SoMemExtTrxByTrxWrap) XXX_DiscardUnknown() {
	xxx_messageInfo_SoMemExtTrxByTrxWrap.DiscardUnknown(m)
}

var xxx_messageInfo_SoMemExtTrxByTrxWrap proto.InternalMessageInfo

func (m *SoMemExtTrxByTrxWrap) GetTrxWrap() *prototype.TransactionWrapper {
	if m != nil {
		return m.TrxWrap
	}
	return nil
}

type SoMemExtTrxByBlockTime struct {
	BlockTime            *prototype.TimePointSec `protobuf:"bytes,1,opt,name=block_time,json=blockTime,proto3" json:"block_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *SoMemExtTrxByBlockTime) Reset()         { *m = SoMemExtTrxByBlockTime{} }
func (m *SoMemExtTrxByBlockTime) String() string { return proto.CompactTextString(m) }
func (*SoMemExtTrxByBlockTime) ProtoMessage()    {}
func (*SoMemExtTrxByBlockTime) Descriptor() ([]byte, []int) {
	return fileDescriptor_76957eaae1f8a1bc, []int{4}
}

func (m *SoMemExtTrxByBlockTime) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoMemExtTrxByBlockTime.Unmarshal(m, b)
}
func (m *SoMemExtTrxByBlockTime) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoMemExtTrxByBlockTime.Marshal(b, m, deterministic)
}
func (m *SoMemExtTrxByBlockTime) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoMemExtTrxByBlockTime.Merge(m, src)
}
func (m *SoMemExtTrxByBlockTime) XXX_Size() int {
	return xxx_messageInfo_SoMemExtTrxByBlockTime.Size(m)
}
func (m *SoMemExtTrxByBlockTime) XXX_DiscardUnknown() {
	xxx_messageInfo_SoMemExtTrxByBlockTime.DiscardUnknown(m)
}

var xxx_messageInfo_SoMemExtTrxByBlockTime proto.InternalMessageInfo

func (m *SoMemExtTrxByBlockTime) GetBlockTime() *prototype.TimePointSec {
	if m != nil {
		return m.BlockTime
	}
	return nil
}

type SoMemExtTrxByTrxCreateOrder struct {
	TrxCreateOrder       *prototype.UserTrxCreateOrder `protobuf:"bytes,1,opt,name=trx_create_order,json=trxCreateOrder,proto3" json:"trx_create_order,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *SoMemExtTrxByTrxCreateOrder) Reset()         { *m = SoMemExtTrxByTrxCreateOrder{} }
func (m *SoMemExtTrxByTrxCreateOrder) String() string { return proto.CompactTextString(m) }
func (*SoMemExtTrxByTrxCreateOrder) ProtoMessage()    {}
func (*SoMemExtTrxByTrxCreateOrder) Descriptor() ([]byte, []int) {
	return fileDescriptor_76957eaae1f8a1bc, []int{5}
}

func (m *SoMemExtTrxByTrxCreateOrder) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoMemExtTrxByTrxCreateOrder.Unmarshal(m, b)
}
func (m *SoMemExtTrxByTrxCreateOrder) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoMemExtTrxByTrxCreateOrder.Marshal(b, m, deterministic)
}
func (m *SoMemExtTrxByTrxCreateOrder) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoMemExtTrxByTrxCreateOrder.Merge(m, src)
}
func (m *SoMemExtTrxByTrxCreateOrder) XXX_Size() int {
	return xxx_messageInfo_SoMemExtTrxByTrxCreateOrder.Size(m)
}
func (m *SoMemExtTrxByTrxCreateOrder) XXX_DiscardUnknown() {
	xxx_messageInfo_SoMemExtTrxByTrxCreateOrder.DiscardUnknown(m)
}

var xxx_messageInfo_SoMemExtTrxByTrxCreateOrder proto.InternalMessageInfo

func (m *SoMemExtTrxByTrxCreateOrder) GetTrxCreateOrder() *prototype.UserTrxCreateOrder {
	if m != nil {
		return m.TrxCreateOrder
	}
	return nil
}

type SoMemExtTrxByBlockId struct {
	BlockId              *prototype.Sha256 `protobuf:"bytes,1,opt,name=block_id,json=blockId,proto3" json:"block_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SoMemExtTrxByBlockId) Reset()         { *m = SoMemExtTrxByBlockId{} }
func (m *SoMemExtTrxByBlockId) String() string { return proto.CompactTextString(m) }
func (*SoMemExtTrxByBlockId) ProtoMessage()    {}
func (*SoMemExtTrxByBlockId) Descriptor() ([]byte, []int) {
	return fileDescriptor_76957eaae1f8a1bc, []int{6}
}

func (m *SoMemExtTrxByBlockId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoMemExtTrxByBlockId.Unmarshal(m, b)
}
func (m *SoMemExtTrxByBlockId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoMemExtTrxByBlockId.Marshal(b, m, deterministic)
}
func (m *SoMemExtTrxByBlockId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoMemExtTrxByBlockId.Merge(m, src)
}
func (m *SoMemExtTrxByBlockId) XXX_Size() int {
	return xxx_messageInfo_SoMemExtTrxByBlockId.Size(m)
}
func (m *SoMemExtTrxByBlockId) XXX_DiscardUnknown() {
	xxx_messageInfo_SoMemExtTrxByBlockId.DiscardUnknown(m)
}

var xxx_messageInfo_SoMemExtTrxByBlockId proto.InternalMessageInfo

func (m *SoMemExtTrxByBlockId) GetBlockId() *prototype.Sha256 {
	if m != nil {
		return m.BlockId
	}
	return nil
}

type SoListExtTrxByTrxId struct {
	TrxId                *prototype.Sha256 `protobuf:"bytes,1,opt,name=trx_id,json=trxId,proto3" json:"trx_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SoListExtTrxByTrxId) Reset()         { *m = SoListExtTrxByTrxId{} }
func (m *SoListExtTrxByTrxId) String() string { return proto.CompactTextString(m) }
func (*SoListExtTrxByTrxId) ProtoMessage()    {}
func (*SoListExtTrxByTrxId) Descriptor() ([]byte, []int) {
	return fileDescriptor_76957eaae1f8a1bc, []int{7}
}

func (m *SoListExtTrxByTrxId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoListExtTrxByTrxId.Unmarshal(m, b)
}
func (m *SoListExtTrxByTrxId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoListExtTrxByTrxId.Marshal(b, m, deterministic)
}
func (m *SoListExtTrxByTrxId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoListExtTrxByTrxId.Merge(m, src)
}
func (m *SoListExtTrxByTrxId) XXX_Size() int {
	return xxx_messageInfo_SoListExtTrxByTrxId.Size(m)
}
func (m *SoListExtTrxByTrxId) XXX_DiscardUnknown() {
	xxx_messageInfo_SoListExtTrxByTrxId.DiscardUnknown(m)
}

var xxx_messageInfo_SoListExtTrxByTrxId proto.InternalMessageInfo

func (m *SoListExtTrxByTrxId) GetTrxId() *prototype.Sha256 {
	if m != nil {
		return m.TrxId
	}
	return nil
}

type SoListExtTrxByBlockHeight struct {
	BlockHeight          uint64            `protobuf:"varint,1,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	TrxId                *prototype.Sha256 `protobuf:"bytes,2,opt,name=trx_id,json=trxId,proto3" json:"trx_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SoListExtTrxByBlockHeight) Reset()         { *m = SoListExtTrxByBlockHeight{} }
func (m *SoListExtTrxByBlockHeight) String() string { return proto.CompactTextString(m) }
func (*SoListExtTrxByBlockHeight) ProtoMessage()    {}
func (*SoListExtTrxByBlockHeight) Descriptor() ([]byte, []int) {
	return fileDescriptor_76957eaae1f8a1bc, []int{8}
}

func (m *SoListExtTrxByBlockHeight) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoListExtTrxByBlockHeight.Unmarshal(m, b)
}
func (m *SoListExtTrxByBlockHeight) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoListExtTrxByBlockHeight.Marshal(b, m, deterministic)
}
func (m *SoListExtTrxByBlockHeight) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoListExtTrxByBlockHeight.Merge(m, src)
}
func (m *SoListExtTrxByBlockHeight) XXX_Size() int {
	return xxx_messageInfo_SoListExtTrxByBlockHeight.Size(m)
}
func (m *SoListExtTrxByBlockHeight) XXX_DiscardUnknown() {
	xxx_messageInfo_SoListExtTrxByBlockHeight.DiscardUnknown(m)
}

var xxx_messageInfo_SoListExtTrxByBlockHeight proto.InternalMessageInfo

func (m *SoListExtTrxByBlockHeight) GetBlockHeight() uint64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

func (m *SoListExtTrxByBlockHeight) GetTrxId() *prototype.Sha256 {
	if m != nil {
		return m.TrxId
	}
	return nil
}

type SoListExtTrxByBlockTime struct {
	BlockTime            *prototype.TimePointSec `protobuf:"bytes,1,opt,name=block_time,json=blockTime,proto3" json:"block_time,omitempty"`
	TrxId                *prototype.Sha256       `protobuf:"bytes,2,opt,name=trx_id,json=trxId,proto3" json:"trx_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *SoListExtTrxByBlockTime) Reset()         { *m = SoListExtTrxByBlockTime{} }
func (m *SoListExtTrxByBlockTime) String() string { return proto.CompactTextString(m) }
func (*SoListExtTrxByBlockTime) ProtoMessage()    {}
func (*SoListExtTrxByBlockTime) Descriptor() ([]byte, []int) {
	return fileDescriptor_76957eaae1f8a1bc, []int{9}
}

func (m *SoListExtTrxByBlockTime) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoListExtTrxByBlockTime.Unmarshal(m, b)
}
func (m *SoListExtTrxByBlockTime) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoListExtTrxByBlockTime.Marshal(b, m, deterministic)
}
func (m *SoListExtTrxByBlockTime) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoListExtTrxByBlockTime.Merge(m, src)
}
func (m *SoListExtTrxByBlockTime) XXX_Size() int {
	return xxx_messageInfo_SoListExtTrxByBlockTime.Size(m)
}
func (m *SoListExtTrxByBlockTime) XXX_DiscardUnknown() {
	xxx_messageInfo_SoListExtTrxByBlockTime.DiscardUnknown(m)
}

var xxx_messageInfo_SoListExtTrxByBlockTime proto.InternalMessageInfo

func (m *SoListExtTrxByBlockTime) GetBlockTime() *prototype.TimePointSec {
	if m != nil {
		return m.BlockTime
	}
	return nil
}

func (m *SoListExtTrxByBlockTime) GetTrxId() *prototype.Sha256 {
	if m != nil {
		return m.TrxId
	}
	return nil
}

type SoListExtTrxByTrxCreateOrder struct {
	TrxCreateOrder       *prototype.UserTrxCreateOrder `protobuf:"bytes,1,opt,name=trx_create_order,json=trxCreateOrder,proto3" json:"trx_create_order,omitempty"`
	TrxId                *prototype.Sha256             `protobuf:"bytes,2,opt,name=trx_id,json=trxId,proto3" json:"trx_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *SoListExtTrxByTrxCreateOrder) Reset()         { *m = SoListExtTrxByTrxCreateOrder{} }
func (m *SoListExtTrxByTrxCreateOrder) String() string { return proto.CompactTextString(m) }
func (*SoListExtTrxByTrxCreateOrder) ProtoMessage()    {}
func (*SoListExtTrxByTrxCreateOrder) Descriptor() ([]byte, []int) {
	return fileDescriptor_76957eaae1f8a1bc, []int{10}
}

func (m *SoListExtTrxByTrxCreateOrder) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoListExtTrxByTrxCreateOrder.Unmarshal(m, b)
}
func (m *SoListExtTrxByTrxCreateOrder) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoListExtTrxByTrxCreateOrder.Marshal(b, m, deterministic)
}
func (m *SoListExtTrxByTrxCreateOrder) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoListExtTrxByTrxCreateOrder.Merge(m, src)
}
func (m *SoListExtTrxByTrxCreateOrder) XXX_Size() int {
	return xxx_messageInfo_SoListExtTrxByTrxCreateOrder.Size(m)
}
func (m *SoListExtTrxByTrxCreateOrder) XXX_DiscardUnknown() {
	xxx_messageInfo_SoListExtTrxByTrxCreateOrder.DiscardUnknown(m)
}

var xxx_messageInfo_SoListExtTrxByTrxCreateOrder proto.InternalMessageInfo

func (m *SoListExtTrxByTrxCreateOrder) GetTrxCreateOrder() *prototype.UserTrxCreateOrder {
	if m != nil {
		return m.TrxCreateOrder
	}
	return nil
}

func (m *SoListExtTrxByTrxCreateOrder) GetTrxId() *prototype.Sha256 {
	if m != nil {
		return m.TrxId
	}
	return nil
}

type SoUniqueExtTrxByTrxId struct {
	TrxId                *prototype.Sha256 `protobuf:"bytes,1,opt,name=trx_id,json=trxId,proto3" json:"trx_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SoUniqueExtTrxByTrxId) Reset()         { *m = SoUniqueExtTrxByTrxId{} }
func (m *SoUniqueExtTrxByTrxId) String() string { return proto.CompactTextString(m) }
func (*SoUniqueExtTrxByTrxId) ProtoMessage()    {}
func (*SoUniqueExtTrxByTrxId) Descriptor() ([]byte, []int) {
	return fileDescriptor_76957eaae1f8a1bc, []int{11}
}

func (m *SoUniqueExtTrxByTrxId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SoUniqueExtTrxByTrxId.Unmarshal(m, b)
}
func (m *SoUniqueExtTrxByTrxId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SoUniqueExtTrxByTrxId.Marshal(b, m, deterministic)
}
func (m *SoUniqueExtTrxByTrxId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SoUniqueExtTrxByTrxId.Merge(m, src)
}
func (m *SoUniqueExtTrxByTrxId) XXX_Size() int {
	return xxx_messageInfo_SoUniqueExtTrxByTrxId.Size(m)
}
func (m *SoUniqueExtTrxByTrxId) XXX_DiscardUnknown() {
	xxx_messageInfo_SoUniqueExtTrxByTrxId.DiscardUnknown(m)
}

var xxx_messageInfo_SoUniqueExtTrxByTrxId proto.InternalMessageInfo

func (m *SoUniqueExtTrxByTrxId) GetTrxId() *prototype.Sha256 {
	if m != nil {
		return m.TrxId
	}
	return nil
}

func init() {
	proto.RegisterType((*SoExtTrx)(nil), "table.so_extTrx")
	proto.RegisterType((*SoMemExtTrxByTrxId)(nil), "table.so_mem_extTrx_by_trx_id")
	proto.RegisterType((*SoMemExtTrxByBlockHeight)(nil), "table.so_mem_extTrx_by_block_height")
	proto.RegisterType((*SoMemExtTrxByTrxWrap)(nil), "table.so_mem_extTrx_by_trx_wrap")
	proto.RegisterType((*SoMemExtTrxByBlockTime)(nil), "table.so_mem_extTrx_by_block_time")
	proto.RegisterType((*SoMemExtTrxByTrxCreateOrder)(nil), "table.so_mem_extTrx_by_trx_create_order")
	proto.RegisterType((*SoMemExtTrxByBlockId)(nil), "table.so_mem_extTrx_by_block_id")
	proto.RegisterType((*SoListExtTrxByTrxId)(nil), "table.so_list_extTrx_by_trx_id")
	proto.RegisterType((*SoListExtTrxByBlockHeight)(nil), "table.so_list_extTrx_by_block_height")
	proto.RegisterType((*SoListExtTrxByBlockTime)(nil), "table.so_list_extTrx_by_block_time")
	proto.RegisterType((*SoListExtTrxByTrxCreateOrder)(nil), "table.so_list_extTrx_by_trx_create_order")
	proto.RegisterType((*SoUniqueExtTrxByTrxId)(nil), "table.so_unique_extTrx_by_trx_id")
}

func init() { proto.RegisterFile("app/table/so_extTrx.proto", fileDescriptor_76957eaae1f8a1bc) }

var fileDescriptor_76957eaae1f8a1bc = []byte{
	// 473 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x94, 0xd1, 0x6b, 0xd3, 0x50,
	0x14, 0xc6, 0xb9, 0x75, 0xed, 0xb6, 0x33, 0x11, 0x0d, 0x82, 0x69, 0xa7, 0xa3, 0xcb, 0x53, 0x91,
	0xd9, 0xc0, 0x44, 0xd1, 0xd7, 0x4d, 0xc4, 0xfa, 0x22, 0x84, 0xe1, 0xc0, 0x97, 0xcb, 0x4d, 0x72,
	0x69, 0x2e, 0x36, 0x39, 0xf1, 0xde, 0x13, 0xcc, 0x5e, 0xfd, 0x33, 0xfc, 0x87, 0xfc, 0xb7, 0x24,
	0x37, 0xb3, 0x66, 0x6b, 0xca, 0x1a, 0x8b, 0x2f, 0x85, 0x7c, 0xf7, 0x9c, 0xef, 0x7c, 0x7c, 0x3f,
	0x28, 0x0c, 0x45, 0x9e, 0xfb, 0x24, 0xc2, 0x85, 0xf4, 0x0d, 0x72, 0x59, 0xd2, 0x85, 0x2e, 0xa7,
	0xb9, 0x46, 0x42, 0xa7, 0x6f, 0xe5, 0x91, 0x6b, 0xbf, 0xe8, 0x2a, 0x97, 0x7e, 0x5a, 0x2c, 0x48,
	0x71, 0x15, 0xd7, 0x03, 0xa3, 0xc7, 0x7f, 0x5f, 0xaa, 0x9f, 0x6b, 0xf5, 0xb0, 0xa1, 0x6a, 0x91,
	0x19, 0x11, 0x91, 0xc2, 0xac, 0x7e, 0xf4, 0x7e, 0xf5, 0x60, 0x7f, 0x79, 0xc7, 0x99, 0xc0, 0x80,
	0x74, 0xc9, 0x55, 0xec, 0xb2, 0x31, 0x9b, 0x1c, 0x9c, 0x3e, 0x9a, 0x2e, 0x77, 0xa7, 0x26, 0x11,
	0xa7, 0xaf, 0x5e, 0x07, 0x7d, 0xd2, 0xe5, 0x2c, 0x76, 0x8e, 0xe1, 0x7e, 0xb8, 0xc0, 0xe8, 0x2b,
	0x4f, 0xa4, 0x9a, 0x27, 0xe4, 0xf6, 0xc6, 0x6c, 0xb2, 0x13, 0x1c, 0x58, 0xed, 0x83, 0x95, 0x9c,
	0xb7, 0xb0, 0x57, 0x99, 0x7d, 0xd7, 0x22, 0x77, 0xef, 0x59, 0xbb, 0xa3, 0x86, 0x5d, 0x23, 0x8a,
	0x1d, 0xc9, 0xa5, 0x0e, 0x76, 0x49, 0x97, 0x97, 0x5a, 0xe4, 0xce, 0x1b, 0x80, 0xda, 0x9d, 0x54,
	0x2a, 0xdd, 0x1d, 0xbb, 0x3c, 0x6c, 0x2e, 0xab, 0x54, 0xf2, 0x1c, 0x55, 0x46, 0xdc, 0xc8, 0x28,
	0xd8, 0xb7, 0xc3, 0x17, 0x2a, 0x95, 0xce, 0x47, 0x78, 0x58, 0x1d, 0x8d, 0xb4, 0x14, 0x24, 0x39,
	0xea, 0x58, 0x6a, 0xb7, 0x6f, 0xf7, 0xc7, 0x8d, 0xfd, 0xc2, 0x48, 0xcd, 0x6f, 0xcf, 0x05, 0x0f,
	0x48, 0x97, 0xe7, 0x56, 0xf8, 0x54, 0x7d, 0x3b, 0x27, 0xb0, 0x57, 0xa7, 0x50, 0xb1, 0x3b, 0x58,
	0xd7, 0xc7, 0xae, 0x1d, 0x99, 0xc5, 0xde, 0x39, 0x3c, 0x31, 0xc8, 0x53, 0x99, 0x5e, 0x97, 0xc9,
	0xc3, 0x2b, 0x5e, 0x97, 0xb9, 0x79, 0xad, 0xde, 0x19, 0x3c, 0x5b, 0x31, 0x69, 0xf6, 0xbc, 0xd2,
	0x3b, 0x5b, 0xe9, 0xdd, 0xfb, 0x0c, 0xc3, 0xd6, 0x20, 0x55, 0xcb, 0x37, 0xa0, 0xb0, 0x4e, 0x50,
	0xbc, 0x4b, 0x38, 0x5c, 0x93, 0xad, 0xc2, 0x71, 0x8b, 0x19, 0xdb, 0x9c, 0x99, 0x87, 0x70, 0xdc,
	0x1a, 0xb8, 0x09, 0xa7, 0x15, 0x2c, 0xfb, 0x37, 0xb0, 0xde, 0xac, 0xa5, 0xa1, 0x3f, 0xa4, 0x6f,
	0x50, 0x67, 0x77, 0x52, 0x7f, 0x07, 0xae, 0x41, 0xbe, 0x50, 0x86, 0xb6, 0xc1, 0x9e, 0xc2, 0xd1,
	0xaa, 0x4b, 0x47, 0xee, 0x8d, 0x73, 0xbd, 0x3b, 0xce, 0xfd, 0x60, 0xf0, 0x74, 0xdd, 0xbd, 0xed,
	0x58, 0x76, 0x08, 0xf1, 0x93, 0x81, 0xd7, 0x5e, 0xdd, 0xff, 0xe2, 0xde, 0x21, 0xdc, 0x7b, 0x18,
	0x19, 0xe4, 0x45, 0xa6, 0xbe, 0x15, 0x72, 0x0b, 0xb0, 0x67, 0x27, 0x5f, 0x9e, 0xcf, 0x15, 0x25,
	0x45, 0x38, 0x8d, 0x30, 0xf5, 0x23, 0x34, 0x51, 0x22, 0x54, 0xe6, 0x47, 0x98, 0x91, 0xcc, 0x08,
	0xcd, 0x8b, 0x39, 0xfa, 0xcb, 0x3f, 0xfc, 0x70, 0x60, 0x6d, 0x5e, 0xfe, 0x0e, 0x00, 0x00, 0xff,
	0xff, 0x87, 0x75, 0xcf, 0x0b, 0x04, 0x06, 0x00, 0x00,
}
