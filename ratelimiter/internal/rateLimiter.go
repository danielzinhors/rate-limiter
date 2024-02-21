package ratelimiter

import (
	"context"
	"time"

	"github.com/danielzinhors/rate-limiter/ratelimiter/internal/config"
	helper "github.com/danielzinhors/rate-limiter/ratelimiter/internal/helpers"
)

func CheckRateLimit(ctx context.Context, keyType string, key string, limitConf *config.LimiterConfig, rateConfig *config.RateConfig) (*time.Time, error) {
	if key == "" {
		return nil, nil
	}

	block, err := limitConf.StorageAdapter.GetBlock(ctx, keyType, key)
	if err != nil {
		return nil, err
	}

	if block == nil {
		success, count, err := limitConf.StorageAdapter.IncrementAccesses(ctx, keyType, key, rateConfig.MaxRequestsPerSecond)
		if err != nil {
			return nil, err
		}

		if success {
			helper.PrintfD(limitConf, "%d of %d (%dms if blocked)", keyType, key, count, rateConfig.MaxRequestsPerSecond, rateConfig.BlockTimeMilliseconds)
		} else {
			helper.PrintfD(limitConf, "adding a block of %dms", keyType, key, rateConfig.BlockTimeMilliseconds)
			block, err = limitConf.StorageAdapter.AddBlock(ctx, keyType, key, rateConfig.BlockTimeMilliseconds)
			if err != nil {
				return nil, err
			}
		}
	}

	if block != nil {
		helper.PrintfD(limitConf, "block time %.2f seconds", keyType, key, helper.GetBlockTime(block))
		return block, nil
	}

	return nil, nil
}
