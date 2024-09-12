package pages

import (
	"context"
	"crypto/rsa"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/BillysBigFileServer/bfsp-go"
)

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
