package test

import (
	"fmt"
	"strconv"

	"github.com/almanzharoz/operation"
)

const (
	ToNumberIsEmpty   operation.OperationResult = "IsEmpty"
	ToNumberNotNumber operation.OperationResult = "NotNumber"
	MultiplyNegative  operation.OperationResult = "Negative"
	MultiplyZero      operation.OperationResult = "Zero"
)

type Number int
type NumberInterface interface {
	operation.OperationInterface
	GetNumber() int
}
type ToNumberResult struct {
	operation.BaseOperation
	Number int

	input string
}

func (op *ToNumberResult) StringToNumber(s string) bool {
	parse := func() bool {
		n, e := strconv.Atoi(s)
		op.Number = n
		return e == nil
	}
	if !op.TryExecute() {
		return false
	}
	op.input = s
	return (s != "" || op.SetError("ToNumberIsEmpty")) &&
		(parse() || op.SetError("ToNumberNotNumber"))
}

func (op *ToNumberResult) GetNumber() int {
	return op.Number
}
func (n Number) GetNumber() int {
	return int(n)
}

type MultiplyResult struct {
	operation.BaseOperation
	Result int

	input1 NumberInterface
	input2 NumberInterface
}

func (op *MultiplyResult) String() string {
	return fmt.Sprintf("{n1: %v, n2: %v, result: %d, results: %v}", op.input1, op.input2, op.Result, op.GetResults())
}

func (op *ToNumberResult) String() string {
	return fmt.Sprintf("{s: \"%s\", result: %d, results: %v}", op.input, op.Number, op.GetResults())
}

// func (op *MultiplyResult) SimpleMultiply(n1, n2 int) bool {
// 	set := func() bool {
// 		op.Result = n1 * n2
// 		return op.Result > 0 || op.SetError("MultiplyNegative")
// 	}
// 	if !op.TryExecute() {
// 		return op.HasErrors()
// 	}
// 	op.input1 = n1
// 	op.input2 = n2
// 	return op.SetErrorIf(n1 == 0, "MultiplyZero") &&
// 		(n2 != 0 || op.SetError("MultiplyZero")) &&
// 		set()
// }

func (op *MultiplyResult) Multiply(n1 NumberInterface, n2 NumberInterface) bool {
	set := func() bool {
		op.Result = op.input1.GetNumber() * op.input2.GetNumber()
		return op.Result > 0 || op.SetError("MultiplyNegative")
	}
	if !op.TryExecute() {
		return op.HasErrors()
	}
	op.input1 = n1
	op.input2 = n2
	return op.UnionWithPrefix(n1, "Number1_") && op.UnionWithPrefix(n2, "Number2_") && set()
}

// wrapper for Multiply
func (op *MultiplyResult) MultiplyStrings(s1, s2 string) bool {
	if op.IsDone() {
		return op.HasErrors()
	}
	op.TryInit()

	n1 := &ToNumberResult{}
	n2 := &ToNumberResult{}

	n1.StringToNumber(s1)
	n2.StringToNumber(s2)

	return op.Multiply(n1, n2)
}

func main() {
	m := &MultiplyResult{}

	r := m.MultiplyStrings("-1a", "2")
	fmt.Println(r, m, m.GetResults())
	//fmt.Println(n2, m)
}
