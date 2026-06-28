const ZERO_DATE = "0001-01-01T00:00:00Z"

export function hasValidDate(date: string): boolean {
  if (!date || date === ZERO_DATE) return false
  return !Number.isNaN(Date.parse(date))
}

export function formatRelativeTime(date: string): string | null {
  if (!hasValidDate(date)) return null

  const diffMs = Date.now() - Date.parse(date)
  const diffSec = Math.max(0, Math.floor(diffMs / 1000))

  const minutes = Math.floor(diffSec / 60)
  if (minutes < 1) return "just now"
  if (minutes < 60) return `${minutes}m ago`

  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}h ago`

  const days = Math.floor(hours / 24)
  if (days < 7) return `${days}d ago`

  const weeks = Math.floor(days / 7)
  if (weeks < 5) return `${weeks}w ago`

  return new Date(date).toLocaleDateString(undefined, {
    day: "numeric",
    month: "short",
    year: "numeric",
  })
}

export function stripHtml(html: string): string {
  if (!html) return ""
  const doc = new DOMParser().parseFromString(html, "text/html")
  return (doc.body.textContent ?? "").replace(/\s+/g, " ").trim()
}
