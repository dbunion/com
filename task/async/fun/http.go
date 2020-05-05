package fun

import (
	"context"
	"github.com/zssky/tc/http"
	"time"
)

// HTTPGet - request http get request
func (m *defaultFuncWrap) HTTPGet(ctx context.Context, url string, deadline, dialTimeout int64) (string, error) {
	results, err := m.wrangler(ctx, func(ctx context.Context) ([]string, error) {
		data, _, err := http.SimpleGet(url, time.Duration(deadline), time.Duration(dialTimeout))
		if err != nil {
			return nil, err
		}
		return []string{string(data[0])}, nil
	})
	if err != nil {
		return "", err
	}

	return results[0], nil
}
