package state

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/BillysBigFileServer/bfsp-go"
)

type AppState struct {
	Files       map[string]*bfsp.FileMetadata
	RwLock      sync.RWMutex
	Initialized atomic.Bool
}

func InitAppState(ctx context.Context) context.Context {
	appState := AppState{}
	go func() {
		client := bfsp.ClientFromContext(ctx)
		masterKey := bfsp.MasterKeyFromContext(ctx)
		fileMetas, err := bfsp.ListFileMetadata(client, []string{}, masterKey)
		if err != nil {
			panic(err)
		}

		appState.Initialized.Store(true)
		appState.RwLock.Lock()
		appState.Files = fileMetas
		appState.RwLock.Unlock()

		time.Sleep(10 * time.Second)
	}()
	return ContextWithAppState(ctx, &appState)
}

type appStateContextKeyType struct{}

var appStateContextKey = appStateContextKeyType{}

func ContextWithAppState(ctx context.Context, appState *AppState) context.Context {
	return context.WithValue(ctx, appStateContextKey, appState)
}

func FromContext(ctx context.Context) *AppState {
	return ctx.Value(appStateContextKey).(*AppState)
}
