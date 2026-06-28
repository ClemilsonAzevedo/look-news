import { Icon } from "./Icon"
import { TIME_RANGES, type TimeRange } from "../types/feed"

interface FeedHeaderProps {
  title: string
  subtitle: string
  activeRange: TimeRange
  onRangeChange: (range: TimeRange) => void
  sortLabel?: string
}

export function FeedHeader({
  title,
  subtitle,
  activeRange,
  onRangeChange,
  sortLabel = "Newest",
}: FeedHeaderProps) {
  return (
    <header className="flex flex-col justify-between gap-6 border-b border-outline-variant/40 pb-6 md:flex-row md:items-end">
      <div>
        <h2 className="mb-2 font-headline text-5xl leading-tight tracking-tight text-on-background md:text-6xl">
          {title}
        </h2>
        <p className="max-w-xl font-body text-sm text-on-surface-variant md:text-base">
          {subtitle}
        </p>
      </div>

      <div className="flex flex-wrap items-center gap-4 font-label text-sm">
        <div className="flex items-center rounded-full bg-surface-container-high p-1">
          {TIME_RANGES.map((range) => {
            const isActive = range === activeRange
            return (
              <button
                key={range}
                type="button"
                onClick={() => onRangeChange(range)}
                className={
                  "rounded-full px-3 py-1 transition-all " +
                  (isActive
                    ? "bg-primary text-on-primary shadow-sm"
                    : "text-on-surface-variant hover:text-on-surface")
                }
              >
                {range}
              </button>
            )
          })}
        </div>

        <button
          type="button"
          className="group flex items-center gap-2 rounded-full border border-outline-variant/60 px-4 py-1.5 transition-colors hover:border-outline"
        >
          <span className="text-on-surface-variant">Sort:</span>
          <span className="font-medium text-on-surface">{sortLabel}</span>
          <Icon name="expand_more" className="text-sm text-on-surface-variant" />
        </button>
      </div>
    </header>
  )
}
