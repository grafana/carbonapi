package timeShift

import (
	"testing"
	"time"

	fconfig "github.com/go-graphite/carbonapi/expr/functions/config"
	"github.com/go-graphite/carbonapi/expr/interfaces"
	"github.com/go-graphite/carbonapi/expr/metadata"
	"github.com/go-graphite/carbonapi/expr/types"
	"github.com/go-graphite/carbonapi/pkg/parser"
	th "github.com/go-graphite/carbonapi/tests"
)

var (
	md []interfaces.FunctionMetadata = New("")
)

func init() {
	for _, m := range md {
		metadata.RegisterFunction(m.Name, m.F)
	}
}

func TestTimeShift(t *testing.T) {
	var startTime int64 = 86400

	tests := []th.EvalTestItemWithRange{
		// TODO(civil): Do not pass `true` resetEnd parameter in 0.15
		// Note: some of these test cases are based off of timeShift tests in Graphite-web to ensure consistency
		{
			Target: `timeShift(metric1, "-10minutes")`,
			M: map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1", From: 600, Until: 1200}: {types.MakeMetricData("metric1", []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 60, 600)},
				{Metric: "metric1", From: 0, Until: 600}:    {types.MakeMetricData("metric1", []float64{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, 60, 0)},
			},
			Want: []*types.MetricData{types.MakeMetricData("timeShift(metric1,'-600')",
				[]float64{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, 60, 600).SetTag("timeshift", "-600")},
			From:  600,
			Until: 1200,
		},
		{
			Target: `timeShift(metric1, "-10minutes")`,
			M: map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1", From: 600, Until: 1200}: {}, // Check handling of empty series
				{Metric: "metric1", From: 0, Until: 600}:    {types.MakeMetricData("metric1", []float64{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, 60, 0)},
			},
			Want:  []*types.MetricData{},
			From:  600,
			Until: 1200,
		},
		{
			Target: `timeShift(metric1, "0s", true)`,
			M: map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1", From: startTime, Until: startTime + 6}: {types.MakeMetricData("metric1", []float64{5, 4, 3, 2, 1, 0}, 1, startTime)},
				{Metric: "metric1", From: startTime, Until: startTime + 6}: {types.MakeMetricData("metric1", []float64{0, 1, 2, 3, 4, 5}, 1, startTime)},
			},
			Want: []*types.MetricData{types.MakeMetricData("timeShift(metric1,'0')",
				[]float64{0, 1, 2, 3, 4, 5}, 1, startTime).SetTag("timeshift", "0")},
			From:  startTime,
			Until: startTime + 6,
		},
		{
			Target: `timeShift(metric1, "1s", false)`,
			M: map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1", From: startTime, Until: startTime + 6}:     {types.MakeMetricData("metric1", []float64{5, 4, 3, 2, 1, 0}, 1, startTime)},
				{Metric: "metric1", From: startTime - 1, Until: startTime + 5}: {types.MakeMetricData("metric1", []float64{-1, 0, 1, 2, 3, 4}, 1, startTime-1)},
			},
			Want: []*types.MetricData{types.MakeMetricData("timeShift(metric1,'-1')",
				[]float64{-1, 0, 1, 2, 3, 4}, 1, startTime).SetTag("timeshift", "-1")},
			From:  startTime,
			Until: startTime + 6,
		},
		{
			Target: `timeShift(metric1, "1s", true)`,
			M: map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1", From: startTime, Until: startTime + 6}:     {types.MakeMetricData("metric1", []float64{3, 2, 1, 0, -1}, 1, startTime)},
				{Metric: "metric1", From: startTime - 1, Until: startTime + 5}: {types.MakeMetricData("metric1", []float64{-1, 0, 1, 2, 3}, 1, startTime-1)},
			},
			Want: []*types.MetricData{types.MakeMetricData("timeShift(metric1,'-1')",
				[]float64{-1, 0, 1, 2, 3}, 1, startTime).SetTag("timeshift", "-1")},
			From:  startTime,
			Until: startTime + 6,
		},
		{
			Target: `timeShift(metric1, "1h", false)`,
			M: map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1", From: startTime, Until: startTime + 6}:                 {types.MakeMetricData("metric1", []float64{4, 3, 2, 1, 0, -1}, 1, startTime)},
				{Metric: "metric1", From: startTime - 60*60, Until: startTime - 60*60 + 6}: {types.MakeMetricData("metric1", []float64{-1, 0, 1, 2, 3, 4}, 1, startTime-60*60)},
			},
			Want: []*types.MetricData{types.MakeMetricData("timeShift(metric1,'-3600')",
				[]float64{-1, 0, 1, 2, 3, 4}, 1, startTime).SetTag("timeshift", "-3600")},
			From:  startTime,
			Until: startTime + 6,
		},
		{
			Target: `timeShift(metric1, "1h", true)`,
			M: map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1", From: startTime, Until: startTime + 6}:                 {types.MakeMetricData("metric1", []float64{4, 3, 2, 1, 0, -1}, 1, startTime)},
				{Metric: "metric1", From: startTime - 60*60, Until: startTime - 60*60 + 6}: {types.MakeMetricData("metric1", []float64{-1, 0, 1, 2, 3, 4}, 1, startTime-60*60)},
			},
			Want: []*types.MetricData{types.MakeMetricData("timeShift(metric1,'-3600')",
				[]float64{-1, 0, 1, 2, 3, 4}, 1, startTime).SetTag("timeshift", "-3600")},
			From:  startTime,
			Until: startTime + 6,
		},
		{
			Target: `timeShift(metric1, "1d", false)`,
			M: map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1", From: startTime, Until: startTime + 6}:                 {types.MakeMetricData("metric1", []float64{4, 3, 2, 1, 0, -1}, 1, startTime)},
				{Metric: "metric1", From: startTime - 86400, Until: startTime - 86400 + 6}: {types.MakeMetricData("metric1", []float64{-1, 0, 1, 2, 3, 4}, 1, startTime-86400)},
			},
			Want: []*types.MetricData{types.MakeMetricData("timeShift(metric1,'-86400')",
				[]float64{-1, 0, 1, 2, 3, 4}, 1, startTime).SetTag("timeshift", "-86400")},
			From:  startTime,
			Until: startTime + 6,
		},
		{
			Target: `timeShift(metric1, "1d", true)`,
			M: map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1", From: startTime, Until: startTime + 6}:                 {types.MakeMetricData("metric1", []float64{4, 3, 2, 1, 0, -1}, 1, startTime)},
				{Metric: "metric1", From: startTime - 86400, Until: startTime - 86400 + 6}: {types.MakeMetricData("metric1", []float64{-1, 0, 1, 2, 3, 4}, 1, startTime-86400)},
			},
			Want: []*types.MetricData{types.MakeMetricData("timeShift(metric1,'-86400')",
				[]float64{-1, 0, 1, 2, 3, 4}, 1, startTime).SetTag("timeshift", "-86400")},
			From:  startTime,
			Until: startTime + 6,
		},
	}

	for _, tt := range tests {
		testName := tt.Target
		t.Run(testName, func(t *testing.T) {
			eval := th.EvaluatorFromFunc(md[0].F)
			th.TestEvalExprWithRange(t, eval, &tt)
		})
	}
}

func TestTimeShift_AlignDST(t *testing.T) {
	tests := []th.EvalTestItemWithRange{
		{
			Target: `timeShift(metric1, "-1d", false, true)`,
			M: map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1", From: 323869020, Until: 323872620}: {types.MakeMetricData("metric1", []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 3600, 323869020)},
				{Metric: "metric1", From: 323786220, Until: 323789820}: {types.MakeMetricData("metric1", []float64{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, 3600, 323786220)},
			},
			Want: []*types.MetricData{types.MakeMetricData("timeShift(metric1,'-86400')",
				[]float64{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, 3600, 323869020).SetTag("timeshift", "-86400")},
			From:  323869020,
			Until: 323872620,
		},
		{
			Target: `timeShift(metric1, "-1h", false, true)`,
			M: map[parser.MetricRequest][]*types.MetricData{
				{Metric: "metric1", From: 323869020, Until: 323872620}: {types.MakeMetricData("metric1", []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 3600, 323869020)},
				{Metric: "metric1", From: 323865420, Until: 323869020}: {types.MakeMetricData("metric1", []float64{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, 3600, 323865420)},
			},
			Want: []*types.MetricData{types.MakeMetricData("timeShift(metric1,'-3600')",
				[]float64{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, 3600, 323869020).SetTag("timeshift", "-3600")},
			From:  323869020,
			Until: 323872620,
		},
	}

	for _, tt := range tests {
		testName := tt.Target
		loc, _ := time.LoadLocation("Europe/Berlin")
		fconfig.Config.DefaultTimeZone = loc
		t.Run(testName, func(t *testing.T) {
			eval := th.EvaluatorFromFunc(md[0].F)
			th.TestEvalExprWithRange(t, eval, &tt)
		})
	}
}
