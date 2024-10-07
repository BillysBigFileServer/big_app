package pages

import (
	"bytes"
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/BillysBigFileServer/app/helper"
	"github.com/BillysBigFileServer/bfsp-go/usage"
	"github.com/vicanso/go-charts/v2"
)

func UsagePage(ctx context.Context, w fyne.Window) fyne.CanvasObject {
	title := container.NewCenter(widget.NewLabel("Usage"))
	usage, err := usage.GetUsage(ctx)
	if err != nil {
		panic(err)
	}

	values := []uint64{
		usage.StorageCap,
		usage.TotalUsage,
	}
	valuesf64 := []float64{
		float64(usage.StorageCap),
		float64(usage.TotalUsage),
	}

	backArrow := resourceBackarrowPng
	backButton := widget.NewButtonWithIcon("", backArrow, func() {
		w.SetContent(FilesPage(ctx, w, false))
	})
	p, err := charts.PieRender(
		valuesf64,
		charts.LegendOptionFunc(charts.LegendOption{
			Orient: charts.OrientVertical,
			Data: []string{
				"Total",
				"Usage",
			},
			Left: charts.PositionLeft,
		}),
		func(opt *charts.ChartOption) {
			for idx := range opt.SeriesList {
				opt.SeriesList[idx].Label.Show = true
				opt.SeriesList[idx].Label.Formatter = helper.HumanSize(values[idx])

			}
		},
		charts.ThemeOptionFunc(charts.ThemeDark),
	)
	if err != nil {
		panic(err)
	}

	buf, err := p.Bytes()
	if err != nil {
		panic(err)
	}
	r := bytes.NewReader(buf)
	image := canvas.NewImageFromReader(r, "image")
	image.FillMode = canvas.ImageFillOriginal

	topBox := container.NewHBox(backButton, title)
	return container.NewVBox(topBox, image)
}
