package middlewares

import (
	"net"
	"net/http"
	"sync"

	"github.com/go-fuego/fuego"
	"golang.org/x/time/rate"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		next.ServeHTTP(w, r)
	})
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var limiterMu sync.Mutex
		limiters := make(map[string]*rate.Limiter)

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		limiterMu.Lock()
		l, ok := limiters[ip]
		if !ok {
			l = rate.NewLimiter(5, 15)
			limiters[ip] = l
		}
		limiterMu.Unlock()

		if !l.Allow() {
			w.Header().Set("Retry-After", "1")
			fuego.SendError(w, r, fuego.HTTPError{
				Title:  "To many request",
				Status: http.StatusTooManyRequests,
				Detail: "you have exceeded the rate limit",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
