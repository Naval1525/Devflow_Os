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
import { Separator } from "@/components/ui/separator"

const nav = [
  { to: "/", label: "Dashboard", icon: LayoutDashboard },
  { to: "/coding-log", label: "Coding log", icon: Terminal },
  { to: "/ideas", label: "Ideas", icon: Lightbulb },
  { to: "/leetcode", label: "LeetCode", icon: Code2 },
  { to: "/opportunities", label: "Opportunities", icon: Briefcase },
  { to: "/finance", label: "Finance", icon: Wallet },
  { to: "/ai", label: "AI Generator", icon: Sparkles },
]

export function AppSidebar() {
  const location = useLocation()

  return (
    <aside className="flex h-full w-56 flex-col border-r border-border bg-card/50">
      <div className="p-4">
        <h1 className="font-display text-lg font-semibold text-foreground">
          DevFlow OS
        </h1>
        <p className="text-xs text-muted-foreground">Control center</p>
      </div>
      <Separator />
      <nav className="flex-1 space-y-0.5 p-2">
        {nav.map((item) => {
          const Icon = item.icon
          const isActive = location.pathname === item.to
          return (
            <Link
              key={item.to}
              to={item.to}
              className={cn(
                "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                isActive
                  ? "bg-primary/10 text-primary"
                  : "text-muted-foreground hover:bg-accent/50 hover:text-accent-foreground"
              )}
            >
              <Icon className="h-4 w-4" />
              {item.label}
            </Link>
          )
        })}
      </nav>
    </aside>
  )
}
