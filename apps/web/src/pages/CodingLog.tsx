import { useEffect, useState } from "react"
import { apiFetch } from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
type CodingLogEntry = {
  id: string
  title: string
  description: string
  created_at: string
}

export function CodingLog() {
  const [logs, setLogs] = useState<CodingLogEntry[]>([])
  const [loading, setLoading] = useState(true)
  const [title, setTitle] = useState("")
  const [description, setDescription] = useState("")

  async function load() {
    setLoading(true)
    try {
      const res = await apiFetch("/coding-logs")
      if (res.ok) {
        const data = await res.json()
        setLogs(data.logs ?? [])
      }
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  async function handleAdd(e: React.FormEvent) {
    e.preventDefault()
    if (!title.trim()) return
    const res = await apiFetch("/coding-logs", {
      method: "POST",
      body: JSON.stringify({ title: title.trim(), description: description.trim() }),
    })
    if (res.ok) {
      setTitle("")
      setDescription("")
      load()
    }
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="font-display text-2xl font-semibold">Coding log</h1>
        <p className="text-muted-foreground">
          Log what you built, fixed, or learned. Use it in AI Generator to create content.
        </p>
      </div>

      <Card className="border-border">
        <CardHeader>
          <CardTitle className="font-display">What I did</CardTitle>
          <CardDescription>Short title and optional details (e.g. &quot;Fixed auth bug&quot;, &quot;Built dashboard API&quot;)</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleAdd} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="title">Title</Label>
              <Input
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="e.g. Fixed login redirect bug"
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">Details (optional)</Label>
              <textarea
                id="description"
                className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                placeholder="What you did, what you learned..."
                value={description}
                onChange={(e) => setDescription(e.target.value)}
              />
            </div>
            <Button type="submit">Add log</Button>
          </form>
        </CardContent>
      </Card>

      <Card className="border-border">
        <CardHeader>
          <CardTitle className="font-display">Your coding logs</CardTitle>
          <CardDescription>Use these in AI Generator to turn your work into tweets, hooks, or reel scripts</CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <p className="text-muted-foreground">Loading…</p>
          ) : logs.length === 0 ? (
            <p className="text-muted-foreground">No logs yet. Add what you did above.</p>
          ) : (
            <ul className="space-y-2">
              {logs.map((log) => (
                <li key={log.id} className="rounded-lg border border-border p-3">
                  <p className="font-medium">{log.title}</p>
                  {log.description && (
                    <p className="mt-1 text-sm text-muted-foreground">{log.description}</p>
                  )}
                  <p className="mt-1 text-xs text-muted-foreground">
                    {new Date(log.created_at).toLocaleDateString()}
                  </p>
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
