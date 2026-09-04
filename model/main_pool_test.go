package model

import (
	"os"
	"testing"

	"github.com/QuantumNous/new-api/common"
)

func TestSQLMaxIdleConnsDefaultsToZeroInDatabaseIdleMode(t *testing.T) {
	originalIdleMode := common.DatabaseIdleMode
	t.Cleanup(func() { common.DatabaseIdleMode = originalIdleMode })

	common.DatabaseIdleMode = true
	t.Setenv("SQL_MAX_IDLE_CONNS", "restored-by-testing")
	if err := os.Unsetenv("SQL_MAX_IDLE_CONNS"); err != nil {
		t.Fatal(err)
	}
	if got := sqlMaxIdleConns(); got != 0 {
		t.Fatalf("sqlMaxIdleConns() = %d, want 0", got)
	}
}

func TestSQLMaxIdleConnsHonorsExplicitOverride(t *testing.T) {
	originalIdleMode := common.DatabaseIdleMode
	t.Cleanup(func() { common.DatabaseIdleMode = originalIdleMode })

	common.DatabaseIdleMode = true
	t.Setenv("SQL_MAX_IDLE_CONNS", "3")

	if got := sqlMaxIdleConns(); got != 3 {
		t.Fatalf("sqlMaxIdleConns() = %d, want 3", got)
	}
}
