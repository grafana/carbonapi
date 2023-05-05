package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNudgedAggregatedValues(t *testing.T) {
	tests := []struct {
		name      string
		values    []float64
		step      int64
		start     int64
		mdp       int64
		want      []float64
		wantStep  int64
		wantStart int64
	}{
		{
			name:      "empty",
			values:    []float64{},
			step:      60,
			mdp:       100,
			want:      []float64{},
			wantStep:  60,
			wantStart: 0,
		},
		{
			name:      "one point",
			values:    []float64{1, 2, 3, 4},
			start:     10,
			step:      10,
			mdp:       1,
			want:      []float64{10},
			wantStep:  40,
			wantStart: 40,
		},
		{
			name:      "can't trim if response ends empty",
			values:    []float64{1, 2, 3, 4},
			start:     7,
			step:      3,
			mdp:       2,
			want:      []float64{3, 7},
			wantStart: 10,
			wantStep:  6,
		},
		{
			name:      "no trim due to not many points",
			values:    []float64{1, 2, 3, 4},
			step:      10,
			start:     20,
			mdp:       1,
			want:      []float64{10},
			wantStep:  40,
			wantStart: 50,
		},

		{
			name:      "should trim the first point",
			values:    []float64{1, 2, 3, 4, 5, 6},
			start:     20,
			step:      10,
			mdp:       3,
			want:      []float64{5, 9, 6},
			wantStep:  20,
			wantStart: 40,
		},
		{
			name:      "should be stable with previous",
			values:    []float64{2, 3, 4, 5, 6, 7},
			start:     30,
			step:      10,
			mdp:       3,
			want:      []float64{5, 9, 13},
			wantStep:  20,
			wantStart: 40,
		},
		{
			name:      "a bit more data",
			values:    []float64{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
			start:     20,
			step:      10,
			mdp:       3,
			want:      []float64{40, 50},
			wantStep:  50,
			wantStart: 100,
		},
		{
			name:      "a bit more data even",
			values:    []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10.0, 11, 12, 13, 14},
			start:     10,
			step:      10,
			mdp:       3,
			want:      []float64{15, 40, 50},
			wantStep:  50,
			wantStart: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := MakeMetricData("test", tt.values, tt.step, tt.start)
			input.ConsolidationFunc = "sum"
			ConsolidateJSON(tt.mdp, true, []*MetricData{input})

			got := input.AggregatedValues()
			gotStep := input.AggregatedTimeStep()
			gotStart := input.AggregatedStartTime()

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantStep, gotStep)
			assert.Equal(t, tt.wantStart, gotStart)
		})
	}
}
