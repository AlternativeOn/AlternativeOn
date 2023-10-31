// Versão: 1.0.0
package main

import (
	_ "embed"
	"errors"
	"fmt"
	"image/color"
	"io"
	"os"
	"path/filepath"
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

var TokenData pgo.Token
var TokenDataPrimitivo pgo.DadosPrimitivos
var LivrosData []pgo.Item

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

	//Sessão do usuário - Verificação nas configurações
	if alternativeOnApp.Preferences().String("config_session") == "" {
		alternativeOnApp.Preferences().SetString("username", "")
		alternativeOnApp.Preferences().SetString("password", "")
	}

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
	loginPainelSalvarSessão := container.NewHBox(loginPainelSalvarSessãoCheck, loginPainelSalvarSessãoAjuda)
	//Fim da opção de sessão

	loginPainelEntrada := container.NewVBox(loginPainelTextoAjuda, loginPainelLoginUsuarioTexto, loginPainelLoginUsuario, loginPainelLoginSenhaTexto, loginPainelLoginSenha)

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
		TokenData = *userToken
		TokenDataPrimitivo = *userData
		interfacePrincipal(alternativeOnWindow, alternativeOnApp)
	})
	loginPainelBtnEnviar.Importance = widget.HighImportance
	loginPainelBtns := container.NewHBox(loginPainelBtnEsqueciSenha, loginPainelBtnEnviar)
	loginPainelEspaçador := canvas.NewLine(color.Transparent)
	loginPainelEspaçador.StrokeWidth = 3

	loginPane := container.NewVBox(loginPainelEntrada, loginPainelSalvarSessão, loginPainelEspaçador, loginPainelBtns)

	alternativeOnWindow.SetContent(loginPane)

	if strings.Contains(alternativeOnApp.Preferences().String("config_session"), "yes") {
		userToken, err := pgo.Login(alternativeOnApp.Preferences().String("username"), alternativeOnApp.Preferences().String("password"))
		if err != nil {
			alternativeOnApp.Preferences().SetString("config_session", "")
			fmt.Println(err)
			main()
		}
		oldUserToken, err := pgo.LegacyLogin(alternativeOnApp.Preferences().String("username"), alternativeOnApp.Preferences().String("password"))
		if err != nil {
			alternativeOnApp.Preferences().SetString("config_session", "")
			fmt.Println(err)
			main()
		}
		userData, err := pgo.DadosUsuario(oldUserToken.AccessToken)
		if err != nil {
			alternativeOnApp.Preferences().SetString("config_session", "")
			fmt.Println(err)
			main()
		}

		alternativeOnWindow.SetTitle("Alternative On")
		alternativeOnApp.SendNotification(fyne.NewNotification("Sessão restaurada!", "Sua sessão foi automaticamente restaurada. Para mudar isso clique em 'Sair'."))
		TokenData = *userToken
		TokenDataPrimitivo = *userData
		interfacePrincipal(alternativeOnWindow, alternativeOnApp)
	}
	alternativeOnWindow.Show()
	alternativeOnApp.Run()
}

func interfacePrincipal(janela fyne.Window, app fyne.App) {
	janela.Resize(fyne.NewSize(800, 600))
	janela.CenterOnScreen()

	links := pgo.ObterRecursos(TokenData.IdEscola, TokenData.Token, TokenData.TokenParceiro)

	//UI APOS O LOGIN
	/* Tab 1: Principal */
	labelTabHub := widget.NewRichTextFromMarkdown(textoHub)
	labelTabHub.Wrapping = fyne.TextWrapWord
	labelAccordionAtividades := widget.Label{
		Text: "Acesse suas atividades do Positivo On.",
	}
	botaoAccordionAtividades := widget.Button{
		Text:       "Ver atividades",
		Icon:       theme.NewThemedResource(resourceHistoryeduSvg),
		OnTapped:   func() { app.OpenURL(parseUrl(links.Studos)) },
		Importance: widget.HighImportance,
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
		Text: "Visualizar livros",
		Icon: theme.NewThemedResource(resourceBookSvg),
		OnTapped: func() {
			labelQuestão := widget.NewLabel("Você quer vizualizar seus livros, ou baixar todos eles?")
			labelQuestão.Wrapping = fyne.TextWrapWord

			questãoVerOuBaixar := dialog.NewCustomWithoutButtons("Ver/Baixar livros", labelQuestão, janela)
			btnVerTodos := widget.NewButton("Vizualizar livros", func() {
				interfaceLivros(janela, app)
				questãoVerOuBaixar.Hide()
			})
			btnBaixarTodos := widget.NewButton("Baixar livros", func() {
				interfaceBaixarTudo(app, janela)
				questãoVerOuBaixar.Hide()
			})
			questãoVerOuBaixar.SetButtons([]fyne.CanvasObject{btnVerTodos, btnBaixarTodos})
			questãoVerOuBaixar.Show()
		},
		Importance: widget.HighImportance,
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
		Text: "Leia as mensagens enviadas a você.",
	}
	botaoAccordionMensagens := widget.Button{
		Text:       "Ler mensagens",
		Icon:       theme.NewThemedResource(resourceChatSvg),
		OnTapped:   func() { app.OpenURL(parseUrl(links.Mensagens)) },
		Importance: widget.HighImportance,
	}

	containerAccordionMensagens := container.NewVBox(&labelAccordionMensagens, &botaoAccordionMensagens)
	accordionMensagens := widget.NewAccordion(
		&widget.AccordionItem{
			Title:  "Mensagens",
			Detail: containerAccordionMensagens,
			Open:   true,
		},
	)
	botaoAccordionLogout := widget.Button{
		Text: "Sair",
		Icon: theme.NewThemedResource(resourceLogoutSvg),
		OnTapped: func() {
			dialog.ShowConfirm("Você tem certeza?", "Você realmente quer sair do app?\nIsso também encerrará sua sessão.", func(b bool) {
				if b {
					app.Preferences().SetString("config_session", "")
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
	labelTabContaNomeUsuario := widget.NewLabel(fmt.Sprintf("Olá, %v!\n(ID: %v, ID legado: %v)", TokenDataPrimitivo.Nome, TokenData.IdUsuario, TokenDataPrimitivo.IdUsuarioEscola)) //TODO: Implementar userinfo no pgo
	labelTabContaNomeUsuario.Wrapping = fyne.TextWrapWord
	labelTabContaNomeEscola := widget.NewLabel("Escola: " + TokenData.NomeEscola + " (ID: " + TokenData.IdEscola + ")")
	labelTabContaNomeEscola.Wrapping = fyne.TextWrapWord
	botaoTabContaMudarSenha := widget.Button{
		Text:       "Mudar senha",
		Icon:       lockResetIcon,
		Importance: widget.HighImportance,
		OnTapped: func() {
			mudarSenhaAntigaLabel := widget.NewLabelWithStyle("Digite a antiga senha:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			mudarSenhaAntigaEntry := widget.NewPasswordEntry()
			mudarSenhaAntigaEntry.AcceptsTab()
			mudarSenhaNovaLabel := widget.NewLabelWithStyle("Digite a nova senha:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			mudarSenhaNovaEntry := widget.NewPasswordEntry()
			mudarSenhaNovaEntry.AcceptsTab()
			mudarSenhaNovaEntry.PlaceHolder = "Coloque uma senha forte!"
			mudarSenhaLayout := container.NewVBox(mudarSenhaAntigaLabel, mudarSenhaAntigaEntry, mudarSenhaNovaLabel, mudarSenhaNovaEntry)
			mudarSenhaDialog := dialog.NewCustomConfirm("Alterar senha - Alternative On", "Mudar", "Fechar", mudarSenhaLayout, func(b bool) {
				if b {
					resultado, err := pgo.AlterarSenha(mudarSenhaAntigaEntry.Text, mudarSenhaNovaEntry.Text, TokenData.Token)
					if err != nil {
						dialog.ShowError(err, janela)
						return
					}
					if resultado.Erro {
						dialog.ShowInformation("Sucesso!", resultado.Mensagem, janela)
						return
					}
				}
			}, janela)

			mudarSenhaDialog.Show()
		},
	}

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
		Text:       "Ir para a página do projeto",
		Icon:       theme.ComputerIcon(),
		OnTapped:   func() { app.OpenURL(parseUrl("https://github.com/AlternativeOn/AlternativeOn")) },
		Importance: widget.HighImportance,
	}
	botaoTabSobreConfig := widget.Button{
		Text:       "Configurações do app",
		Icon:       theme.SettingsIcon(),
		Importance: widget.MediumImportance,
	}
	botaoTabSobreConfig.OnTapped = func() {
		btnMudarTemaClaro := widget.NewButtonWithIcon("Tema claro", theme.NewInvertedThemedResource(resourceLightmodeSvg), func() {
			app.Settings().SetTheme(theme.LightTheme())
		})
		btnMudarTemaClaro.Importance = widget.HighImportance
		btnMudarTemaEscuro := widget.NewButtonWithIcon("Tema escuro", theme.NewInvertedThemedResource(resourceDarkmodeSvg), func() {
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
		container.NewTabItemWithIcon("Sobre", theme.NewThemedResource(resourceInfoSvg), conteudoTabSobre),
	)
	abasUiAposLogin.SetTabLocation(container.TabLocationTop)
	livros, err := pgo.ObterLivros(TokenData.Token)
	if err != nil {
		dialog.ShowError(err, janela)
		return
	}
	LivrosData = livros

	janela.SetContent(abasUiAposLogin)
}

func interfaceLivros(win fyne.Window, app fyne.App) {
	livros, err := pgo.ObterLivros(TokenData.Token)
	if err != nil {
		dialog.ShowError(err, win)
		return
	}
	LivrosData = livros

	/* BOTÕES DA INTERFACE DOS LIVROS */
	interfaceLivrosVoltarBtn := widget.NewButtonWithIcon("Voltar", theme.NavigateBackIcon(), func() { interfacePrincipal(win, app) })

	interfaceLivrosAjudaBtn := widget.NewButtonWithIcon("Ajuda", theme.HelpIcon(), func() {
		ajudaTexto := widget.NewLabel(textoLivrosAjuda)
		ajudaTexto.Wrapping = fyne.TextWrapWord
		ajuda := dialog.NewCustom("Ajuda dos livros - Alternative On", "Fechar", ajudaTexto, win)
		ajuda.Show()
	})

	livrosCont := widget.NewLabel(fmt.Sprintf("Livros: %v", len(livros)))
	livrosCont.Wrapping = fyne.TextWrapBreak

	btnBar := container.NewHBox(interfaceLivrosVoltarBtn, interfaceLivrosAjudaBtn)

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

		linkLivro := widget.NewHyperlink("Baixar", parseUrl(book.URL))
		linkLivro.OnTapped = func() {
			fmt.Println("tapped!", linkLivro.URL)
			downloadPdf(fmt.Sprint(linkLivro.URL), txtLivro.Text, win)
		}
		btnLivro := widget.NewButtonWithIcon("", theme.NewThemedResource(resourceDownloadsingleSvg), func() {
			downloadPdf(fmt.Sprint(linkLivro.URL), txtLivro.Text, win)
		})

		hAlign := container.NewHBox(btnLivro, txtLivro)
		hyperlinkContainer.Add(hAlign)
	}
	fmt.Println(len(livros))
	bookStartDialog.Hide()
	labelLivros.SetText("Baixe seus livros por aqui!")
}

func interfaceBaixarTudo(app fyne.App, win fyne.Window) {
	btnVoltarInterfaceBaixarTudo := widget.NewButtonWithIcon("Voltar ao início", theme.NavigateBackIcon(), func() {
		interfacePrincipal(win, app)
	})
	btnVoltarInterfaceBaixarTudo.Importance = widget.HighImportance

	labelInterfaceBaixarTudo := widget.NewLabel(textoInterfaceBaixarTudo)
	labelInterfaceBaixarTudo.Wrapping = fyne.TextWrapWord
	labelStatusInterfaceBaixarTudo := widget.NewLabel("Status: Esperando o usuário escolher uma pasta...")
	labelStatusInterfaceBaixarTudo.Wrapping = fyne.TextWrapWord

	status := func(s string) {
		labelStatusInterfaceBaixarTudo.SetText(s)
	}

	barraInterfaceBaixarTudo := widget.NewProgressBar()
	totalLivros := len(LivrosData)
	barraInterfaceBaixarTudo.Max = float64(totalLivros)
	barraInterfaceBaixarTudo.Min = 1.
	barraInterfaceBaixarTudo.SetValue(1)

	//var SalvarPara fyne.ListableURI

	btnIniciarDownload := widget.NewButtonWithIcon("Iniciar download", theme.NewThemedResource(resourceDownloadSvg), nil)
	btnIniciarDownload.Disable()
	dlgSalvarLivros := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		if lu == nil {
			return
		}

		dialog.ShowInformation("Sir", lu.Path(), win)
		btnIniciarDownload.Enable()
		app.Preferences().SetString("path", lu.Path())
		//SalvarPara = lu
	}, win)
	dlgSalvarLivros.SetConfirmText("Selecionar")
	dlgSalvarLivros.SetDismissText("Cancelar")

	btnEscolherPasta := widget.NewButtonWithIcon("Escolher pasta", theme.FolderOpenIcon(), func() {
		dlgSalvarLivros.Show()
	})
	btnEscolherPasta.Importance = widget.SuccessImportance

	btnIniciarDownload.OnTapped = func() {
		btnVoltarInterfaceBaixarTudo.Disable()
		btnIniciarDownload.Disable()
		btnEscolherPasta.Disable()
		downloadPath := app.Preferences().String("path")
		fmt.Println(downloadPath)

		status(fmt.Sprintf("Baixando livros... (0/%v)", totalLivros))

		g := got.New()

		nãoBaixados := make([]string, 0)

		for baixados, book := range LivrosData {
			caminho := fmt.Sprintf("%v/%v - %v - %v [%v].pdf", downloadPath, book.ComponenteCurricular, book.Volume, book.Serie, book.Tipo)

			err := g.Download(book.URL, caminho)
			if err != nil {
				dialog.ShowError(err, win)
				btnVoltarInterfaceBaixarTudo.Enable()
				btnIniciarDownload.Enable()
				btnEscolherPasta.Enable()
				status("Download falhou, Motivo: " + err.Error())
				break
			}

			config := model.NewDefaultConfiguration()
			config.DecodeAllStreams = true

			config.OwnerPW = "@rc0Tech"
			err = api.DecryptFile(caminho, caminho, config)
			if err != nil {
				fmt.Println(err)
				status("Falha ao processar arquivo: " + err.Error())
				nãoBaixados = append(nãoBaixados, filepath.Base(caminho))
			}

			barraInterfaceBaixarTudo.Value++
			barraInterfaceBaixarTudo.Refresh()
			status(fmt.Sprintf("Baixando livros... Seja paciente :)\nBaixados/Total: %v/%v", baixados, totalLivros))
		}

		status("Todos os livros baixados com sucesso.")
		btnVoltarInterfaceBaixarTudo.Enable()
		btnIniciarDownload.Enable()
		btnEscolherPasta.Enable()

		if len(nãoBaixados) > 0 {
			status(fmt.Sprintf("Não foi possível processar %v arquivos.\nDownload finalizado\nSe necessário, utilize a senha @rc0Tech (já copiada para a área de transferencia)", len(nãoBaixados)))
			win.Clipboard().SetContent("@rc0Tech")
		}
	}

	btnBar := container.NewHBox(btnVoltarInterfaceBaixarTudo, widget.NewSeparator())
	botoes := container.NewHBox(btnEscolherPasta, btnIniciarDownload)
	resto := container.NewVBox(labelInterfaceBaixarTudo, barraInterfaceBaixarTudo, widget.NewSeparator(), botoes, widget.NewSeparator(), labelStatusInterfaceBaixarTudo)
	tudo := container.NewBorder(btnBar, nil, nil, nil, resto)
	win.SetContent(tudo)
}

func downloadPdf(url string, nome string, win fyne.Window) {
	downloadDialog := dialog.NewCustomWithoutButtons("Baixando livro....", widget.NewProgressBarInfinite(), win)
	downloadDialog.Show()
	g := got.New()

	tmpFile, err := os.CreateTemp("", "downloaded_*.pdf")
	if err != nil {
		dialog.ShowError(err, win)
		downloadDialog.Hide()
		return
	}
	defer tmpFile.Close()

	err = g.Download(url, tmpFile.Name())
	if err != nil {
		dialog.ShowError(err, win)
		downloadDialog.Hide()
		return
	}

	config := model.NewDefaultConfiguration()
	config.OwnerPW = "@rc0Tech"
	err = api.DecryptFile(tmpFile.Name(), tmpFile.Name()+"out.pdf", config)
	if err != nil {
		fmt.Println(err)
		downloadDialog.Hide()
		dialog.ShowError(errors.New(err.Error()+"\nO Arquivo foi salvo com senha!"), win)
		saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if writer == nil {
				return
			}

			_ = os.Rename(tmpFile.Name(), tmpFile.Name()+"-senha=@rc0Tech.pdf")
			newSrc, _ := os.Open(tmpFile.Name() + "-senha=@rc0Tech.pdf")
			_, err = io.Copy(writer, newSrc)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			writer.Close()
			os.Remove(tmpFile.Name())
			os.Remove(tmpFile.Name() + "-senha=@rc0Tech.pdf")
			dialog.ShowInformation("Pronto", "O livro foi salvo com sucesso!", win)

		}, win)

		saveDialog.SetFileName(nome + "-senha=@rc0Tech.pdf")
		saveDialog.SetFilter(&storage.ExtensionFileFilter{Extensions: []string{".pdf"}})
		saveDialog.SetConfirmText("Salvar")
		saveDialog.SetDismissText("Fechar")
		saveDialog.Show()
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

	saveDialog.SetFileName(nome + ".pdf")
	saveDialog.SetFilter(&storage.ExtensionFileFilter{Extensions: []string{".pdf"}})
	saveDialog.SetConfirmText("Salvar")
	saveDialog.SetDismissText("Fechar")
	saveDialog.Show()
}
