package statics

import (
	"fmt"
	"strings"
)

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
		s = string(num%10+'0') + s
	}
	return s
}
func IntToString(num int) string {
	var s string = ""
	for ; num > 0; num /= 10 {
		s = string(num%10+'0') + s
	}
	return s
}

// 检查是否包含数字
func ContainNum(ss string) bool {
	for _, v := range ss {
		if v >= '0' && v <= '9' {
			return true
		}
	}
	return false
}

// 检查是否全为数字
func AllNum(ss string) bool {
	for _, v := range ss {
		if v < '0' || v > '9' {
			return false
		}
	}
	return true
}

// 检查是否包含英文字符
func ContainAlpha(ss string) bool {
	for _, v := range ss {
		if v >= 'a' && v <= 'z' || v >= 'A' && v <= 'Z' {
			return true
		}
	}
	return false
}

// 取出字符串中数字并拼接
func CatchNumber(ss string) string {
	s := ""
	for _, v := range ss {
		if v >= '0' && v <= '9' {
			s = s + string(v)
		}
	}
	return s
}

// 得到文件名称及后缀
func GetFileName(path string) string {
	index := strings.LastIndex(path, "/")
	path = path[index+1:]
	return path
}

// 计算数据大小
func FormatSize(raw int64) string {
	r := float64(raw)
	unit := []string{"B", "KB", "MB", "GB", "TB", "EB"}
	i := 0
	for i = 0; i < len(unit); i++ {
		if r < 1024 {
			if r >= 100 {
				return fmt.Sprintf("%.0f %s", r, unit[i])
			} else if r >= 10 {
				return fmt.Sprintf("%.1f %s", r, unit[i])
			} else {
				return fmt.Sprintf("%.2f %s", r, unit[i])
			}
		}
		r /= float64(int64(1) << 10)
	}
	if r >= 100 {
		return fmt.Sprintf("%.0f %s", r, unit[5])
	} else if r >= 10 {
		return fmt.Sprintf("%.1f %s", r, unit[5])
	} else {
		return fmt.Sprintf("%.2f %s", r, unit[5])
	}
	return "0B"
}
