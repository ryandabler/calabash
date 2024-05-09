package value

type Bottom struct{}

func (v *Bottom) v() vtype {
	return value
}

func (v *Bottom) Hash() string {
	return "btm"
}

func (v *Bottom) Proto() *Proto {
	return nil
}

func (v *Bottom) Inherit(_ *Proto) Value {
	return v
}

// Compile time checks
var _ Value = (*Bottom)(nil)
