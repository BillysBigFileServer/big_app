package pages

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/BillysBigFileServer/app/helper"
	"github.com/BillysBigFileServer/app/state"
	"github.com/BillysBigFileServer/bfsp-go"
)

func FilesPage(ctx context.Context, w fyne.Window, directory []string, updateFiles bool) fyne.CanvasObject {
	appState := state.FromContext(ctx)
	if appState == nil {
		ctx = state.InitAppState(ctx)
		appState = state.FromContext(ctx)
	}
	fileMetaList := helper.GetFileList(ctx, appState, directory, updateFiles)
	dirList := helper.GetDirList(ctx, appState, directory, updateFiles)

	fileButtons := []*widget.Button{}
	for _, dir := range dirList {
		fileButtons = append(fileButtons, DirectoryButton(ctx, w, dir, func(dir []string) fyne.CanvasObject { return FilesPage(ctx, w, dir, updateFiles) }))
	}

	for _, fileMeta := range fileMetaList {
		fileButtons = append(fileButtons, FileButton(ctx, w, fileMeta))
	}

	clickableList := widget.NewList(func() int {
		return len(fileButtons)
	}, func() fyne.CanvasObject {
		return widget.NewButton("template", func() {})
	}, func(i widget.ListItemID, o fyne.CanvasObject) {
		o.(*widget.Button).SetText(fileButtons[i].Text)
		o.(*widget.Button).OnTapped = fileButtons[i].OnTapped
	})

	uploadButton := widget.NewButton("Upload", func() {
		fileDialog := dialog.NewFileOpen(func(uri fyne.URIReadCloser, err error) {
			if err != nil {
				panic(err)
			}

			go func() {
				err = bfsp.UploadFile(ctx, &bfsp.FileInfo{
					Name:   uri.URI().Name(),
					Reader: uri,
				}, 100)
				if err != nil {
					panic(err)
				}

				w.SetContent(FilesPage(ctx, w, directory, true))
			}()

		}, w)
		fileDialog.Show()
	})
	usageButton := widget.NewButton("Usage", func() {
		w.SetContent(UsagePage(ctx, w))
	})

	buttons := container.NewGridWithColumns(2, uploadButton, usageButton)

	var top *fyne.Container

	switch len(directory) {
	case 0:
		top = container.NewHBox()
	default:
		backIcon := resourceBackarrowPng
		backButton := widget.NewButtonWithIcon("", backIcon, func() {
			page := FilesPage(ctx, w, directory[:len(directory)-1], false)
			w.SetContent(page)
		})
		top = container.NewHBox(backButton)
	}

	return container.NewBorder(top, buttons, nil, nil, clickableList)
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
