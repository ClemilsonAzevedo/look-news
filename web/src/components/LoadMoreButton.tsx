interface LoadMoreButtonProps {
  onClick?: () => void
  label?: string
}

export function LoadMoreButton({
  onClick,
  label = "Load More Articles",
}: LoadMoreButtonProps) {
  return (
    <div className="mt-8 flex justify-center">
      <button
        type="button"
        onClick={onClick}
        className="rounded-lg border border-outline-variant/60 px-8 py-3 font-label font-medium text-primary transition-colors duration-300 hover:bg-surface-container focus:ring-2 focus:ring-primary/50 focus:outline-none"
      >
        {label}
      </button>
    </div>
  )
}
