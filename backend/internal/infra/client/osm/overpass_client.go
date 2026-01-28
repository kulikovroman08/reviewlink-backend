package osm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout = 20 * time.Second
	defaultLimit   = 50
	maxLimit       = 200
)

type OverpassClient struct {
	httpClient *http.Client
	baseURLs   []string
}

func NewOverpassClient() *OverpassClient {
	return &OverpassClient{
		httpClient: &http.Client{Timeout: defaultTimeout},
		baseURLs: []string{
			"https://overpass-api.de/api/interpreter",
			"https://overpass.kumi.systems/api/interpreter",
			"https://overpass.nchc.org.tw/api/interpreter",
		},
	}
}

func (c *OverpassClient) SearchPlaces(ctx context.Context, p SearchParams) ([]Place, error) {
	if p.City == "" {
		return nil, fmt.Errorf("city is required")
	}

	limit := p.Limit
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	query := buildOverpassQuery(p, limit)

	form := url.Values{}
	form.Set("data", query)
	body := form.Encode()

	var lastErr error

	for _, baseURL := range c.baseURLs {
		// 1 попытка + 1 retry для перегруженного upstream
		for attempt := 0; attempt < 2; attempt++ {
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, strings.NewReader(body))
			if err != nil {
				return nil, fmt.Errorf("create overpass request: %w", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

			resp, err := c.httpClient.Do(req)
			if err != nil {
				lastErr = fmt.Errorf("overpass request failed (%s): %w", baseURL, err)
				// network error -> попробуем retry/следующий endpoint
				if attempt == 0 {
					sleepBackoff(ctx, 250*time.Millisecond)
					continue
				}
				break
			}

			// читаем тело полностью (нужно и для ошибок, и для json.Unmarshal)
			raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				msg := strings.TrimSpace(string(raw))
				lastErr = fmt.Errorf("overpass status %d (%s): %s", resp.StatusCode, baseURL, msg)

				// перегрузка/лимиты — ретраим или пробуем другой инстанс
				if resp.StatusCode == 429 || resp.StatusCode == 502 || resp.StatusCode == 503 || resp.StatusCode == 504 {
					if attempt == 0 {
						sleepBackoff(ctx, 400*time.Millisecond)
						continue
					}
					break
				}

				// любые другие коды — это не “перегрузка”, возвращаем сразу
				return nil, lastErr
			}

			var data overpassResp
			if err := json.Unmarshal(raw, &data); err != nil {
				return nil, fmt.Errorf("decode overpass response: %w", err)
			}

			out := make([]Place, 0, len(data.Elements))
			for _, el := range data.Elements {
				pl, ok := mapElementToPlace(el)
				if !ok {
					continue
				}
				out = append(out, pl)
			}

			return out, nil
		}
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return nil, fmt.Errorf("overpass request failed: no endpoints available")
}

func sleepBackoff(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-ctx.Done():
		return
	case <-t.C:
		return
	}
}
