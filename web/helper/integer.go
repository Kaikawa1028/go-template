package helper

// IntToIntPtr 数値リテラルをポインタ型に変換します
func IntToIntPtr(s int) *int {
	tmp := s
	return &tmp
}
