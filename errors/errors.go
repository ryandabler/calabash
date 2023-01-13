package errors

type ScanError struct {
	Msg string
}

func (e ScanError) Error() string {
	return e.Msg
}

type ParseError struct {
	Msg string
}

func (e ParseError) Error() string {
	return e.Msg
}
