package reqlog

type ContextHeader string

const (
	XCorrelationIdLogKey  = "x-correlation-id"
	XSessionEntropyLogKey = "x-session-entropy"
	XCorrelationId        = "X-Correlation-Id"
	XSessionEntropy       = "X-Session-Entropy"

	XCorrelationIdKey  = ContextHeader(XCorrelationId)
	XSessionEntropyKey = ContextHeader(XSessionEntropy)
)
