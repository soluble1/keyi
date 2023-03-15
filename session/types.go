package session

import (
	"context"
	"net/http"
)

type Session interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	ID() string
}

// Store 管理session
type Store interface {
	// Generate 生成一个session
	Generate(ctx context.Context, id string) (Session, error)

	// Refresh 刷新一个session
	Refresh(ctx context.Context, id string) error
	// Remove 删除一个session
	Remove(ctx context.Context, id string) error
	// Get 获得session
	Get(ctx context.Context, id string) (Session, error)
}

// Propagator 操作session
type Propagator interface {
	// Inject 将session id注入到http.ResponseWriter中
	Inject(id string, writer http.ResponseWriter) error
	// Extract 将session id从响应中提取出
	Extract(req *http.Request) (string, error)
	// Remove 将session id从http.ResponseWriter中删除
	Remove(writer http.ResponseWriter) error
}
