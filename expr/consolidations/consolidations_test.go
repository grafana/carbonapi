package consolidations

import (
	"math"
	"testing"
)

func TestSummarizeValues(t *testing.T) {
	epsilon := math.Nextafter(1, 2) - 1
	tests := []struct {
		name         string
		function     string
		values       []float64
		xFilesFactor float32
		expected     float64
	}{
		{
			name:         "no values",
			function:     "sum",
			values:       []float64{},
			xFilesFactor: 0,
			expected:     math.NaN(),
		},
		{
			name:         "sum",
			function:     "sum",
			values:       []float64{1, 2, 3},
			xFilesFactor: 0,
			expected:     6,
		},
		{
			name:         "sum alias",
			function:     "total",
			values:       []float64{1, 2, 3},
			xFilesFactor: 0,
			expected:     6,
		},
		{
			name:         "avg",
			function:     "avg",
			values:       []float64{1, 2, 3, 4},
			xFilesFactor: 0,
			expected:     2.5,
		},
		{
			name:         "avg with nones",
			function:     "avg",
			values:       []float64{1, 2, 3, 4, math.NaN()},
			xFilesFactor: 0,
			expected:     2.5,
		},
		{
			name:         "avg xFilesFactor",
			function:     "avg",
			values:       []float64{1, 2, 3, 4, math.NaN()},
			xFilesFactor: 0.9,
			expected:     math.NaN(),
		},
		{
			name:         "max",
			function:     "max",
			values:       []float64{1, 2, 3, 4},
			xFilesFactor: 0,
			expected:     4,
		},
		{
			name:         "min",
			function:     "min",
			values:       []float64{1, 2, 3, 4},
			xFilesFactor: 0,
			expected:     1,
		},
		{
			name:         "last",
			function:     "last",
			values:       []float64{1, 2, 3, 4},
			xFilesFactor: 0,
			expected:     4,
		},
		{
			name:         "range",
			function:     "range",
			values:       []float64{1, 2, 3, 4},
			xFilesFactor: 0,
			expected:     3,
		},
		{
			name:         "median",
			function:     "median",
			values:       []float64{1, 2, 3, 10, 11},
			xFilesFactor: 0,
			expected:     3,
		},
		{
			name:         "multiply",
			function:     "multiply",
			values:       []float64{1, 2, 3, 4},
			xFilesFactor: 0,
			expected:     24,
		},
		{
			name:         "diff",
			function:     "diff",
			values:       []float64{1, 2, 3, 4},
			xFilesFactor: 0,
			expected:     -8,
		},
		{
			name:         "count",
			function:     "count",
			values:       []float64{1, 2, 3, 4},
			xFilesFactor: 0,
			expected:     4,
		},
		{
			name:         "stddev",
			function:     "stddev",
			values:       []float64{1, 2, 3, 4},
			xFilesFactor: 0,
			expected:     1.118033988749895,
		},
		{
			name:         "p50 (fallback)",
			function:     "p50",
			values:       []float64{1, 2, 3, 10, 11},
			xFilesFactor: 0,
			expected:     3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := SummarizeValues(tt.function, tt.values, tt.xFilesFactor)
			if math.Abs(actual-tt.expected) > epsilon {
				t.Errorf("actual %v, expected %v", actual, tt.expected)
			}
		})
	}

}

func TestIsValidConsolidationFunc(t *testing.T) {
	tests := []struct {
		name           string
		expectedResult bool
	}{
		{
			name:           "sum",
			expectedResult: true,
		},
		{
			name:           "avg",
			expectedResult: true,
		},
		{
			name:           "p50",
			expectedResult: true,
		},
		{
			name:           "p99.9",
			expectedResult: true,
		},
		{
			name:           "test",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidConsolidationFunc(tt.name)
			if result != tt.expectedResult {
				t.Errorf("actual %v, expected %v", result, tt.expectedResult)
			}
		})
	}

}
