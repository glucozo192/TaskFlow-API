package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/glu-project/internal/taks/models"
	elastic_client "github.com/glu-project/pkg/elasticsearch_client"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const indexName = "tasks"

type TaskESRepository struct{}

// IndexTask indexes a task document into Elasticsearch.
// It upserts the document using the task ID as the document ID.
func (r *TaskESRepository) IndexTask(ctx context.Context, client *elastic_client.ElasticClient, task *models.Task) error {
	body := map[string]any{
		"id":          task.ID,
		"title":       task.Title.String,
		"description": task.Description.String,
		"status":      task.Status.String,
		"user_id":     task.UserID.String,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: task.ID,
		Body:       bytes.NewReader(b),
		Refresh:    "true",
	}
	res, err := req.Do(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("esapi.IndexRequest.Do: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch index error [%d]", res.StatusCode)
	}
	return nil
}

// SearchTasks performs a multi-match full-text search on title and description fields.
func (r *TaskESRepository) SearchTasks(ctx context.Context, client *elastic_client.ElasticClient, q string) ([]*models.Task, error) {
	query := map[string]any{
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":  q,
				"fields": []string{"title", "description"},
			},
		},
	}

	b, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  bytes.NewReader(b),
	}
	res, err := req.Do(ctx, client.Client)
	if err != nil {
		return nil, fmt.Errorf("esapi.SearchRequest.Do: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch search error [%d]", res.StatusCode)
	}

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source struct {
					ID          string `json:"id"`
					Title       string `json:"title"`
					Description string `json:"description"`
					Status      string `json:"status"`
					UserID      string `json:"user_id"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	tasks := make([]*models.Task, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		s := hit.Source
		t := &models.Task{ID: s.ID}
		t.Title.String = s.Title
		t.Title.Valid = true
		t.Description.String = s.Description
		t.Description.Valid = true
		t.Status.String = s.Status
		t.Status.Valid = true
		t.UserID.String = s.UserID
		t.UserID.Valid = true
		tasks = append(tasks, t)
	}
	return tasks, nil
}
