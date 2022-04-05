package contracts

import (
	"errors"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/arwen"
	mock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/context"
	test "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/testcommon"
	"github.com/stretchr/testify/require"
)

// TransferToThirdPartyAsyncChildMock is an exposed mock contract method
func TransferToThirdPartyAsyncChildMock(instanceMock *mock.InstanceMock, config interface{}) {
	instanceMock.AddMockMethod("transferToThirdParty", func() *mock.InstanceMock {
		testConfig := config.(*test.TestConfig)
		host := instanceMock.Host
		instance := mock.GetMockInstance(host)
		t := instance.T

		metering := host.Metering()
		err := metering.UseGasBounded(testConfig.GasUsedByChild)
		if err != nil {
			host.Runtime().SetRuntimeBreakpointValue(arwen.BreakpointOutOfGas)
			return instance
		}

		arguments := host.Runtime().Arguments()
		outputContext := host.Output()

		if len(arguments) != 3 {
			host.Runtime().SignalUserError("wrong num of arguments")
			return instance
		}

		behavior := byte(0)
		if len(arguments[2]) != 0 {
			behavior = arguments[2][0]
		}
		err = handleChildBehaviorArgument(host, behavior)
		if err != nil {
			return instance
		}

		scAddress := host.Runtime().GetSCAddress()
		valueToTransfer := big.NewInt(0).SetBytes(arguments[0])
		err = outputContext.Transfer(
			testConfig.GetThirdPartyAddress(),
			scAddress,
			0,
			0,
			valueToTransfer,
			arguments[1],
			0)
		require.Nil(t, err)
		outputContext.Finish([]byte("thirdparty"))

		valueToTransfer = big.NewInt(testConfig.TransferToVault)
		err = outputContext.Transfer(
			testConfig.GetVaultAddress(),
			scAddress,
			0,
			0,
			valueToTransfer,
			[]byte{},
			0)
		require.Nil(t, err)
		outputContext.Finish([]byte("vault"))

		host.Storage().SetStorage(test.ChildKey, test.ChildData)

		return instance
	})
}

func handleChildBehaviorArgument(host arwen.VMHost, behavior byte) error {
	if behavior == 1 {
		host.Runtime().SignalUserError("child error")
		return errors.New("behavior / child error")
	}
	if behavior == 2 {
		for {
			host.Output().Finish([]byte("loop"))
		}
	}
	host.Output().Finish([]byte{behavior})
	return nil
}
