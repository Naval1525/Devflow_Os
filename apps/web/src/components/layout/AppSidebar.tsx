import { Link, useLocation } from "react-router-dom"
import { cn } from "@/lib/utils"
import {
  LayoutDashboard,
  Lightbulb,
  Code2,
  Terminal,
  Briefcase,
  Wallet,
  Sparkles,
} from "lucide-react"
const nav = [
  { to: "/", label: "Dashboard", icon: LayoutDashboard },
  { to: "/coding-log", label: "Coding log", icon: Terminal },
  { to: "/ideas", label: "Ideas", icon: Lightbulb },
  { to: "/leetcode", label: "LeetCode", icon: Code2 },
  { to: "/opportunities", label: "Opportunities", icon: Briefcase },
  { to: "/finance", label: "Finance", icon: Wallet },
  { to: "/ai", label: "AI Generator", icon: Sparkles },
]

type AppSidebarProps = {
  open?: boolean
  onClose?: () => void
  mobile?: boolean
  className?: string
}

export function AppSidebar({ open = true, onClose, mobile, className }: AppSidebarProps) {
  const location = useLocation()

  return (
    <aside
      className={cn(
        "flex h-full flex-col border-r border-border bg-card/50",
        "w-56 shrink-0",
        mobile && [
          "fixed left-0 top-0 z-50 h-full w-64 max-w-[85vw] transform transition-transform duration-200 ease-out",
          open ? "translate-x-0" : "-translate-x-full",
        ],
        className
      )}
    >
      {/* Same height + border as AppLayout main header so the divider is one continuous line */}
      <div className="flex h-14 shrink-0 items-center justify-between gap-2 border-b border-border px-4">
        <div className="min-w-0 flex flex-col justify-center gap-0.5 leading-tight">
          <h1 className="font-display truncate text-base font-semibold text-foreground sm:text-lg">
            DevFlow OS
          </h1>
          <p className="truncate text-[11px] text-muted-foreground sm:text-xs">Control center</p>
        </div>
        {mobile && onClose && (
          <button
            type="button"
            onClick={onClose}
            className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-muted-foreground hover:bg-accent/50 hover:text-foreground min-h-[44px] min-w-[44px] sm:min-h-0 sm:min-w-0"
            aria-label="Close menu"
          >
            <span className="text-lg font-medium leading-none">×</span>
          </button>
        )}
      </div>
      <nav className="flex-1 space-y-0.5 overflow-y-auto p-2">
        {nav.map((item) => {
          const Icon = item.icon
          const isActive = location.pathname === item.to
          return (
            <Link
              key={item.to}
              to={item.to}
              onClick={mobile ? onClose : undefined}
              className={cn(
                "flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors min-h-[44px] sm:min-h-0 sm:py-2",
                isActive
                  ? "bg-primary/10 text-primary"
                  : "text-muted-foreground hover:bg-accent/50 hover:text-accent-foreground"
              )}
            >
              <Icon className="h-4 w-4 shrink-0" />
              {item.label}
            </Link>
          )
        })}
      </nav>
    </aside>
  )
}
