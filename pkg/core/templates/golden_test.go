package template

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/golden"
)

func getTestdataDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to get current test file path")
	}

	return filepath.Join(filepath.Dir(filename), "testdata")
}

func assertGolden(t *testing.T, namespace, actual, relPath string) {
	t.Helper()

	goldenPath := filepath.Join(
		getTestdataDir(),
		namespace,
		relPath+".golden",
	)

	if golden.FlagUpdate() {
		require.NoError(
			t,
			os.MkdirAll(filepath.Dir(goldenPath), 0755),
		)

		require.NoError(
			t,
			os.WriteFile(goldenPath, []byte(actual), 0644),
		)

		t.Logf("updated golden: %s", goldenPath)
		return
	}

	expected, err := os.ReadFile(goldenPath)
	require.NoError(
		t,
		err,
		"golden file not found: %s",
		goldenPath,
	)

	assert.Equal(
		t,
		string(expected),
		actual,
		"golden mismatch: %s",
		relPath,
	)
}
