package helper

import (
	"context"
	"slices"
	"sort"
	"time"

	"github.com/BillysBigFileServer/app/state"
	"github.com/BillysBigFileServer/bfsp-go"
)

func GetDirList(ctx context.Context, appState *state.AppState, directory []string, update bool) []string {
	if update {
		client := bfsp.ClientFromContext(ctx)
		masterKey := bfsp.MasterKeyFromContext(ctx)
		appState.UpdateAppState(client, masterKey)
	}

	// we use a map to get a unique list of directories
	directories := map[string]bool{}

	for {
		appState.RwLock.RLock()
		for _, meta := range appState.Files {
			if IsSubdirectory(directory, meta.Directory) {
				directories[SliceToDirectory(meta.Directory)] = true
			}
		}
		appState.RwLock.RUnlock()

		// wait until the file metadata has been populated
		if appState.Initialized.Load() {
			break
		}

		time.Sleep(1 * time.Second)
	}

	directoriesList := []string{}
	for dir, _ := range directories {
		directoriesList = append(directoriesList, dir)
	}

	sort.Strings(directoriesList)
	return directoriesList
}

func GetFileList(ctx context.Context, appState *state.AppState, directory []string, update bool) []*bfsp.FileMetadata {
	if update {
		client := bfsp.ClientFromContext(ctx)
		masterKey := bfsp.MasterKeyFromContext(ctx)
		appState.UpdateAppState(client, masterKey)
	}

	fileMetaList := []*bfsp.FileMetadata{}

	for {
		if !appState.Initialized.Load() {
			time.Sleep(1 * time.Second)
		}

		break
	}

	appState.RwLock.RLock()
	for _, meta := range appState.Files {
		if IsDirEqual(directory, meta.Directory) {
			fileMetaList = append(fileMetaList, meta)
		}
	}
	appState.RwLock.RUnlock()

	// sort the files alphabetically
	slices.SortStableFunc(fileMetaList, func(a, b *bfsp.FileMetadata) int {
		switch {
		case a.FileName == b.FileName:
			return 0
		case a.FileName < b.FileName:
			return -1
		default:
			return 1
		}
	})
	return fileMetaList
}
