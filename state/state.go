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

func (state *AppState) UpdateAppState(client bfsp.FileServerClient, masterKey bfsp.MasterKey) {
	fileMetas, err := bfsp.ListFileMetadata(client, []string{}, masterKey)
	if err != nil {
		panic(err)
	}

	state.Initialized.Store(true)
	state.RwLock.Lock()
	state.Files = fileMetas
	state.RwLock.Unlock()
}

func InitAppState(ctx context.Context) context.Context {
	appState := &AppState{}
	go func() {
		client := bfsp.ClientFromContext(ctx)
		masterKey := bfsp.MasterKeyFromContext(ctx)
		appState.UpdateAppState(client, masterKey)

		time.Sleep(10 * time.Second)
	}()
	return ContextWithAppState(ctx, appState)
}

type appStateContextKeyType struct{}

var appStateContextKey = appStateContextKeyType{}

func ContextWithAppState(ctx context.Context, appState *AppState) context.Context {
	return context.WithValue(ctx, appStateContextKey, appState)
}

func FromContext(ctx context.Context) *AppState {
	val := ctx.Value(appStateContextKey)
	switch val {
	case nil:
		return nil
	default:
		return val.(*AppState)
	}
}
