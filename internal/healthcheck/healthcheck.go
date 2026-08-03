// Package healthcheck polls an HTTP endpoint until it reports healthy or a
// deadline passes — used to confirm a deploy actually came up before we
// consider it successful.
package healthcheck

import (
	"fmt"
	"net/http"
	"time"
)

// WaitHealthy polls url (expects HTTP 200, e.g. Spring's /actuator/health)
// every interval until it succeeds or timeout elapses.
func WaitHealthy(url string, timeout, interval time.Duration) error {
	client := &http.Client{Timeout: 5 * time.Second}

	deadline := time.Now().Add(timeout)
	var lastErr error

	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err != nil {
			lastErr = err
		} else {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
			lastErr = fmt.Errorf("status %d", resp.StatusCode)
		}
		time.Sleep(interval)
	}

	return fmt.Errorf("timed out after %s waiting for %s to become healthy: %w", timeout, url, lastErr)
}
