package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func Random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	if (max - min) == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}

func SplitTimeInRange(a int, b int, n int, duration time.Duration) []time.Time {
	times := make([]time.Time, n)
	f := float32(float32((b - a)) / float32(n))
	currentTime := time.Now()
	for i := 0; i < n; i++ {
		x := float32(a) + (float32(i) * f)
		y := float32(a) + float32(i+1)*f
		r := Random(int(x), int(y))
		times[i] = currentTime.Add(time.Duration(r) * duration)
	}
	return times
}

func Contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsStr(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Set Difference String: A - B
func Difference(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

// get keys and values of map
func UnpackMap(m map[string]string) ([]string, []string) {
	keys := make([]string, 0, len(m))
	values := make([]string, 0, len(m))

	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

func GetCurrentDateKeys() (string, string, string, string) {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	t := time.Now().In(loc)
	dayKey := t.Format("20060102")
	_, weekKey := t.ISOWeek()
	weekKeyString := strconv.Itoa(weekKey)
	if weekKey < 10 {
		weekKeyString = "0" + weekKeyString
	}
	monthKey := t.Format("200601")
	yearKey := t.Format("2006")
	return "D" + dayKey, "W" + yearKey + weekKeyString, "M" + monthKey, "Y" + yearKey
}

func GetStartOfMonth() time.Time {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	t := time.Now().In(loc)
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

func GetEndOfMonth() time.Time {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	t := time.Now().In(loc)
	firstday := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	lastday := firstday.AddDate(0, 1, 0).Add(time.Nanosecond * -1)
	return lastday
}

func WeekStart(year, week int) time.Time {
	// Start from the middle of the year:
	loc, _ := time.LoadLocation("Asia/Kolkata")
	t := time.Date(year, 7, 1, 0, 0, 0, 0, loc)

	// Roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}

func WeekRange() (start, end time.Time) {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	t := time.Now().In(loc)
	_, week := t.ISOWeek()
	start = WeekStart(t.Year(), week)
	end = start.AddDate(0, 0, 6)
	return
}

func GetStartOfYear() time.Time {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	t := time.Now().In(loc)
	y, _, _ := t.Date()
	return time.Date(y, time.January, 1, 0, 0, 0, 0, t.Location())
}

func GetEndOfYear() time.Time {
	return GetStartOfYear().AddDate(1, 0, 0).Add(-time.Nanosecond)
}

func GetStartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func GetEndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999, t.Location())
}

func ISOFormat(t time.Time) string {
	return t.Format("2006-01-02T15:04:05.000Z")
}

func ParseISOToTime(isoString string) time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05.000Z", isoString)
	return t
}

func TimeStampInMillis(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func IncludesInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func TokenizeText(s string, size int) []string {
	n := len(s)
	results := []string{s}
	for i := 0; i <= n-size; i++ {
		results = append(results, s[i:i+4])
	}
	return results
}
