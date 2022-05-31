package main

import (
	"github.com/alternativeon/pgo"
	"github.com/gonutz/w32/v2"
	"github.com/gonutz/wui/v2"
	llgg "github.com/rs/zerolog/log"
)

func mostrarUiAposLogin() {
	windowFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	window := wui.NewWindow()
	window.SetFont(windowFont)
	window.SetTitle("Alternative On")
	window.SetHasMaxButton(false)
	window.SetResizable(false)

	txtEstudanteFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: 18,
	})

	txtEstudante := wui.NewLabel()
	txtEstudante.SetFont(txtEstudanteFont)
	txtEstudante.SetAnchors(wui.AnchorCenter, wui.AnchorCenter)
	txtEstudante.SetBounds(230, 3, 150, 20)
	txtEstudante.SetText("Olá, " + Username + "!")
	txtEstudante.SetAlignment(wui.AlignCenter)
	window.Add(txtEstudante)

	boxAvaliacoesFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	boxAvaliacoes := wui.NewPanel()
	boxAvaliacoes.SetFont(boxAvaliacoesFont)
	boxAvaliacoes.SetBounds(10, 51, 570, 110)
	boxAvaliacoes.SetBorderStyle(wui.PanelBorderSunken)
	window.Add(boxAvaliacoes)

	btnAvaliacoes := wui.NewButton()
	btnAvaliacoes.SetBounds(5, 6, 560, 25)
	btnAvaliacoes.SetText("Acessar suas atividades")
	btnAvaliacoes.SetOnClick(func() {
		err := w32.ShellExecute(0, "open", pgo.GetHomework(Usertoken), "", "", w32.SW_SHOW)
		if err != nil {
			llgg.Error().Str("Erro", err.Error()).Msg("Não foi possivel abrir o link!")
		}
	})

	boxAvaliacoes.Add(btnAvaliacoes)

	txtAvaliacoes2 := wui.NewLabel()
	txtAvaliacoes2.SetBounds(5, 35, 555, 15)
	txtAvaliacoes2.SetText("Acesse suas atividades, trilhas, avaliações, treinos clicando no botão acima.")
	boxAvaliacoes.Add(txtAvaliacoes2)

	txtAvaliacoes3 := wui.NewLabel()
	txtAvaliacoes3.SetBounds(5, 74, 521, 15)
	txtAvaliacoes3.SetText("A página será aberta de acordo com as configurações da aplicação. ")
	boxAvaliacoes.Add(txtAvaliacoes3)

	txtAvaliacoes := wui.NewLabel()
	txtAvaliacoes.SetBounds(10, 36, 194, 15)
	txtAvaliacoes.SetText("Trabalhos, avaliações e atividades")
	window.Add(txtAvaliacoes)

	txtLivros := wui.NewLabel()
	txtLivros.SetBounds(10, 184, 150, 15)
	txtLivros.SetText("Conteúdo didático")
	window.Add(txtLivros)

	boxLivrosFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	boxLivros := wui.NewPanel()
	boxLivros.SetFont(boxLivrosFont)
	boxLivros.SetBounds(10, 201, 570, 70)
	boxLivros.SetBorderStyle(wui.PanelBorderSunken)
	window.Add(boxLivros)

	txtLivros2 := wui.NewLabel()
	txtLivros2.SetBounds(5, 33, 531, 14)
	txtLivros2.SetText("Clique no botão acima para ver seus livros digitais em PDF. Você também pode salvar-los no seu computador.")
	boxLivros.Add(txtLivros2)

	btnLivros := wui.NewButton()
	btnLivros.SetEnabled(false)
	btnLivros.SetBounds(5, 6, 560, 25)
	btnLivros.SetText("Ver seus livros digitais")
	boxLivros.Add(btnLivros)

	txtLivrosNAOPRONTOFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -12,
		Bold:   true,
	})

	txtLivrosNAOPRONTO := wui.NewLabel()
	txtLivrosNAOPRONTO.SetFont(txtLivrosNAOPRONTOFont)
	txtLivrosNAOPRONTO.SetBounds(200, 52, 200, 15)
	txtLivrosNAOPRONTO.SetText("Função ainda não implementada.")
	txtLivrosNAOPRONTO.SetAlignment(wui.AlignCenter)
	boxLivros.Add(txtLivrosNAOPRONTO)

	txtMsgs := wui.NewLabel()
	txtMsgs.SetBounds(10, 285, 150, 15)
	txtMsgs.SetText("Mensagens da escola")
	window.Add(txtMsgs)

	boxMsgFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -11,
	})

	boxMsg := wui.NewPanel()
	boxMsg.SetFont(boxMsgFont)
	boxMsg.SetBounds(10, 302, 570, 85)
	boxMsg.SetBorderStyle(wui.PanelBorderSunken)
	window.Add(boxMsg)

	btnMsg := wui.NewButton()
	btnMsg.SetBounds(5, 6, 560, 25)
	btnMsg.SetText("Visualizar novas mensagens")
	btnMsg.SetOnClick(func() {
		err := w32.ShellExecute(0, "open", pgo.GetMessages(Usertoken), "", "", w32.SW_SHOW)
		if err != nil {
			llgg.Error().Str("Erro", err.Error()).Msg("Não foi possivel abrir o link!")
		}
	})
	boxMsg.Add(btnMsg)

	txtMsgs2 := wui.NewLabel()
	txtMsgs2.SetBounds(5, 40, 553, 15)
	txtMsgs2.SetText("Aqui, você consegue ver as mensagens enviadas pelos professores e direção da escola, apenas clicando no botão.")
	boxMsg.Add(txtMsgs2)

	txtMsgs3 := wui.NewLabel()
	txtMsgs3.SetBounds(5, 59, 543, 15)
	txtMsgs3.SetText("A página será aberta como foi selecionado nas configurações do programa.")
	boxMsg.Add(txtMsgs3)

	txtVersaoFont, _ := wui.NewFont(wui.FontDesc{
		Name:   "Tahoma",
		Height: -10,
	})

	txtVersao := wui.NewLabel()
	txtVersao.SetFont(txtVersaoFont)
	txtVersao.SetBounds(10, 387, 200, 12)
	txtVersao.SetText("Alternative On - Beta (versão " + version + ")")
	window.Add(txtVersao)

	window.Show()
}
