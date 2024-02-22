package ratelimiter

import (
	"io"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (suite *UtilsTestSuite) TestDebugPrintf_DebugTrue() {
	rateLimiterConfig := &LimiterConfig{
		Debug: true,
	}
	keyType := "IP"
	keyValue := "127.0.0.1"
	message := "Test message %d"
	messageParam := 1
	outputRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[RATE LIMITER]\[IP\]\[127\.0\.0\.1\] Test message 1\n$`)

	output, err := captureOutput(func() error {
		_, err := PrintfD(rateLimiterConfig, message, keyType, keyValue, messageParam)
		return err
	})

	assert.Nil(suite.T(), err)
	assert.Regexp(suite.T(), outputRegex, output)
}

func (suite *UtilsTestSuite) TestDebugPrintf_DebugFalse() {
	rateLimiterConfig := &LimiterConfig{
		Debug: false,
	}
	keyType := "IP"
	keyValue := "127.0.0.1"
	message := "Test message %d"
	messageParam := 1

	output, err := captureOutput(func() error {
		_, err := PrintfD(rateLimiterConfig, message, keyType, keyValue, messageParam)
		return err
	})

	assert.Nil(suite.T(), err)
	assert.Empty(suite.T(), output)
}

func (suite *UtilsTestSuite) TestDebugPrinDKey_DebugTrue() {
	rateLimiterConfig := &LimiterConfig{
		Debug: true,
	}
	message := "Test message %d"
	messageParam := 1
	outputRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[RATE LIMITER] Test message 1\n$`)

	output, err := captureOutput(func() error {
		_, err := PrintfWD(rateLimiterConfig, message, messageParam)
		return err
	})

	assert.Nil(suite.T(), err)
	assert.Regexp(suite.T(), outputRegex, output)
}

func (s *UtilsTestSuite) TestDebugPrintfWithoutKey_DebugFalse() {
	rateLimiterConfig := &LimiterConfig{
		Debug: false,
	}
	message := "Test message %d"
	messageParam := 1

	output, err := captureOutput(func() error {
		_, err := PrintfWD(rateLimiterConfig, message, messageParam)
		return err
	})

	assert.Nil(s.T(), err)
	assert.Empty(s.T(), output)
}

func (s *UtilsTestSuite) TestGetRemainingBlockTime() {
	block := time.Now().Add(time.Second * 6)
	diff := GetBlockTime(&block)
	assert.Equal(s.T(), 5, int(diff))
}

func (s *UtilsTestSuite) TestGetStringEnv() {
	os.Setenv("MY_ENV", "ENV_VALUE")
	value, ok := GetEnvString("MY_ENV")
	assert.True(s.T(), ok)
	assert.Equal(s.T(), "ENV_VALUE", value)
}

func (s *UtilsTestSuite) TestGetStringEnv_NoValue() {
	os.Setenv("MY_ENV", "ENV_VALUE")
	value, ok := GetEnvString("MY_ANOTHER_ENV")
	assert.False(s.T(), ok)
	assert.Equal(s.T(), "", value)
}

func (s *UtilsTestSuite) TestGetStringEnv_EmptyValue() {
	os.Setenv("MY_ENV", "")
	value, ok := GetEnvString("MY_ENV")
	assert.False(s.T(), ok)
	assert.Equal(s.T(), "", value)
}

func (s *UtilsTestSuite) TestGetBoolEnv_TrueValue() {
	os.Setenv("MY_ENV", "true")
	value, ok := GetEnvBoolean("MY_ENV")
	assert.True(s.T(), ok)
	assert.Equal(s.T(), true, value)
}

func (s *UtilsTestSuite) TestGetBoolEnv_FalseValue() {
	os.Setenv("MY_ENV", "false")
	value, ok := GetEnvBoolean("MY_ENV")
	assert.True(s.T(), ok)
	assert.Equal(s.T(), false, value)
}

func (s *UtilsTestSuite) TestGetBoolEnv_NoValue() {
	os.Setenv("MY_ENV", "false")
	value, ok := GetEnvBoolean("MY_ANOTHER_ENV")
	assert.False(s.T(), ok)
	assert.Equal(s.T(), false, value)
}

func (s *UtilsTestSuite) TestGetBoolEnv_EmptyValue() {
	os.Setenv("MY_ENV", "")
	value, ok := GetEnvBoolean("MY_ENV")
	assert.False(s.T(), ok)
	assert.Equal(s.T(), false, value)
}

func (s *UtilsTestSuite) TestGetBoolEnv_InvalidValue() {
	os.Setenv("MY_ENV", "NOT_A_BOOL")
	value, ok := GetEnvBoolean("MY_ENV")
	assert.False(s.T(), ok)
	assert.Equal(s.T(), false, value)
}

func (s *UtilsTestSuite) TestGetInt64Env() {
	os.Setenv("MY_ENV", "567")
	value, ok := GetEnvLargeint("MY_ENV")
	assert.True(s.T(), ok)
	assert.Equal(s.T(), int64(567), value)
}

func (s *UtilsTestSuite) TestGetInt64Env_NoValue() {
	os.Setenv("MY_ENV", "567")
	value, ok := GetEnvLargeint("MY_ANOTHER_ENV")
	assert.False(s.T(), ok)
	assert.Equal(s.T(), int64(0), value)
}

func (s *UtilsTestSuite) TestGetInt64Env_EmptyValue() {
	os.Setenv("MY_ENV", "")
	value, ok := GetEnvLargeint("MY_ENV")
	assert.False(s.T(), ok)
	assert.Equal(s.T(), int64(0), value)
}

func (s *UtilsTestSuite) TestGetInt64Env_InvalidValue() {
	os.Setenv("MY_ENV", "NOT_A_INT64")
	value, ok := GetEnvLargeint("MY_ENV")
	assert.False(s.T(), ok)
	assert.Equal(s.T(), int64(0), value)
}

func captureOutput(f func() error) (string, error) {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := f()
	os.Stdout = orig
	w.Close()
	out, _ := io.ReadAll(r)
	return string(out), err
}
