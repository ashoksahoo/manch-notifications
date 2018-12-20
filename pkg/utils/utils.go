package utils

import(
    "math/rand"
    "time"
)

func Random(min, max int) int {
    rand.Seed(time.Now().UTC().UnixNano())
    if (max-min) == 0 {
        return max
    }
    return rand.Intn(max - min) + min
}

func SplitTimeInRange(a int, b int, n int, duration time.Duration) []time.Time {
    times := make([]time.Time, n)
    f := (b - a) / n
    currentTime := time.Now()
    for i := 0; i < n; i++ {
        r := Random((a + (i * f)), (a + (i + 1) * f))
        times[i] = currentTime.Add(time.Duration(r) * duration);
    }
    return times;
}
