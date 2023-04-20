package smartSummarize

import (
	"math"
	"testing"

	"github.com/go-graphite/carbonapi/expr/helper"
	"github.com/go-graphite/carbonapi/expr/metadata"
	"github.com/go-graphite/carbonapi/expr/types"
	"github.com/go-graphite/carbonapi/pkg/parser"
	th "github.com/go-graphite/carbonapi/tests"
)

func init() {
	md := New("")
	evaluator := th.EvaluatorFromFunc(md[0].F)
	metadata.SetEvaluator(evaluator)
	helper.SetEvaluator(evaluator)
	for _, m := range md {
		metadata.RegisterFunction(m.Name, m.F)
	}
}

func TestSummarizeEmptyData(t *testing.T) {
	tests := []th.EvalTestItem{
		{
			"smartSummarize(metric1,'1hour','sum','1y')",
			map[parser.MetricRequest][]*types.MetricData{
				{"foo.bar", 0, 1}: {},
			},
			[]*types.MetricData{},
		},
	}

	for _, tt := range tests {
		testName := tt.Target
		t.Run(testName, func(t *testing.T) {
			th.TestEvalExpr(t, &tt)
		})
	}

}

func TestEvalSummarize(t *testing.T) {
	tests := []th.SummarizeEvalTestItem{
		{
			"smartSummarize(metric1,'1hour','sum','1y')",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1", 0, 1}: {types.MakeMetricData("metric1", generateValues(0, 1800+3.5*3600, 1), 1, 0)},
			},
			[]float64{6478200, 19438200, 32398200, 45358200},
			"smartSummarize(metric1,'1hour','sum','1y')",
			3600,
			0,
			14400,
		},
		{
			"smartSummarize(metric1,'1hour','sum','y')",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1", 0, 1}: {types.MakeMetricData("metric1", generateValues(0, 1800+3.5*3600, 1), 1, 0)},
			},
			[]float64{6478200, 19438200, 32398200, 45358200},
			"smartSummarize(metric1,'1hour','sum','y')",
			3600,
			0,
			14400,
		},
		{
			"smartSummarize(metric1,'1hour','sum','month')",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1", 0, 1}: {types.MakeMetricData("metric1", generateValues(0, 1800+3.5*3600, 1), 1, 0)},
			},
			[]float64{6478200, 19438200, 32398200, 45358200},
			"smartSummarize(metric1,'1hour','sum','month')",
			3600,
			0,
			14400,
		},
		{
			"smartSummarize(metric1,'1hour','sum','1month')",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1", 0, 1}: {types.MakeMetricData("metric1", generateValues(0, 1800+3.5*3600, 1), 1, 0)},
			},
			[]float64{6478200, 19438200, 32398200, 45358200},
			"smartSummarize(metric1,'1hour','sum','1month')",
			3600,
			0,
			14400,
		},
		{
			"smartSummarize(metric1,'1minute','sum','minute')",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1", 0, 1}: {types.MakeMetricData("metric1", generateValues(0, 240, 1), 1, 0)},
			},
			[]float64{1770, 5370, 8970, 12570},
			"smartSummarize(metric1,'1minute','sum','minute')",
			60,
			0,
			240,
		},
		{
			"smartSummarize(metric1,'1minute','sum','1minute')",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1", 0, 1}: {types.MakeMetricData("metric1", generateValues(0, 240, 1), 1, 0)},
			},
			[]float64{1770, 5370, 8970, 12570},
			"smartSummarize(metric1,'1minute','sum','1minute')",
			60,
			0,
			240,
		},
		{
			"smartSummarize(metric1,'1minute','avg','minute')",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1", 0, 1}: {types.MakeMetricData("metric1", generateValues(0, 240, 1), 1, 0)},
			},
			[]float64{29.5, 89.5, 149.5, 209.5},
			"smartSummarize(metric1,'1minute','avg','minute')",
			60,
			0,
			240,
		},
		{
			"smartSummarize(metric1,'1minute','last','minute')",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1", 0, 1}: {types.MakeMetricData("metric1", generateValues(0, 240, 1), 1, 0)},
			},
			[]float64{59, 119, 179, 239},
			"smartSummarize(metric1,'1minute','last','minute')",
			60,
			0,
			240,
		},
		{
			"smartSummarize(metric1,'4hours','sum','weeks')",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1", 0, 1}: {types.MakeMetricData("metric1", generateValues(0, 14400, 1), 1, 0)},
			},
			[]float64{103672800},
			"smartSummarize(metric1,'4hours','sum','weeks')",
			14400,
			0,
			14400,
		},
		{
			"smartSummarize(metric1,'1d','sum','days')",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1", 0, 1}: {types.MakeMetricData("metric1", generateValues(0, 86400, 60), 60, 0)},
			},
			[]float64{62164800},
			"smartSummarize(metric1,'1d','sum','days')",
			86400,
			0,
			86400,
		},
		{
			"smartSummarize(metric1,'1minute','sum','seconds')",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1", 0, 1}: {types.MakeMetricData("metric1", generateValues(0, 240, 1), 1, 0)},
			},
			[]float64{1770, 5370, 8970, 12570},
			"smartSummarize(metric1,'1minute','sum','seconds')",
			60,
			0,
			240,
		},
		{
			"smartSummarize(metric1,'1hour','max','hours')",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1", 0, 1}: {types.MakeMetricData("metric1", generateValues(0, 14400, 1), 1, 0)},
			},
			[]float64{3599, 7199, 10799, 14399},
			"smartSummarize(metric1,'1hour','max','hours')",
			3600,
			0,
			14400,
		},
		{
			"smartSummarize(metric1,'6m','sum', 'minutes')", // Test having a smaller interval than the data's step
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1", 0, 1}: {types.MakeMetricData("metric1", []float64{
					2, 4, 6}, 600, 1410345000)},
			},
			[]float64{2, 4, math.NaN(), 6, math.NaN()},
			"smartSummarize(metric1,'6m','sum','minutes')",
			360,
			1410345000,
			1410345000 + 3*600,
		},
		{
			"smartSummarize(metric2,'2minute','sum')",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric2", 0, 1}: {types.MakeMetricData("metric2", []float64{1, 2, 3, 4}, 60, 0)},
			},
			[]float64{3, 7},
			"smartSummarize(metric2,'2minute','sum')",
			120,
			0,
			240,
		},
	}

	for _, tt := range tests {
		th.TestSummarizeEvalExpr(t, &tt)
	}
}

func TestFunctionUseNameWithWildcards(t *testing.T) {
	tests := []th.MultiReturnEvalTestItem{
		{
			"smartSummarize(metric1.*,'1minute','last')",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1.*", 0, 1}: {
					types.MakeMetricData("metric1.foo", generateValues(0, 240, 1), 1, 0),
					types.MakeMetricData("metric1.bar", generateValues(0, 240, 1), 1, 0),
				},
			},
			"smartSummarize",
			map[string][]*types.MetricData{
				"smartSummarize(metric1.foo,'1minute','last')": {types.MakeMetricData("smartSummarize(metric1.foo,'1minute','last')",
					[]float64{59, 119, 179, 239}, 60, 0).SetTag("smartSummarize", "60").SetTag("smartSummarizeFunction", "last")},
				"smartSummarize(metric1.bar,'1minute','last')": {types.MakeMetricData("smartSummarize(metric1.bar,'1minute','last')",
					[]float64{59, 119, 179, 239}, 60, 0).SetTag("smartSummarize", "60").SetTag("smartSummarizeFunction", "last")},
			},
		},
	}

	for _, tt := range tests {
		testName := tt.Target
		t.Run(testName, func(t *testing.T) {
			th.TestMultiReturnEvalExpr(t, &tt)
		})
	}
}

func generateValues(start, stop, step int64) (values []float64) {
	for i := start; i < stop; i += step {
		values = append(values, float64(i))
	}
	return
}
