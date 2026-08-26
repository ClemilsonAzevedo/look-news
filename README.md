<p align="center">
  <img src="tray/assets/icon.png" alt="Look News" width="512" />
</p>

<h1 align="center">Look News</h1>

<p align="center">
  Agregador de notícias que acompanha várias fontes RSS/Atom/RDF, filtra por relevância usando IA e entrega tudo pronto através de uma API — consumido por um app de bandeja (tray) em Electron.
</p>

## Sumário

- [Instalação rápida](#instalação-rápida-recomendado)
- [Instalação manual](#instalação-manual)
- [Dependências](#dependências)
- [Instalando o front (tray)](#instalando-o-front-tray)
- [Instalando a API](#instalando-a-api)
- [Usando a API](#usando-a-api)
- [Documentação interativa (Swagger)](#documentação-interativa-swagger)
- [Exemplos de chamada](#exemplos-de-chamada)
- [Ambientes](#ambientes)
- [Contribuindo](#contribuindo)

## Instalação rápida (recomendado)

Baixa o binário pronto da última release — sem precisar de dependências.

1. macOS e Linux:
```Bash
curl -fsSL https://raw.githubusercontent.com/clemilsonazevedo/look-news/main/install.sh | bash
```

2. Windows (PowerShell):
```Bash
irm https://raw.githubusercontent.com/clemilsonazevedo/look-news/main/install.ps1 | iex
```
O script consulta a API do GitHub, baixa o asset correto e instala o app. 

Basta executar e acompanhar a instalação.

## Instalação manual

### Dependências

| Ferramenta | Pra quê | Download |
|---|---|---|
| Git | Clonar o repositório | [git-scm.com/downloads](https://git-scm.com/downloads) |
| Go | Rodar a API (versão exata em `api/go.mod`) | [go.dev/dl](https://go.dev/dl/) |
| make | Rodar os comandos do Makefile | [gnu.org/software/make](https://www.gnu.org/software/make/) |
| Node.js + npm | Rodar o front (versão em `tray/package.json`) | [nodejs.org/en/download](https://nodejs.org/en/download) |
| Chave de API da Groq | Filtro de relevância por IA | [console.groq.com/keys](https://console.groq.com/keys) |

> `make` no macOS já vem com as Command Line Tools do Xcode. No Windows, use via WSL ou Git Bash.

### Instalando o front (tray)

O front é um app de bandeja em Electron: fica rodando discretamente, buscando e mostrando notícias sem precisar de janela aberta.

1. Entre na pasta e instale as dependências:

   ```bash
   cd tray
   npm install
   ```

2. Rode em modo de desenvolvimento:

   ```bash
   npm start
   ```

O app abre com o ícone já na bandeja do sistema. Por padrão, ele consome a API hospedada em produção — não precisa da API local rodando pra testar.

Pra apontar pro backend local em vez do hospedado, crie um `.env` dentro de `tray/`:

```bash
LOOK_NEWS_BACKEND_URL=http://localhost:8080
```

#### Gerando o executável

Pra empacotar o app num instalador nativo:

```bash
cd tray
npm run make
```

O [Electron Forge](https://www.electronforge.io/) compila e empacota o app em `tray/out/make/`.

> O comando gera o instalador só pra plataforma onde ele é rodado (mac → `.dmg`, Windows → `.exe`, Linux → `.deb`/`.rpm`). Não há cross-compilação automática. Se já tinha um build antigo, rode `rm -rf out` antes pra garantir que compila do zero.

### Instalando a API

1. Clone o repositório e entre na pasta da API:

   ```bash
   git clone https://github.com/clemilsonazevedo/look-news.git
   cd look-news/api
   ```

2. Instale as dependências:

   ```bash
   make download
   ```

3. Crie um `.env` com sua chave da Groq:

   ```bash
   GROQ_API_KEY=sua-chave-aqui
   ```

4. Suba o servidor:

   ```bash
   make run
   ```

Pronto — a API sobe em `http://localhost:8080`.


### Usando a API

#### `GET /news`

| Query param | Tipo | Obrigatório |
|---|---|---|
| `criterion` | string — critério de relevância pro filtro | sim |
| `sources` | array de strings — uma ou mais URLs de feed (`?sources=A&sources=B`) | sim |

**Resposta — `200 OK`**

```json
{
  "articles": [
    {
      "title": "string",
      "summary": "string",
      "link": "string",
      "date": "2026-08-20T12:00:00Z",
      "source": "string",
      "author": "string",
      "published": "string",
      "terms": ["string"]
    }
  ],
  "total": 0
}
```

**Outros status possíveis**

| Status | Quando acontece |
|---|---|
| `304 Not Modified` | O `ETag` enviado em `If-None-Match` bate com o atual — corpo vazio |
| `400 Bad Request` | Nenhuma fonte (`sources`) foi enviada |
| `429 Too Many Requests` | Limite de requisições por IP excedido (5 req/s, rajada de 15) |
| `503 Service Unavailable` | Todas as fontes pedidas falharam e não há cache disponível |

A API suporta cache HTTP padrão via `ETag` + `Cache-Control: max-age=60, stale-while-revalidate=300` — uma requisição condicional com `If-None-Match` evita baixar a resposta inteira quando nada mudou.

### Documentação interativa (Swagger)

A API gera sua própria documentação OpenAPI automaticamente:

| Ambiente | Swagger UI |
|---|---|
| Local | `http://localhost:8080/swagger` |
| Produção | `https://look-news-api.onrender.com/swagger` |

Lá dá pra ver todas as rotas, parâmetros e testar chamadas direto pelo navegador.

### Exemplos de chamada

**Uma única fonte:**

```bash
curl "http://localhost:8080/news?criterion=tecnologia&sources=https://exemplo.com/feed.xml"
```

**Múltiplas fontes de uma vez:**

```bash
curl "http://localhost:8080/news?criterion=ciência+e+tecnologia&sources=https://exemplo.com/feed.xml&sources=https://outroexemplo.com/rss"
```

**Requisição condicional (revalida sem baixar tudo de novo):**

```bash
# primeira chamada — guarda o ETag da resposta
curl -i "http://localhost:8080/news?criterion=tecnologia&sources=https://exemplo.com/feed.xml"
# -> ETag: "a1b2c3..."

# chamada seguinte — usa o ETag guardado
curl -i "http://localhost:8080/news?criterion=tecnologia&sources=https://exemplo.com/feed.xml" \
  -H 'If-None-Match: "a1b2c3..."'
# -> 304 Not Modified, sem corpo, se nada mudou
```

**Requisição inválida (sem fontes):**

```bash
curl -i "http://localhost:8080/news?criterion=tecnologia"
# -> 400 Bad Request
```

### Ambientes

| Ambiente | URL |
|---|---|
| Local | `http://localhost:8080` |
| Produção | `https://look-news-api.onrender.com` |

> A API de produção roda no plano gratuito do Render, que "dorme" após um tempo sem uso — a primeira requisição depois disso pode demorar mais que o normal pra responder.

### Contribuindo

1. Faça um fork do repositório.
2. Crie uma branch a partir da `main`: `git checkout -b feat/nome-da-feature` ou `fix/nome-do-bug`.
3. Instale as [dependências](#dependências) e siga os passos de instalação da [API](#instalando-a-api) e do [front](#instalando-o-front-tray).
4. Antes de abrir o PR, confirme que tudo builda:
   - API: `go build ./...` (dentro de `api/`)
   - Front: `npm start` ou `npm run make` (dentro de `tray/`)
5. Abra o Pull Request descrevendo o que mudou e por quê.
6. Mudanças grandes? Abra uma issue antes pra alinhar o approach.

Bugs e sugestões também são bem-vindos na aba de Issues.
