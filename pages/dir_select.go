package pages

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/BillysBigFileServer/app/helper"
	"github.com/BillysBigFileServer/app/state"
	"github.com/BillysBigFileServer/bfsp-go"
)

func DirectorySelectPage(ctx context.Context, appState *state.AppState, w fyne.Window, dir []string, fileMeta *bfsp.FileMetadata, updateFiles bool) fyne.CanvasObject {
	dirList := helper.GetDirList(ctx, appState, dir, updateFiles)
	dirButtons := []*widget.Button{}
	for _, dir := range dirList {
		dirButtons = append(dirButtons, DirectoryButton(ctx, w, dir, func(dir []string) fyne.CanvasObject {
			return DirectorySelectPage(ctx, appState, w, dir, fileMeta, false)
		}))
	}

	list := widget.NewList(func() int {
		return len(dirButtons)
	}, func() fyne.CanvasObject {
		return widget.NewButton("template", func() {})
	}, func(i widget.ListItemID, o fyne.CanvasObject) {
		o.(*widget.Button).SetText(dirButtons[i].Text)
		o.(*widget.Button).OnTapped = dirButtons[i].OnTapped
	})

	cancelIcon := resourceCancelPng
	cancelButton := widget.NewButtonWithIcon("Cancel", cancelIcon, func() {
		page := FilePage(ctx, fileMeta, w)
		w.SetContent(page)
	})

	var top *fyne.Container
	switch len(dir) {
	case 0:
		top = container.NewHBox(cancelButton)
	default:
		backIcon := resourceBackarrowPng
		backButton := widget.NewButtonWithIcon("", backIcon, func() {
			w.SetContent(DirectorySelectPage(ctx, appState, w, dir[:len(dir)-1], fileMeta, false))
		})

		top = container.NewHBox(backButton, cancelButton)
	}

	fileMoveFolderIcon := resourceFilemovefolderPng
	bottom := widget.NewButtonWithIcon("Move File", fileMoveFolderIcon, func() {
		cli := bfsp.ClientFromContext(ctx)
		masterKey := bfsp.MasterKeyFromContext(ctx)

		fileMeta.Directory = dir
		if err := bfsp.UpdateFileMetadata(cli, fileMeta, masterKey); err != nil {
			panic(err)
		}

		w.SetContent(FilesPage(ctx, w, fileMeta.Directory, false))
	})

	return container.NewBorder(top, bottom, nil, nil, list)
}
