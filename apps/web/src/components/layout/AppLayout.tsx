import { useState } from "react"
import { Outlet } from "react-router-dom"
import { AppSidebar } from "./AppSidebar"
import { useAuth } from "@/contexts/AuthContext"
import { Button } from "@/components/ui/button"
import { LogOut, Menu } from "lucide-react"
import { cn } from "@/lib/utils"

export function AppLayout() {
  const { logout } = useAuth()
  const [sidebarOpen, setSidebarOpen] = useState(false)

  return (
    <div className="flex h-screen w-full min-w-0 overflow-hidden bg-background">
      {/* Mobile overlay when sidebar is open */}
      <button
        type="button"
        aria-label="Close menu"
        className={cn(
          "fixed inset-0 z-40 bg-black/50 lg:hidden",
          sidebarOpen ? "block" : "hidden"
        )}
        onClick={() => setSidebarOpen(false)}
      />
      <AppSidebar
        open={sidebarOpen}
        onClose={() => setSidebarOpen(false)}
        className="hidden lg:flex"
      />
      <AppSidebar
        open={sidebarOpen}
        onClose={() => setSidebarOpen(false)}
        className="flex lg:hidden"
        mobile
      />
      <main className="flex min-w-0 flex-1 flex-col overflow-auto">
        <header className="flex h-14 shrink-0 items-center justify-between gap-2 border-b border-border px-3 py-2 sm:px-4">
          <Button
            variant="ghost"
            size="icon"
            className="h-10 w-10 min-h-[44px] min-w-[44px] lg:hidden"
            onClick={() => setSidebarOpen(true)}
            aria-label="Open menu"
          >
            <Menu className="h-5 w-5" />
          </Button>
          <div className="flex flex-1 justify-end">
            <Button
              variant="ghost"
              size="sm"
              onClick={logout}
              className="gap-2 min-h-[44px] px-3 sm:min-h-0 sm:px-2"
            >
              <LogOut className="h-4 w-4" />
              <span className="hidden sm:inline">Log out</span>
            </Button>
          </div>
        </header>
        <div className="flex-1 p-4 sm:p-6">
          <Outlet />
        </div>
      </main>
    </div>
  )
}
