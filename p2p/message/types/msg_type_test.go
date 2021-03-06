package types

import (
	"fmt"
	"testing"

	"github.com/coschain/contentos-go/prototype"
	"github.com/gogo/protobuf/proto"
)

func Test_Serialize(t *testing.T) {
	obj1 := new(TransferMsg)
	obj2 := new(TransferMsg)

	trx := &prototype.Transaction{
		RefBlockNum:    1,
		RefBlockPrefix: 2,
	}

	sigtrx := new(prototype.SignedTransaction)
	sigtrx.Trx = trx
	msg := new(BroadcastSigTrx)
	msg.SigTrx = sigtrx

	obj1.MakeMsg(msg)
	obj1.cmdType()
	fmt.Printf("before marshal BroadcastSigTrx message, message data: +%v\n",obj1)
	data, err := proto.Marshal(obj1)
	if err != nil {
		t.Error("BroadcastSigTrx Marshal failed")
	}
	err = proto.Unmarshal(data, obj2)
	if err != nil {
		t.Error("BroadcastSigTrx Unmarshal failed")
	}
	fmt.Printf("after unmarshal BroadcastSigTrx message, message data: +%v\n\n",obj2)

	sigBlk := new(prototype.SignedBlock)
	sigBlkHdr := new(prototype.SignedBlockHeader)
	sigBlkHdr.Header = new(prototype.BlockHeader)
	sigBlkHdr.Header.BlockProducer = new(prototype.AccountName)
	sigBlkHdr.Header.BlockProducer.Value = "alice"
	sigBlk.SignedHeader = sigBlkHdr

	msg2 := new(SigBlkMsg)
	msg2.SigBlk = sigBlk

	obj1.MakeMsg(msg2)
	obj1.cmdType()
	fmt.Printf("before marshal SigBlkMsg message, message data: +%v\n",obj1)
	data, err = proto.Marshal(obj1)
	if err != nil {
		t.Error("SigBlkMsg Marshal failed")
	}
	err = proto.Unmarshal(data, obj2)
	if err != nil {
		t.Error("SigBlkMsg Unmarshal failed")
	}
	fmt.Printf("after unmarshal SigBlkMsg message, message data: +%v\n\n",obj2)
}

func (m *TransferMsg) MakeMsg(msg interface{}) {
	switch msg.(type){
	case *BroadcastSigTrx:
		ptr := msg.(*BroadcastSigTrx)
		m.Msg = &TransferMsg_Msg1{Msg1:ptr}
	case *IdMsg:
		ptr := msg.(*IdMsg)
		m.Msg = &TransferMsg_Msg2{Msg2:ptr}
	case *SigBlkMsg:
		ptr := msg.(*SigBlkMsg)
		m.Msg = &TransferMsg_Msg3{Msg3:ptr}
	default:
		panic("error omsg type")
	}
}

func (m *TransferMsg) cmdType() {
	switch m.Msg.(type) {
	case *TransferMsg_Msg1:
		fmt.Println("sigtrx")
	case *TransferMsg_Msg3:
		fmt.Println("sigblk")
	}
}