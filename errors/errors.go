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

type RuntimeError struct {
	Msg string
}

func (e RuntimeError) Error() string {
	return e.Msg
}

type StaticError struct {
	Msg string
}

func (e StaticError) Error() string {
	return e.Msg
}
