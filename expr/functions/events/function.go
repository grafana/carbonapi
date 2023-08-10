package events

import (
	"context"

	"github.com/go-graphite/carbonapi/expr/interfaces"
	"github.com/go-graphite/carbonapi/expr/types"
	"github.com/go-graphite/carbonapi/pkg/parser"
)

type events struct {
	interfaces.FunctionBase
}

func GetOrder() interfaces.Order {
	return interfaces.Any
}

func New(_ string) []interfaces.FunctionMetadata {
	return []interfaces.FunctionMetadata{
		{
			Name: "events",
			F:    &events{},
		},
	}
}

func (f *events) Do(ctx context.Context, e parser.Expr, from, until int64, values map[parser.MetricRequest][]*types.MetricData) ([]*types.MetricData, error) {
	tags, err := e.GetStringArg(0)
	if err != nil {
		return nil, err
	}

	fetchTarget := parser.NewEventTagsExpr(tags)
	targetValues, err := f.GetEvaluator().Fetch(ctx, []parser.Expr{fetchTarget}, from, until, values)
	if err != nil {
		return nil, err
	}

	return f.GetEvaluator().Eval(ctx, fetchTarget, from, until, targetValues)
}

func (f *events) Description() map[string]types.FunctionDescription {
	return map[string]types.FunctionDescription{
		"events": {
			Description: "Returns the number of events at this point in time. Usable with drawAsInfinite.",
			Function:    "events(*tags)",
			Group:       "Special",
			Module:      "graphite.render.functions",
			Name:        "events",
			Params: []types.FunctionParam{
				{
					Name:     "tags",
					Required: false,
					Type:     types.String,
				},
			},
			ValuesChange: true,
		},
	}
}
