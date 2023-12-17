package statics

func StringToInt64(s string) int64 {
	var num int64 = 0
	for i := 0; i < len(s); i++ {
		num = num*10 + int64(s[i]-'0')
	}
	return num
}
func Int64ToString(num int64) string {
	var s string = ""
	for ; num > 0; num /= 10 {
		s = s + string(num%10+'0')
	}
	return s
}
func IntToString(num int) string {
	var s string = ""
	for ; num > 0; num /= 10 {
		s = s + string(num%10+'0')
	}
	return s
}
func ContainNum(ss string) bool {
	for _, v := range ss {
		if v >= '0' && v <= '9' {
			return true
		}
	}
	return false
}
func AllNum(ss string) bool {
	for _, v := range ss {
		if v < '0' || v > '9' {
			return false
		}
	}
	return true
}
func ContainAlpha(ss string) bool {
	for _, v := range ss {
		if v >= 'a' && v <= 'z' || v >= 'A' && v <= 'Z' {
			return true
		}
	}
	return false
}
