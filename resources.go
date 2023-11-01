package main

import (
	_ "embed"
)

//go:embed assets/Icon.png
var windowIcon []byte

var textoHub = `## Aqui você pode acessar suas atividades, livros e mensagens`
var textoFuncLivros = `Baixe seus livros digitais no formato de PDF. É possível baixar todos os livros ou escolher livros individuais.`
var textoTabConta = `## Aqui você pode controlar e alterar sua conta do Positivo.`
var textoTabSobre = `## Alternative On - v1.0.1

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

var textoInterfaceBaixarTudo = "Primeiro selecione uma pasta para baixar os livros, depois clique em iniciar."
