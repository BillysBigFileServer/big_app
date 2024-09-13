package pages

import (
	"context"
	"crypto/rsa"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/BillysBigFileServer/bfsp-go"
	bfspConfig "github.com/BillysBigFileServer/bfsp-go/config"
)

func AuthPage(ctx context.Context, w fyne.Window, url *url.URL, dlToken string, tempPrivKey *rsa.PrivateKey) fyne.CanvasObject {
	waiting := widget.NewLabel("Waiting for you to login at:")
	pleaseOpen := widget.NewHyperlink(url.String(), url)

	go func() {
		tokenInfo, err := bfsp.GetToken("http://localhost:4000/", dlToken, tempPrivKey)
		if err != nil {
			panic(err)
		}

		client, err := bfsp.NewHTTPFileServerClient(tokenInfo.Token, "localhost:9998", false)
		if err != nil {
			panic(err)
		}
		ctx = bfsp.ContextWithMasterKey(ctx, tokenInfo.MasterKey)
		ctx = bfsp.ContextWithClient(ctx, client)

		//TODO: use the fyne preference api for easy cross platform preferenes
		configFile, err := bfspConfig.OpenDefaultConfigFile()
		defer configFile.Close()
		if err != nil {
			panic(err)
		}
		config, err := bfspConfig.ReadConfig(configFile)
		if err != nil {
			panic(err)
		}
		config.SetEncryptionKey(tokenInfo.MasterKey)
		config.Token = tokenInfo.Token

		err = bfspConfig.WriteConfigToFile(configFile, config)
		if err != nil {
			panic(err)
		}

		w.SetContent(FilesPage(ctx, w, false))
	}()

	return container.NewVBox(waiting, pleaseOpen)
}
