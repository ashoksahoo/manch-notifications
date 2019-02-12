package utils

import (
	"math/rand"
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
