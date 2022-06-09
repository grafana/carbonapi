package aliasQuery

import (
	"math"
	"testing"
	"time"

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

func TestAliasByNode(t *testing.T) {
	now32 := int64(time.Now().Unix())

	tests := []th.EvalTestItem{
		{
			"aliasQuery(chan.pow.1,'chan\\.pow\\.([0-9]+)','chan.freq.\\1',\"Channel %g MHz\")",
			map[parser.MetricRequest][]*types.MetricData{
				{"chan.pow.1", 0, 1}: {types.MakeMetricData("chan.pow.1", []float64{1, 2, 3, 4, 5}, 1, now32)},
			},
			[]*types.MetricData{types.MakeMetricData("Channel 5 MHz",
				[]float64{1, 2, 3, 4, 5}, 1, now32)},
		},
		{
			"aliasQuery(chan.pow.1,'chan\\.pow\\.([0-9]+)','chan.freq.\\1',\"Channel %g MHz\")",
			map[parser.MetricRequest][]*types.MetricData{
				{"chan.pow.1", 0, 1}: {types.MakeMetricData("chan.pow.1", []float64{1, 2, 3, 4, math.NaN()}, 1, now32)},
			},
			[]*types.MetricData{
				types.MakeMetricData("Channel 4 MHz", []float64{1, 2, 3, 4, math.NaN()}, 1, now32),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		testName := tt.Target
		t.Run(testName, func(t *testing.T) {
			th.TestEvalExpr(t, &tt)
		})
	}

}
