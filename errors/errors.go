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

type ReturnError struct{}

func (e ReturnError) Error() string {
	return ""
}

type BreakError struct{}

func (e BreakError) Error() string {
	return ""
}

type ContinueError struct{}

func (e ContinueError) Error() string {
	return ""
}
