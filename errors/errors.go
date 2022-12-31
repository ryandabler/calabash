package errors

type ScanError struct {
	Msg string
}

func (e ScanError) Error() string {
	return e.Msg
}
