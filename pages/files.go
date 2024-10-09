package pages

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"

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
	fileList := getFileList(ctx, appState, w, directory, updateFiles)

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

				w.SetContent(FilesPage(ctx, w, directory, true))

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

	return container.NewBorder(top, buttons, nil, nil, fileList)
}

func getFileList(ctx context.Context, appState *state.AppState, w fyne.Window, directory []string, update bool) fyne.CanvasObject {
	if update {
		client := bfsp.ClientFromContext(ctx)
		masterKey := bfsp.MasterKeyFromContext(ctx)
		appState.UpdateAppState(client, masterKey)
	}

	directories := map[string]bool{}
	fileMetaList := []*bfsp.FileMetadata{}

	for {
		appState.RwLock.RLock()
		for _, meta := range appState.Files {
			if isDirEqual(directory, meta.Directory) {
				fileMetaList = append(fileMetaList, meta)
			}

			if isSubdirectory(directory, meta.Directory) {
				directories[sliceToDirectory(meta.Directory)] = true
			}
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
	directoriesList := []string{}
	for dir, _ := range directories {
		directoriesList = append(directoriesList, dir)
	}
	sort.Strings(directoriesList)

	buttons := []*widget.Button{}

	for _, dir := range directoriesList {
		button := widget.NewButton(dir, func() {
			dir := directoryToSlice(dir)
			page := FilesPage(ctx, w, dir, false)
			w.SetContent(page)
		})

		buttons = append(buttons, button)

	}

	for _, meta := range fileMetaList {
		button := widget.NewButton(meta.FileName, func() {
			page := FilePage(ctx, meta, w)
			w.SetContent(page)
		})
		buttons = append(buttons, button)
	}

	list := widget.NewList(func() int {
		return len(buttons)
	}, func() fyne.CanvasObject {
		return widget.NewButton("template", func() {})
	}, func(i widget.ListItemID, o fyne.CanvasObject) {
		o.(*widget.Button).SetText(buttons[i].Text)
		o.(*widget.Button).OnTapped = buttons[i].OnTapped
	})

	return list
}

func isDirEqual(dir1 []string, dir2 []string) bool {
	if len(dir1) != len(dir2) {
		return false
	}

	for idx := range len(dir1) {
		if dir1[idx] != dir2[idx] {
			return false
		}
	}

	return true
}

// This isn't too hard. We check that our current directory (dir1) is < dir2's length by 1, then we just check that each item in dir1 matches each item in dir2 (to make sure it's actually a subdirectory)
func isSubdirectory(dir1 []string, dir2 []string) bool {
	if len(dir1)+1 != len(dir2) {
		return false
	}

	for idx := range dir1 {
		if dir1[idx] != dir2[idx] {
			return false
		}
	}

	return true
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

func sliceToDirectory(dir []string) string {
	return "/" + strings.Join(dir, "/")
}
func directoryToSlice(dir string) []string {
	if dir == "/" {
		return []string{}
	}
	dir = strings.TrimPrefix(dir, "/")
	slice := strings.Split(dir, "/")
	for idx := range slice {
		if slice[idx] == " " {
			slice[idx] = ""
		}
	}
	return slice
}
