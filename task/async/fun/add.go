package fun

import (
	"context"
)

// AddEval - request http get request
func (m *defaultFuncWrap) AddEval(ctx context.Context, a, b int64) (int64, error) {
	return a + b, nil
}
