import type { Article } from "../types/article"
import { formatRelativeTime, stripHtml } from "../lib/format"
import { Icon } from "./Icon"

interface ArticleCardProps {
  article: Article
}

function StandardCard({ article }: { article: Article }) {
  const timeAgo = formatRelativeTime(article.date)
  const summary = stripHtml(article.summary)

  return (
    <a
      href={article.link}
      target="_blank"
      rel="noopener noreferrer"
      className="group flex cursor-pointer flex-col gap-2 rounded-xl border border-outline-variant/30 bg-surface-container-low p-5 transition-colors hover:border-outline-variant/80"
    >
      <div className="flex items-center justify-between font-label text-xs text-on-surface-variant">
        <span className="font-semibold tracking-widest text-tertiary uppercase">
          {article.source}
        </span>
        {timeAgo && <span>{timeAgo}</span>}
      </div>

      <h3 className="font-headline text-xl leading-snug text-on-background transition-colors group-hover:text-primary lg:text-2xl">
        {article.title}
      </h3>

      {summary && (
        <p className="line-clamp-2 font-body text-sm leading-relaxed text-on-surface-variant">
          {summary}
        </p>
      )}

      <div className="mt-auto flex flex-wrap gap-2 pt-4 font-label text-xs">
        {article.terms.map((term) => (
          <span
            key={term}
            className="rounded bg-surface-container-high px-2 py-1 text-on-surface-variant"
          >
            {term}
          </span>
        ))}
        {article.author && (
          <span className="ml-auto flex items-center gap-1 text-on-surface-variant">
            <Icon name="person" className="text-[14px]" /> {article.author}
          </span>
        )}
      </div>
    </a>
  )
}



export function ArticleCard({ article }: ArticleCardProps) {
  return (
    <StandardCard article={article} />
  )
}
