package pages

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/BillysBigFileServer/app/helper"
	"github.com/BillysBigFileServer/bfsp-go"
)

func FilePage(ctx context.Context, fileMeta *bfsp.FileMetadata, w fyne.Window) fyne.CanvasObject {
	backArrow, err := fyne.LoadResourceFromPath("resources/back_arrow.png")
	if err != nil {
		panic(err)
	}

	downloadIcon, err := fyne.LoadResourceFromPath("resources/download.png")
	if err != nil {
		panic(err)
	}

	downloadButton := widget.NewButtonWithIcon("", downloadIcon, func() {
		fileSaver := dialog.NewFileSave(func(uri fyne.URIWriteCloser, err error) {
			err = bfsp.DownloadFile(ctx, fileMeta, uri, "")
			if err != nil {
				panic(err)
			}
		}, w)
		fileSaver.SetFileName(fileMeta.FileName)
		fileSaver.Show()
	})
	buttonSeparator := widget.NewSeparator()
	backButton := widget.NewButtonWithIcon("", backArrow, func() {
		w.SetContent(FilesPage(ctx, w, false))
	})
	buttonBox := container.NewHBox(backButton, buttonSeparator, downloadButton)

	fileName := widget.NewLabel(fileMeta.FileName)
	fileName.Wrapping = fyne.TextWrapWord
	fileName.TextStyle.Bold = true

	humanReadableSize := widget.NewLabel("Size: " + helper.HumanSize(fileMeta.FileSize))

	return container.NewVBox(buttonBox, fileName, humanReadableSize)
}
