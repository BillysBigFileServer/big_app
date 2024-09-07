package main

import (
	"context"
	"fmt"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/BillysBigFileServer/app/state"
	"github.com/BillysBigFileServer/bfsp-go"
)

func main() {
	ctx := context.Background()
	a := app.New()
	w := a.NewWindow("BBFS")

	masterKey, err := bfsp.CreateMasterEncKey("johnbon")
	if err != nil {
		panic(err)
	}
	ctx = bfsp.ContextWithMasterKey(ctx, masterKey)

	token := "Eu8BCoQBCgI0MAoGcmlnaHRzCgZkZWxldGUKB3BheW1lbnQKD3JlYWRfbWFzdGVyX2tleQoFdXNhZ2UKEHdyaXRlX21hc3Rlcl9rZXkYAyIJCgcIChIDGIAIIi4KLAiBCBInOiUKAhgACgIYAQoCGBsKAxiCCAoDGIMICgMYhAgKAxiFCAoDGIYIEiQIABIghw3MltWZUp_BM8_S1j9ewSJWw0dTYHyBksu06aIJLeMaQJxE5jVYJwJ5DmapX5SigmbtVtKVjC3hOJ50EOrlR0XnmslDyK244423BSHSHNedVvpPSCWQ3bB6Jjd748JY2QkiIgogfvkAKWRSh315yZdcCQmnhYCvjEXKTOq0uR4RnH99Ip4="
	client, err := bfsp.NewHTTPFileServerClient(token, "localhost:9998", false)
	if err != nil {
		panic(err)
	}
	ctx = bfsp.ContextWithClient(ctx, client)

	appState := state.AppState{}
	go func() {
		client := bfsp.ClientFromContext(ctx)
		masterKey := bfsp.MasterKeyFromContext(ctx)
		fileMetas, err := bfsp.ListFileMetadata(client, []string{}, masterKey)
		if err != nil {
			panic(err)
		}

		appState.RwLock.Lock()
		appState.Files = fileMetas
		appState.RwLock.Unlock()

		time.Sleep(10 * time.Second)
	}()
	ctx = state.ContextWithAppState(ctx, &appState)

	w.SetContent(FilesPage(ctx, w))
	w.ShowAndRun()
}

func FilesPage(ctx context.Context, w fyne.Window) fyne.CanvasObject {
	fileMetaList := []*bfsp.FileMetadata{}

	appState := state.FromContext(ctx)
	for {
		appState.RwLock.RLock()
		for _, meta := range appState.Files {
			fileMetaList = append(fileMetaList, meta)
		}
		appState.RwLock.RUnlock()

		// wait until the file metadata has been populated
		if len(fileMetaList) > 0 {
			break
		}

		time.Sleep(1 * time.Second)
	}

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

	fileWidgets := []fyne.CanvasObject{}
	for _, fileMeta := range fileMetaList {
		fileWidgets = append(fileWidgets, FileWidget(ctx, fileMeta, w))
	}

	return container.NewVBox(
		fileWidgets...,
	)
}

func FilePage(ctx context.Context, fileMeta *bfsp.FileMetadata, w fyne.Window) fyne.CanvasObject {
	backArrow, err := fyne.LoadResourceFromPath("resources/back_arrow.png")
	if err != nil {
		panic(err)
	}

	backButton := widget.NewButtonWithIcon("", backArrow, func() {
		w.SetContent(FilesPage(ctx, w))

	})

	fileName := widget.NewLabel(fileMeta.FileName)
	fileName.Wrapping = fyne.TextWrapWord
	fileName.TextStyle.Bold = true

	humanReadableSize := widget.NewLabel("Size: " + HumanSize(fileMeta.FileSize))

	return container.NewVBox(backButton, fileName, humanReadableSize)
}

func FileWidget(ctx context.Context, fileMeta *bfsp.FileMetadata, w fyne.Window) fyne.CanvasObject {
	checkbox := widget.NewCheck("", func(isToggled bool) {
		fmt.Println(isToggled)
	})
	fileName := widget.NewButton(AbridgedFileName(fileMeta.FileName), func() {
		filePage := FilePage(ctx, fileMeta, w)
		w.SetContent(filePage)

	})

	return container.NewHBox(checkbox, fileName)
}
