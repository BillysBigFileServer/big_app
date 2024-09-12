package pages

import (
	"context"
	"fmt"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/BillysBigFileServer/app/helper"
	"github.com/BillysBigFileServer/app/state"
	"github.com/BillysBigFileServer/bfsp-go"
)

func FilesPage(ctx context.Context, w fyne.Window) fyne.CanvasObject {
	ctx = state.InitAppState(ctx)
	fileMetaList := []*bfsp.FileMetadata{}

	appState := state.FromContext(ctx)
	for {
		appState.RwLock.RLock()
		for _, meta := range appState.Files {
			fileMetaList = append(fileMetaList, meta)
		}
		appState.RwLock.RUnlock()

		// wait until the file metadata has been populated
		if appState.Initialized.Load() {
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
		fileWidgets = append(fileWidgets, minimalFileWidget(ctx, fileMeta, w))
	}

	return container.NewVScroll(container.NewVBox(
		fileWidgets...,
	))
}

func minimalFileWidget(ctx context.Context, fileMeta *bfsp.FileMetadata, w fyne.Window) fyne.CanvasObject {
	checkbox := widget.NewCheck("", func(isToggled bool) {
		fmt.Println(isToggled)
	})
	fileName := widget.NewButton(helper.AbridgedFileName(fileMeta.FileName), func() {
		filePage := FilePage(ctx, fileMeta, w)
		w.SetContent(filePage)

	})

	return container.NewHBox(checkbox, fileName)
}
