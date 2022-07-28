package linearRegression

import (
	"math"
	"testing"
	"time"

	"github.com/grafana/carbonapi/expr/helper"
	"github.com/grafana/carbonapi/expr/metadata"
	"github.com/grafana/carbonapi/expr/types"
	"github.com/grafana/carbonapi/pkg/parser"
	th "github.com/grafana/carbonapi/tests"
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

func TestFunction(t *testing.T) {
	now32 := int64(time.Now().Unix())

	tests := []th.EvalTestItem{
		{
			"linearRegression(metric1)",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1", 0, 1}: {
					types.MakeMetricData("metric1",
						[]float64{1, 2, math.NaN(), math.NaN(), 5, 6}, 1, now32),
				},
			},
			[]*types.MetricData{
				types.MakeMetricData("linearRegression(metric1)",
					[]float64{1, 2, 3, 4, 5, 6}, 1, now32),
			},
		},
		{
			"linearRegression(metric1,'1h')",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1", 0, 1}: {
					types.MakeMetricData("metric1",
						[]float64{1, 2, math.NaN(), math.NaN(), 5, 6}, 1, now32),
				},
			},
			[]*types.MetricData{
				types.MakeMetricData("linearRegression(metric1,'1h')",
					[]float64{1, 2, 3, 4, 5, 6}, 1, now32-3600),
			},
		},
		{
			"linearRegression(metric1,'1h','1m')",
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1", 0, 1}: {
					types.MakeMetricData("metric1",
						[]float64{1, 2, math.NaN(), math.NaN(), 5, 6}, 1, now32),
				},
			},
			[]*types.MetricData{
				types.MakeMetricData("linearRegression(metric1,'1h','1m')",
					[]float64{1, 2, 3, 4, 5, 6}, 1, now32-3600),
			},
		},
	}

	for _, tt := range tests {
		testName := tt.Target
		t.Run(testName, func(t *testing.T) {
			th.TestEvalExpr(t, &tt)
		})
	}

}
