import {Icon} from "./Icon"

interface FeedHeaderProps {
    title: string
    subtitle: string
    onManageSources: () => void
    onRefresh: () => void
    isLoading?: boolean
}

export function FeedHeader({
   title,
   subtitle,
   onManageSources,
   onRefresh,
   isLoading = false,
}: FeedHeaderProps) {
    return (
        <header
            className="flex flex-col justify-between gap-6 border-b border-outline-variant/40 pb-6 md:flex-row md:items-end">
            <div>
                <h2 className="mb-2 font-headline text-5xl leading-tight tracking-tight text-on-background md:text-6xl">
                    {title}
                </h2>
                <p className="max-w-xl font-body text-sm text-on-surface-variant md:text-base">
                    {subtitle}
                </p>
            </div>

            <div className="flex items-center gap-3">
                <button
                    type="button"
                    onClick={onRefresh}
                    disabled={isLoading}
                    className="flex items-center gap-2 rounded-full border border-outline-variant/60 px-4 py-1.5 text-sm transition-colors hover:border-outline disabled:opacity-50"
                >
                    <Icon
                        name="refresh"
                        className={`text-sm ${isLoading ? "animate-spin" : ""}`}
                    />
                    Atualizar
                </button>

                <button
                    type="button"
                    onClick={onManageSources}
                    className="flex items-center gap-2 rounded-full bg-primary px-4 py-1.5 text-sm font-medium text-on-primary transition-opacity hover:opacity-90"
                >
                    <Icon name="link" className="text-sm"/>
                    Fontes
                </button>
            </div>
        </header>
    )
}