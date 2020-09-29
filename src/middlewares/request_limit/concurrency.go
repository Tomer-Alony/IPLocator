package request_limit

import (
	"net/http"
	"time"
)

func ConcurrencyHandler(m RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if m.IsLimitExceeded() {
				http.Error(w, "To many requests", 429)
			}

			token, err := m.Acquire()
			if err == nil {
				next.ServeHTTP(w, r)
			}
			m.Release(token)
		})
	}
}

func NewConcurrencyManager(maxRateLimit, maxTokenDuration int) RateLimiter {
	if maxRateLimit <= 0 {
		return nil
	}

	m := NewManager(&Config{
		Limit:           maxRateLimit,
		TokenResetAfter: time.Duration(maxTokenDuration),
	})

	//m.runResetTokenTask(time.Duration(maxTokenDuration))

	// max concurrency await function
	await := func() {
		go func() {
			for {
				select {
					case <-m.inChan:
						m.tryGenerateToken()
					case t := <-m.releaseChan:
						m.releaseToken(t)
					}
			}
		}()
	}

	await()
	return m
}

