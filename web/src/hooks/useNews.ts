import { useCallback, useEffect, useState } from "react"
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

export function useNews(query: NewsQuery): UseNewsResult {
  const [articles, setArticles] = useState<Article[]>([])
  const [total, setTotal] = useState(0)
  const [status, setStatus] = useState<FetchStatus>("loading")
  const [error, setError] = useState<string | null>(null)
  const [reloadToken, setReloadToken] = useState(0)

  const refetch = useCallback(() => setReloadToken((n) => n + 1), [])

  const queryKey = JSON.stringify(query)

  useEffect(() => {
    const controller = new AbortController()
    setStatus("loading")
    setError(null)

    fetchNews(JSON.parse(queryKey) as NewsQuery, controller.signal)
      .then((result) => {
        setArticles(result.articles)
        setTotal(result.total)
        setStatus("success")
      })
      .catch((err: unknown) => {
        if (controller.signal.aborted) return
        setError(err instanceof Error ? err.message : "Unknown error")
        setStatus("error")
      })

    return () => controller.abort()
  }, [queryKey, reloadToken])

  return { articles, total, status, error, refetch }
}
