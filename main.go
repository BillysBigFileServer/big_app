package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"net/url"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/BillysBigFileServer/app/pages"
	"github.com/BillysBigFileServer/app/preferences"
	"github.com/BillysBigFileServer/bfsp-go"
	bfspConfig "github.com/BillysBigFileServer/bfsp-go/config"
	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	a := app.NewWithID("io.bbfs.app")
	w := a.NewWindow("BBFS")

	ctx = preferences.ContextWithPreferences(ctx, a.Preferences())

	encryptionKey := a.Preferences().String("master_key")
	token := a.Preferences().String("token")

	if encryptionKey == "" && token == "" {
		w.SetContent(StartPage(ctx, w))
	} else {
		client, err := bfsp.NewHTTPFileServerClient(token, bfspConfig.FileServerBaseURL(), bfspConfig.FileServerHTTPS())
		if err != nil {
			panic(err)
		}

		masterKey, err := base64.StdEncoding.DecodeString(encryptionKey)
		if err != nil {
			panic(err)
		}
		ctx = bfsp.ContextWithMasterKey(ctx, masterKey)
		ctx = bfsp.ContextWithClient(ctx, client)

		w.SetContent(pages.FilesPage(ctx, w, []string{}, false))
	}

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
		signupURL, _ := url.Parse(bfspConfig.BigCentralBaseURL() + "/signup?dl_token=" + dlToken.String() + "#" + tempPubKey)
		fyne.CurrentApp().OpenURL(signupURL)
		w.SetContent(pages.AuthPage(ctx, w, signupURL, dlToken.String(), tempPrivKey.key))
	})

	loginButton := widget.NewButton("Login", func() {
		tempPubKey := tempPrivKey.Encode()
		loginURL, _ := url.Parse(bfspConfig.BigCentralBaseURL() + "/auth?dl_token=" + dlToken.String() + "#" + tempPubKey)
		fyne.CurrentApp().OpenURL(loginURL)
		w.SetContent(pages.AuthPage(ctx, w, loginURL, dlToken.String(), tempPrivKey.key))
	})

	return container.NewVBox(title, signupButton, loginButton)
}
