package operation

import (
	"fmt"
	"log"
)

type OperationResult string

type OperationInterface interface {
	GetResults() map[OperationResult]bool
	IsDone() bool
	String() string
}

type BaseOperation struct {
	OperationInterface
	hasError bool
	isDone   bool
	names    map[OperationResult]bool
}

//var OperationNotCreated = errors.New("OperationNotCreated")

const (
	OperationNotExecuted OperationResult = "OperationNotExecuted"
)

func (op *BaseOperation) GetResults() map[OperationResult]bool {
	// if op.names == nil {
	// 	return nil
	// }
	// result := make([]OperationResult, len(op.names))
	// i := 0
	// for r, _ := range op.names {
	// 	result[i] = r
	// 	i++
	// }
	return op.names
}

func (op *BaseOperation) TryInit() bool {
	// if op == nil {
	// 	panic("OperationNotCreated")
	// }
	if op.isDone {
		return false
	}
	if op.names == nil {
		op.names = make(map[OperationResult]bool)
	}
	return true
}

func (op *BaseOperation) TryExecute() bool {
	// if op == nil {
	// 	panic("OperationNotCreated")
	// }
	if op.isDone {
		return false
	}
	op.isDone = true
	if op.names == nil {
		op.names = make(map[OperationResult]bool)
	}
	return true
}

func (op *BaseOperation) HasErrors() bool {
	return !op.isDone || op.hasError
}

func (op *BaseOperation) IsDone() bool {
	return op.isDone
}

func (op *BaseOperation) SetResult(result OperationResult) bool {
	op.names[result] = false
	return op.hasError
}

func (op *BaseOperation) SetResultIf(value bool, result OperationResult) bool {
	if value {
		op.names[result] = false
	}
	return op.hasError
}

func (op *BaseOperation) SetResultIfNot(value bool, result OperationResult) bool {
	if !value {
		op.names[result] = false
	}
	return op.hasError
}

func (op *BaseOperation) String() string {
	return fmt.Sprintf("%v", op.names)
}

func (op *BaseOperation) SetError(result OperationResult) bool {
	op.hasError = true
	op.names[result] = true
	log.Println((*op).String())
	return false
}

func (op *BaseOperation) SetErrorIf(value bool, result OperationResult) bool {
	if !value {
		return true
	}
	op.hasError = true
	op.names[result] = true
	return false
}

func (op *BaseOperation) SetErrorIfNot(value bool, result OperationResult) bool {
	if value {
		return true
	}
	op.hasError = true
	op.names[result] = true
	return false
}

func (op *BaseOperation) UnionResult(op1 OperationInterface) bool {
	if op1 == nil || !op1.IsDone() {
		return op.SetError(OperationNotExecuted)
	}
	for n := range op1.GetResults() {
		op.names[n] = false
	}
	return !op.hasError
}

func (op *BaseOperation) Union(op1 OperationInterface) bool {
	if op1 == nil || !op1.IsDone() {
		return op.SetError(OperationNotExecuted)
	}
	for n, r := range op1.GetResults() {
		op.names[n] = r
		if r {
			op.hasError = true
		}
	}
	return !op.hasError
}

func (op *BaseOperation) UnionWithPrefix(op1 OperationInterface, prefix string) bool {
	if op1 == nil || !op1.IsDone() {
		return op.SetError(OperationResult(prefix + string(OperationNotExecuted)))
	}
	for n, r := range op1.GetResults() {
		op.names[OperationResult(prefix+string(n))] = r
		if r {
			op.hasError = true
		}
	}
	return !op.hasError
}
