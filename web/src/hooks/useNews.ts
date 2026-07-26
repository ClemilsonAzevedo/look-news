import { useCallback, useEffect, useRef, useState } from "react"
import { fetchNews, type NewsQuery } from "../api/news"
import type { Article } from "../types/article"

export type FetchStatus = "loading" | "success" | "error"

export interface UseNewsResult {
  articles: Article[]
  total: number
  status: FetchStatus
  error: string | null
  refetch: () => void
}

export function useNews(
    sources: string[] = [],
    query: NewsQuery = {},
): UseNewsResult {
  const [articles, setArticles] = useState<Article[]>([])
  const [total, setTotal] = useState(0)
  const [status, setStatus] = useState<FetchStatus>("loading")
  const [error, setError] = useState<string | null>(null)

  const sourcesRef = useRef(sources)
  const queryRef = useRef(query)

  sourcesRef.current = sources
  queryRef.current = query

  const load = useCallback(async (signal?: AbortSignal) => {
    setStatus("loading")
    setError(null)

    try {
      const result = await fetchNews(
          sourcesRef.current,
          queryRef.current,
          signal,
      )
      setArticles(result.articles)
      setTotal(result.total)
      setStatus("success")
    } catch (err: unknown) {
      if (signal?.aborted) return
      setError(err instanceof Error ? err.message : "Unknown error")
      setStatus("error")
    }
  }, [])

  const refetch = useCallback(() => {
    load()
  }, [load])

  useEffect(() => {
    const controller = new AbortController()
    load(controller.signal)
    return () => controller.abort()
  }, [JSON.stringify(sources), JSON.stringify(query), load])

  return { articles, total, status, error, refetch }
}