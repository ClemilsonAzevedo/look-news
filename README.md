<p align="center">
  <img src="tray/assets/icon.png" alt="Look News" width="512" />
</p>

<h1 align="center">Look News</h1>

<p align="center">
  Agregador de notícias que acompanha várias fontes RSS/Atom/RDF, filtra por relevância usando IA e entrega tudo pronto através de uma API — consumido por um app de bandeja (tray) em Electron.
</p>

## Sumário

- [Instalando a API](#instalando-a-api)
- [Instalando o front (tray)](#instalando-o-front-tray)
- [Usando a API](#usando-a-api)
- [Documentação interativa (Swagger)](#documentação-interativa-swagger)
- [Exemplos de chamada](#exemplos-de-chamada)
- [Ambientes](#ambientes)

## Instalando a API

Clone o repositório e instale as dependências do Go:

```bash
git clone https://github.com/clemilsonazevedo/look-news.git
cd look-news
cd api
make download
```

Crie um arquivo `.env` na raiz com sua chave da Groq:

```bash
<!-- TODO: confirmar o nome exato da variável usada para a chave da Groq -->
GROQ_API_KEY=sua-chave-aqui
```

Suba o servidor:

```bash
<!-- TODO: ajustar para o comando/caminho real do entrypoint -->
make run
```

Por padrão, a API sobe em `http://localhost:8080`.

## Instalando o front (tray)

O front é um app de bandeja feito em Electron — fica rodando discretamente no sistema, buscando e mostrando as notícias sem precisar de uma janela aberta o tempo todo.

```bash
cd tray
npm install
npm start
```

Isso abre o app em modo de desenvolvimento, já com o ícone na bandeja do sistema. Por padrão, ele consome a API já hospedada em produção — não precisa da API local rodando pra testar.

Se quiser apontar pro backend local em vez do hospedado, crie um `.env` dentro de `tray/`:

```bash
LOOK_NEWS_BACKEND_URL=http://localhost:8080
```

## Usando a API

### `GET /news`

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

## Documentação interativa (Swagger)

A API gera sua própria documentação OpenAPI automaticamente. Com o servidor rodando, acesse:

| Ambiente | Swagger UI |
|---|---|
| Local | `http://localhost:8080/swagger` |
| Produção | `https://look-news-api.onrender.com/swagger` |

Lá dá pra ver todas as rotas, parâmetros e testar chamadas direto pelo navegador.

## Exemplos de chamada

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

## Ambientes

| Ambiente | URL |
|---|---|
| Local | `http://localhost:8080` |
| Produção | `https://look-news-api.onrender.com` |

> A API de produção roda no plano gratuito do Render, que "dorme" após um tempo sem uso — a primeira requisição depois disso pode demorar mais que o normal pra responder.
