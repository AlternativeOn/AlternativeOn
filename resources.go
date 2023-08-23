package main

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed assets/Icon.png
var windowIcon []byte

//go:embed assets/lock_reset.svg
var lockReset []byte
var lockResetIconResource = fyne.NewStaticResource("Reset lock icon", lockReset)
var lockResetIcon = theme.NewThemedResource(lockResetIconResource)

var textoHub = `## Aqui você pode acessar suas atividades, livros e mensagens`
var textoFuncLivros = `Acesse seus livros digitais no formato de PDF. 

**Função ainda em desenvolvimento**, portanto ela pode não estar completa e/ou apresentar bugs.
A senha dos livros é copiada para sua área de transferencia.`
var textoTabConta = `## Aqui você pode controlar e alterar sua conta do Positivo.

Atualmente é possivel:
- Ver o nome e escola

Ainda não foi adicionado a possibilidade de mudança de senha pois o servidor está retonando o erro 500 para todas as mudanças de senha.`
var textoTabSobre = `## Alternative On - v0.1.0

Alternative On é um projeto que visa ser um cliente alternativo a plataforma de estudos Positivo On, contendo aplicativos próprios para Android, Linux e Windows.

**Veja mais sobre o projeto no página do Github.**



## Créditos:
- Fyne
- Postman
- Material Icons
- Yngrid, Luis Fernado
- Tamiris, Nicole, Maria Cecilia
- Todo mundo que contribui no desenvolvimento programacional (<3)
- E você, por usar o aplicativo!

---
> Aplicativo licenciado para você sobre a licença BSD-3. O icone do aplicativo está sobre a licença do Creative Commons BY-NC (CC BY-NC)

> Os outros icones dentro do aplicativo estão lincenciados sobre a licença do Fyne.
`

var textoPainelLogin = `# Alternative On
Seja bem-vindo! Para continuar faça seu login utilizando o mesmo usuário e senha do Positivo On.

Caso tenha esquecido sua senha, clique no botão "Recuperar senha".`

var textoPainelSessão = "Essa opção faz com que suas credenciais sejam salvadas para que na proxima vez que abrir o app não ser necessário inserir suas credenciais novamente"

var textoLivrosAjuda = "Baixe os seus livros digitais por aqui!\n-Baixar qualquer livro manualmente NÃO tera a senha removida automaticamente;\n\n-Baixar todos os livros remove a senha automaticamente;\n\n-A senha será copiada para a área de transferencia ao baixar manualmente."
