package ratelimiter

import (
	"context"
	"time"
)

func CheckRateLimit(ctx context.Context, keyType string, key string, limitConf *LimiterConfig, rateConfig *RateConfig) (*time.Time, error) {
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
			PrintfD(limitConf, "%d of %d (%dms if blocked)", keyType, key, count, rateConfig.MaxRequestsPerSecond, rateConfig.BlockTimeMilliseconds)
		} else {
			PrintfD(limitConf, "adding a block of %dms", keyType, key, rateConfig.BlockTimeMilliseconds)
			block, err = limitConf.StorageAdapter.AddBlock(ctx, keyType, key, rateConfig.BlockTimeMilliseconds)
			if err != nil {
				return nil, err
			}
		}
	}

	if block != nil {
		PrintfD(limitConf, "block time %.2f seconds", keyType, key, GetBlockTime(block))
		return block, nil
	}

	return nil, nil
}
