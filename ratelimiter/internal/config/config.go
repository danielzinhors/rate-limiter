package config

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/danielzinhors/rate-limiter/ratelimiter/internal/adapters"
	"github.com/danielzinhors/rate-limiter/ratelimiter/internal/helpers"
	//response_writers "github.com/danielzinhors/rate-limiter/ratelimiter/internal/response_writers"
)

const envKeyIPMaxRequestsPerSecond = "RATE_LIMITER_IP_MAX_REQUESTS"
const envKeyIPBlockTimeMilliseconds = "RATE_LIMITER_IP_BLOCK_TIME"
const envKeyTokenMaxRequestsPerSecond = "RATE_LIMITER_TOKEN_MAX_REQUESTS"
const envKeyTokenBlockTimeMilliseconds = "RATE_LIMITER_TOKEN_BLOCK_TIME"
const envKeyDebug = "RATE_LIMITER_DEBUG"
const envUseRedis = "RATE_LIMITER_USE_REDIS"
const envRedisAddress = "RATE_LIMITER_REDIS_ADDRESS"
const envRedisPassword = "RATE_LIMITER_REDIS_PASSWORD"
const envRedisDB = "RATE_LIMITER_REDIS_DB"

type RateConfig struct {
	MaxRequestsPerSecond  int64 `json:"maxRequestsPerSecond"`
	BlockTimeMilliseconds int64 `json:"blockTimeMilliseconds"`
}

type LimiterConfig struct {
	IP             *RateConfig                                `json:"ip"`
	Token          *RateConfig                                `json:"token"`
	CustomTokens   *map[string]*RateConfig                    `json:"tokens"`
	StorageAdapter adapters.RateLimitStorageAdapter           `json:"-"`
	ResponseWriter response_writers.RateLimiterResponseWriter `json:"-"`
	Debug          bool                                       `json:"debug"`
	DisableEnvs    bool                                       `json:"disableEnvs"`
}

func (c *LimiterConfig) GetRateLimiterRateConfigForToken(token string) (*RateConfig, bool) {
	customTokenConfig, ok := (*c.CustomTokens)[token]
	if ok {
		return customTokenConfig, true
	} else {
		return c.Token, false
	}
}

func getDefaultConfiguration() *LimiterConfig {
	return &LimiterConfig{
		IP: &RateConfig{
			MaxRequestsPerSecond:  100,
			BlockTimeMilliseconds: 1000,
		},
		Token: &RateConfig{
			MaxRequestsPerSecond:  200,
			BlockTimeMilliseconds: 500,
		},
		CustomTokens:   &map[string]*RateConfig{},
		StorageAdapter: adapters.NewRateLimitMemoryStorageAdapter(),
		ResponseWriter: response_writers.NewRateLimiterDefaultResponseWriter(),
		Debug:          false,
	}
}

func SetConfiguration(config *LimiterConfig) *LimiterConfig {
	defaultConfiguration := getDefaultConfiguration()

	if config == nil {
		config = defaultConfiguration
	}

	if !config.DisableEnvs {
		debug, ok := helpers.GetEnvBoolean(envKeyDebug)
		if ok {
			config.Debug = debug
			helpers.PrintfWD(config, "using env %s", envKeyDebug)
		}
	}

	configureIP(config, defaultConfiguration)
	configureToken(config, defaultConfiguration)
	configureCustomTokens(config, defaultConfiguration)
	configureStorageAdapter(config, defaultConfiguration)
	configureResponseWriter(config, defaultConfiguration)

	if config.Debug {
		jsonConfiguration, err := json.Marshal(config)
		if err == nil {
			helpers.PrintfWD(config, "using configuration: %s", jsonConfiguration)
		}
	}

	return config
}

func configureIP(config *LimiterConfig, defaultConfiguration *LimiterConfig) {
	if config.IP == nil {
		config.IP = defaultConfiguration.IP
	}

	if !config.DisableEnvs {
		mrps, ok := helpers.GetEnvLargeint(envKeyIPMaxRequestsPerSecond)
		if ok {
			config.IP.MaxRequestsPerSecond = mrps
			helpers.PrintfWD(config, "using env %s", envKeyIPMaxRequestsPerSecond)
		}

		bt, ok := helpers.GetEnvLargeint(envKeyIPBlockTimeMilliseconds)
		if ok {
			config.IP.BlockTimeMilliseconds = bt
			helpers.PrintfWD(config, "using env %s", envKeyIPBlockTimeMilliseconds)
		}
	}
}

func configureToken(config *LimiterConfig, defaultConfiguration *LimiterConfig) {
	if config.Token == nil {
		config.Token = defaultConfiguration.Token
	}

	if !config.DisableEnvs {
		mrps, ok := helpers.GetEnvLargeint(envKeyTokenMaxRequestsPerSecond)
		if ok {
			config.Token.MaxRequestsPerSecond = mrps
			helpers.PrintfWD(config, "using env %s", envKeyTokenMaxRequestsPerSecond)
		}

		bt, ok := helpers.GetEnvLargeint(envKeyTokenBlockTimeMilliseconds)
		if ok {
			config.Token.BlockTimeMilliseconds = bt
			helpers.PrintfWD(config, "using env %s", envKeyTokenBlockTimeMilliseconds)
		}
	}
}

func configureCustomTokens(config *LimiterConfig, defaultConfiguration *LimiterConfig) {
	if config.CustomTokens == nil {
		config.CustomTokens = defaultConfiguration.CustomTokens
	}

	for key := range *config.CustomTokens {
		value, ok := (*config.CustomTokens)[key]
		if !ok || value == nil {
			(*config.CustomTokens)[key] = config.Token
		}
	}

	customTokens := getCustomTokenList()
	for _, customToken := range *customTokens {
		configureCustomToken(config, defaultConfiguration, customToken)
	}
}

func getCustomTokenList() *[]string {
	envKeyRegex := regexp.MustCompile("^RATE_LIMITER_TOKEN_(.*)_(MAX_REQUESTS|BLOCK_TIME)$")

	foundTokens := map[string]bool{}

	envs := os.Environ()
	for _, env := range envs {
		envPair := strings.SplitN(env, "=", 2)
		envKey := envPair[0]
		if envKeyRegex.Match([]byte(envKey)) {
			foundTokens[envKeyRegex.FindStringSubmatch(envKey)[1]] = true
		}
	}

	tokens := []string{}
	for k := range foundTokens {
		tokens = append(tokens, k)
	}

	return &tokens
}

func configureCustomToken(config *LimiterConfig, defaultConfiguration *LimiterConfig, customToken string) {

	helpers.PrintfWD(config, "configuring custom token \"%s\"", customToken)

	maxRequestsPerSecondEnvKey := fmt.Sprintf("RATE_LIMITER_TOKEN_%s_MAX_REQUESTS", customToken)
	maxRequestsPerSecond, ok := helpers.GetEnvLargeint(maxRequestsPerSecondEnvKey)
	if !ok {
		defaultValue := config.Token.MaxRequestsPerSecond
		helpers.PrintfWD(config, "env \"%s\" not found: using default value %d", maxRequestsPerSecondEnvKey, defaultValue)
		maxRequestsPerSecond = defaultValue
	}

	blockTimeMillisecondEnvKey := fmt.Sprintf("RATE_LIMITER_TOKEN_%s_BLOCK_TIME", customToken)
	blockTimeMilliseconds, ok := helpers.GetEnvLargeint(blockTimeMillisecondEnvKey)
	if !ok {
		defaultValue := config.Token.BlockTimeMilliseconds
		helpers.PrintfWD(config, "env \"%s\" not found: using default value %d", blockTimeMillisecondEnvKey, defaultValue)
		blockTimeMilliseconds = defaultValue
	}

	(*config.CustomTokens)[customToken] = &RateConfig{
		MaxRequestsPerSecond:  maxRequestsPerSecond,
		BlockTimeMilliseconds: blockTimeMilliseconds,
	}
}

func configureStorageAdapter(config *LimiterConfig, defaultConfiguration *LimiterConfig) {
	if config.StorageAdapter == nil {
		config.StorageAdapter = defaultConfiguration.StorageAdapter
	}

	useRedis, ok := helpers.GetEnvBoolean(envUseRedis)
	if ok && useRedis {
		configureRedisStorageAdapter(config)
	} else if config.StorageAdapter != defaultConfiguration.StorageAdapter {
		helpers.PrintfWD(config, "using StorageAdapter Custom")
	} else {
		helpers.PrintfWD(config, "using StorageAdapter Default")
	}
}

func configureRedisStorageAdapter(config *LimiterConfig) {
	helpers.PrintfWD(config, "using StorageAdapter Redis")

	redisAddress, ok := helpers.GetEnvString(envRedisAddress)
	if !ok {
		panic(fmt.Sprintf("%s env is required when using redis adapter with env configuration", envRedisAddress))
	}

	redisPassword, ok := helpers.GetEnvString(envRedisPassword)
	if !ok {
		redisPassword = ""
	}

	redisDB, ok := helpers.GetEnvLargeint(envRedisDB)
	if !ok {
		redisDB = 0
	}

	config.StorageAdapter = adapters.NewRateLimitRedisStorageAdapter(redisAddress, redisPassword, redisDB)
}

func configureResponseWriter(config *LimiterConfig, defaultConfiguration *LimiterConfig) {
	if config.ResponseWriter == nil {
		config.ResponseWriter = defaultConfiguration.ResponseWriter
	}

	if config.ResponseWriter != defaultConfiguration.ResponseWriter {
		helpers.PrintfWD(config, "using ResponseWriter Custom")
	} else {
		helpers.PrintfWD(config, "using ResponseWriter Default")
	}
}
