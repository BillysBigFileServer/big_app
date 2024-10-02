package pages

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/BillysBigFileServer/app/helper"
	"github.com/BillysBigFileServer/app/preferences"
	"github.com/BillysBigFileServer/bfsp-go"
	"github.com/BillysBigFileServer/bfsp-go/config"
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
	shareIcon, err := fyne.LoadResourceFromPath("resources/share.png")
	if err != nil {
		panic(err)
	}

	backButton := widget.NewButtonWithIcon("", backArrow, func() {
		w.SetContent(FilesPage(ctx, w, false))
	})
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
	copiedToClipboardText := widget.NewLabel("")
	shareButton := widget.NewButtonWithIcon("", shareIcon, func() {
		token := preferences.GetPreferences(ctx).String("token")
		encryptionKey := preferences.GetPreferences(ctx).String("master_key")
		masterKey, err := base64.StdEncoding.DecodeString(encryptionKey)
		if err != nil {
			panic(err)
		}

		viewInfo, err := bfsp.ShareFile(fileMeta, token, masterKey)
		if err != nil {
			panic(err)
		}

		viewInfoStr, err := bfsp.EncodeViewFileInfo(viewInfo)
		if err != nil {
			panic(err)
		}

		baseURL := config.BigCentralBaseURL()
		url := fmt.Sprintf("%s/files/view_file/#z:%s", baseURL, viewInfoStr)

		w.Clipboard().SetContent(url)

		go func() {
			copiedToClipboardText.SetText("Copied to clipboard")
			time.Sleep(3 * time.Second)
			copiedToClipboardText.SetText("")
		}()
	})

	buttonSeparator := widget.NewSeparator()
	buttonBox := container.NewHBox(backButton, buttonSeparator, downloadButton, buttonSeparator, shareButton, copiedToClipboardText)

	fileName := widget.NewLabel(fileMeta.FileName)
	fileName.Wrapping = fyne.TextWrapWord
	fileName.TextStyle.Bold = true

	humanReadableSize := widget.NewLabel("Size: " + helper.HumanSize(fileMeta.FileSize))

	return container.NewVBox(buttonBox, fileName, humanReadableSize)
}
