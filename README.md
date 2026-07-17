# Look-News

Agregador de notícias em tempo real. Consome feeds RSS/Atom/RDF de múltiplas fontes, expõe uma API HTTP simples e serve um leitor editorial com filtros interactivos.

---

## O que é

Look-News é um servidor Go que agrega feeds de notícias tech periodicamente e os serve através de uma API REST. O frontend consome essa API directamente — sem frameworks, sem build step — e apresenta os artigos num leitor com filtros por fonte, período e ordenação.

**Fontes actuais:**
- [The Hacker News](https://thehackernews.com)
- [TabNews](https://www.tabnews.com.br)
- ...

---

## Stack

| Camada | Tecnologia              |
|--------|-------------------------|
| Servidor | Go (net/http)           |
| Parsing de feeds | RSS / Atom / RDF        |
| Cache | In-memory, TTL 24h      |
| Frontend | Vite + React + Tailwind |

---

## Arrancar

```bash
git clone https://github.com/teu_usuario/look-news
cd look-news

# Para o frontend (Precisa do Node.js e do Bun instalado)
cd web
bun install
bun dev

# Para o Backend (Precisa do Golang instalado)
cd api
go mod tidy
make run
```

A página web sobe em `http://localhost:5173`.

O servidor sobe em `http://localhost:8080`.

Os feeds são agregados automaticamente na inicialização e actualizados a cada hora.

---

## API

### `GET /news`

Devolve um array de artigos. Todos os parâmetros são opcionais.

| Parâmetro | Tipo | Descrição |
|-----------|------|-----------|
| `source` | string | Filtrar por fonte (case-insensitive, exacta). Ex: `TabNews` |
| `term` | string | Filtrar por tag/categoria (case-insensitive, exacta). Ex: `golang` |
| `since` | duração | Artigos mais recentes que X. Aceita `30m`, `2h`, `7d` |
| `sort` | string | `asc` para mais antigos primeiro. Por omissão: `desc` |
| `limit` | int | Máximo de resultados (por omissão: `50`) |

O header `X-Total-Count` da resposta contém o total de artigos após filtros, antes do `limit`.

**Formato de cada artigo:**

```json
{
  "title": "Título do artigo",
  "summary": "Resumo (pode conter HTML do feed)",
  "link": "https://fonte.com/artigo",
  "date": "2026-06-25T14:30:00Z",
  "source": "TabNews",
  "author": "Nome do autor",
  "published": "Wed, 25 Jun 2026 14:30:00 +0000",
  "terms": ["tag1", "tag2"]
}
```

> `summary` pode conter HTML — sanitizar antes de renderizar.  
> `date` pode ser `0001-01-01T00:00:00Z` se o feed não tiver data parseável.  
> `terms` pode ser um array vazio dependendo do feed.

---

### `GET /health`

```json
{
  "status": "ok",
  "articles": 142,
  "newest": "2026-06-25T14:30:00Z"
}
```

---

## Exemplos

```bash
# Últimas 20 notícias
curl "http://localhost:8080/news?limit=20"

# Notícias do TabNews das últimas 24h
curl "http://localhost:8080/news?source=TabNews&since=24h"

# Total disponível (via header)
curl -I "http://localhost:8080/news?term=golang" | grep X-Total-Count

# Health check
curl "http://localhost:8080/health"
```

---

## Notas

- CORS aberto para qualquer origem (`Access-Control-Allow-Origin: *`)
- Cache actualizado de hora em hora, expiração após 24h
- Filtros `source` e `term` são correspondência exacta — não suportam pesquisa parcial
- Sem autenticação
- Porta e fontes configuradas directamente no código
