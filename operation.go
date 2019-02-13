package operation

// Operation Test comment
type Operation interface {
	SetErrorIfTrue(result bool, error OperationResult) bool
	SetErrorIfFalse(result bool, error OperationResult) bool
	SetResult(result bool, error OperationResult, good OperationResult) bool
	SetDone() bool
	Done() (bool, bool, OperationResult)
	HasResult(result OperationResult) bool
	Invoke() bool
}

type OperationResult = uint

const (
	opNotExecuted = OperationResult(0)
	opGood        = OperationResult(1)
	opError       = OperationResult(2)
)

type operationData struct {
	result OperationResult
}

type BaseOperation struct {
	operationData
}

func (op *operationData) SetResult(result bool, error OperationResult, good OperationResult) bool {
	if error < 4 || good < 4 {
		panic("Result must be greater than or equal 4")
	}
	if result {
		op.result |= (((op.result & 3) >> 1) + 1) | good
	} else {
		op.result = (op.result &^ 1) | opError | error
	}
	return (op.result & 1) > 0
}

func (op *operationData) SetErrorIfTrue(result bool, error OperationResult) bool {
	if error < 4 {
		panic("Result must be greater than or equal 4")
	}
	if !result {
		op.result |= ((op.result & 3) >> 1) + 1
	} else {
		op.result = (op.result &^ opGood) | opError | error
	}
	return result
}

func (op *operationData) SetErrorIfFalse(result bool, error OperationResult) bool {
	if error < 4 {
		panic("Result must be greater than or equal 4")
	}
	if result {
		op.result |= ((op.result & 3) >> 1) + 1
	} else {
		op.result = (op.result &^ opGood) | opError | error
	}
	return result
}

func (op *operationData) SetDone() bool {
	op.result |= ((op.result & 3) >> 1) + 1
	return op.result&1 > 0
}

func (op *operationData) Done() (bool, bool, OperationResult) {
	return op.result > 0, op.result&2 > 0, op.result
}

func (op *operationData) HasResult(result OperationResult) bool {
	return op.result&result > 0
}
