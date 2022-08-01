package percentileOfSeries

import (
	"context"

	"github.com/grafana/carbonapi/expr/consolidations"
	"github.com/grafana/carbonapi/expr/helper"
	"github.com/grafana/carbonapi/expr/interfaces"
	"github.com/grafana/carbonapi/expr/types"
	"github.com/grafana/carbonapi/pkg/parser"
)

type percentileOfSeries struct {
	interfaces.FunctionBase
}

func GetOrder() interfaces.Order {
	return interfaces.Any
}

func New(configFile string) []interfaces.FunctionMetadata {
	res := make([]interfaces.FunctionMetadata, 0)
	f := &percentileOfSeries{}
	functions := []string{"percentileOfSeries"}
	for _, n := range functions {
		res = append(res, interfaces.FunctionMetadata{Name: n, F: f})
	}
	return res
}

// percentileOfSeries(seriesList, n, interpolate=False)
func (f *percentileOfSeries) Do(ctx context.Context, e parser.Expr, from, until int64, values map[parser.MetricRequest][]*types.MetricData) ([]*types.MetricData, error) {
	// TODO(dgryski): make sure the arrays are all the same 'size'
	args, err := helper.GetSeriesArg(ctx, e.Args()[0], from, until, values)
	if err != nil {
		return nil, err
	}

	percent, err := e.GetFloatArg(1)
	if err != nil {
		return nil, err
	}

	interpolate, err := e.GetBoolNamedOrPosArgDefault("interpolate", 2, false)
	if err != nil {
		return nil, err
	}

	xFilesFactor := args[0].XFilesFactor

	return helper.AggregateSeries(e, args, func(values []float64) float64 {
		return consolidations.Percentile(values, percent, interpolate)
	}, float64(xFilesFactor))
}

// Description is auto-generated description, based on output of https://github.com/graphite-project/graphite-web
func (f *percentileOfSeries) Description() map[string]types.FunctionDescription {
	return map[string]types.FunctionDescription{
		"percentileOfSeries": {
			Description: "percentileOfSeries returns a single series which is composed of the n-percentile\nvalues taken across a wildcard series at each point. Unless `interpolate` is\nset to True, percentile values are actual values contained in one of the\nsupplied series.",
			Function:    "percentileOfSeries(seriesList, n, interpolate=False)",
			Group:       "Combine",
			Module:      "graphite.render.functions",
			Name:        "percentileOfSeries",
			Params: []types.FunctionParam{
				{
					Name:     "seriesList",
					Required: true,
					Type:     types.SeriesList,
				},
				{
					Name:     "n",
					Required: true,
					Type:     types.Integer,
				},
				{
					Default: types.NewSuggestion(false),
					Name:    "interpolate",
					Type:    types.Boolean,
				},
			},
		},
	}
}
