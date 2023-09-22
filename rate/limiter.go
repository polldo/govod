package rate

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Limiter leverages the std rate package to limit every client request.
// It assumes that the not allowed requests are DROPPED !
// The alternative, to not drop requests, would be to use `Reserve()` instead of `Allow()`.
type Limiter struct {
	Expiry   int     // Expressed in minutes. Clients older than this value will be cleaned.
	Burst    int     // How much tokens can be consumed by a client in one request.
	LimitRPS float64 // Expressed in requests per second.
	clients  map[string]*clientLimiter
	mu       sync.RWMutex
}

// clientLimiter represents the rate limiter for an individual client.
type clientLimiter struct {
	limiter    *rate.Limiter
	lastAccess time.Time
}

// NewLimiter constructs a new rate limiter with the passed options and returns it.
func NewLimiter(burst int, expiry int, limitRPS float64) *Limiter {
	clients := make(map[string]*clientLimiter)
	lm := &Limiter{
		Expiry:   expiry,
		LimitRPS: limitRPS,
		Burst:    burst,
		clients:  clients,
	}
	go lm.refresh()
	return lm
}

// Check updates the limiter for the passed client and checks
// whether the client request can be processed or it has exceeded the limit.
func (l *Limiter) Check(id string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	cl, ok := l.clients[id]
	if !ok {
		l.clients[id] = &clientLimiter{
			limiter:    rate.NewLimiter(rate.Limit(l.LimitRPS), l.Burst),
			lastAccess: time.Now(),
		}
		return l.clients[id].limiter.Allow()
	}
	cl.lastAccess = time.Now()
	return cl.limiter.Allow()
}

// refresh deletes clients older than a threshold (l.Expiry).
// This is neccessary because, as long as the application is running,
// if a cleanup process is not in place the clients map will continue to grow unbounded.
func (l *Limiter) refresh() {
	for {
		time.Sleep(time.Minute)

		l.mu.Lock()
		for id, v := range l.clients {
			if time.Since(v.lastAccess) > time.Duration(l.Expiry)*time.Minute {
				delete(l.clients, id)
			}
		}
		l.mu.Unlock()
	}
}

// Every converts a minimum time interval between events to a requests-per-second value.
func Every(interval time.Duration) float64 {
	return float64(rate.Every(interval))
}
