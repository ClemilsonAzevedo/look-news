import { useState } from "react"
import { FeedHeader } from "./components/FeedHeader"
import { ArticleCard } from "./components/ArticleCard"
import { BottomNavBar } from "./components/BottomNavBar"
import { SourcesModal } from "./components/SourcesModal"
import { useNews } from "./hooks/useNews"
import { useSources } from "./hooks/useSources"

export function App() {
    const [activeTab, setActiveTab] = useState("feed")
    const [showSources, setShowSources] = useState(false)

    const { sources, add, remove } = useSources()
    const { articles, status, error, refetch } = useNews(sources, { limit: 50 })

    return (
        <div className="flex min-h-screen flex-col bg-background font-body text-on-background selection:bg-primary-container selection:text-on-primary-container">
            <main className="mx-auto flex w-full max-w-5xl grow flex-col gap-12 px-4 pt-12 pb-32 md:px-8">
                <FeedHeader
                    title="Global Feed"
                    subtitle="Curated insights from across the network. Updated in real-time."
                    onManageSources={() => setShowSources(true)}
                    onRefresh={refetch}
                    isLoading={status === "loading"}
                />

                {status === "error" && (
                    <p className="text-center text-sm text-error">{error}</p>
                )}

                {status === "loading" && articles.length === 0 && (
                    <p className="text-center text-sm text-on-surface-variant">Carregando...</p>
                )}

                {articles.length === 0 && status === "success" && (
                    <p className="text-center text-sm text-on-surface-variant">
                        Nenhuma notícia. Adicione fontes para começar.
                    </p>
                )}

                <div className="grid grid-cols-1 gap-8">
                    {articles.map((article) => (
                        <ArticleCard key={article.link} article={article} />
                    ))}
                </div>
            </main>

            <BottomNavBar activeTab={activeTab} onTabChange={setActiveTab} />

            {showSources && (
                <SourcesModal
                    sources={sources}
                    onAdd={add}
                    onRemove={remove}
                    onClose={() => setShowSources(false)}
                />
            )}
        </div>
    )
}