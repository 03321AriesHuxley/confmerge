package merger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMerge_OverrideScalar(t *testing.T) {
	dst := map[string]any{"key": "base"}
	src := map[string]any{"key": "override"}

	result, err := Merge(dst, src, StrategyOverride)
	require.NoError(t, err)
	assert.Equal(t, "override", result["key"])
}

func TestMerge_KeepScalar(t *testing.T) {
	dst := map[string]any{"key": "base"}
	src := map[string]any{"key": "override"}

	result, err := Merge(dst, src, StrategyKeep)
	require.NoError(t, err)
	assert.Equal(t, "base", result["key"])
}

func TestMerge_AddsNewKey(t *testing.T) {
	dst := map[string]any{"existing": 1}
	src := map[string]any{"new": 2}

	result, err := Merge(dst, src, StrategyOverride)
	require.NoError(t, err)
	assert.Equal(t, 1, result["existing"])
	assert.Equal(t, 2, result["new"])
}

func TestMerge_DeepNestedMap(t *testing.T) {
	dst := map[string]any{
		"database": map[string]any{
			"host": "localhost",
			"port": 5432,
		},
	}
	src := map[string]any{
		"database": map[string]any{
			"host": "prod.db",
			"name": "mydb",
		},
	}

	result, err := Merge(dst, src, StrategyOverride)
	require.NoError(t, err)

	db := result["database"].(map[string]any)
	assert.Equal(t, "prod.db", db["host"])
	assert.Equal(t, 5432, db["port"])
	assert.Equal(t, "mydb", db["name"])
}

func TestMerge_EmptySrc(t *testing.T) {
	dst := map[string]any{"key": "value"}
	src := map[string]any{}

	result, err := Merge(dst, src, StrategyOverride)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"key": "value"}, result)
}

func TestMerge_EmptyDst(t *testing.T) {
	dst := map[string]any{}
	src := map[string]any{"key": "value"}

	result, err := Merge(dst, src, StrategyOverride)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"key": "value"}, result)
}
