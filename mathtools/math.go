package mathtools

import (
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

// Centered rolling average
func RollingAverage(in []float64, points int) []float64 {
	out := make([]float64, len(in))

	if points%2 == 0 {
		points = points + 1
	}

	buckets := make([]float64, points)

	for i := range in {
		for j := 0; j < points; j++ {
			index := Mod(((i + j) - points/2), len(in))
			buckets[j] = in[index]
		}
		out[i] = Sum(buckets) / float64(points)
	}

	return out
}

// Python style mod, rather than the remainder that go '%' operator provides
func Mod(d, m int) int {
	var res int = d % m
	if (res < 0 && m > 0) || (res > 0 && m < 0) {
		return res + m
	}

	return res
}

func Sum(in []float64) float64 {
	return floats.Sum(in)
}

func Min(in []float64) float64 {
	return floats.Min(in)
}

func Max(in []float64) float64 {
	return floats.Max(in)
}

func Histogram(in []float64, nBuckets int) []int {
	buckets := make([]int, nBuckets)
	max := Max(in)

	for _, d := range in {
		bucket := int(math.Floor(d / (max / float64(nBuckets))))
		// Max case cases an index-out-of-range depending on floating
		// point math
		if bucket >= nBuckets {
			bucket -= 1
		}
		buckets[bucket] += 1
	}

	return buckets
}

func Mean(in []float64) float64 {
	return stat.Mean(in, nil)
}

func StdDev(in []float64) float64 {
	return stat.StdDev(in, nil)
}

// squelch will zero all values less than floor
func Squelch(in []float64, floor float64) []float64 {
	ret := make([]float64, len(in))
	for i, d := range in {
		if d > floor {
			ret[i] = d
		} else {
			ret[i] = 0
		}
	}

	return ret
}

// AbsAroundMean creates a zero point at the mean, and converts all values to
// their absolute around this zero
func AbsAroundMean(in []float64) []float64 {
	ret := make([]float64, len(in))
	mean := Mean(in)

	for i, d := range in {
		if d < mean {
			ret[i] = mean - d
		} else {
			ret[i] = d - mean
		}
	}

	return ret
}

// Denoise removes spikes that only last one sample
func Denoise(in []float64) []float64 {
	ret := make([]float64, len(in))

	for i, d := range in {
		if len(in)-2 < i {
			ret[i] = d
			continue
		}
		if i < 2 {
			ret[i] = d
			continue
		}
		if in[i-1] < d && in[i+1] < d {
			ret[i] = in[i-1] + in[i+1]/2
			continue
		}
		if in[i-1] > d && in[i+1] > d {
			ret[i] = in[i-1] + in[i+1]/2
			continue
		}
		ret[i] = d
	}

	return ret
}

// Edge finder tries to square up pulses
func EdgeFinder(in []float64, buckets int) []float64 {
	ret := make([]float64, len(in))

	max := Max(in)
	stddev := StdDev(in)

	high := false
	for i, d := range in {
		if d > stddev {
			ret[i] = max
			high = true
			continue
		}
		if high {
			end := i + buckets
			if end >= len(in) {
				end = len(in) - 1
			}
			next := in[i:end]
			localmax := Max(next)
			if localmax > stddev {
				ret[i] = max
			} else {
				high = false
				ret[i] = 0
			}
			continue
		}
		ret[i] = 0
	}

	return ret
}
