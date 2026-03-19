import { useEffect, useState } from "react"
import { apiFetch } from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Play, Square } from "lucide-react"

type Task = {
  id: string
  type: string
  date: string
  completed: boolean
}

type Session = {
  id: string
  start_time: string
  end_time: string | null
}

export function Dashboard() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [activeSession, setActiveSession] = useState<Session | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  async function load() {
    setError(null)
    try {
      const [tasksRes, sessionRes] = await Promise.all([
        apiFetch("/tasks/today"),
        apiFetch("/sessions/active"),
      ])
      if (tasksRes.ok) {
        const data = await tasksRes.json()
        setTasks(data.tasks ?? [])
      }
      if (sessionRes.ok) {
        const data = await sessionRes.json()
        setActiveSession(data)
      }
    } catch {
      setError("Failed to load dashboard")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  async function toggleTask(type: string, completed: boolean) {
    const res = await apiFetch("/tasks/update", {
      method: "POST",
      body: JSON.stringify({ type, completed }),
    })
    if (res.ok) load()
  }

  async function startSession() {
    const res = await apiFetch("/sessions/start", { method: "POST" })
    if (res.ok) load()
  }

  async function endSession() {
    const res = await apiFetch("/sessions/end", { method: "POST" })
    if (res.ok) load()
  }

  const labels: Record<string, string> = {
    coding: "Coding",
    leetcode: "LeetCode",
    content: "Content",
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <p className="text-muted-foreground">Loading…</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="space-y-4">
        <h1 className="font-display text-2xl font-semibold">Dashboard</h1>
        <p className="text-destructive">{error}</p>
        <button
          type="button"
          onClick={() => load()}
          className="text-sm text-primary hover:underline"
        >
          Retry
        </button>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="font-display text-2xl font-semibold">Dashboard</h1>
        <p className="text-muted-foreground">Today&apos;s checklist and focus</p>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card className="border-border">
          <CardHeader>
            <CardTitle className="font-display">Today&apos;s checklist</CardTitle>
            <CardDescription>Mark done as you go</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            {tasks.map((task) => (
              <label
                key={task.id}
                className="flex cursor-pointer items-center gap-3 rounded-lg border border-border p-3 transition-colors hover:bg-accent/30"
              >
                <Checkbox
                  checked={task.completed}
                  onCheckedChange={(checked) =>
                    toggleTask(task.type, !!checked)
                  }
                />
                <span className="text-sm font-medium">
                  {labels[task.type] ?? task.type}
                </span>
              </label>
            ))}
          </CardContent>
        </Card>

        <Card className="border-border">
          <CardHeader>
            <CardTitle className="font-display">Deep work</CardTitle>
            <CardDescription>Start or stop a focus session</CardDescription>
          </CardHeader>
          <CardContent>
            {activeSession ? (
              <div className="flex items-center gap-4">
                <p className="text-sm text-muted-foreground">Session in progress</p>
                <Button size="sm" variant="destructive" onClick={endSession}>
                  <Square className="h-4 w-4" />
                  Stop
                </Button>
              </div>
            ) : (
              <Button onClick={startSession}>
                <Play className="h-4 w-4" />
                Start session
              </Button>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
