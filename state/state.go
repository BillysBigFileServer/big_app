package state

import (
	"context"
	"sync"

	"github.com/BillysBigFileServer/bfsp-go"
)

type AppState struct {
	Files  map[string]*bfsp.FileMetadata
	RwLock sync.RWMutex
}

type appStateContextKeyType struct{}

var appStateContextKey = appStateContextKeyType{}

func ContextWithAppState(ctx context.Context, appState *AppState) context.Context {
	return context.WithValue(ctx, appStateContextKey, appState)
}

func FromContext(ctx context.Context) *AppState {
	return ctx.Value(appStateContextKey).(*AppState)
}
