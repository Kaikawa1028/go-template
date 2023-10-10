package date

//formatとして使用できる数字や文字は決まっているため、以下のような指定になっている。
//指定できるformat https://go.dev/src/time/format.go
const (
	FormatTypeDayTime         = "2006-01-02 15:04:05"
	FormatTypeDay             = "2006-01-02"
	FormatTypeTime            = "15:04"
	FormatTypeDayTimeFileName = "20060102150405"
)
