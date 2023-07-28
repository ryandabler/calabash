package value

type Bottom struct{}

func (v *Bottom) v() vtype {
	return value
}

func (v *Bottom) Hash() string {
	return "btm"
}

// Compile time checks
var _ Value = (*Bottom)(nil)
