package utils

import (
	"fmt"
	"reflect"
	"time"
)

func Ptr[T any](v T) *T {
	return &v
}

func IfNilReturnStr(value any, valueIfNil string, valueIfNotNil string) string {
	// If the interface itself is nil
	if value == nil {
		return valueIfNil
	}

	v := reflect.ValueOf(value)

	// Only certain types can be nil in Go
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		if v.IsNil() {
			return valueIfNil
		}
	}

	return valueIfNotNil
}

// ex. 1 ธันวาคม 2568
func GetFullThaiDate(t time.Time) string {
	return fmt.Sprintf("%v %v %v", int(t.Day()), FullThaiMonthMap[int(t.Month())], t.Year()+543)
}

// ex. 1 ธ.ค. 68
func GetShortThaiDate(t time.Time) string {
	return fmt.Sprintf("%v %v %v", int(t.Day()), ShortThaiMonthMap[int(t.Month())], (t.Year()+543)%100)
}

var FullThaiMonthMap map[int]string = map[int]string{
	1:  "มกราคม",
	2:  "กุมภาพันธ์",
	3:  "มีนาคม",
	4:  "เมษายน",
	5:  "พฤษภาคม",
	6:  "มิถุนายน",
	7:  "กรกฎาคม",
	8:  "สิงหาคม",
	9:  "กันยายน",
	10: "ตุลาคม",
	11: "พฤศจิกายน",
	12: "ธันวาคม",
}

var ShortThaiMonthMap map[int]string = map[int]string{
	1:  "ม.ค.",
	2:  "ก.พ.",
	3:  "มี.ค.",
	4:  "เม.ย.",
	5:  "พ.ค.",
	6:  "มิ.ย.",
	7:  "ก.ค.",
	8:  "ส.ค.",
	9:  "ก.ย.",
	10: "ต.ค.",
	11: "พ.ย.",
	12: "ธ.ค.",
}

func ParseTime(data any) (parsedTime *time.Time, err error) {
	if data != nil {
		strDate, ok := data.(string)
		if ok && strDate != "" {
			t, err := time.Parse(time.RFC3339, strDate)
			if err == nil {
				parsedTime = &t
			}
		}
	}
	return parsedTime, err
}
func ParseString(data string) (parsedStr *string) {
	if data != "" {
		parsedStr = &data
	}
	return parsedStr
}
