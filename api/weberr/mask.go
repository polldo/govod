package weberr

// maskError wraps an error masking all its behaviors.
type maskError struct {
	Err error
}

func (e *maskError) Error() string {
	return e.Err.Error()
}
