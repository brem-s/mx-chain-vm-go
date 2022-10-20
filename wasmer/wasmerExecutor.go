package wasmer

import (
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/executor"
)

// WasmerExecutor oversees the creation of Wasmer instances and execution.
type WasmerExecutor struct {
	eiFunctionNames vmcommon.FunctionNames
}

// NewExecutor creates a new wasmer executor.
func NewExecutor() (*WasmerExecutor, error) {
	functionNames, err := injectCgoFunctionPointers()
	if err != nil {
		return nil, err
	}
	return &WasmerExecutor{
		eiFunctionNames: functionNames,
	}, nil
}

// SetOpcodeCosts sets gas costs globally inside the Wasmer executor.
func (wasmerExecutor *WasmerExecutor) SetOpcodeCosts(opcodeCosts *[executor.OpcodeCount]uint32) {
	SetOpcodeCosts(opcodeCosts)
}

// SetRkyvSerializationEnabled controls a Wasmer flag.
func (wasmerExecutor *WasmerExecutor) SetRkyvSerializationEnabled(enabled bool) {
	SetRkyvSerializationEnabled(enabled)
}

// SetSIGSEGVPassthrough controls a Wasmer flag.
func (wasmerExecutor *WasmerExecutor) SetSIGSEGVPassthrough() {
	SetSIGSEGVPassthrough()
}

func (wasmerExecutor *WasmerExecutor) FunctionNames() vmcommon.FunctionNames {
	return wasmerExecutor.eiFunctionNames
}

// NewInstanceWithOptions creates a new Wasmer instance from WASM bytecode,
// respecting the provided options
func (wasmerExecutor *WasmerExecutor) NewInstanceWithOptions(
	contractCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	return NewInstanceWithOptions(contractCode, options)
}

// NewInstanceFromCompiledCodeWithOptions creates a new Wasmer instance from
// precompiled machine code, respecting the provided options
func (wasmerExecutor *WasmerExecutor) NewInstanceFromCompiledCodeWithOptions(
	compiledCode []byte,
	options executor.CompilationOptions,
) (executor.Instance, error) {
	return NewInstanceFromCompiledCodeWithOptions(compiledCode, options)
}