package helper

// StringToStringPtr リテラル文字列をポインタに変換します
func StringToStringPtr(s string) *string {
	tmp := s
	return &tmp
}
