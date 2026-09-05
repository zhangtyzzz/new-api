package service

import (
	"context"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"

	"github.com/bytedance/gopkg/util/gopool"
)

var (
	databaseIdleMaintenanceOnce sync.Once
	databaseActivityWakeup      = make(chan struct{}, 1)
)

// StartDatabaseIdleMaintenance replaces independent maintenance timers with a
// worker that runs due work only after real application traffic has already
// woken the database. This preserves business maintenance without preventing a
// serverless database from scaling to zero between requests.
func StartDatabaseIdleMaintenance() {
	databaseIdleMaintenanceOnce.Do(func() {
		if !common.IsMasterNode || !common.DatabaseIdleMode {
			return
		}

		gopool.Go(func() {
			logger.LogInfo(context.Background(), "database maintenance is wake-driven in database idle mode")
			var lastSubscriptionReset time.Time
			var lastCodexRefresh time.Time
			var lastAuthCleanup time.Time

			for range databaseActivityWakeup {
				now := time.Now()
				if now.Sub(lastSubscriptionReset) >= subscriptionResetTickInterval {
					lastSubscriptionReset = now
					runSubscriptionQuotaResetOnce()
				}
				if now.Sub(lastCodexRefresh) >= codexCredentialRefreshTickInterval {
					lastCodexRefresh = now
					runCodexCredentialAutoRefreshOnce()
				}
				if now.Sub(lastAuthCleanup) >= authArtifactCleanupInterval {
					lastAuthCleanup = now
					cleanupAuthArtifacts()
				}
			}
		})
	})
}

// NotifyDatabaseActivity is non-blocking. Coalescing is intentional: one pass
// performs all due maintenance and the system-task runner drains runnable work.
func NotifyDatabaseActivity() {
	if !common.IsMasterNode || !common.DatabaseIdleMode {
		return
	}
	notifySystemTaskRunner()
	select {
	case databaseActivityWakeup <- struct{}{}:
	default:
	}
}
