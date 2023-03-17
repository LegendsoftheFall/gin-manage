package hooks

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func CurrentDate(layout string) string {
	return time.Now().Format(layout)
}

func DateDay(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}

//func Date(date time.Time) string {
//	return date.Format("2006-01-02")
//}

func Date(date time.Time) string {
	//return date.Format("2006-01-02")
	y, m, d := date.Date()
	year := strconv.Itoa(y)
	month := strconv.Itoa(int(m))
	day := strconv.Itoa(d)
	return year + "年" + month + "月" + day + "日"
}

func TimeSubDay(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)

	return int(t1.Sub(t2).Hours() / 24)
}

func TimeSubHour(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), t2.Hour(), 0, 0, 0, time.Local)

	return int(t1.Sub(t2).Minutes() / 60)
}

func TimeSubMinute(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), t2.Hour(), t2.Minute(), 0, 0, time.Local)

	return int(t1.Sub(t2).Seconds() / 60)
}

func TimeSub(t1, t2 time.Time) (format string) {
	day := TimeSubDay(t1, t2)
	hour := TimeSubHour(t1, t2)
	minute := TimeSubMinute(t1, t2)

	if day <= 7 && day != 0 {
		format = strconv.Itoa(day) + "天前"
	} else if day == 0 && hour != 0 {
		format = strconv.Itoa(hour) + "小时前"
	} else if hour == 0 && minute != 0 {
		format = strconv.Itoa(minute) + "分钟前"
	} else {
		format = Date(t2)
	}
	return
}

func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}
