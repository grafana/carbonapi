package exponentialMovingAverage

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/lomik/zapwriter"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/grafana/carbonapi/expr/helper"
	"github.com/grafana/carbonapi/expr/interfaces"
	"github.com/grafana/carbonapi/expr/types"
	"github.com/grafana/carbonapi/pkg/parser"
)

type exponentialMovingAverage struct {
	interfaces.FunctionBase

	config movingConfig
}

func GetOrder() interfaces.Order {
	return interfaces.Any
}

type movingConfig struct {
	ReturnNaNsIfStepMismatch *bool
}

func New(configFile string) []interfaces.FunctionMetadata {
	logger := zapwriter.Logger("functionInit").With(zap.String("function", "moving"))
	res := make([]interfaces.FunctionMetadata, 0)
	f := &exponentialMovingAverage{}
	functions := []string{"exponentialMovingAverage"}
	for _, n := range functions {
		res = append(res, interfaces.FunctionMetadata{Name: n, F: f})
	}
	cfg := movingConfig{}
	v := viper.New()
	v.SetConfigFile(configFile)
	err := v.ReadInConfig()
	if err != nil {
		logger.Info("failed to read config file, using default",
			zap.Error(err),
		)
	} else {
		err = v.Unmarshal(&cfg)
		if err != nil {
			logger.Fatal("failed to parse config",
				zap.Error(err),
			)
			return nil
		}
		f.config = cfg
	}

	if cfg.ReturnNaNsIfStepMismatch == nil {
		v := true
		f.config.ReturnNaNsIfStepMismatch = &v
	}
	return res
}

func (f *exponentialMovingAverage) Do(ctx context.Context, e parser.Expr, from, until int64, values map[parser.MetricRequest][]*types.MetricData) ([]*types.MetricData, error) {
	var n int
	var err error
	var constant float64
	var scaleByStep bool

	var argstr string

	if len(e.Args()) < 2 {
		return nil, parser.ErrMissingArgument
	}

	switch e.Args()[1].Type() {
	case parser.EtConst:
		// In this case, zipper does not request additional retrospective points,
		// and leading `n` values, that used to calculate window, become NaN
		n, err = e.GetIntArg(1)
		argstr = strconv.Itoa(n)
	case parser.EtString:
		var n32 int32
		n32, err = e.GetIntervalArg(1, 1)
		argstr = fmt.Sprintf("%q", e.Args()[1].StringValue())
		n = int(n32)
		scaleByStep = true
	default:
		err = parser.ErrBadType
	}
	if err != nil {
		return nil, err
	}

	windowSize := n

	start := from
	if scaleByStep {
		start -= int64(n)
	}

	arg, err := helper.GetSeriesArg(ctx, e.Args()[0], start, until, values)
	if err != nil {
		return nil, err
	}

	var results []*types.MetricData

	if len(arg) == 0 {
		return results, nil
	}

	var offset int

	if scaleByStep {
		windowSize /= int(arg[0].StepTime)
		offset = windowSize
	}

	constant = float64(2 / (float64(windowSize) + 1))

	for _, a := range arg {
		r := a.CopyLink()
		r.Name = fmt.Sprintf("%s(%s,%s)", e.Target(), a.Name, argstr)

		if windowSize == 0 {
			if *f.config.ReturnNaNsIfStepMismatch {
				r.Values = make([]float64, len(a.Values))
				for i := range a.Values {
					r.Values[i] = math.NaN()
				}
			}
			results = append(results, r)
			continue
		}

		r.Values = make([]float64, len(a.Values))
		r.StartTime = (from + r.StepTime - 1) / r.StepTime * r.StepTime // align StartTime to closest >= StepTime
		r.StopTime = r.StartTime + int64(len(r.Values))*r.StepTime

		w := types.NewExpMovingAverage(windowSize, float64(constant))
		for i, v := range a.Values {
			if ridx := i - offset; ridx >= 0 {
				if math.IsNaN(v) {
					r.Values[i] = math.NaN()
					continue
				}

				r.Values[i] = w.Mean()

				if i < windowSize || math.IsNaN(r.Values[ridx]) {
					r.Values[ridx] = math.NaN()
				}
			}
			w.Push(v)
		}
		results = append(results, r)
	}
	return results, nil
}

func (f *exponentialMovingAverage) Description() map[string]types.FunctionDescription {
	return map[string]types.FunctionDescription{
		"exponentialMovingAverage": {
			Description: "Takes a series of values and a window size and produces an exponential moving average utilizing the following formula:\n\n ema(current) = constant * (Current Value) + (1 - constant) * ema(previous)\n The Constant is calculated as:\n constant = 2 / (windowSize + 1) \n The first period EMA uses a simple moving average for its value.\n Example:\n\n code-block:: none\n\n  &target=exponentialMovingAverage(*.transactions.count, 10) \n\n &target=exponentialMovingAverage(*.transactions.count, '-10s')",
			Function:    "exponentialMovingAverage(seriesList, windowSize)",
			Group:       "Calculate",
			Module:      "graphite.render.functions",
			Name:        "exponentialMovingAverage",
			Params: []types.FunctionParam{
				{
					Name:     "seriesList",
					Required: true,
					Type:     types.SeriesList,
				},
				{
					Name:     "windowSize",
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
