package aliasQuery

import (
	"context"
	"fmt"
	"math"

	"github.com/go-graphite/carbonapi/expr/helper"
	"github.com/go-graphite/carbonapi/expr/interfaces"
	"github.com/go-graphite/carbonapi/expr/types"
	"github.com/go-graphite/carbonapi/pkg/parser"

	"regexp"
)

type aliasQuery struct {
	interfaces.FunctionBase
}

func GetOrder() interfaces.Order {
	return interfaces.Any
}

func New(configFile string) []interfaces.FunctionMetadata {
	res := make([]interfaces.FunctionMetadata, 0)
	f := &aliasQuery{}
	for _, n := range []string{"aliasQuery"} {
		res = append(res, interfaces.FunctionMetadata{Name: n, F: f})
	}
	return res
}

func (f *aliasQuery) Do(ctx context.Context, e parser.Expr, from, until int64, values map[parser.MetricRequest][]*types.MetricData) ([]*types.MetricData, error) {
	args, err := helper.GetSeriesArg(ctx, e.Args()[0], from, until, values)
	if err != nil {
		return nil, err
	}

	search, err := e.GetStringArg(1)
	if err != nil {
		return nil, err
	}

	replace, err := e.GetStringArg(2)
	if err != nil {
		return nil, err
	}

	newName, err := e.GetStringArg(3)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(search)
	if err != nil {
		return nil, err
	}

	replace = helper.Backref.ReplaceAllString(replace, "$${$1}")

	var results []*types.MetricData

	for _, a := range args {
		metric := helper.ExtractMetric(a.Name)

		r := *a
		r.Name = re.ReplaceAllString(metric, replace)
		r.Name = re.ReplaceAllString(newName, r.Name)
		// Get the last value in the series and substitute it into the new name
		r.Name = fmt.Sprintf(r.Name, getLastValue(a.Values))
		results = append(results, &r)
	}

	return results, nil
}

func getLastValue(v []float64) float64 {
	if len(v) > 0 {
		i := len(v)
		for i != 0 {
			i--
			if !math.IsNaN(v[i]) {
				return v[i]
			}
		}
	}
	return math.NaN()
}

// Description is auto-generated description, based on output of https://github.com/graphite-project/graphite-web
func (f *aliasQuery) Description() map[string]types.FunctionDescription {
	return map[string]types.FunctionDescription{
		"aliasQuery": {
			Description: "Performs a query to alias the metrics in seriesList.\n\n.. code-block:: none\n\n  &target=aliasQuery(channel.power.*,\"channel\\.power\\.([0-9]+)\",\"channel.frequency.\\1\", \"Channel %d MHz\")",
			Function:    "aliasQuery(seriesList, search, replace, newName)",
			Group:       "Alias",
			Module:      "graphite.render.functions",
			Name:        "aliasQuery",
			Params: []types.FunctionParam{
				{
					Name:     "seriesList",
					Required: true,
					Type:     types.SeriesList,
				},
				{
					Name:     "search",
					Required: true,
					Type:     types.String,
				},
				{
					Name:     "replace",
					Required: true,
					Type:     types.String,
				},
				{
					Name:     "newName",
					Required: true,
					Type:     types.String,
				},
			},
		},
	}
}
