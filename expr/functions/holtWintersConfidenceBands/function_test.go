package holtWintersConfidenceBands

import (
	"github.com/go-graphite/carbonapi/expr/helper"
	"github.com/go-graphite/carbonapi/expr/metadata"
	"github.com/go-graphite/carbonapi/expr/types"
	"github.com/go-graphite/carbonapi/pkg/parser"
	th "github.com/go-graphite/carbonapi/tests"
	"testing"
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

func TestHoltWintersConfidenceBands(t *testing.T) {
	var startTime int64 = 2678400
	//var endTime int64 = 2678940
	var step int64 = 600
	var points int64 = 10

	// expected from 2073600
	// until 2684400
	// actual from: 2592000
	// actual until: 2598000
	// bootstrap interval is 604800 so this is subtracted from 'from' before grabbing the series args

	tests := []th.EvalTestItemWithRange{
		{
			Target: "holtWintersConfidenceBands(metric1)",
			M: map[parser.MetricRequest][]*types.MetricData{
				{"metric1", startTime - 604800, (startTime + 6000)}: {types.MakeMetricData("metric1", []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, step, startTime)},
			},
			Want: []*types.MetricData{
				types.MakeMetricData("holtWintersConfidenceLower(metric1)", []float64{0.2841206166091448, 1.0581027098774411, 0.3338172102994683, 0.5116859493263242, -0.18199175514936972, 0.2366173792019426, -1.2941554508809152, -0.513426806531049, -0.7970905542723132, 0.09868900726536012}, step, startTime).SetTag("holtWintersConfidenceLower", "1"),
				types.MakeMetricData("holtWintersConfidenceUpper(metric1)", []float64{8.424944558327624, 9.409422251880809, 10.607070189221787, 10.288439865038768, 9.491556863132963, 9.474595784593738, 8.572310478053845, 8.897670449095346, 8.941566968508148, 9.409728797779282}, step, startTime).SetTag("holtWintersConfidenceUpper", "1"),
			},
			From:  startTime,
			Until: startTime + 6000,
		},
	}

	for _, tt := range tests {
		testName := tt.Target
		t.Run(testName, func(t *testing.T) {
			th.TestEvalExprWithRange(t, &tt)
		})
	}
}
