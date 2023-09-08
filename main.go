// Versão: 1.0.0-RC2 (Candidato 2)
package main

import (
	_ "embed"
	"errors"
	"fmt"
	"image/color"
	"io"
	"net/url"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/alternativeon/pgo/v2"
	"github.com/melbahja/got"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
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
			dialog.ShowError(errors.New("usuário e senha não podem ficar vazios"), alternativeOnWindow)
			return
		}
		LoginDialog := dialog.NewCustomWithoutButtons("Realizando login....", widget.NewProgressBarInfinite(), alternativeOnWindow)
		LoginDialog.Show()
		userToken, err := pgo.Login(loginPainelLoginUsuario.Text, loginPainelLoginSenha.Text)
		if err != nil {
			LoginDialog.Hide()
			dialog.ShowError(err, alternativeOnWindow)
			return
		}
		oldUserToken, err := pgo.LegacyLogin(loginPainelLoginUsuario.Text, loginPainelLoginSenha.Text)
		if err != nil {
			LoginDialog.Hide()
			dialog.ShowError(err, alternativeOnWindow)
			return
		}
		userData, err := pgo.DadosUsuario(oldUserToken.AccessToken)
		if err != nil {
			LoginDialog.Hide()
			dialog.ShowError(err, alternativeOnWindow)
			return
		}

		if loginPainelSalvarSessãoCheck.Checked {
			alternativeOnApp.Preferences().SetString("config_session", "yes")
			alternativeOnApp.Preferences().SetString("username", loginPainelLoginUsuario.Text)
			alternativeOnApp.Preferences().SetString("password", loginPainelLoginSenha.Text)
		}
		LoginDialog.Hide()
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
	janela.Resize(fyne.NewSize(800, 600))
	janela.CenterOnScreen()
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
		Importance: widget.DangerImportance,
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
	botaoTabSobreConfig.OnTapped = func() {
		btnMudarTemaClaro := widget.NewButtonWithIcon("Tema claro", theme.NewPrimaryThemedResource(resourceLightmodeSvg), func() {
			app.Settings().SetTheme(theme.LightTheme())
		})
		btnMudarTemaClaro.Importance = widget.HighImportance
		btnMudarTemaEscuro := widget.NewButtonWithIcon("Tema escuro", resourceDarkmodeSvg, func() {
			app.Settings().SetTheme(theme.DarkTheme())
		})
		btnMudarTemaEscuro.Importance = widget.HighImportance
		lblConfigApp := widget.NewLabel("Aqui você pode trocar o tema do aplicativo")
		lblConfigApp.Wrapping = fyne.TextWrapWord
		containerAppSettings := container.NewVBox(btnMudarTemaClaro, btnMudarTemaEscuro, lblConfigApp)
		temaApp := dialog.NewCustom("Configurações do aplicativo", "Fechar", containerAppSettings, janela)
		temaApp.Show()
	}

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

	/* INICIO DOS LIVROS */
	bookStartDialog := dialog.NewCustomWithoutButtons("Carregando livros....", widget.NewProgressBarInfinite(), win)
	bookStartDialog.Show()
	labelLivros := widget.NewLabel("Carregando seus livros... POR FAVOR AGUARDE")
	labelLivros.Wrapping = fyne.TextWrapWord
	hyperlinkContainer := container.NewVBox(labelLivros, widget.NewSeparator())

	toolsBar := container.NewVBox(btnBar, livrosCont)
	livrosInterface := container.NewBorder(toolsBar, nil, nil, nil, container.NewScroll(hyperlinkContainer))
	win.SetContent(livrosInterface)
	for _, book := range livros {
		txtLivro := widget.NewLabel(fmt.Sprintf("%v - %v (%v)", book.ComponenteCurricular, book.Volume, book.Tipo))
		txtLivro.Wrapping = fyne.TextWrapOff
		btnLivro := widget.NewButton("Baixar...", func() {
			downloadPdf(book.URL, txtLivro.Text, win)
		})
		hAlign := container.NewHBox(btnLivro, txtLivro)

		hyperlinkContainer.Add(hAlign)
	}
	bookStartDialog.Hide()
	labelLivros.SetText("Baixe seus livros por aqui!")
}

func downloadPdf(url string, nome string, win fyne.Window) {
	downloadDialog := dialog.NewCustomWithoutButtons("Baixando livro....", widget.NewProgressBarInfinite(), win)
	downloadDialog.Show()
	g := got.New()

	tmpFile, err := os.CreateTemp("", "downloaded_*.pdf")
	if err != nil {
		dialog.ShowError(err, win)
		return
	}
	defer tmpFile.Close()

	err = g.Download(url, tmpFile.Name())
	if err != nil {
		dialog.ShowError(err, win)
		return
	}

	//fmt.Println(tmpFile.Name(), url)
	config := model.NewDefaultConfiguration()
	config.OwnerPW = "@rc0Tech"
	err = api.DecryptFile(tmpFile.Name(), tmpFile.Name()+"out.pdf", config)
	if err != nil {
		fmt.Println(err)
		dialog.ShowError(err, win)
		return
	}
	downloadDialog.Hide()

	// Create a file save dialog to choose where to save the downloaded file.
	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		if writer == nil {
			return
		}

		// Copy the decrypted PDF to the selected location.
		newSrc, _ := os.Open(tmpFile.Name() + "out.pdf")
		_, err = io.Copy(writer, newSrc)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		writer.Close()
		os.Remove(tmpFile.Name())
		os.Remove(tmpFile.Name() + "out.pdf")
		dialog.ShowInformation("Pronto", "O livro foi salvo com sucesso!", win)

	}, win)

	// Set the default file name for the save dialog.
	saveDialog.SetFileName(nome + ".pdf")
	saveDialog.SetFilter(&storage.ExtensionFileFilter{Extensions: []string{".pdf"}})

	// Show the save dialog.
	saveDialog.Show()
}
