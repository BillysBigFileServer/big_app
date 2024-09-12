package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/url"
	"slices"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/BillysBigFileServer/app/state"
	"github.com/BillysBigFileServer/bfsp-go"
	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	a := app.New()
	w := a.NewWindow("BBFS")

	masterKey, err := bfsp.CreateMasterEncKey("johnbon")
	if err != nil {
		panic(err)
	}
	ctx = bfsp.ContextWithMasterKey(ctx, masterKey)

	w.SetContent(StartPage(ctx, w))
	w.ShowAndRun()
}

type rsaKey struct {
	key   *rsa.PrivateKey
	mutex sync.Mutex
}

func (k *rsaKey) initRsaKey() error {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	// TODO: i'm sure there's a faster algorithm we can use
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	privKey.Precompute()
	k.key = privKey
	return nil
}

func (k *rsaKey) Encode() string {
	k.mutex.Lock()
	defer k.mutex.Unlock()

	tempPubEncKeyBytes := x509.MarshalPKCS1PublicKey(&k.key.PublicKey)
	return base64.URLEncoding.EncodeToString(tempPubEncKeyBytes)
}

func StartPage(ctx context.Context, w fyne.Window) fyne.CanvasObject {
	title := widget.NewLabel("BBFS")
	title.Alignment = fyne.TextAlignCenter
	dlToken, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	tempPrivKey := &rsaKey{}
	// we generate the rsa key in the background, to not slow down load times
	go tempPrivKey.initRsaKey()

	signupButton := widget.NewButton("Signup", func() {
		tempPubKey := tempPrivKey.Encode()
		signupURL, _ := url.Parse("http://bbfs.io/signup?dl_token=" + dlToken.String() + "#" + tempPubKey)
		fyne.CurrentApp().OpenURL(signupURL)
		w.SetContent(AuthPage(ctx, w, signupURL, dlToken.String(), tempPrivKey.key))
	})

	loginButton := widget.NewButton("Login", func() {
		tempPubKey := tempPrivKey.Encode()
		loginURL, _ := url.Parse("https://bbfs.io/auth?dl_token=" + dlToken.String() + "#" + tempPubKey)
		fyne.CurrentApp().OpenURL(loginURL)
		w.SetContent(AuthPage(ctx, w, loginURL, dlToken.String(), tempPrivKey.key))
	})

	return container.NewVBox(title, signupButton, loginButton)
}

func AuthPage(ctx context.Context, w fyne.Window, url *url.URL, dlToken string, tempPrivKey *rsa.PrivateKey) fyne.CanvasObject {
	waiting := widget.NewLabel("Waiting for you to login at:")
	pleaseOpen := widget.NewHyperlink(url.String(), url)

	go func() {
		tokenInfo, err := bfsp.GetToken("https://bbfs.io/", dlToken, tempPrivKey)
		if err != nil {
			panic(err)
		}

		client, err := bfsp.NewHTTPFileServerClient(tokenInfo.Token, "big-file-server.fly.dev:9998", true)
		if err != nil {
			panic(err)
		}
		ctx = bfsp.ContextWithMasterKey(ctx, tokenInfo.MasterKey)
		ctx = bfsp.ContextWithClient(ctx, client)

		w.SetContent(FilesPage(ctx, w))
	}()

	return container.NewVBox(waiting, pleaseOpen)
}

func InitAppState(ctx context.Context) context.Context {
	appState := state.AppState{}
	go func() {
		client := bfsp.ClientFromContext(ctx)
		masterKey := bfsp.MasterKeyFromContext(ctx)
		fileMetas, err := bfsp.ListFileMetadata(client, []string{}, masterKey)
		if err != nil {
			panic(err)
		}

		appState.Initialized.Store(true)
		appState.RwLock.Lock()
		appState.Files = fileMetas
		appState.RwLock.Unlock()

		time.Sleep(10 * time.Second)
	}()
	return state.ContextWithAppState(ctx, &appState)
}

func FilesPage(ctx context.Context, w fyne.Window) fyne.CanvasObject {
	ctx = InitAppState(ctx)
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
		fileWidgets = append(fileWidgets, FileWidget(ctx, fileMeta, w))
	}

	return container.NewVScroll(container.NewVBox(
		fileWidgets...,
	))
}

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

	humanReadableSize := widget.NewLabel("Size: " + HumanSize(fileMeta.FileSize))

	return container.NewVBox(backButton, fileName, humanReadableSize)
}

func FileWidget(ctx context.Context, fileMeta *bfsp.FileMetadata, w fyne.Window) fyne.CanvasObject {
	checkbox := widget.NewCheck("", func(isToggled bool) {
		fmt.Println(isToggled)
	})
	fileName := widget.NewButton(AbridgedFileName(fileMeta.FileName), func() {
		filePage := FilePage(ctx, fileMeta, w)
		w.SetContent(filePage)

	})

	return container.NewHBox(checkbox, fileName)
}
