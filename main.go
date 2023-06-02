// Version: 0.0.7 (Beta 7)
package main

import (
	_ "embed"
	"errors"
	"fmt"
	"net/url"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/alternativeon/pgo/v2"
)

//go:embed winres/icon.png
var windowIcon []byte
var textoHub = `## Aqui você pode acessar suas atividades, livros e mensagens`
var textoFuncLivros = `Acesse seus livros digitais no formato de PDF. 

**Função ainda em desenvolvimento**, e não implementada no aplicativo ainda.`
var textoTabConta = `## Aqui você pode controlar e alterar sua conta do Positivo.

Atualmente é possivel:
- Ver o nome e escola

É planejado poder mudar a sua senha nesta aba.`
var textoTabSobre = `## Alternative On - v0.1.0

Alternative On é um projeto que visa ser um cliente alternativo a plataforma de estudos Positivo On, contendo aplicativos próprios para Android, Linux e Windows.

**Veja mais sobre o projeto no página do Github.**



## Créditos:
- Fyne
- Postman
- Yngrid, Luis Fernado
- Tamiris, Nicole, Maria Cecilia
- Todo mundo que contribui no desenvolvimento programacional (<3)
- E você, por usar o aplicativo!

---
> Aplicativo licenciado para você sobre a licença BSD-3. O icone do aplicativo está sobre a licença do Creative Commons BY-NC (CC BY-NC)

> Os outros icones dentro do aplicativo estão lincenciados sobre a licença do Fyne.
`

func main() {
	alternativeOnApp := app.NewWithID("link.princessmortix.aon")
	alternativeOnWindow := alternativeOnApp.NewWindow("Alternative ON")
	alternativeOnWindow.Resize(fyne.NewSize(800, 600))
	//alternativeOnWindow.SetFixedSize(true)
	alternativeOnWindow.SetTitle("Alternative On - Login")
	alternativeOnWindow.SetIcon(&fyne.StaticResource{StaticName: "Icon", StaticContent: windowIcon})

	loginPaneSparator := widget.NewSeparator()
	loginPaneText := widget.NewLabel("Faça seu login para continuarmos. Utilize o mesmo usuário e senha do Positivo On.")
	loginPaneText.Wrapping = fyne.TextWrapWord
	loginPaneUser := widget.NewEntry()
	loginPanePass := widget.NewPasswordEntry()
	loginPaneEntry := &widget.Form{
		Items: []*widget.FormItem{{Text: "Usuário", Widget: loginPaneUser}, {Text: "Senha", Widget: loginPanePass}},
		OnSubmit: func() {
			loginPaneUser.Disable()
			loginPanePass.Disable()
			if loginPaneUser.Text == "" || loginPanePass.Text == "" {
				dialog.ShowError(errors.New("O usuário nem a senha podem ficar vazios."), alternativeOnWindow)
				loginPaneUser.Enable()
				loginPanePass.Enable()
				return
			}
			userData, err := pgo.Login(loginPaneUser.Text, loginPanePass.Text)
			if err != nil {
				fmt.Println(err)
				dialog.ShowError(errors.New("Não foi possivel fazer o login\nVerifique o usuário e a senha e tente novamente."), alternativeOnWindow)
				loginPaneUser.Enable()
				loginPanePass.Enable()
				return
			}
			alternativeOnWindow.SetTitle("Alternative On")
			//links := pgo.ObterRecursos(userData.IdEscola, userData.Token)

			//dialog.ShowInformation("Sucesso!", "O login foi feito com sucesso!", alternativeOnWindow)
			mudarConteudoAposLogin(alternativeOnWindow, alternativeOnApp, *userData)
		},
		SubmitText: "Fazer Login",
	}
	loginPane := container.New(layout.NewVBoxLayout(), loginPaneText, loginPaneSparator, loginPaneEntry)

	alternativeOnWindow.SetContent(loginPane)
	alternativeOnWindow.Show()
	alternativeOnApp.Run()
}

func mudarConteudoAposLogin(janela fyne.Window, app fyne.App, dadosUsuario pgo.Token) {
	links := pgo.ObterRecursos(dadosUsuario.IdEscola, dadosUsuario.Token)

	//UI APOS O LOGIN
	/* Tab 1: Principal */
	labelTabHub := widget.NewRichTextFromMarkdown(textoHub)
	labelTabHub.Wrapping = fyne.TextWrapWord
	labelAccordionAtividades := widget.Label{
		Text: "Acesse suas atividades do Positivo On.",
	}
	botaoAccordionAtividades := widget.Button{
		Text:     "Acessar",
		Icon:     theme.ComputerIcon(),
		OnTapped: func() { app.OpenURL(parseUrl(links.Studos)) },
	}

	containerAccordionAtividades := container.NewVBox(&labelAccordionAtividades, &botaoAccordionAtividades)
	accordionAtividades := widget.NewAccordion(
		&widget.AccordionItem{
			Title:  "Atividades",
			Detail: containerAccordionAtividades,
			Open:   true,
		},
	)

	labelAccordionLivros := widget.NewRichTextFromMarkdown(textoFuncLivros)
	labelAccordionLivros.Wrapping = fyne.TextWrapWord
	botaoAccordionLivros := widget.Button{
		Text:     "Acessar livros",
		Icon:     theme.DocumentIcon(),
		OnTapped: func() { dialog.NewInformation("Cuidado...", "Função ainda não implementada", janela) },
	}
	botaoAccordionLivros.Disable()

	containerAccordionLivros := container.NewVBox(labelAccordionLivros, &botaoAccordionLivros)
	accordionLivros := widget.NewAccordion(
		&widget.AccordionItem{
			Title:  "Livros",
			Detail: containerAccordionLivros,
			Open:   false,
		},
	)

	labelAccordionMensagens := widget.Label{
		Text: "Acesse as mensagens enviadas a você.",
	}
	botaoAccordionMensagens := widget.Button{
		Text:     "Acessar",
		Icon:     theme.FileTextIcon(),
		OnTapped: func() { app.OpenURL(parseUrl(links.Mensagens)) },
	}

	containerAccordionMensagens := container.NewVBox(&labelAccordionMensagens, &botaoAccordionMensagens)
	accordionMensagens := widget.NewAccordion(
		&widget.AccordionItem{
			Title:  "Ver Mensagens",
			Detail: containerAccordionMensagens,
			Open:   true,
		},
	)
	botaoAccordionLogout := widget.Button{
		Text: "Fazer logout",
		Icon: theme.LogoutIcon(),
		OnTapped: func() {
			dialog.ShowConfirm("Você tem certeza?", "Você realmente quer sair do app?", func(b bool) {
				if b {
					app.Quit()
					os.Exit(0)
				}
			}, janela)
		},
	}

	conteudoAccordionAtividades := container.New(layout.NewVBoxLayout(), labelTabHub, accordionAtividades, accordionLivros, accordionMensagens, &botaoAccordionLogout)
	/* Tab 1: Principal */

	/* Tab 2: Conta do usuário */
	labelTabContaPrincipal := widget.NewRichTextFromMarkdown(textoTabConta)
	labelTabContaPrincipal.Wrapping = fyne.TextWrapWord
	labelTabContaNomeUsuario := widget.NewLabel("Seu nome é: ???????" + " (ID: " + dadosUsuario.IdUsuario + ")") //TODO: Implementar userinfo no pgo
	labelTabContaNomeUsuario.Wrapping = fyne.TextWrapWord
	labelTabContaNomeEscola := widget.NewLabel("Escola: " + dadosUsuario.NomeEscola + " (ID: " + dadosUsuario.IdEscola + ")")
	labelTabContaNomeEscola.Wrapping = fyne.TextWrapWord
	botaoTabContaMudarSenha := widget.Button{
		Text: "Mudar senha",
		Icon: theme.AccountIcon(),
	}
	botaoTabContaMudarSenha.Disable()

	conteudoTabConta := container.New(layout.NewVBoxLayout(), labelTabContaPrincipal, labelTabContaNomeUsuario, labelTabContaNomeEscola, &botaoTabContaMudarSenha)
	/* Tab 2: Conta do usuário */

	/* Tab 3: Sobre */
	labelTabSobre := widget.NewRichTextFromMarkdown(textoTabSobre)
	labelTabSobre.Wrapping = fyne.TextWrapWord
	botaoTabSobreGH := widget.Button{
		Text:     "Ir para a página do projeto",
		Icon:     theme.ComputerIcon(),
		OnTapped: func() { app.OpenURL(parseUrl("https://github.com/AlternativeOn/AlternativeOn")) },
	}
	botaoTabSobreConfig := widget.Button{
		Text: "Configurações do app",
		Icon: theme.SettingsIcon(),
	}
	botaoTabSobreConfig.Disable()

	conteudoTabSobre := container.New(layout.NewVBoxLayout(), labelTabSobre, &botaoTabSobreGH, &botaoTabSobreConfig)
	/* Tab 3: Sobre */

	abasUiAposLogin := container.NewAppTabs(
		container.NewTabItemWithIcon("Hub", theme.HomeIcon(), conteudoAccordionAtividades),
		container.NewTabItemWithIcon("Conta", theme.AccountIcon(), conteudoTabConta),
		container.NewTabItemWithIcon("Sobre", theme.InfoIcon(), conteudoTabSobre),
	)
	abasUiAposLogin.SetTabLocation(container.TabLocationLeading)

	janela.SetContent(abasUiAposLogin)
}

func parseUrl(link string) *url.URL {
	parseLink, _ := url.Parse(link)
	return parseLink
}
