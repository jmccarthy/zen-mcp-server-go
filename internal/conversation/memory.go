package conversation

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Turn represents a single conversation exchange.
type Turn struct {
	Role    string   `json:"role"`
	Content string   `json:"content"`
	Files   []string `json:"files,omitempty"`
	Tool    string   `json:"tool,omitempty"`
}

// Thread represents a conversation thread.
type Thread struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Turns     []Turn    `json:"turns"`
}

// Store persists threads in Redis.
type Store struct {
	client *redis.Client
	ttl    time.Duration
}

func NewStore(redisURL string, ttl time.Duration) *Store {
	opt, _ := redis.ParseURL(redisURL)
	return &Store{client: redis.NewClient(opt), ttl: ttl}
}

func (s *Store) CreateThread(ctx context.Context) (*Thread, error) {
	id := uuid.NewString()
	th := &Thread{ID: id, CreatedAt: time.Now(), Turns: []Turn{}}
	if err := s.save(ctx, th); err != nil {
		return nil, err
	}
	return th, nil
}

func (s *Store) AddTurn(ctx context.Context, id string, t Turn) error {
	th, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	th.Turns = append(th.Turns, t)
	return s.save(ctx, th)
}

func (s *Store) Get(ctx context.Context, id string) (*Thread, error) {
	data, err := s.client.Get(ctx, key(id)).Bytes()
	if err != nil {
		return nil, err
	}
	var th Thread
	if err := json.Unmarshal(data, &th); err != nil {
		return nil, err
	}
	return &th, nil
}

func (s *Store) save(ctx context.Context, th *Thread) error {
	data, err := json.Marshal(th)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key(th.ID), data, s.ttl).Err()
}

func key(id string) string { return "thread:" + id }
