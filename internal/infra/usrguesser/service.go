package usrguesser

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chains-lab/cities-svc/internal/domain/models"
	"github.com/chains-lab/profiles-svc/resources"
	"github.com/google/uuid"
)

type Service struct {
	baseURL      string
	httpClient   *http.Client
	maxBatchSize int
}

func NewService(baseURL string, httpClient *http.Client) *Service {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 5 * time.Second}
	}
	return &Service{
		baseURL:      strings.TrimRight(baseURL, "/"),
		httpClient:   httpClient,
		maxBatchSize: 100,
	}
}

func (s *Service) Guess(ctx context.Context, userIDs ...uuid.UUID) (map[uuid.UUID]models.Profile, error) {
	ids := dedup(userIDs)
	out := make(map[uuid.UUID]models.Profile, len(ids))
	if len(ids) == 0 {
		return out, nil
	}

	fetchedAll := make(map[uuid.UUID]models.Profile, len(ids))

	for _, part := range chunk(ids, s.maxBatchSize) {
		u, err := url.Parse(s.baseURL)
		if err != nil {
			return nil, fmt.Errorf("bad baseURL: %w", err)
		}

		q := u.Query()
		for _, id := range part {
			q.Add("user_id", id.String())
		}

		u.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}

		resp, err := s.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("profiles resolve request failed: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("profiles resolve read body: %w", err)
		}
		resp.Body.Close()

		if resp.StatusCode/100 != 2 {
			return nil, fmt.Errorf("profiles resolve non-2xx: %s", resp.Status)
		}

		var flat resources.ProfilesCollection
		if err := json.Unmarshal(body, &flat); err == nil && len(flat.Data) > 0 {
			for _, p := range flat.Data {
				if p.Id != uuid.Nil {
					fetchedAll[p.Id] = models.Profile{
						ID:          p.Id,
						Username:    p.Attributes.Username,
						Pseudonym:   p.Attributes.Pseudonym,
						Description: p.Attributes.Description,
						Avatar:      p.Attributes.Avatar,
						Official:    p.Attributes.Official,
						Sex:         p.Attributes.Sex,
						BirthDate:   p.Attributes.BirthDate,
						UpdatedAt:   p.Attributes.UpdatedAt,
						CreatedAt:   p.Attributes.CreatedAt,
					}
				}
			}
			continue
		}

	}

	for _, id := range ids {
		if p, ok := fetchedAll[id]; ok && strings.TrimSpace(p.Username) != "" {
			out[id] = p
		} else {
			out[id] = models.Profile{
				Username: "deleted",
			}
		}
	}

	return out, nil
}

func chunk[T any](in []T, n int) [][]T {
	if n <= 0 || len(in) == 0 {
		return [][]T{in}
	}
	out := make([][]T, 0, (len(in)+n-1)/n)
	for i := 0; i < len(in); i += n {
		j := i + n
		if j > len(in) {
			j = len(in)
		}
		out = append(out, in[i:j])
	}
	return out
}

func dedup(ids []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{}, len(ids))
	out := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
