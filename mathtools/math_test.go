package mathtools

import (
	"testing"
)

func TestRollingAverage(t *testing.T) {
	// FIXME: These are the calculated values. They look wrong, which means
	//        that the function is wrong. Check it.
	var tests = []struct {
		in []float64
		buckets int
		expected []float64
	}{
		{
			[]float64{0.0, 2.0, 2.0, 2.0, 4.0, 10.0, 10.0, 10.0, -10.0, 8.0},
			4,
			[]float64{0.4, 2.8, 2.0, 4.0, 5.6, 7.2, 4.8, 5.6, 3.6, 2.0},
		},
	}

	for _, test := range tests {
		actual := RollingAverage(test.in, test.buckets)

		for i := range(test.in) {
			if test.expected[i] != actual[i] {
				t.Errorf("Error at index %d expected %f got %f",
					i, test.expected[i], actual[i])
			}
		}
	}
}

func TestMod(t *testing.T) {
	var tests = []struct {
		in int
		modulo int
		expected int
	}{
		{
			1,
			1,
			0,
		},
		{
			10,
			6,
			4,
		},
		{
			-10,
			6,
			2,
		},
		{
			10,
			-6,
			-2,
		},
	}

	for _, test := range tests {
		actual := Mod(test.in, test.modulo)
		if test.expected != actual {
			t.Errorf("Error: %d %% %d = %d got %d",
				test.in, test.modulo, test.expected, actual)
		}
	}
}

func TestMin(t *testing.T) {
	var tests = []struct {
		in []float64
		expected float64
	}{
		{
			[]float64{1,2,3},
			1,
		},
		{
			[]float64{-1,2,3},
			-1,
		},
		{
			[]float64{-1,-2,-3},
			-3,
		},
	}

	for _, test := range tests {
		actual := Min(test.in)
		if test.expected != actual {
			t.Errorf("Error: min(%v) expected %f got %f",
				test.in, test.expected, actual)
		}
	}
}
