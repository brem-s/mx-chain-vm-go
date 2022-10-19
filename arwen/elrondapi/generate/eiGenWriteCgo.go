package elrondapigenerate

import (
	"fmt"
	"os"
)

type cgoWriter struct {
	goPackage string
	cgoPrefix string
}

func cgoType(goType string) string {
	if goType == "int32" {
		return "int32_t"
	}
	if goType == "int64" {
		return "long long"
	}
	return goType
}

func (writer *cgoWriter) cgoFuncName(funcMetadata *EIFunction) string {
	return writer.cgoPrefix + lowerInitial(funcMetadata.Name)
}

func (writer *cgoWriter) cgoImportName(funcMetadata *EIFunction) string {
	return fmt.Sprintf("C.%s", writer.cgoFuncName(funcMetadata))
}

func WriteWasmer1Cgo(out *os.File, eiMetadata *EIMetadata) {
	writer := &cgoWriter{
		goPackage: "wasmer",
		cgoPrefix: "v1_5_",
	}
	writer.writeHeader(out, eiMetadata)
	writer.writeCgoFunctions(out, eiMetadata)
	writer.writePopulateImports(out, eiMetadata)
	writer.writeGoExports(out, eiMetadata)
}

func WriteWasmer2Cgo(out *os.File, eiMetadata *EIMetadata) {
	writer := &cgoWriter{
		goPackage: "wasmer2",
		cgoPrefix: "w2_",
	}
	writer.writeHeader(out, eiMetadata)
	writer.writeCgoFunctions(out, eiMetadata)
	writer.writePopulateFuncPointers(out, eiMetadata)
	writer.writeGoExports(out, eiMetadata)
}

func (writer *cgoWriter) writeHeader(out *os.File, eiMetadata *EIMetadata) {
	out.WriteString(fmt.Sprintf(`package %s

// Code generated by elrondapi generator. DO NOT EDIT.

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// !!!!!!!!!!!!!!!!!!!!!! AUTO-GENERATED FILE !!!!!!!!!!!!!!!!!!!!!!
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef int int32_t;
//
`,
		writer.goPackage))
}

func (writer *cgoWriter) writeCgoFunctions(out *os.File, eiMetadata *EIMetadata) {
	for _, funcMetadata := range eiMetadata.AllFunctions {
		out.WriteString(fmt.Sprintf("// extern %-9s %s(void* context",
			externResult(funcMetadata.Result),
			writer.cgoFuncName(funcMetadata),
		))
		for _, arg := range funcMetadata.Arguments {
			out.WriteString(fmt.Sprintf(", %s %s", cgoType(arg.Type), arg.Name))
		}
		out.WriteString(");\n")
	}

	out.WriteString(`import "C"

import (
	"unsafe"
)

`)
}

func (writer *cgoWriter) writePopulateImports(out *os.File, eiMetadata *EIMetadata) {
	out.WriteString(`// populateWasmerImports populates imports with the ElrondEI API methods
func populateWasmerImports(imports *wasmerImports) error {
	var err error
`)

	for _, funcMetadata := range eiMetadata.AllFunctions {
		out.WriteString(fmt.Sprintf("\terr = imports.append(\"%s\", %s, %s)\n",
			lowerInitial(funcMetadata.Name),
			writer.cgoFuncName(funcMetadata),
			writer.cgoImportName(funcMetadata),
		))
		out.WriteString("\tif err != nil {\n")
		out.WriteString("\t\treturn err\n")
		out.WriteString("\t}\n\n")
	}
	out.WriteString("\treturn nil\n")
	out.WriteString("}\n")
}

func (writer *cgoWriter) writePopulateFuncPointers(out *os.File, eiMetadata *EIMetadata) {
	out.WriteString(`// populateCgoFunctionPointers populates imports with the ElrondEI API methods
func populateCgoFunctionPointers() *cWasmerVmHookPointers {
	return &cWasmerVmHookPointers{`)

	for _, funcMetadata := range eiMetadata.AllFunctions {
		out.WriteString(fmt.Sprintf("\n\t\t%s: funcPointer(%s),",
			cgoFuncPointerFieldName(funcMetadata),
			writer.cgoFuncName(funcMetadata),
		))
	}
	out.WriteString(`
	}
}
`)
}

func (writer *cgoWriter) writeGoExports(out *os.File, eiMetadata *EIMetadata) {
	for _, funcMetadata := range eiMetadata.AllFunctions {
		out.WriteString(fmt.Sprintf("\n//export %s\n",
			writer.cgoFuncName(funcMetadata),
		))
		out.WriteString(fmt.Sprintf("func %s(context unsafe.Pointer",
			writer.cgoFuncName(funcMetadata),
		))
		for _, arg := range funcMetadata.Arguments {
			out.WriteString(fmt.Sprintf(", %s %s", arg.Name, arg.Type))
		}
		out.WriteString(")")
		if funcMetadata.Result != nil {
			out.WriteString(fmt.Sprintf(" %s", funcMetadata.Result.Type))
		}
		out.WriteString(" {\n")
		out.WriteString("\tvmHooks := getVMHooksFromContextRawPtr(context)\n")
		out.WriteString("\t")
		if funcMetadata.Result != nil {
			out.WriteString("return ")
		}
		out.WriteString(fmt.Sprintf("vmHooks.%s(",
			upperInitial(funcMetadata.Name),
		))
		for argIndex, arg := range funcMetadata.Arguments {
			if argIndex > 0 {
				out.WriteString(", ")
			}
			out.WriteString(arg.Name)
		}
		out.WriteString(")\n")

		out.WriteString("}\n")
	}
}
