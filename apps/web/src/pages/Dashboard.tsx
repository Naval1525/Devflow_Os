import { useEffect, useState } from "react"
import { apiFetch } from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Checkbox } from "@/components/ui/checkbox"
import { Play, Square, Plus, Pencil, Trash2, Calendar } from "lucide-react"

type Task = {
  id: string
  type: string
  date: string
  completed: boolean
  created_at?: string
}

type Session = {
  id: string
  start_time: string
  end_time: string | null
}

function formatDisplayDate(iso: string): string {
  if (!iso) return ""
  const d = new Date(iso + "T12:00:00")
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const dNorm = new Date(d)
  dNorm.setHours(0, 0, 0, 0)
  if (dNorm.getTime() === today.getTime()) return "Today"
  const yesterday = new Date(today)
  yesterday.setDate(yesterday.getDate() - 1)
  if (dNorm.getTime() === yesterday.getTime()) return "Yesterday"
  return d.toLocaleDateString("en-IN", { weekday: "short", day: "numeric", month: "short" })
}

function toISODate(d: Date): string {
  return d.toISOString().slice(0, 10)
}

export function Dashboard() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [selectedDate, setSelectedDate] = useState<string>(() => toISODate(new Date()))
  const [activeSession, setActiveSession] = useState<Session | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [newTitle, setNewTitle] = useState("")
  const [adding, setAdding] = useState(false)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editTitle, setEditTitle] = useState("")
  const [editDate, setEditDate] = useState("")

  async function loadTasks() {
    setError(null)
    try {
      const url = selectedDate ? `/tasks?date=${selectedDate}` : "/tasks/today"
      const res = await apiFetch(url)
      if (res.ok) {
        const data = await res.json()
        setTasks(data.tasks ?? [])
      }
    } catch {
      setError("Failed to load tasks")
    } finally {
      setLoading(false)
    }
  }

  async function loadSession() {
    try {
      const res = await apiFetch("/sessions/active")
      if (res.ok) {
        const data = await res.json()
        setActiveSession(data)
      }
    } catch {
      // ignore
    }
  }

  function load() {
    setLoading(true)
    loadTasks()
    loadSession()
  }

  useEffect(() => {
    load()
  }, [selectedDate])

  async function toggleTask(task: Task) {
    const res = await apiFetch(`/tasks/${task.id}`, {
      method: "PATCH",
      body: JSON.stringify({ completed: !task.completed }),
    })
    if (res.ok) loadTasks()
  }

  async function addTask(e: React.FormEvent) {
    e.preventDefault()
    const title = newTitle.trim()
    if (!title) return
    setAdding(true)
    try {
      const res = await apiFetch("/tasks", {
        method: "POST",
        body: JSON.stringify({ title, date: selectedDate || toISODate(new Date()) }),
      })
      if (res.ok) {
        setNewTitle("")
        loadTasks()
      }
    } finally {
      setAdding(false)
    }
  }

  function startEdit(task: Task) {
    setEditingId(task.id)
    setEditTitle(task.type)
    setEditDate(task.date)
  }

  function cancelEdit() {
    setEditingId(null)
    setEditTitle("")
    setEditDate("")
  }

  async function saveEdit(e: React.FormEvent) {
    e.preventDefault()
    if (!editingId) return
    const title = editTitle.trim()
    if (!title) return
    const res = await apiFetch(`/tasks/${editingId}`, {
      method: "PATCH",
      body: JSON.stringify({ title, date: editDate || selectedDate }),
    })
    if (res.ok) {
      cancelEdit()
      loadTasks()
    }
  }

  async function deleteTask(id: string) {
    if (!confirm("Delete this task?")) return
    const res = await apiFetch(`/tasks/${id}`, { method: "DELETE" })
    if (res.ok) loadTasks()
  }

  async function startSession() {
    const res = await apiFetch("/sessions/start", { method: "POST" })
    if (res.ok) load()
  }

  async function endSession() {
    const res = await apiFetch("/sessions/end", { method: "POST" })
    if (res.ok) load()
  }

  if (loading && tasks.length === 0) {
    return (
      <div className="flex items-center justify-center py-12">
        <p className="text-muted-foreground">Loading…</p>
      </div>
    )
  }

  if (error && tasks.length === 0) {
    return (
      <div className="space-y-4">
        <h1 className="font-display text-xl font-semibold sm:text-2xl">Dashboard</h1>
        <p className="text-sm text-destructive sm:text-base">{error}</p>
        <button
          type="button"
          onClick={() => load()}
          className="min-h-[44px] rounded-md px-3 text-sm text-primary hover:underline sm:min-h-0"
        >
          Retry
        </button>
      </div>
    )
  }

  return (
    <div className="space-y-6 sm:space-y-8">
      <div>
        <h1 className="font-display text-xl font-semibold sm:text-2xl">Dashboard</h1>
        <p className="text-sm text-muted-foreground sm:text-base">
          Your personal productivity — tasks and focus
        </p>
      </div>

      <div className="grid gap-4 sm:gap-6 md:grid-cols-2">
        <Card className="border-border md:col-span-2">
          <CardHeader className="flex flex-col gap-3 pb-2 sm:flex-row sm:items-center sm:justify-between sm:space-y-0">
            <div className="min-w-0">
              <CardTitle className="font-display flex items-center gap-2 text-lg sm:text-base">
                <Calendar className="h-5 w-5 shrink-0" />
                <span className="truncate">{formatDisplayDate(selectedDate)}</span>
              </CardTitle>
              <CardDescription>Add and edit tasks for this day</CardDescription>
            </div>
            <div className="flex items-center gap-2 shrink-0">
              <Label htmlFor="dashboard-date" className="text-sm text-muted-foreground">
                Date
              </Label>
              <Input
                id="dashboard-date"
                type="date"
                value={selectedDate}
                onChange={(e) => setSelectedDate(e.target.value)}
                className="w-full min-w-0 sm:w-40 min-h-[44px] sm:min-h-[2.5rem]"
              />
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <form onSubmit={addTask} className="flex flex-col gap-3 sm:flex-row sm:flex-wrap sm:items-end sm:gap-2">
              <div className="min-w-0 flex-1 space-y-1">
                <Label htmlFor="new-task-title">New task</Label>
                <Input
                  id="new-task-title"
                  placeholder="e.g. Coding, Call mom, Review PR"
                  value={newTitle}
                  onChange={(e) => setNewTitle(e.target.value)}
                  className="min-h-[44px] sm:min-h-[2.5rem]"
                />
              </div>
              <Button
                type="submit"
                disabled={adding || !newTitle.trim()}
                className="min-h-[44px] w-full sm:min-h-0 sm:w-auto"
              >
                <Plus className="h-4 w-4" />
                Add
              </Button>
            </form>

            <ul className="space-y-2">
              {tasks.length === 0 && (
                <li className="rounded-lg border border-dashed border-border p-4 text-center text-sm text-muted-foreground">
                  No tasks for this day. Add one above.
                </li>
              )}
              {tasks.map((task) => (
                <li
                  key={task.id}
                  className="flex flex-wrap items-center gap-2 rounded-lg border border-border bg-card p-3 transition-colors hover:bg-accent/20 sm:gap-3"
                >
                  <Checkbox
                    checked={task.completed}
                    onCheckedChange={() => toggleTask(task)}
                    aria-label={`Mark ${task.type} complete`}
                    className="mt-0.5"
                  />
                  {editingId === task.id ? (
                    <form onSubmit={saveEdit} className="flex min-w-0 flex-1 flex-wrap items-center gap-2">
                      <Input
                        value={editTitle}
                        onChange={(e) => setEditTitle(e.target.value)}
                        placeholder="Task title"
                        className="min-w-[120px] flex-1 min-h-[44px] sm:min-h-[2.5rem]"
                        autoFocus
                      />
                      <Input
                        type="date"
                        value={editDate}
                        onChange={(e) => setEditDate(e.target.value)}
                        className="w-full min-w-0 sm:w-40 min-h-[44px] sm:min-h-[2.5rem]"
                      />
                      <div className="flex gap-1 w-full sm:w-auto">
                        <Button type="submit" size="sm" className="min-h-[44px] flex-1 sm:min-h-0 sm:flex-none">
                          Save
                        </Button>
                        <Button type="button" size="sm" variant="ghost" onClick={cancelEdit} className="min-h-[44px] sm:min-h-0">
                          Cancel
                        </Button>
                      </div>
                    </form>
                  ) : (
                    <>
                      <span
                        className={`min-w-0 flex-1 text-sm font-medium break-words ${task.completed ? "text-muted-foreground line-through" : ""}`}
                      >
                        {task.type}
                      </span>
                      <div className="flex shrink-0 gap-1">
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          className="h-9 w-9 min-h-[44px] min-w-[44px] sm:h-8 sm:w-8 sm:min-h-0 sm:min-w-0"
                          onClick={() => startEdit(task)}
                          aria-label="Edit task"
                        >
                          <Pencil className="h-4 w-4" />
                        </Button>
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          className="h-9 w-9 min-h-[44px] min-w-[44px] text-destructive hover:text-destructive sm:h-8 sm:w-8 sm:min-h-0 sm:min-w-0"
                          onClick={() => deleteTask(task.id)}
                          aria-label="Delete task"
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </>
                  )}
                </li>
              ))}
            </ul>
          </CardContent>
        </Card>

        <Card className="border-border">
          <CardHeader className="p-4 sm:p-6">
            <CardTitle className="font-display text-base sm:text-base">Deep work</CardTitle>
            <CardDescription>Start or stop a focus session</CardDescription>
          </CardHeader>
          <CardContent className="p-4 pt-0 sm:p-6 sm:pt-0">
            {activeSession ? (
              <div className="flex flex-wrap items-center gap-3">
                <p className="text-sm text-muted-foreground">Session in progress</p>
                <Button size="sm" variant="destructive" onClick={endSession} className="min-h-[44px] sm:min-h-0">
                  <Square className="h-4 w-4" />
                  Stop
                </Button>
              </div>
            ) : (
              <Button onClick={startSession} className="min-h-[44px] w-full sm:min-h-0 sm:w-auto">
                <Play className="h-4 w-4" />
                Start session
              </Button>
            )}
          </CardContent>
        </Card>

        <Card className="border-border">
          <CardHeader className="p-4 sm:p-6">
            <CardTitle className="font-display text-base sm:text-base">Quick links</CardTitle>
            <CardDescription>Ideas, LeetCode, content, finances</CardDescription>
          </CardHeader>
          <CardContent className="p-4 pt-0 text-sm text-muted-foreground sm:p-6 sm:pt-0">
            Use the sidebar to open Ideas, LeetCode log, AI Generator, Opportunities, and Finance.
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
