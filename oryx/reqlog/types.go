package reqlog

import "net/http"

type XCorrelationIdContextKey struct{}

var (
	XCorrelationIdLogKey = "x-correlation-id"
	XCorrelationIdKey    = http.CanonicalHeaderKey(XCorrelationIdLogKey)
)
