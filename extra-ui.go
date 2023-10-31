package main

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/alternativeon/pgo/v2"
)

func recuperarSenha(win fyne.Window) {
	recuperarSenhaTextoAjuda := widget.NewLabel("Informe seu e-mail, usuário ou cpf para continuar")
	recuperarSenhaTextoEntry := widget.NewEntry()
	recuperarSenhaTextoEntry.PlaceHolder = "CPF, E-mail ou usuário..."
	recuperarSenhaContainer := container.New(layout.NewVBoxLayout(), recuperarSenhaTextoAjuda, recuperarSenhaTextoEntry)
	recuperarSenhaDlg := dialog.NewCustomConfirm("Recuperar senha", "Enviar", "Fechar", recuperarSenhaContainer, func(b bool) {
		if b {
			ok, err := pgo.RecuperarSenha(recuperarSenhaTextoEntry.Text)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			dialog.ShowInformation("Sucesso!", ok.Mensagem, win)
		}
	}, win)
	recuperarSenhaDlg.Show()
}

func parseUrl(link string) *url.URL {
	parseLink, _ := url.Parse(link)
	return parseLink
}
