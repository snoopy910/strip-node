package sequencer

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/go-pg/pg/v10"
)

// Tests the database connection pool configuration and behavior under various scenarios
func TestDBConnectionPool(t *testing.T) {
	// Verifies that the connection pool maintains the correct number of connections
	// - Checks if idle connections meet minimum pool size
	// - Ensures total connections don't exceed maximum pool size
	t.Run("Pool Size Constraints", func(t *testing.T) {
		InitialiseDB("localhost:5432", "test_db", "test_user", "test_password")

		// Initial pool stats check
		stats := client.PoolStats()
		if stats.IdleConns < minPoolSize {
			t.Errorf("Expected minimum %d idle connections, got %d", minPoolSize, stats.IdleConns)
		}
		if stats.TotalConns > maxPoolSize {
			t.Errorf("Total connections %d exceeds maximum pool size %d", stats.TotalConns, maxPoolSize)
		}
	})

	// Validates that connections are properly returned to the pool after use
	// - Executes 100 sequential queries
	// - Verifies that connections are released back to idle state
	// - Checks idle connection count doesn't decrease after operations
	t.Run("Connection Release", func(t *testing.T) {
		// Get initial pool stats
		initialStats := client.PoolStats()

		// Execute a batch of queries
		const numQueries = 100
		for i := 0; i < numQueries; i++ {
			_, err := client.Exec("SELECT 1")
			if err != nil {
				t.Fatalf("Query failed: %v", err)
			}
		}

		// Get stats after queries
		afterStats := client.PoolStats()

		// Verify connections are properly released back to the pool
		if afterStats.IdleConns < initialStats.IdleConns {
			t.Errorf("Expected connections to be released back to pool. Initial idle: %d, After idle: %d",
				initialStats.IdleConns, afterStats.IdleConns)
		}
	})

	// Tests connection pool behavior under concurrent load
	// - Spawns 10 concurrent workers
	// - Each worker executes 50 queries (mix of fast and slow)
	// - Verifies pool remains stable under concurrent access
	// - Ensures connection limits are respected during high load
	t.Run("Pool Under Load", func(t *testing.T) {
		// Create a workload that simulates real usage patterns
		const (
			numWorkers = 10
			numQueries = 50
		)

		var wg sync.WaitGroup
		errors := make(chan error, numWorkers*numQueries)

		// Launch workers
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < numQueries; i++ {
					// Mix of short and long queries
					query := "SELECT 1"
					if i%5 == 0 {
						query = "SELECT pg_sleep(0.1)"
					}
					_, err := client.Exec(query)
					if err != nil {
						errors <- err
					}
				}
			}()
		}

		// Wait for all workers to complete
		wg.Wait()
		close(errors)

		// Check for errors
		for err := range errors {
			t.Errorf("Error during pool stress test: %v", err)
		}

		// Verify pool stats are still within bounds
		finalStats := client.PoolStats()
		if finalStats.TotalConns > maxPoolSize {
			t.Errorf("Pool size exceeded maximum: got %d, want <= %d",
				finalStats.TotalConns, maxPoolSize)
		}
	})

	// Tests behavior when connection pool is exhausted
	// - Manually holds all available connections from the pool
	// - Attempts to acquire an additional connection
	// - Verifies that the request times out appropriately
	// - Ensures proper error handling when pool is full
	t.Run("Pool Exhaustion", func(t *testing.T) {
		// First, occupy all connections in the pool
		var conns []*pg.Conn
		for i := 0; i < maxPoolSize; i++ {
			conn := client.Conn()
			if err := conn.Ping(context.Background()); err != nil {
				t.Fatalf("Failed to ping connection: %v", err)
			}
			conns = append(conns, conn)
		}
		defer func() {
			for _, conn := range conns {
				_ = conn.Close()
			}
		}()

		// Now try to get one more connection with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		_, err := client.ExecContext(ctx, "SELECT 1")
		if err == nil {
			t.Error("Expected error due to pool exhaustion, got none")
		}
	})
}
