package exponentialMovingAverage

import (
	"context"
	"math"

	"github.com/go-graphite/carbonapi/expr/helper"
	"github.com/go-graphite/carbonapi/expr/interfaces"
	"github.com/go-graphite/carbonapi/expr/types"
	"github.com/go-graphite/carbonapi/pkg/parser"
)

type exponentialMovingAverage struct {
	interfaces.FunctionBase
}

func GetOrder() interfaces.Order {
	return interfaces.Any
}

func New(configFile string) []interfaces.FunctionMetadata {
	res := make([]interfaces.FunctionMetadata, 0)
	f := &exponentialMovingAverage{}
	functions := []string{"exponentialMovingAverage"}
	for _, n := range functions {
		res = append(res, interfaces.FunctionMetadata{Name: n, F: f})
	}
	return res
}

func (f *exponentialMovingAverage) Do(ctx context.Context, e parser.Expr, from, until int64, values map[parser.MetricRequest][]*types.MetricData) ([]*types.MetricData, error) {
	arg, err := helper.GetSeriesArg(ctx, e.Args()[0], from, until, values)
	if err != nil {
		return nil, err
	}

	windowSize, err := e.GetFloatArg(1)
	if err != nil {
		return nil, err
	}

	e.SetTarget("ema")

	// ugh, helper.ForEachSeriesDo does not handle arguments properly
	var results []*types.MetricData

	var constant = 2 / (windowSize + 1)
	for _, a := range arg {
		r := *a

		r.Values = make([]float64, len(a.Values))

		w := types.NewExpMovingAverage(constant)

		for i, v := range a.Values {
			if math.IsNaN(v) {
				r.Values[i] = math.NaN()
				continue
			}

			w.Push(v)
			r.Values[i] = w.Mean()
		}
		results = append(results, &r)
	}
	return results, nil
}

func (f *exponentialMovingAverage) Description() map[string]types.FunctionDescription {
	return map[string]types.FunctionDescription{
		"exponentialMovingAverage": {
			Description: "Takes a series of values and a window size and produces an exponential moving average utilizing the following formula:\n\n ema(current) = constant * (Current Value) + (1 - constant) * ema(previous)\n The Constant is calculated as:\n constant = 2 / (windowSize + 1) \n The first period EMA uses a simple moving average for its value.\n Example:\n\n code-block:: none\n\n  &target=exponentialMovingAverage(*.transactions.count, 10) \n\n &target=exponentialMovingAverage(*.transactions.count, '-10s')",
			Function:    "exponentialWeightedMovingAverage(seriesList, alpha)",
			Group:       "Calculate",
			Module:      "graphite.render.functions.custom",
			Name:        "exponentialWeightedMovingAverage",
			Params: []types.FunctionParam{
				{
					Name:     "seriesList",
					Required: true,
					Type:     types.SeriesList,
				},
				{
					Name:     "alpha",
					Required: true,
					Suggestions: types.NewSuggestions(
						0.1,
						0.5,
						0.7,
					),
					Type: types.Float,
				},
			},
		},
	}
}
