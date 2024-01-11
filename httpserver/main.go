package httpserver

import (
	"context"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// The HttpServer interface...
type HttpServer interface {
	Serve(ctx context.Context) error
}
