package op

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/coschain/contentos-go/app/table"
	"github.com/coschain/contentos-go/common/constants"
	. "github.com/coschain/contentos-go/dandelion"
	"github.com/coschain/contentos-go/prototype"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

type ContractTester struct {}

func (tester *ContractTester) Test(t *testing.T, d *Dandelion) {
	t.Run("sha256", d.Test(tester.sha256))
	t.Run("contractInfo", d.Test(tester.contractInfo))
	t.Run("requireAuth", d.Test(tester.requireAuth))
	t.Run("chainInfo", d.Test(tester.chainInfo))
	t.Run("transfer", d.Test(tester.transfer))
}

func (tester *ContractTester) sha256(t *testing.T, d *Dandelion) {
	data := make([]byte, 16)
	// sha256 for random bytes
	for i := 0; i < 10; i++ {
		_, _ = rand.Reader.Read(data)
		sum := sha256.Sum256(data)
		ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.sha256 %v, %v", BytesToJson(data), BytesToJson(sum[:])))
	}
	// sha256 for nil
	sum := sha256.Sum256(nil)
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.sha256 [], %v", BytesToJson(sum[:])))
}

func (tester *ContractTester) contractInfo(t *testing.T, d *Dandelion) {
	//
	// scenario #1, called by a user
	//

	// is_contract_called_by_user() == true
	ApplyNoError(t, d, "actor1: actor1.native_tester.is_contract_called_by_user true")

	// get_contract_caller() always returns caller account
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.get_contract_caller %q", "actor1"))
	ApplyNoError(t, d, fmt.Sprintf("actor0: actor1.native_tester.get_contract_caller %q", "actor0"))

	// get_contract_caller_contract() returns empty since the caller is not a contract
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.get_contract_caller_contract %q, %q", "", ""))

	// get_contract_name() returns name of contract
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.get_contract_name %q, %q", "actor1", "native_tester"))

	// get_contract_method() returns name of method
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.get_contract_method %q", "get_contract_method"))

	//
	// scenario #2, called by a contract
	// actor1.native_tester calls actor0.native_tester, and we test results of the callee, actor0.native_tester.
	//

	// is_contract_called_by_user() == false
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.call_is_contract_called_by_user %q, %q, false", "actor0", "native_tester"))

	// the caller account is still actor1, who sent the original operation
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.call_get_contract_caller %q, %q, %q", "actor0", "native_tester", "actor1"))

	// the caller contract is 'actor1.native_tester'
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.call_get_contract_caller_contract %q, %q, %q, %q", "actor0", "native_tester", "actor1", "native_tester"))

}

func (tester *ContractTester) requireAuth(t *testing.T, d *Dandelion) {
	//
	// scenario #1, called by a user
	//

	// require_auth of the caller account succeeds
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.require_auth %q", "actor1"))
	ApplyError(t, d, fmt.Sprintf("actor1: actor1.native_tester.require_auth %q", "actor0"))

	// require_auth of any contracts must fail, since the contract is called by user, not by contract.
	ApplyError(t, d, fmt.Sprintf("actor1: actor1.native_tester.require_auth_contract %q, %q", "actor1", "native_tester"))
	ApplyError(t, d, fmt.Sprintf("actor1: actor1.native_tester.require_auth_contract %q, %q", "actor1", ""))

	//
	// scenario #2, called by a contract
	// actor1.native_tester calls actor0.native_tester, and we test results of the callee, actor0.native_tester.
	//

	// require_auth of the caller account succeeds
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.call_require_auth %q, %q, %q", "actor0", "native_tester", "actor1"))
	ApplyError(t, d, fmt.Sprintf("actor1: actor1.native_tester.call_require_auth %q, %q, %q", "actor0", "native_tester", "actor0"))

	// require_auth of the caller contract succeeds
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.call_require_auth_contract %q, %q, %q, %q", "actor0", "native_tester", "actor1", "native_tester"))
	ApplyError(t, d, fmt.Sprintf("actor1: actor1.native_tester.call_require_auth_contract %q, %q, %q, %q", "actor0", "native_tester", "actor0", "native_tester"))
}

func (tester *ContractTester) chainInfo(t *testing.T, d *Dandelion) {
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.current_block_number %d", d.GlobalProps().HeadBlockNumber))
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.current_timestamp %d", d.GlobalProps().Time.UtcSeconds))
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.current_witness %q", d.GlobalProps().CurrentWitness.Value))
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.block_producers %s", StringsToJson(tester.blockProducers(d))))
}

func (tester *ContractTester) blockProducers(d *Dandelion) (names []string) {
	nameList := table.SWitnessOwnerWrap{Dba:d.Database()}
	_ = nameList.ForEachByOrder(nil, nil, nil, nil, func(mVal *prototype.AccountName, sVal *prototype.AccountName, idx uint32) bool {
		if table.NewSoWitnessWrap(d.Database(), mVal).GetActive() {
			names = append(names, mVal.Value)
		}
		return true
	})
	return
}

func (tester *ContractTester) transfer(t *testing.T, d *Dandelion) {
	t.Run("user_and_contract", d.Test(tester.transferBetweenUserAndContract))
	t.Run("contract_and_contract", d.Test(tester.transferBetweenContractAndContract))
}

func (tester *ContractTester) transferBetweenUserAndContract(t *testing.T, d *Dandelion) {
	a := assert.New(t)
	userBalance := d.Account("actor0").GetBalance().Value
	contractBalance := d.Contract("actor1", "native_tester").GetBalance().Value

	// user->contract: normal
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.get_user_balance %q, %d", "actor0", userBalance))
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.get_contract_balance %q, %q, %d", "actor1", "native_tester", contractBalance))
	ApplyNoError(t, d, fmt.Sprintf("actor0: %d actor1.native_tester.get_contract_sender_value %d", 123, 123))
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.get_user_balance %q, %d", "actor0", userBalance - 123))
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.get_contract_balance %q, %q, %d", "actor1", "native_tester", contractBalance + 123))
	userBalance -= 123
	contractBalance += 123

	// user->contract: too much
	NoApply(t, d, fmt.Sprintf("actor0: %d actor1.native_tester.get_contract_sender_value %d", userBalance + 1, userBalance + 1))
	NoApply(t, d, fmt.Sprintf("actor0: %d actor1.native_tester.get_contract_sender_value %d", userBalance + 100, userBalance + 100))
	NoApply(t, d, fmt.Sprintf("actor0: %d actor1.native_tester.get_contract_sender_value %d", uint64(math.MaxUint64), uint64(math.MaxUint64)))

	// user->unknown contract
	NoApply(t, d, fmt.Sprintf("actor0: %d initminer.native_tester.get_contract_sender_value %d", 1, 1))
	NoApply(t, d, fmt.Sprintf("actor0: %d xxx.native_tester.get_contract_sender_value %d", 1, 1))

	// contract->unknown user
	ApplyError(t, d, fmt.Sprintf("actor1: actor1.native_tester.transfer_to_user %q, %d", "xxxxxxx", 1))

	// contract->user: too much
	ApplyError(t, d, fmt.Sprintf("actor1: actor1.native_tester.transfer_to_user %q, %d", "actor0", contractBalance + 1))
	ApplyError(t, d, fmt.Sprintf("actor1: actor1.native_tester.transfer_to_user %q, %d", "actor0", contractBalance + 100))
	ApplyError(t, d, fmt.Sprintf("actor1: actor1.native_tester.transfer_to_user %q, %d", "actor0", uint64(math.MaxUint64)))

	// contract->user: normal
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.transfer_to_user %q, %d", "actor0", 123))
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.get_contract_balance %q, %q, %d", "actor1", "native_tester", contractBalance - 123))
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.get_user_balance %q, %d", "actor0", userBalance + 123))
	userBalance += 123
	contractBalance -= 123

	a.Equal(userBalance, d.Account("actor0").GetBalance().Value)
	a.Equal(contractBalance, d.Contract("actor1", "native_tester").GetBalance().Value)
}

func (tester *ContractTester) transferBetweenContractAndContract(t *testing.T, d *Dandelion) {
	a := assert.New(t)

	// first, fund 2 contracts
	ApplyNoError(t, d, fmt.Sprintf("%s: %d actor0.native_tester.get_contract_sender_value %d", constants.COSInitMiner, 123456, 123456))
	ApplyNoError(t, d, fmt.Sprintf("%s: %d actor1.native_tester.get_contract_sender_value %d", constants.COSInitMiner, 123456, 123456))
	contractBalance0 := d.Contract("actor0", "native_tester").GetBalance().Value
	contractBalance1 := d.Contract("actor1", "native_tester").GetBalance().Value
	a.True(contractBalance0 >= 123456 && contractBalance1 >= 123456)

	// contract -> contract: normal
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.call_get_contract_sender_value %q, %q, %d", "actor0", "native_tester", 123))
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.get_contract_balance %q, %q, %d", "actor0", "native_tester", contractBalance0 + 123))
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.get_contract_balance %q, %q, %d", "actor1", "native_tester", contractBalance1 - 123))
	contractBalance0 += 123
	contractBalance1 -= 123

	// contract -> contract: too much
	ApplyError(t, d, fmt.Sprintf("actor1: actor1.native_tester.call_get_contract_sender_value %q, %q, %d", "actor0", "native_tester", contractBalance1 + 1))
	ApplyError(t, d, fmt.Sprintf("actor1: actor1.native_tester.call_get_contract_sender_value %q, %q, %d", "actor0", "native_tester", uint64(math.MaxUint64)))

	// contract -> unknown contract
	ApplyError(t, d, fmt.Sprintf("actor1: actor1.native_tester.call_get_contract_sender_value %q, %q, %d", "xxxxxxx", "native_tester", 10))

	// contract -> self
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.call_get_contract_sender_value %q, %q, %d", "actor1", "native_tester", 123))
	ApplyNoError(t, d, fmt.Sprintf("actor1: actor1.native_tester.get_contract_balance %q, %q, %d", "actor1", "native_tester", contractBalance1))

	a.Equal(contractBalance0, d.Contract("actor0", "native_tester").GetBalance().Value)
	a.Equal(contractBalance1, d.Contract("actor1", "native_tester").GetBalance().Value)
}