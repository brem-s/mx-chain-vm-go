package arwen

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/v1_4/math"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// AsyncCall holds the information about an individual async call
type AsyncCall struct {
	CallID          []byte
	Status          AsyncCallStatus
	ExecutionMode   AsyncCallExecutionMode
	Source          []byte
	Destination     []byte
	Data            []byte
	GasLimit        uint64
	ExtraGasLocked  uint64
	ValueBytes      []byte
	SuccessCallback string
	ErrorCallback   string
}

// Clone creates a deep clone of the AsyncCall
func (ac *AsyncCall) Clone() *AsyncCall {
	clone := &AsyncCall{
		CallID:          ac.CallID,
		Status:          ac.Status,
		ExecutionMode:   ac.ExecutionMode,
		Source:          make([]byte, len(ac.Source)),
		Destination:     make([]byte, len(ac.Destination)),
		Data:            make([]byte, len(ac.Data)),
		GasLimit:        ac.GasLimit,
		ExtraGasLocked:  ac.ExtraGasLocked,
		ValueBytes:      make([]byte, len(ac.ValueBytes)),
		SuccessCallback: ac.SuccessCallback,
		ErrorCallback:   ac.ErrorCallback,
	}

	copy(clone.Source, ac.Source)
	copy(clone.Destination, ac.Destination)
	copy(clone.Data, ac.Data)
	copy(clone.ValueBytes, ac.ValueBytes)

	return clone
}

// GetIdentifier returns the identifier of an async call
func (ac *AsyncCall) GetIdentifier() []byte {
	return ac.CallID
}

// GetSource returns the destination of an async call
func (ac *AsyncCall) GetSource() []byte {
	return ac.Source
}

// GetDestination returns the destination of an async call
func (ac *AsyncCall) GetDestination() []byte {
	return ac.Destination
}

// GetData returns the transaction data of the async call
func (ac *AsyncCall) GetData() []byte {
	return ac.Data
}

// GetGasLimit returns the gas limit of the current async call
func (ac *AsyncCall) GetGasLimit() uint64 {
	return ac.GasLimit
}

// GetGasLocked returns the gas locked for the async callback
func (ac *AsyncCall) GetGasLocked() uint64 {
	return ac.ExtraGasLocked
}

// GetTotalGas returns the sum of the gas limit and gas locked
func (ac *AsyncCall) GetTotalGas() uint64 {
	return math.AddUint64(ac.GasLimit, ac.ExtraGasLocked)
}

// GetValue returns the byte representation of the value of the async call
func (ac *AsyncCall) GetValue() []byte {
	return ac.ValueBytes
}

// IsLocal returns true if the async call allows for local execution
func (ac *AsyncCall) IsLocal() bool {
	return !ac.IsRemote()
}

// IsRemote returns true if the async call must be sent remotely
func (ac *AsyncCall) IsRemote() bool {
	remote := (ac.ExecutionMode == AsyncUnknown) || (ac.ExecutionMode == AsyncBuiltinFuncCrossShard)
	return remote
}

// HasCallback returns true if there is a callback to execute, depending on the status of the async call
func (ac *AsyncCall) HasCallback() bool {
	callback := ac.GetCallbackName()
	if len(callback) == 0 {
		return false
	}
	return true
}

// HasDefinedAnyCallback returns true if this AsyncCall defines at least one non-empty callback name
func (ac *AsyncCall) HasDefinedAnyCallback() bool {
	return len(ac.SuccessCallback) > 0 || len(ac.ErrorCallback) > 0
}

// UpdateStatus sets the status of the async call depending on the provided ReturnCode
func (ac *AsyncCall) UpdateStatus(returnCode vmcommon.ReturnCode) {
	ac.Status = AsyncCallResolved
	if returnCode != vmcommon.Ok {
		ac.Status = AsyncCallRejected
	}
}

// Reject sets the rejected status for this async call
func (ac *AsyncCall) Reject() {
	ac.Status = AsyncCallRejected
}

// GetCallbackName returns the name of the callback to execute, depending on
// the status of the async call
func (ac *AsyncCall) GetCallbackName() string {
	if ac.Status == AsyncCallResolved {
		return ac.SuccessCallback
	}

	return ac.ErrorCallback
}

// IsInterfaceNil returns true if there is no value under the interface
func (ac *AsyncCall) IsInterfaceNil() bool {
	return ac == nil
}

func (ac *AsyncCall) toSerializable() *SerializableAsyncCall {
	return &SerializableAsyncCall{
		CallID:          ac.CallID,
		Status:          SearializableAsyncCallStatus(ac.Status),
		ExecutionMode:   SearializableAsyncCallExecutionMode(ac.ExecutionMode),
		Source:          ac.Source,
		Destination:     ac.Destination,
		Data:            ac.Data,
		GasLimit:        ac.GasLimit,
		ExtraGasLocked:  ac.ExtraGasLocked,
		ValueBytes:      ac.ValueBytes,
		SuccessCallback: ac.SuccessCallback,
		ErrorCallback:   ac.ErrorCallback,
	}
}

func fromSerializableAsyncCalls(serializableAsyncCalls []*SerializableAsyncCall) []*AsyncCall {
	var asyncCalls = make([]*AsyncCall, 0)
	for _, serAsyncCall := range serializableAsyncCalls {
		asyncCalls = append(asyncCalls, serAsyncCall.FromSerializable())
	}
	return asyncCalls
}

// FromSerializable -
func (serAsyncCall *SerializableAsyncCall) FromSerializable() *AsyncCall {
	return &AsyncCall{
		CallID:          serAsyncCall.CallID,
		Status:          AsyncCallStatus(serAsyncCall.Status),
		ExecutionMode:   AsyncCallExecutionMode(serAsyncCall.ExecutionMode),
		Source:          serAsyncCall.Source,
		Destination:     serAsyncCall.Destination,
		Data:            serAsyncCall.Data,
		GasLimit:        serAsyncCall.GasLimit,
		ExtraGasLocked:  serAsyncCall.ExtraGasLocked,
		ValueBytes:      serAsyncCall.ValueBytes,
		SuccessCallback: serAsyncCall.SuccessCallback,
		ErrorCallback:   serAsyncCall.ErrorCallback,
	}
}
