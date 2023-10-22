package httpclient

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/easy-techno-lab/proton/logger"
)

// The RoundTripper type is an adapter to allow the use of ordinary functions as HTTP round trippers.
// If f is a function with the appropriate signature, Func(f) is a RoundTripper that calls f.
type RoundTripper func(*http.Request) (*http.Response, error)

// RoundTrip calls f(r).
func (f RoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// RoundTripperSequencer chains http.RoundTrippers in a chain.
func RoundTripperSequencer(baseRoundTripper http.RoundTripper, rts ...func(http.RoundTripper) http.RoundTripper) http.RoundTripper {
	for _, f := range rts {
		baseRoundTripper = f(baseRoundTripper)
	}
	return baseRoundTripper
}

// Timer measures the time taken by http.RoundTripper.
func Timer(logLevel logger.Level) func(http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripper(func(r *http.Request) (*http.Response, error) {
			if logger.InLevel(logLevel) {
				defer func(start time.Time) {
					logLevel.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
				}(time.Now())
			}
			return next.RoundTrip(r)
		})
	}
}

// PanicCatcher handles panics in http.RoundTripper.
func PanicCatcher(next http.RoundTripper) http.RoundTripper {
	return RoundTripper(func(r *http.Request) (*http.Response, error) {
		defer func() {
			if rec := recover(); rec != nil {
				if logger.InLevel(logger.LevelError) {
					logger.Errorf("%s\n%s", rec, debug.Stack())
				}
			}
		}()
		return next.RoundTrip(r)
	})
}

// DumpHttp dumps the HTTP request and response, and prints out with logFunc.
func DumpHttp(logLevel logger.Level) func(http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripper(func(r *http.Request) (*http.Response, error) {
			if logger.InLevel(logLevel) {
				logger.DumpHttpRequest(r, logLevel)

				resp, err := next.RoundTrip(r)
				if err != nil {
					return nil, err
				}

				logger.DumpHttpResponse(resp, logLevel)

				return resp, nil
			}
			return next.RoundTrip(r)
		})
	}
}
