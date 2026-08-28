<p align="center">
  <img src="tray/assets/icon.png" alt="Look News" width="520" />
</p>

<h1 align="center">Look News</h1>

<p align="center">
  Agregador de notícias que acompanha várias fontes RSS/Atom/RDF, filtra por relevância usando IA e entrega tudo pronto — direto na sua bandeja do sistema.
</p>

<p align="center">
  <a href="https://github.com/clemilsonazevedo/look-news/releases/latest">
    <img src="https://img.shields.io/github/v/release/clemilsonazevedo/look-news" alt="Última versão">
  </a>
</p>

## Sumário

* [Instalação rápida](#instalação-rápida)
* [Como funciona](#como-funciona)
* [Fontes RSS para começar](#fontes-rss-para-começar)
* [Como encontrar fontes RSS](#como-encontrar-fontes-rss)
* [O que tem aqui](#o-que-tem-aqui)
* [Contribuindo](#contribuindo)
* [Licença](#licença)

## Instalação rápida

Sem precisar instalar nada além do que já vem no seu sistema — sem Git, sem Node, sem Go:

**macOS / Linux**

```bash
curl -fsSL https://raw.githubusercontent.com/clemilsonazevedo/look-news/main/install.sh | bash
```

**Windows** (PowerShell)

```powershell
irm https://raw.githubusercontent.com/clemilsonazevedo/look-news/main/install.ps1 | iex
```

O script busca a última versão publicada e instala sozinho.

## Como funciona

Depois de instalado, o Look News roda discretamente como um ícone na bandeja do sistema — sem janela aberta ocupando espaço. Clique no ícone para ver as notícias mais recentes, já filtradas por relevância.

<p align="center">
  <img src="tray/assets/tray-diagram.png" alt="Onde o app aparece no macOS e no Windows" width="700" />
</p>

> A ilustração acima mostra só o posicionamento — não é um screenshot real da interface. Veja mais detalhes na [documentação do tray](tray/README.md#onde-encontrar-o-app).

## Fontes RSS para começar

Você pode adicionar praticamente qualquer site que disponibilize um feed RSS, Atom ou RDF.

Para começar, aqui estão algumas fontes de tecnologia e programação:

| Fonte            | Feed                                      |
| ---------------- | ----------------------------------------- |
| **Huncoding**    | `https://huncoding.com/feed.xml`          |
| **TabNews**      | `https://www.tabnews.com.br/recentes/rss` |
| **Tecnoblog**    | `https://tecnoblog.net/feed/`             |
| **TechCrunch**   | `https://techcrunch.com/feed/`            |
| **SecurityWeek** | `https://www.securityweek.com/feed/`      |

Essas são apenas sugestões. O objetivo do Look News é permitir que você monte sua própria seleção de fontes de acordo com os assuntos que acompanha.

## Como encontrar fontes RSS

Existem várias formas de encontrar feeds RSS.

Muitos sites já possuem um feed e indicam isso no rodapé, na página de notícias ou através de um ícone/link de RSS. Também vale procurar pela palavra `RSS` ou `feed` no próprio site.

Outra opção é usar uma IA para encontrar feeds de um determinado assunto. Por exemplo, na **Perplexity**, você pode pedir:

> "Encontre fontes RSS sobre desenvolvimento em Go, inteligência artificial e startups. Retorne os links diretos dos feeds RSS."

Você também pode fazer o mesmo com qualquer tema de interesse — tecnologia, segurança, ciência, economia, programação, notícias locais etc.

Depois, basta adicionar os feeds encontrados ao Look News.

## O que tem aqui

Este é um monorepo com duas partes, cada uma com sua própria documentação:

| Parte        | O que faz                                                                  | Documentação                       |
| ------------ | -------------------------------------------------------------------------- | ---------------------------------- |
| 🖥️ **Tray** | App de bandeja em Electron — o que o usuário final abre e usa no dia a dia | [`tray/README.md`](tray/README.md) |
| ⚙️ **API**   | Backend em Go — busca, filtra por IA e serve as notícias                   | [`api/README.md`](api/README.md)   |

Só quer *usar* o app? A instalação rápida acima já é tudo que você precisa. Quer rodar ou desenvolver localmente? Entra na documentação da parte que te interessa.

## Contribuindo

1. Fork → branch a partir da `main` (`feat/...` ou `fix/...`).
2. Siga o guia de instalação da parte que você vai mexer: [tray](tray/README.md) ou [API](api/README.md).
3. Confirme que builda antes de abrir o PR.
4. Descreva claramente o que mudou. Mudança grande? Abre uma issue antes pra alinhar o approach.

Bugs e sugestões também são bem-vindos na aba de Issues.
