package operation

import (
	"errors"
	"fmt"
	"log"
)

type OperationResult string

type OperationInterface interface {
	GetResults() map[OperationResult]bool
	IsDone() bool
	//String() string
	HasGoodResult(result OperationResult) bool
	HasErrorResult(result OperationResult) bool
	GetError() error
	//SetInfo(log string)
	HasErrors() bool
}

type BaseOperation struct {
	OperationInterface
	hasError bool
	isDone   bool
	names    map[OperationResult]bool // true - предупреждение, false - ошибка
	error    error
}

//var OperationNotCreated = errors.New("OperationNotCreated")

const (
	OperationNotExecuted OperationResult = "OperationNotExecuted"
	OperationPanic                       = "Panic"
)

func (op *BaseOperation) GetResults() map[OperationResult]bool {
	return op.names
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

	op.names[result] = true
	return true
}

func (op *BaseOperation) SetResultIf(value bool, result OperationResult) bool {
	if value {
		op.names[result] = true
	}
	return !value
}

func (op *BaseOperation) SetResultIfNot(value bool, result OperationResult) bool {
	if !value {
		op.names[result] = true
	}
	return value
}

func (op *BaseOperation) String() string {
	return fmt.Sprintf("%v", op.names)
}

func (op *BaseOperation) SetError(result OperationResult) bool {
	op.hasError = true
	op.names[result] = false
	log.Println((*op).String())
	return false
}

func (op *BaseOperation) SetErrorIf(value bool, result OperationResult) bool {
	if !value {
		return true
	}
	op.hasError = true
	op.names[result] = false
	return false
}

func (op *BaseOperation) SetErrorIfNot(value bool, result OperationResult) bool {
	if value {
		return true
	}
	op.hasError = true
	op.names[result] = false
	return false
}

func (op *BaseOperation) TryAll(values ...bool) bool {
	result := true
	for _, v := range values {
		result = result && v
	}
	return result
}

func (op *BaseOperation) UnionResult(op1 OperationInterface) bool {
	if op1 == nil || !op1.IsDone() {
		return op.SetError(OperationNotExecuted)
	}
	result := true
	for n, v := range op1.GetResults() {
		if !v {
			result = false
		}
		op.names[n] = true
	}
	return result
}

func (op *BaseOperation) UnionResultWithPrefix(op1 OperationInterface, prefix string) bool {
	if op1 == nil || !op1.IsDone() {
		return op.SetError(OperationNotExecuted)
	}
	result := true
	for n, v := range op1.GetResults() {
		if !v {
			result = false
		}
		op.names[OperationResult(prefix+string(n))] = true
	}
	return result
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

func (op *BaseOperation) HasGoodResult(result OperationResult) bool {
	if !op.isDone {
		op.error = errors.New("HasGoodResult under not executed operation")
		return false
	}
	v, ok := op.names[result]
	return ok && v
}

func (op *BaseOperation) HasErrorResult(result OperationResult) bool {
	if !op.isDone {
		op.error = errors.New("HasErrorResult under not executed operation")
		return false
	}
	if !op.hasError {
		return false
	}
	v, ok := op.names[result]
	return ok && !v
}

func (op *BaseOperation) GetPanic() error {
	return op.error
}
func (op *BaseOperation) SetPanic(error error) bool {
	op.error = error
	op.hasError = true
	op.names[OperationPanic] = false
	log.Println(error)
	return false
}

func (op *BaseOperation) SetPanicWithError(err OperationResult, error error) bool {
	op.error = error
	op.hasError = true
	op.names[err] = false
	log.Println(error)
	return false
}
