[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/alternativeon/alternativeon?style=flat)](https://go.dev/dl/)
[![Open Source Helpers](https://www.codetriage.com/alternativeon/alternativeon/badges/users.svg)](https://www.codetriage.com/alternativeon/alternativeon)
![GitHub](https://img.shields.io/github/license/alternativeon/alternativeon)

# Como contribuir para o Alternative On

Estamos muito felizes que você esteja lendo isto. Aqui tem alguns recursos úteis para você:

- [Download do Go](https://go.dev/dl/), a lingaguem que usamos é essa. Caso queira aprender sobre, [clique aqui](https://go.dev/learn).
- Nosso [roadmap](https://github.com/AlternativeOn/AlternativeOn/projects?type=beta) é aonde tem o que queremos implementar na aplicação.
- Bugs? Saiba como reportar-los [aqui](#Bugs)

## Testando

Por favor, teste o seu código antes de enviar uma pull request. [Saiba como criar testes](https://edsoncelio.github.io/posts/testes-unitarios-golang-/)

## Enviando sua Pull Request

Na sua Pull Request, sempre inclua uma explicação de o que você tenha mudado (e por que, se aplicavel). Por favor, também siga nosso [estilo de código](#estilo-de-código) e também tenha certeza que todas suas commits sejam atômicas (um recurso por commit).

### Você escreveu um patch que arruma um bug?

- Abra uma Pull Request com o patch
- Certifique-se de que a descrição da Pull Request descreva claramente o problema e a solução. Inclua o número do problema (issue) relevante, se aplicável.
- Antes de enviar, lembre-se de ler todo esse arquivo de como contribuir para o Alternative On.

### Você quer adicionar um novo recurso ou mudar um existente?

- Sugira sua mudança nas [discussões](https://github.com/AlternativeOn/AlternativeOn/discussions) e comece a escrever seu código
- Se receber algum feedback positivo, abra um issue sobre sua mudança.

### Você tem perguntas sobre o código-fonte?

- Sinta-se livre para perguntar [nas discussões](https://github.com/AlternativeOn/AlternativeOn/discussions).
  - Não abra um issue sobre isso. Issues são reservados para outras questões, como bugs.

### Você quer contribuir para a documentação do Alternative On?

- Artigos que recomendamos ler:
  - [Sintaxe do markdown do GitHub](https://docs.github.com/pt/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax)
  - [Como editar arquivos no GitHub](https://docs.github.com/pt/repositories/working-with-files/managing-files/editing-files#editing-files-in-another-users-repository)

###### _[Leia mais sobre as Pull Requests](http://help.github.com/pull-requests/)_

## Estilo de código

Para que o código fique legivel para todos que queiram ler-lo, nós seguimos alguns padrões, e eles são bem simples:

- Após terminar de escrever, use a ferramenta `go fmt`, assim ela deixa tudo com o estilo único e padrão.
- Lembre-se que isso é um software de código aberto. Considere que outras pessoas vão ler seu código, e deixe-o legivel para elas. É como cuidar de um irmão: talvez vocês briguem quando o seu responsável não está por perto, mas com ele perto vocês devem se comportar.

## Bugs

Encontrou um bug? É aqui que você vai saber como enviar ele para a gente!

- Encontrou uma vunerabilidade de segurança? **NÃO abra um issue sobre isso**, [leia este arquivo](./SECURITY.md) para saber o que fazer nesses casos.

- Tenha certeza que o bug já não tenha sido enviado anteriormente, você pode pesquisar isso nos [issues](https://github.com/AlternativeOn/AlternativeOn/issues).
  - Se não tiver encontrado algum issue com o seu bug, [abra um novo issue](https://github.com/AlternativeOn/AlternativeOn/issues/new) e certifique-se de incluir um título e uma descrição clara, o máximo possível de informações relevantes e uma amostra de código ou um caso de teste executável demonstrando o comportamento esperado que não está ocorrendo.
  - Recomendamos que você siga os templetes de issues.
  
