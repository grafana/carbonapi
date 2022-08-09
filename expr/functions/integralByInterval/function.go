package integralByInterval

import (
	"context"
	"math"

	"github.com/grafana/carbonapi/expr/helper"
	"github.com/grafana/carbonapi/expr/interfaces"
	"github.com/grafana/carbonapi/expr/types"
	"github.com/grafana/carbonapi/pkg/parser"
)

type integralByInterval struct {
	interfaces.FunctionBase
}

func GetOrder() interfaces.Order {
	return interfaces.Any
}

func New(configFile string) []interfaces.FunctionMetadata {
	res := make([]interfaces.FunctionMetadata, 0)
	f := &integralByInterval{}
	functions := []string{"integralByInterval"}
	for _, n := range functions {
		res = append(res, interfaces.FunctionMetadata{Name: n, F: f})
	}
	return res
}

// integralByInterval(seriesList, intervalString)
func (f *integralByInterval) Do(ctx context.Context, e parser.Expr, from, until int64, values map[parser.MetricRequest][]*types.MetricData) ([]*types.MetricData, error) {
	args, err := helper.GetSeriesArg(ctx, e.Arg(0), from, until, values)
	if err != nil {
		return nil, err
	}
	if len(args) == 0 {
		return nil, nil
	}

	bucketSizeInt32, err := e.GetIntervalArg(1, 1)
	if err != nil {
		return nil, err
	}
	bucketSize := int64(bucketSizeInt32)
<<<<<<< HEAD
	intervalString, err := e.GetStringArg(1)
	if err != nil {
		return nil, err
	}
=======
	bucketSizeStr := e.Arg(1).StringValue()
>>>>>>> upstream/main

	startTime := from
	results := make([]*types.MetricData, len(args))
	for j, arg := range args {
		current := 0.0
		currentTime := arg.StartTime

<<<<<<< HEAD
		name := fmt.Sprintf("integralByInterval(%s,'%s')", arg.Name, e.Args()[1].StringValue())
		result := arg.CopyLink()
		result.Name = name
		result.PathExpression = name
		result.Values = make([]float64, len(arg.Values))

		result.Tags["integralByInterval"] = intervalString

=======
		name := "integralByInterval(" + arg.Name + ",'" + bucketSizeStr + "')"
		result := &types.MetricData{
			FetchResponse: pb.FetchResponse{
				Name:              name,
				Values:            make([]float64, len(arg.Values)),
				StepTime:          arg.StepTime,
				StartTime:         arg.StartTime,
				StopTime:          arg.StopTime,
				XFilesFactor:      arg.XFilesFactor,
				PathExpression:    name,
				ConsolidationFunc: arg.ConsolidationFunc,
			},
			Tags: arg.Tags,
		}
>>>>>>> upstream/main
		for i, v := range arg.Values {
			if (currentTime-startTime)/bucketSize != (currentTime-startTime-arg.StepTime)/bucketSize {
				current = 0
			}
			if math.IsNaN(v) {
				v = 0
			}
			current += v
			result.Values[i] = current
			currentTime += arg.StepTime
		}

<<<<<<< HEAD
		results = append(results, result)
=======
		results[j] = result
>>>>>>> upstream/main
	}

	return results, nil
}

// Description is auto-generated description, based on output of https://github.com/graphite-project/graphite-web
func (f *integralByInterval) Description() map[string]types.FunctionDescription {
	return map[string]types.FunctionDescription{
		"integralByInterval": {
			Description: "This will do the same as integralByInterval() funcion, except resetting the total to 0 at the given time in the parameter “from” Useful for finding totals per hour/day/week/..",
			Function:    "integralByInterval(seriesList, intervalString)",
			Group:       "Transform",
			Module:      "graphite.render.functions",
			Name:        "integralByInterval",
			Params: []types.FunctionParam{
				{
					Name:     "seriesList",
					Required: true,
					Type:     types.SeriesList,
				}, {
					Name:     "intervalString",
					Required: true,
					Suggestions: types.NewSuggestions(
						"10min",
						"1h",
						"1d",
					),
					Type: types.Interval,
				},
			},
			NameChange:   true, // name changed
			ValuesChange: true, // values changed
		},
	}
}
