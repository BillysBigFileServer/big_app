package pages

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/BillysBigFileServer/app/helper"
	"github.com/BillysBigFileServer/bfsp-go"
)

func DirectoryButton(ctx context.Context, w fyne.Window, dir string, pageFunc func(dir []string) fyne.CanvasObject) *widget.Button {
	return widget.NewButton(dir, func() {
		dir := helper.DirectoryToSlice(dir)
		w.SetContent(pageFunc(dir))
	})
}

func FileButton(ctx context.Context, w fyne.Window, fileMeta *bfsp.FileMetadata) *widget.Button {
	return widget.NewButton(fileMeta.FileName, func() {
		page := FilePage(ctx, fileMeta, w)
		w.SetContent(page)
	})
}
