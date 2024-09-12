package pages

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/BillysBigFileServer/app/helper"
	"github.com/BillysBigFileServer/bfsp-go"
)

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

	humanReadableSize := widget.NewLabel("Size: " + helper.HumanSize(fileMeta.FileSize))

	return container.NewVBox(backButton, fileName, humanReadableSize)
}
