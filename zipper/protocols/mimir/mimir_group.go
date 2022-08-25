package mimir

import (
	"context"
	"github.com/ansel1/merry"

	"github.com/grafana/carbonapi/limiter"
	"github.com/grafana/carbonapi/zipper/helper"
	"github.com/grafana/carbonapi/zipper/metadata"
	"github.com/grafana/carbonapi/zipper/protocols/prometheus"
	"github.com/grafana/carbonapi/zipper/types"

	"go.uber.org/zap"
)

func init() {
	metadata.Metadata.Lock()
	defer metadata.Metadata.Unlock()

	metadata.Metadata.SupportedProtocols["mimir"] = struct{}{}
	metadata.Metadata.ProtocolInits["mimir"] = New
	metadata.Metadata.ProtocolInitsWithLimiter["mimir"] = NewWithLimiter
}

func New(logger *zap.Logger, config types.BackendV2, tldCacheDisabled bool) (types.BackendServer, merry.Error) {
	if config.ConcurrencyLimit == nil {
		return nil, types.ErrConcurrencyLimitNotSet
	}
	if len(config.Servers) == 0 {
		return nil, types.ErrNoServersSpecified
	}
	l := limiter.NewServerLimiter([]string{config.GroupName}, *config.ConcurrencyLimit)

	return NewWithLimiter(logger, config, tldCacheDisabled, l)
}

func NewWithLimiter(logger *zap.Logger, config types.BackendV2, tldCacheDisabled bool, limiter limiter.ServerLimiter) (types.BackendServer, merry.Error) {
	pg, err := prometheus.NewPrometheusGroupWithLimiter(logger, config, tldCacheDisabled, limiter)
	if err != nil {
		logger.Fatal("problem creating prometheus group with limiter", zap.Error(err))
	}

	rawTenantID, ok := config.BackendOptions["tenant_id"]
	if !ok {
		logger.Fatal("missing required option tenant ID")
	}
	tenantID, ok := rawTenantID.(string)
	if !ok {
		logger.Fatal("failed to cast tenant ID as a string")
	}

	pg.QueryExecutor = func(httpQuery *helper.HttpQuery, ctx context.Context, logger *zap.Logger, uri string, r types.Request) (*helper.ServerResponse, merry.Error) {
		headers := map[string]string{
			"X-Scope-OrgID": tenantID,
		}

		return httpQuery.DoQueryWithAdditionalHeaders(ctx, logger, uri, r, headers)
	}

	return pg, nil
}
