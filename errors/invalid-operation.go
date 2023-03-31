package errors

type InvalidOperation struct {
	Operation string
}

func NewInvalidOperation(operation string) InvalidOperation {
	return InvalidOperation{Operation: operation}
}

func (e InvalidOperation) Error() string {
	return "非法的操作: " + e.Operation
}
