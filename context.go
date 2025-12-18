package synapse

import (
	"context"
)

type contextKey int

const (
	namespaceKey contextKey = iota
	metadataKey
)

// WithNamespace adds a namespace to the context
func WithNamespace(ctx context.Context, namespace string) context.Context {
	return context.WithValue(ctx, namespaceKey, namespace)
}

// GetNamespace retrieves the namespace from the context
func GetNamespace(ctx context.Context) string {
	if ns, ok := ctx.Value(namespaceKey).(string); ok {
		return ns
	}
	return ""
}

// WithMetadata adds metadata to the context
func WithMetadata(ctx context.Context, key string, value any) context.Context {
	md := getMetadata(ctx)
	if md == nil {
		md = make(map[string]any)
	}
	newMd := make(map[string]any, len(md)+1)
	for k, v := range md {
		newMd[k] = v
	}
	newMd[key] = value
	return context.WithValue(ctx, metadataKey, newMd)
}

// GetMetadata retrieves a metadata value from the context
func GetMetadata(ctx context.Context, key string) (any, bool) {
	md := getMetadata(ctx)
	if md == nil {
		return nil, false
	}
	val, ok := md[key]
	return val, ok
}

func getMetadata(ctx context.Context) map[string]any {
	if md, ok := ctx.Value(metadataKey).(map[string]any); ok {
		return md
	}
	return nil
}
