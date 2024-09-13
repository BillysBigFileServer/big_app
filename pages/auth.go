package pages

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/BillysBigFileServer/app/preferences"
	"github.com/BillysBigFileServer/bfsp-go"
	bfspConfig "github.com/BillysBigFileServer/bfsp-go/config"
)

func AuthPage(ctx context.Context, w fyne.Window, url *url.URL, dlToken string, tempPrivKey *rsa.PrivateKey) fyne.CanvasObject {
	waiting := widget.NewLabel("Waiting for you to login at:")
	pleaseOpen := widget.NewHyperlink(url.String(), url)

	go func() {
		tokenInfo, err := bfsp.GetToken(bfspConfig.BigCentralBaseURL(), dlToken, tempPrivKey)
		if err != nil {
			panic(err)
		}

		client, err := bfsp.NewHTTPFileServerClient(tokenInfo.Token, bfspConfig.FileServerBaseURL(), bfspConfig.FileServerHTTPS())
		if err != nil {
			panic(err)
		}
		ctx = bfsp.ContextWithMasterKey(ctx, tokenInfo.MasterKey)
		ctx = bfsp.ContextWithClient(ctx, client)

		preferences := preferences.GetPreferences(ctx)
		preferences.SetString("master_key", base64.StdEncoding.EncodeToString(tokenInfo.MasterKey))
		preferences.SetString("token", string(tokenInfo.Token))

		w.SetContent(FilesPage(ctx, w, false))
	}()

	return container.NewVBox(waiting, pleaseOpen)
}
