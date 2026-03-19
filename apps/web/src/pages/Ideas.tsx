import { useEffect, useState } from "react"
import { apiFetch } from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select } from "@/components/ui/select"
import { cn } from "@/lib/utils"

type Idea = {
  id: string
  hook: string
  idea: string
  type: string
  status: string
  created_at: string
}

const IDEA_TYPES = ["tweet", "reel", "thread", "linkedin"] as const
const IDEA_STATUSES = ["idea", "ready", "posted"] as const

export function Ideas() {
  const [ideas, setIdeas] = useState<Idea[]>([])
  const [loading, setLoading] = useState(true)
  const [filterType, setFilterType] = useState<string>("")
  const [filterStatus, setFilterStatus] = useState<string>("")
  const [hook, setHook] = useState("")
  const [ideaDetail, setIdeaDetail] = useState("")
  const [addType, setAddType] = useState<string>("tweet")
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editStatus, setEditStatus] = useState<string>("")

  async function load() {
    setLoading(true)
    try {
      const params = new URLSearchParams()
      if (filterType) params.set("type", filterType)
      if (filterStatus) params.set("status", filterStatus)
      const res = await apiFetch(`/ideas?${params}`)
      if (res.ok) {
        const data = await res.json()
        setIdeas(data.ideas ?? [])
      }
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [filterType, filterStatus])

  async function handleAdd(e: React.FormEvent) {
    e.preventDefault()
    if (!hook.trim()) return
    const res = await apiFetch("/ideas", {
      method: "POST",
      body: JSON.stringify({ hook: hook.trim(), idea: ideaDetail.trim(), type: addType, status: "idea" }),
    })
    if (res.ok) {
      setHook("")
      setIdeaDetail("")
      load()
    }
  }

  async function handleUpdateStatus(id: string, status: string) {
    const res = await apiFetch(`/ideas/${id}`, {
      method: "PATCH",
      body: JSON.stringify({ status }),
    })
    if (res.ok) {
      setEditingId(null)
      load()
    }
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="font-display text-2xl font-semibold">Idea Bank</h1>
        <p className="text-muted-foreground">Quick add and filter content ideas</p>
      </div>

      <Card className="border-border">
        <CardHeader>
          <CardTitle className="font-display">Quick add</CardTitle>
          <CardDescription>Hook (punch line) and optional details</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleAdd} className="flex flex-col gap-3">
            <div className="grid gap-2 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="hook">Hook</Label>
                <Input
                  id="hook"
                  placeholder="Short punch line..."
                  value={hook}
                  onChange={(e) => setHook(e.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="add-type">Type</Label>
                <Select id="add-type" value={addType} onChange={(e) => setAddType(e.target.value)}>
                  {IDEA_TYPES.map((t) => (
                    <option key={t} value={t}>{t}</option>
                  ))}
                </Select>
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="idea">Details (optional)</Label>
              <Input
                id="idea"
                placeholder="More context..."
                value={ideaDetail}
                onChange={(e) => setIdeaDetail(e.target.value)}
              />
            </div>
            <Button type="submit">Add idea</Button>
          </form>
        </CardContent>
      </Card>

      <Card className="border-border">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <div>
            <CardTitle className="font-display">Ideas</CardTitle>
            <CardDescription>Filter by type and status</CardDescription>
          </div>
          <div className="flex gap-2">
            <Select
              value={filterType}
              onChange={(e) => setFilterType(e.target.value)}
              className="w-32"
            >
              <option value="">All types</option>
              {IDEA_TYPES.map((t) => (
                <option key={t} value={t}>{t}</option>
              ))}
            </Select>
            <Select
              value={filterStatus}
              onChange={(e) => setFilterStatus(e.target.value)}
              className="w-32"
            >
              <option value="">All statuses</option>
              {IDEA_STATUSES.map((s) => (
                <option key={s} value={s}>{s}</option>
              ))}
            </Select>
          </div>
        </CardHeader>
        <CardContent>
          {loading ? (
            <p className="text-muted-foreground">Loading…</p>
          ) : ideas.length === 0 ? (
            <p className="text-muted-foreground">No ideas yet. Add one above.</p>
          ) : (
            <ul className="space-y-2">
              {ideas.map((i) => (
                <li
                  key={i.id}
                  className={cn(
                    "flex flex-wrap items-center justify-between gap-2 rounded-lg border border-border p-3",
                    editingId === i.id && "ring-2 ring-primary/20"
                  )}
                >
                  <div className="min-w-0 flex-1">
                    <p className="font-medium">{i.hook}</p>
                    {i.idea && <p className="text-sm text-muted-foreground">{i.idea}</p>}
                    <div className="mt-1 flex gap-2">
                      <span className="rounded bg-muted px-2 py-0.5 text-xs">{i.type}</span>
                      <span className="rounded bg-muted px-2 py-0.5 text-xs">{i.status}</span>
                    </div>
                  </div>
                  {editingId === i.id ? (
                    <div className="flex items-center gap-2">
                      <Select
                        value={editStatus}
                        onChange={(e) => setEditStatus(e.target.value)}
                        className="w-28"
                      >
                        {IDEA_STATUSES.map((s) => (
                          <option key={s} value={s}>{s}</option>
                        ))}
                      </Select>
                      <Button size="sm" onClick={() => handleUpdateStatus(i.id, editStatus)}>Save</Button>
                      <Button size="sm" variant="ghost" onClick={() => setEditingId(null)}>Cancel</Button>
                    </div>
                  ) : (
                    <Button size="sm" variant="outline" onClick={() => { setEditingId(i.id); setEditStatus(i.status); }}>Edit status</Button>
                  )}
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
