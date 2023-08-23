// Versão: 1.0.0-RC1 (Candidato 1)
package main

import (
	_ "embed"
	"errors"
	"fmt"
	"image/color"
	"net/url"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/alternativeon/pgo/v2"
)

func main() {
	//Inicialização da UI
	alternativeOnApp := app.NewWithID("link.princessmortix.aon")
	alternativeOnWindow := alternativeOnApp.NewWindow("Alternative ON")
	alternativeOnWindow.Resize(fyne.NewSize(800, 200))
	alternativeOnWindow.SetTitle("Alternative On - Login")
	alternativeOnWindow.SetIcon(&fyne.StaticResource{StaticName: "Icon", StaticContent: windowIcon})
	alternativeOnMetadata := fyne.AppMetadata{
		Build:   int(fyne.BuildRelease),
		Name:    "Alternative On",
		Release: true,
	}
	app.SetMetadata(alternativeOnMetadata)

	//Sessão do usuário
	fmt.Println(strings.Contains(alternativeOnApp.Preferences().String("config_session"), "yes"))
	if alternativeOnApp.Preferences().String("config_session") == "" {
		alternativeOnApp.Preferences().SetString("username", "")
		alternativeOnApp.Preferences().SetString("password", "")
	}
	//Sessão do usuário

	//Login do usuário
	loginPainelTextoAjuda := widget.NewRichTextFromMarkdown(textoPainelLogin)
	loginPainelTextoAjuda.Wrapping = fyne.TextWrapWord
	loginPainelLoginUsuarioTexto := widget.NewLabel("Usuário")
	loginPainelLoginUsuarioTexto.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	loginPainelLoginUsuario := widget.NewEntry()
	loginPainelLoginUsuario.PlaceHolder = "Coloque seu usuário aqui"
	loginPainelLoginSenhaTexto := widget.NewLabel("Senha")
	loginPainelLoginSenhaTexto.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	loginPainelLoginSenha := widget.NewPasswordEntry()
	loginPainelLoginSenha.PlaceHolder = "Coloque sua senha aqui"
	//Parte da sessão do usuário
	loginPainelSalvarSessãoCheck := widget.NewCheck("Salvar credenciais?", nil)
	loginPainelSalvarSessãoCheck.SetChecked(true)

	loginPainelSalvarSessãoAjuda := widget.NewHyperlink("O que é isso?", nil)
	dialogSalvarSessãoTexto := widget.NewLabel(textoPainelSessão)
	dialogSalvarSessãoTexto.Wrapping = fyne.TextWrapWord
	loginPainelSalvarSessãoAjuda.OnTapped = func() {
		lgnPainelSlvSajd := dialog.NewCustom("Salvar sessão - Ajuda", "Fechar", dialogSalvarSessãoTexto, alternativeOnWindow)
		lgnPainelSlvSajd.Show()
	}
	loginPainelSalvarSessão := container.New(layout.NewHBoxLayout(), loginPainelSalvarSessãoCheck, loginPainelSalvarSessãoAjuda)
	//Fim da opção de sessão

	loginPainelEntrada := container.New(layout.NewVBoxLayout(), loginPainelTextoAjuda, loginPainelLoginUsuarioTexto, loginPainelLoginUsuario, loginPainelLoginSenhaTexto, loginPainelLoginSenha)

	//Botão de Recuperar senha
	loginPainelBtnEsqueciSenha := widget.NewButtonWithIcon("Recuperar senha", lockResetIcon, func() { recuperarSenha(alternativeOnWindow) })
	loginPainelBtnEsqueciSenha.Importance = widget.MediumImportance
	loginPainelBtnEnviar := widget.NewButtonWithIcon("Entrar", theme.LoginIcon(), func() {
		if loginPainelLoginUsuario.Text == "" || loginPainelLoginSenha.Text == "" {
			dialog.ShowError(errors.New("Usuário e senha não podem ficar vazios"), alternativeOnWindow)
			return
		}
		userToken, err := pgo.Login(loginPainelLoginUsuario.Text, loginPainelLoginSenha.Text)
		if err != nil {
			dialog.ShowError(err, alternativeOnWindow)
			return
		}
		oldUserToken, err := pgo.LegacyLogin(loginPainelLoginUsuario.Text, loginPainelLoginSenha.Text)
		if err != nil {
			dialog.ShowError(err, alternativeOnWindow)
			return
		}
		userData, err := pgo.DadosUsuario(oldUserToken.AccessToken)
		if err != nil {
			dialog.ShowError(err, alternativeOnWindow)
			return
		}

		if loginPainelSalvarSessãoCheck.Checked {
			alternativeOnApp.Preferences().SetString("config_session", "yes")
			alternativeOnApp.Preferences().SetString("username", loginPainelLoginUsuario.Text)
			alternativeOnApp.Preferences().SetString("password", loginPainelLoginSenha.Text)
		}
		mudarConteudoAposLogin(alternativeOnWindow, alternativeOnApp, *userToken, *userData)
	})
	loginPainelBtnEnviar.Importance = widget.HighImportance
	loginPainelBtns := container.New(layout.NewHBoxLayout(), loginPainelBtnEsqueciSenha, loginPainelBtnEnviar)
	loginPainelEspaçador := canvas.NewLine(color.Transparent)
	loginPainelEspaçador.StrokeWidth = 3

	loginPane := container.New(layout.NewVBoxLayout(), loginPainelEntrada, loginPainelSalvarSessão, loginPainelEspaçador, loginPainelBtns)

	alternativeOnWindow.SetContent(loginPane)
	if strings.Contains(alternativeOnApp.Preferences().String("config_session"), "yes") {
		userToken, err := pgo.Login(alternativeOnApp.Preferences().String("username"), alternativeOnApp.Preferences().String("password"))
		if err != nil {
			alternativeOnApp.Preferences().SetString("config_session", "")
			fmt.Println(err)
			main()
			return
		}
		oldUserToken, err := pgo.LegacyLogin(alternativeOnApp.Preferences().String("username"), alternativeOnApp.Preferences().String("password"))
		if err != nil {
			alternativeOnApp.Preferences().SetString("config_session", "")
			fmt.Println(err)
			main()
			return
		}
		userData, err := pgo.DadosUsuario(oldUserToken.AccessToken)
		if err != nil {
			alternativeOnApp.Preferences().SetString("config_session", "")
			fmt.Println(err)
			main()
			return
		}

		alternativeOnWindow.SetTitle("Alternative On")
		alternativeOnApp.SendNotification(fyne.NewNotification("Sessão restaurada!", "Sua sessão foi automaticamente restaurada. Para mudar isso clique em 'Fazer logout'"))
		mudarConteudoAposLogin(alternativeOnWindow, alternativeOnApp, *userToken, *userData)
	}
	alternativeOnWindow.Show()
	alternativeOnApp.Run()
}

func mudarConteudoAposLogin(janela fyne.Window, app fyne.App, tokenUsuario pgo.Token, dadosUsuario pgo.DadosPrimitivos) {
	links := pgo.ObterRecursos(tokenUsuario.IdEscola, tokenUsuario.Token, tokenUsuario.TokenParceiro)

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
		OnTapped: func() { livrosUI(janela, app, tokenUsuario, dadosUsuario) },
	}

	containerAccordionLivros := container.NewVBox(labelAccordionLivros, &botaoAccordionLivros)
	accordionLivros := widget.NewAccordion(
		&widget.AccordionItem{
			Title:  "Livros",
			Detail: containerAccordionLivros,
			Open:   true,
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
			dialog.ShowConfirm("Você tem certeza?", "Você realmente quer sair do app?\nIsso também encerrará sua sessão.", func(b bool) {
				if b {
					app.Preferences().SetString("config_session", "")
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
	labelTabContaNomeUsuario := widget.NewLabel(fmt.Sprintf("Olá, %v!\n(ID: %v, ID legado: %v)", dadosUsuario.Nome, tokenUsuario.IdUsuario, dadosUsuario.IdUsuarioEscola)) //TODO: Implementar userinfo no pgo
	labelTabContaNomeUsuario.Wrapping = fyne.TextWrapWord
	labelTabContaNomeEscola := widget.NewLabel("Escola: " + tokenUsuario.NomeEscola + " (ID: " + tokenUsuario.IdEscola + ")")
	labelTabContaNomeEscola.Wrapping = fyne.TextWrapWord
	botaoTabContaMudarSenha := widget.Button{
		Text: "Mudar senha",
		Icon: lockResetIcon,
	}
	botaoTabContaMudarSenha.Disable()
	labelSessãoStatus := widget.NewLabel("Sua sessão **não** está sendo salva para o usuário atual.")
	labelSessãoStatus.Wrapping = fyne.TextWrapWord
	if strings.Contains(app.Preferences().String("config_session"), "yes") {
		labelSessãoStatus.Text = "Sua sessão **está** sendo salva para o usuário atual."
		labelSessãoStatus.Refresh()
	}

	conteudoTabConta := container.New(layout.NewVBoxLayout(), labelTabContaPrincipal, labelTabContaNomeUsuario, labelTabContaNomeEscola, labelSessãoStatus, &botaoTabContaMudarSenha)
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

func recuperarSenha(win fyne.Window) {
	recuperarSenhaTextoAjuda := widget.NewLabel("Informe seu e-mail, usuário ou cpf para continuar")
	recuperarSenhaTextoEntry := widget.NewEntry()
	recuperarSenhaTextoEntry.PlaceHolder = "CPF, E-mail ou usuário..."
	recuperarSenhaContainer := container.New(layout.NewVBoxLayout(), recuperarSenhaTextoAjuda, recuperarSenhaTextoEntry)
	recuperarSenhaDlg := dialog.NewCustomConfirm("Recuperar senha", "Enviar", "Fechar", recuperarSenhaContainer, func(b bool) {
		if b {
			ok, err := pgo.ResetarSenha(recuperarSenhaTextoEntry.Text)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			dialog.ShowInformation("Sucesso!", ok.Mensagem, win)
		}
	}, win)
	recuperarSenhaDlg.Show()
}

func livrosUI(win fyne.Window, app fyne.App, userToken pgo.Token, oldData pgo.DadosPrimitivos) {
	livros, err := pgo.ObterLivros(userToken.Token)
	if err != nil {
		dialog.ShowError(err, win)
		return
	}

	/* BOTÕES DA INTERFACE DOS LIVROS */
	interfaceLivrosVoltarBtn := widget.NewButtonWithIcon("Voltar", theme.NavigateBackIcon(), func() { mudarConteudoAposLogin(win, app, userToken, oldData) })
	interfaceLivrosAjudaBtn := widget.NewButtonWithIcon("Ajuda", theme.HelpIcon(), func() {
		ajudaTexto := widget.NewLabel(textoLivrosAjuda)
		ajudaTexto.Wrapping = fyne.TextWrapWord
		ajuda := dialog.NewCustom("Ajuda dos livros - Alternative On", "Fechar", ajudaTexto, win)
		ajuda.Show()
	})

	configLivrosCaminho := widget.NewLabelWithStyle("Pasta para salvar:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	configLivrosCaminhoLabel := widget.NewLabel("Nenhuma.")
	configLivrosCaminhoLabel.Wrapping = fyne.TextWrapBreak
	if app.Preferences().String("save_folder") != "" {
		configLivrosCaminhoLabel.Text = app.Preferences().String("save_folder")
	}
	configLivrosBtnMudarPasta := widget.NewButtonWithIcon("Mudar pasta", theme.FolderOpenIcon(), func() {
		pasta := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
			if err != nil || lu == nil {
				return
			}

			fmt.Printf("New path: %v, Old path: %v\n", lu.String(), app.Preferences().String("save_folder"))
			app.Preferences().SetString("save_folder", lu.String())
			configLivrosCaminhoLabel.Text = lu.Path()
			configLivrosCaminhoLabel.Refresh()

		}, win)
		pasta.SetConfirmText("Selecionar pasta")
		pasta.Show()
	})
	configLivrosBtnMudarPasta.Importance = widget.HighImportance
	configLivrosLayout := container.NewVBox(configLivrosCaminho, configLivrosCaminhoLabel, configLivrosBtnMudarPasta)

	configLivrosDlg := dialog.NewCustom("Configurações para baixar os livros", "Salvar", configLivrosLayout, win)
	configLivrosDlgBtn := widget.NewButtonWithIcon("Configurações", theme.SettingsIcon(), func() {
		configLivrosDlg.Show()
	})
	livrosCont := widget.NewLabel(fmt.Sprintf("Livros: %v", len(livros)))
	livrosCont.Wrapping = fyne.TextWrapBreak
	btnBar := container.NewHBox(interfaceLivrosVoltarBtn, interfaceLivrosAjudaBtn, configLivrosDlgBtn)

	/* FIM DOS BOTÕES */

	/*if fyne.CurrentDevice().IsMobile() {
		dialog.ShowConfirm("Não é possivel baixar os livros", "Em celulares essa função não é suportada,\ndeseja baixar todos os livros?", func(b bool) {
			if b {
				configLivrosDlg.Show()
				progress := widget.NewProgressBar()
				dlgLivrosDowload := dialog.NewCustom("Baixando....", "Fechar", progress, win)
				dlgLivrosDowload.Show()
			}
		}, win)
	}*/

	/* INICIO DOS LIVROS */

	/*for _, book := range livros {
		fmt.Println("Componente Curricular:", book.ComponenteCurricular)
		fmt.Println("Volume:", book.Volume)
		fmt.Println("Tipo:", book.Tipo)
		fmt.Println("URL:", book.URL)
		fmt.Println()
	}*/

	table := container.New(layout.NewGridLayout(4))
	cc := widget.NewLabel("Componente Curricular")
	cc.TextStyle = fyne.TextStyle{Bold: true}
	vol := widget.NewLabel("Volume")
	vol.TextStyle = fyne.TextStyle{Bold: true}
	tipo := widget.NewLabel("Tipo")
	tipo.TextStyle = fyne.TextStyle{Bold: true}
	action := widget.NewLabel("Ação")
	action.TextStyle = fyne.TextStyle{Bold: true}
	table.Add(cc)
	table.Add(vol)
	table.Add(tipo)
	table.Add(action)

	for _, book := range livros {
		table.Add(widget.NewLabel(book.ComponenteCurricular))
		table.Add(widget.NewLabel(book.Volume))
		table.Add(widget.NewLabel(book.Tipo))
		table.Add(widget.NewButton("Baixar", func() {
			win.Clipboard().SetContent("@rc0Tech")
			app.SendNotification(fyne.NewNotification("Senha copiada!", "A senha foi copiada para a área de transferencia."))
			app.OpenURL(parseUrl(book.URL))
		}))
	}
	tableContainer := fyne.NewContainerWithLayout(layout.NewMaxLayout(), table)
	scrollContainer := container.NewScroll(tableContainer)

	// Crie um container para armazenar a tabela
	/*content := container.NewVBox(
		widget.NewLabel("Seus livros digitais"),
		scrollContainer,
	)*/

	toolsBar := container.NewVBox(btnBar, livrosCont)
	//everything := container.NewVBox(toolsBar, scrollContainer)
	livrosInterface := container.NewBorder(toolsBar, nil, nil, nil, container.NewMax(scrollContainer))
	win.SetContent(livrosInterface)
}
