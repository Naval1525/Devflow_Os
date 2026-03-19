import { Outlet } from "react-router-dom"
import { AppSidebar } from "./AppSidebar"
import { useAuth } from "@/contexts/AuthContext"
import { Button } from "@/components/ui/button"
import { LogOut } from "lucide-react"

export function AppLayout() {
  const { logout } = useAuth()

  return (
    <div className="flex h-screen w-full overflow-hidden bg-background">
      <AppSidebar />
      <main className="flex flex-1 flex-col overflow-auto">
        <header className="flex h-14 shrink-0 items-center justify-end border-b border-border px-4">
          <Button variant="ghost" size="sm" onClick={logout} className="gap-2">
            <LogOut className="h-4 w-4" />
            Log out
          </Button>
        </header>
        <div className="flex-1 p-6">
          <Outlet />
        </div>
      </main>
    </div>
  )
}
