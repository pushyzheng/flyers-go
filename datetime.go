package times

import "time"

const (
	StandardDatetimeLayout = "2006-01-02 15:04:05"
	StandardDateLayout     = "2006-01-02"
)

func GetTimeStamp() int64 {
	return time.Now().UnixMilli()
}

func ParseDatetime(s string) (time.Time, error) {
	return time.Parse(StandardDatetimeLayout, s)
}

func FormatDatetime(t time.Time) string {
	return t.Format(StandardDatetimeLayout)
}

func ParseDate(s string) (time.Time, error) {
	return time.Parse(StandardDateLayout, s)
}

func FormatDate(t time.Time) string {
	return t.Format(StandardDateLayout)
}
