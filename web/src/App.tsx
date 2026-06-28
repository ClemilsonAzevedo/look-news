import { useState } from "react"
import { FeedHeader } from "./components/FeedHeader"
import { type TimeRange } from "./types/feed"
import { LoadMoreButton } from "./components/LoadMoreButton"
import { BottomNavBar } from "./components/BottomNavBar"
import {useNews} from "./hooks/useNews.ts";
import {ArticleCard} from "./components/ArticleCard.tsx";

export function App() {
  const [activeRange, setActiveRange] = useState<TimeRange>("Last 30m")
  const [activeTab, setActiveTab] = useState("feed")
    const result = useNews({})

  return (
    <div className="flex min-h-screen flex-col bg-background font-body text-on-background selection:bg-primary-container selection:text-on-primary-container">
      <main className="mx-auto flex w-full max-w-5xl grow flex-col gap-12 px-4 pt-12 pb-32 md:px-8">
        <FeedHeader
          title="Global Feed"
          subtitle="Curated insights from across the network. Updated in real-time."
          activeRange={activeRange}
          onRangeChange={setActiveRange}
        />

        {/*<div className="grid grid-cols-1 gap-8 md:grid-cols-1 lg:gap-12">*/}
        <div className="grid grid-cols-1 gap-8">
          {result.articles.map((article) => (
            <ArticleCard key={article.link} article={article} />
          ))}
        </div>

        <LoadMoreButton />
      </main>

      <BottomNavBar activeTab={activeTab} onTabChange={setActiveTab} />
    </div>
  )
}
