package middleware

import (
	"net/http"
	"sync"
)

func Limit(l Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			if l.TryAcquire() {
				next.ServeHTTP(rw, r)
				l.Release()
				return
			}
			rw.WriteHeader(http.StatusTooManyRequests)
		}

		return http.HandlerFunc(fn)
	}
}

type Limiter interface {
	TryAcquire() bool
	Release()
}

type MutexLimiter struct {
	mutex   *sync.Mutex
	current int
	limit   int
}

func NewMutexLimiter(count int) *MutexLimiter {
	var m sync.Mutex
	mutexLimiter := MutexLimiter{
		mutex:   &m,
		current: 0,
		limit:   count,
	}
	return &mutexLimiter
}

func (l *MutexLimiter) TryAcquire() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.current >= l.limit {
		return false
	}

	l.current += 1

	return true
}

func (l *MutexLimiter) Release() {
	l.mutex.Lock()
	l.current -= 1
	l.mutex.Unlock()
}

type ChanLimiter struct {
	limiter chan struct{}
}

func NewChanLimiter(count int) *ChanLimiter {
	chanLimiter := ChanLimiter{limiter: make(chan struct{}, count)}
	return &chanLimiter
}

func (l *ChanLimiter) TryAcquire() bool {
	select {
	case l.limiter <- struct{}{}:
		return true
	default:
		return false
	}
}

func (l *ChanLimiter) Release() {
	_ = <-l.limiter
}
