package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func getListData(podData *[]string) (binding.ExternalStringList, *widget.List) {
	// list binding, bind pod list data to data
	data := binding.BindStringList(
		podData,
	)

	list := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	return data, list
}

func setupErrorUI(stringErrorResponse string, list *widget.List) (*widget.Label, *widget.List, *widget.Button) {
	title := widget.NewLabel("")
	title.TextStyle = fyne.TextStyle{Monospace: true}
	title.Text = stringErrorResponse
	title.Wrapping = fyne.TextWrapBreak
	title.TextStyle = fyne.TextStyle{Italic: true, Bold: true}
	title.Refresh()
	list.Hide()
	refresh := widget.NewButton("Refresh", func() {})
	refresh.Refresh()
	return title, list, refresh
}
