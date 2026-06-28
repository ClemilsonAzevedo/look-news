import { Icon } from "./Icon"

interface NavTab {
  id: string
  label: string
  icon: string
}

const TABS: NavTab[] = [
  { id: "feed", label: "Feed", icon: "newspaper" },
  { id: "discover", label: "Discover", icon: "explore" },
  { id: "saved", label: "Saved", icon: "bookmark" },
]

interface BottomNavBarProps {
  activeTab?: string
  onTabChange?: (id: string) => void
}

export function BottomNavBar({
  activeTab = "feed",
  onTabChange,
}: BottomNavBarProps) {
  return (
    <nav className="fixed bottom-0 left-0 z-50 flex w-full items-center justify-around border-t border-outline-variant/60 bg-surface/80 px-4 py-3 backdrop-blur-md md:hidden">
      {TABS.map((tab) => {
        const isActive = tab.id === activeTab
        return (
          <button
            key={tab.id}
            type="button"
            onClick={() => onTabChange?.(tab.id)}
            className={
              "flex flex-col items-center justify-center rounded-xl px-3 py-1 transition-all " +
              (isActive
                ? "scale-90 bg-primary-container/10 text-primary"
                : "text-on-surface-variant hover:bg-surface-container-high")
            }
          >
            <Icon name={tab.icon} filled={isActive} className="mb-1 text-2xl" />
            <span className="font-body text-[10px] font-medium tracking-wider uppercase">
              {tab.label}
            </span>
          </button>
        )
      })}
    </nav>
  )
}
