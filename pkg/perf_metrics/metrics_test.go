package perfmetrics

import (
	"sync"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/require"
)

func TestRecordSkipsCollectionInDatabaseIdleMode(t *testing.T) {
	originalIdleMode := common.DatabaseIdleMode
	t.Cleanup(func() {
		common.DatabaseIdleMode = originalIdleMode
		hotBuckets = sync.Map{}
	})

	common.DatabaseIdleMode = true
	hotBuckets = sync.Map{}

	Record(Sample{Model: "test-model", Group: "default", Success: true})

	count := 0
	hotBuckets.Range(func(_, _ any) bool {
		count++
		return true
	})
	require.Zero(t, count)
}
