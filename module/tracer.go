package module

import "go.opentelemetry.io/otel"

var (
	tracer = otel.Tracer("github.com/datachainlab/besu-ibc-relay-prover/module")
)
