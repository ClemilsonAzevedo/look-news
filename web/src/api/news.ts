import type { Article } from "../types/article"

const BASE_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080"

export type SortOrder = "asc" | "desc"

export interface NewsQuery {
  term?: string
  source?: string
  since?: string
  sort?: SortOrder
  limit?: number
}

export interface NewsResult {
  articles: Article[]
  total: number
}

function buildQuery(query: NewsQuery): string {
  const params = new URLSearchParams()
  if (query.term) params.set("term", query.term)
  if (query.source) params.set("source", query.source)
  if (query.since) params.set("since", query.since)
  if (query.sort) params.set("sort", query.sort)
  if (query.limit !== undefined) params.set("limit", String(query.limit))
  const qs = params.toString()
  return qs ? `?${qs}` : ""
}

export async function fetchNews(
  query: NewsQuery = {},
  signal?: AbortSignal,
): Promise<NewsResult> {
  const res = await fetch(`${BASE_URL}/news${buildQuery(query)}`, { signal })
  if (!res.ok) {
    throw new Error(`Failed to load news (${res.status})`)
  }

  const articles = (await res.json()) as Article[]
  const totalHeader = res.headers.get("X-Total-Count")
  const total = totalHeader ? Number(totalHeader) : articles.length

  return { articles, total: Number.isNaN(total) ? articles.length : total }
}

// export async function fetchHealth(signal?: AbortSignal): Promise<Health> {
//   const res = await fetch(`${BASE_URL}/health`, { signal })
//   if (!res.ok) {
//     throw new Error(`Health check failed (${res.status})`)
//   }
//   return (await res.json()) as Health
// }
