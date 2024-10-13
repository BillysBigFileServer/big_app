//go:generate fyne bundle -o bundled.go --pkg pages ../resources/back_arrow.png
//go:generate fyne bundle -o bundled.go --pkg pages -append ../resources/download.png
//go:generate fyne bundle -o bundled.go --pkg pages -append ../resources/share.png
//go:generate fyne bundle -o bundled.go --pkg pages -append ../resources/delete.png
//go:generate fyne bundle -o bundled.go --pkg pages -append ../resources/file_move_folder.png
//go:generate fyne bundle -o bundled.go --pkg pages -append ../resources/cancel.png

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
	"github.com/BillysBigFileServer/app/state"
	"github.com/BillysBigFileServer/bfsp-go"
	"github.com/BillysBigFileServer/bfsp-go/config"
)

func FilePage(ctx context.Context, fileMeta *bfsp.FileMetadata, w fyne.Window) fyne.CanvasObject {
	backArrow := resourceBackarrowPng
	downloadIcon := resourceDownloadPng
	shareIcon := resourceSharePng
	deleteIcon := resourceDeletePng
	fileMoveFolderIcon := resourceFilemovefolderPng

	backButton := widget.NewButtonWithIcon("", backArrow, func() {
		w.SetContent(FilesPage(ctx, w, fileMeta.Directory, false))
	})
	downloadButton := widget.NewButtonWithIcon("", downloadIcon, func() {
		fileSaver := dialog.NewFileSave(func(uri fyne.URIWriteCloser, err error) {
			go func() {
				err = bfsp.DownloadFile(ctx, fileMeta, uri, "")
				if err != nil {
					panic(err)
				}
			}()
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
		url := fmt.Sprintf("%s/files/view_file#z:%s", baseURL, viewInfoStr)

		w.Clipboard().SetContent(url)

		go func() {
			copiedToClipboardText.SetText("Copied to clipboard")
			time.Sleep(3 * time.Second)
			copiedToClipboardText.SetText("")
		}()
	})
	deleteButton := widget.NewButtonWithIcon("", deleteIcon, func() {
		cli := bfsp.ClientFromContext(ctx)

		if err := bfsp.DeleteFileMetadata(cli, fileMeta.Id); err != nil {
			err = fmt.Errorf("error deleting file metadata: %w", err)
			panic(err)
		}

		var chunks []string
		for _, chunkID := range fileMeta.Chunks {
			chunks = append(chunks, chunkID)
		}

		if err := bfsp.DeleteChunks(cli, chunks); err != nil {
			err = fmt.Errorf("error deleting file chunks: %w", err)
			panic(err)
		}

		// we should go back to whatever directory we were looking at
		w.SetContent(FilesPage(ctx, w, fileMeta.Directory, true))
	})
	fileMoveFolderButton := widget.NewButtonWithIcon("", fileMoveFolderIcon, func() {
		appState := state.FromContext(ctx)
		w.SetContent(DirectorySelectPage(ctx, appState, w, fileMeta.Directory, fileMeta, false))
	})

	buttonSeparator := widget.NewSeparator()
	leftButtonBox := container.NewHBox(backButton, buttonSeparator, downloadButton, buttonSeparator, shareButton, copiedToClipboardText)
	rightButtonBox := container.NewHBox(fileMoveFolderButton, deleteButton)
	buttonBox := container.NewBorder(nil, nil, leftButtonBox, rightButtonBox)

	fileName := widget.NewLabel(fileMeta.FileName)
	fileName.Wrapping = fyne.TextWrapWord
	fileName.TextStyle.Bold = true

	humanReadableSize := widget.NewLabel("Size: " + helper.HumanSize(fileMeta.FileSize))

	return container.NewVBox(buttonBox, fileName, humanReadableSize)
}
