package holtWintersAberration

import (
	"testing"

	"github.com/go-graphite/carbonapi/expr/helper"
	"github.com/go-graphite/carbonapi/expr/holtwinters"
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

func TestHoltWintersAberration(t *testing.T) {
	points := int64(10)
	step := int64(600)
	origStartTime := int64(2678400) // 1970-02-01
	weekSeconds := int64(7 * 86400)
	startTime := origStartTime - weekSeconds // The parser adjusts 'from' by 'bootstrapInterval' (which defaults to 7d).

	tests := []th.EvalTestItemWithRange{
		{
			Target: "holtWintersAberration(metric.*)",
			M: map[parser.MetricRequest][]*types.MetricData{
				{"metric.*", startTime, startTime + step*points}: {
					types.MakeMetricData("metric.foo", holtwinters.GenerateTestRange(0, ((weekSeconds/step)+points)*step, step, 0), step, startTime-weekSeconds),
					types.MakeMetricData("metric.bar", holtwinters.GenerateTestRange(0, ((weekSeconds/step)+points)*step, step, 10), step, startTime-weekSeconds),
				},
			},
			Want: []*types.MetricData{
				types.MakeMetricData("holtWintersAberration(metric.foo)", []float64{0, 0, -0.33381721029946854, 0, 0, 0, 0, 0, 0, 0}, step, startTime).SetTag("holtWintersAberration", "1"),
				types.MakeMetricData("holtWintersAberration(metric.bar)", []float64{0, 0, -0.33381721029947187, 0, 0, 0, 0, 0, 0, 0}, step, startTime).SetTag("holtWintersAberration", "1"),
			},
			From:  startTime,
			Until: startTime + step*points,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Target, func(t *testing.T) {
			th.TestEvalExprWithRange(t, &tt)
		})
	}
}
