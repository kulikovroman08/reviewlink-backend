package news

import (
	"context"
	"time"

	newsclient "github.com/kulikovroman08/reviewlink-backend/internal/infra/client/news"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
)

type PublicNewsClient interface {
	ListNews(ctx context.Context, limit int) ([]newsclient.Item, *time.Time, error)
}

type Service struct {
	newsClient PublicNewsClient
}

func NewNewsService(newsClient PublicNewsClient) *Service {
	return &Service{
		newsClient: newsClient,
	}
}

func (s *Service) ListNews(ctx context.Context, limit int) ([]model.NewsItem, *time.Time, error) {
	items, cachedAt, err := s.newsClient.ListNews(ctx, limit)
	if err != nil {
		return nil, nil, err
	}

	out := make([]model.NewsItem, 0, len(items))
	for _, it := range items {
		out = append(out, model.NewsItem{
			Title:       it.Title,
			Link:        it.Link,
			PublishedAt: it.PublishedAt,
			Source:      it.Source,
		})
	}

	return out, cachedAt, nil
}
