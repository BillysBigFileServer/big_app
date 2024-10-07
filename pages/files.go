package pages

import (
	"context"
	"fmt"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/BillysBigFileServer/app/helper"
	"github.com/BillysBigFileServer/app/state"
	"github.com/BillysBigFileServer/bfsp-go"
)

func FilesPage(ctx context.Context, w fyne.Window, updateFiles bool) fyne.CanvasObject {
	appState := state.FromContext(ctx)
	if appState == nil {
		ctx = state.InitAppState(ctx)
		appState = state.FromContext(ctx)
	}
	fileList := getFileList(ctx, appState, w, updateFiles)

	uploadButton := widget.NewButton("Upload", func() {
		fileDialog := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
			if err != nil {
				panic(err)
			}

			err = bfsp.UploadFile(ctx, &bfsp.FileInfo{
				Name:   uri.URI().Name(),
				Reader: uri,
			}, 100)
			if err != nil {
				panic(err)
			}

			w.SetContent(FilesPage(ctx, w, true))

		}, w)
		fileDialog.Show()
	})
	usageButton := widget.NewButton("Usage", func() {
		w.SetContent(UsagePage(ctx, w))

	})

	buttons := container.NewGridWithColumns(2, uploadButton, usageButton)

	return container.NewBorder(nil, buttons, nil, nil, fileList)
}

func getFileList(ctx context.Context, appState *state.AppState, w fyne.Window, update bool) fyne.CanvasObject {
	if update {
		client := bfsp.ClientFromContext(ctx)
		masterKey := bfsp.MasterKeyFromContext(ctx)
		appState.UpdateAppState(client, masterKey)
	}

	fileMetaList := []*bfsp.FileMetadata{}
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

	list := widget.NewList(func() int {
		return len(fileMetaList)
	}, func() fyne.CanvasObject {
		return widget.NewButton("template", func() {

		})

	}, func(i widget.ListItemID, o fyne.CanvasObject) {
		fileMeta := fileMetaList[i]
		o.(*widget.Button).SetText(fileMeta.FileName)
		o.(*widget.Button).OnTapped = func() {
			page := FilePage(ctx, fileMeta, w)
			w.SetContent(page)
		}
	})
	return list

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
