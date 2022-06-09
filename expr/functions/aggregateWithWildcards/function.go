package aggregateWithWildcards

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-graphite/carbonapi/expr/consolidations"
	"github.com/go-graphite/carbonapi/expr/helper"
	"github.com/go-graphite/carbonapi/expr/interfaces"
	"github.com/go-graphite/carbonapi/expr/types"
	"github.com/go-graphite/carbonapi/pkg/parser"
)

type aggregateWithWildcards struct {
	interfaces.FunctionBase
}

func GetOrder() interfaces.Order {
	return interfaces.Any
}

func New(configFile string) []interfaces.FunctionMetadata {
	f := &aggregateWithWildcards{}
	res := make([]interfaces.FunctionMetadata, 0)
	for _, n := range []string{"aggregateWithWildcards"} {
		res = append(res, interfaces.FunctionMetadata{Name: n, F: f})
	}

	// Also register aliases for each and every summarizer
	for _, n := range consolidations.AvailableSummarizers {
		res = append(res,
			interfaces.FunctionMetadata{Name: n, F: f},
			interfaces.FunctionMetadata{Name: n + "Series", F: f},
		)
	}
	return res
}

func (f *aggregateWithWildcards) Do(ctx context.Context, e parser.Expr, from, until int64, values map[parser.MetricRequest][]*types.MetricData) ([]*types.MetricData, error) {
	var args []*types.MetricData
	isAggregateFunc := true

	callback, err := e.GetStringArg(1)
	if err != nil {
		if e.Target() == "aggregate" {
			return nil, err
		} else {
			args, err = helper.GetSeriesArgsAndRemoveNonExisting(ctx, e, from, until, values)
			if err != nil {
				return nil, err
			}
			callback = strings.Replace(e.Target(), "Series", "", 1)
			isAggregateFunc = false
		}
	} else {
		args, err = helper.GetSeriesArg(ctx, e.Args()[0], from, until, values)
		if err != nil {
			return nil, err
		}
	}
	positions, err := e.GetIntArgs(2)
	if err != nil {
		return nil, err
	}

	aggFunc, ok := consolidations.ConsolidationToFunc[callback]
	if !ok {
		return nil, fmt.Errorf("unsupported consolidation function %s", callback)
	}
	target := fmt.Sprintf("%sSeries", callback)

	e.SetTarget(target)
	if isAggregateFunc {
		e.SetRawArgs(e.Args()[0].Target())
	}
	//name := fmt.Sprintf("%s(%s)", e.Target(), e.RawArgs())

	groups := make(map[string][]*types.MetricData)
	var keys []string

	for _, a := range args {
		key := filterNodesByPositions(a.Name, positions)
		_, ok := groups[key]
		if !ok {
			keys = append(keys, key)
		}
		groups[key] = append(groups[key], a)
	}

	results := make([]*types.MetricData, 0, len(groups))

	for _, key := range keys {
		res, err := helper.AggregateSeries(e, groups[key], aggFunc)
		if err != nil {
			return nil, err
		}
		res[0].Name = key
		results = append(results, res...)
	}
	return results, nil
}

func filterNodesByPositions(name string, nodes []int) string {
	parts := strings.Split(name, ".")
	var newName []string
	for i, word := range parts {
		var found bool
		for _, n := range nodes {
			found = (n == i) || found
		}
		if !found {
			newName = append(newName, word)
		}
	}
	return strings.Join(newName, ".")
}

// Description is auto-generated description, based on output of https://github.com/graphite-project/graphite-web
func (f *aggregateWithWildcards) Description() map[string]types.FunctionDescription {
	return map[string]types.FunctionDescription{
		"aggregateWithWildcards": {
			Name:        "aggregateWithWildcards",
			Function:    "aggregateWithWildcards(seriesList, func, *positions)",
			Description: "Call aggregator after inserting wildcards at the given position(s).\n\nExample:\n\n.. code-block:: none\n\n &target=aggregateWithWildcards(host.cpu-[0-7].cpu-{user,system}.value, 'sum', 1)",
			Module:      "graphite.render.functions",
			Group:       "Calculate",
			Params: []types.FunctionParam{
				{
					Name:     "seriesList",
					Type:     types.SeriesList,
					Required: true,
				},
				{
					Name:     "func",
					Type:     types.AggFunc,
					Required: false,
					Options:  types.StringsToSuggestionList(consolidations.AvailableConsolidationFuncs()),
					Default: &types.Suggestion{
						Value: "average",
						Type:  types.SString,
					},
				},
				{
					Name: "keepStep",
					Type: types.Boolean,
					Default: &types.Suggestion{
						Value: false,
						Type:  types.SBool,
					},
				},
			},
		},
	}
}
