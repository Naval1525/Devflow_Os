import { useEffect, useState } from "react"
import { apiFetch } from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select } from "@/components/ui/select"

type Opportunity = {
  id: string
  name: string
  type: string
  stage: string
  source: string
  notes: string
  created_at: string
}

const OPP_TYPES = ["job", "freelance"] as const
const STAGES = ["applied", "interview", "closed"] as const

export function Opportunities() {
  const [opportunities, setOpportunities] = useState<Opportunity[]>([])
  const [loading, setLoading] = useState(true)
  const [name, setName] = useState("")
  const [oppType, setOppType] = useState<string>("job")
  const [stage, setStage] = useState<string>("applied")
  const [source, setSource] = useState("")
  const [notes, setNotes] = useState("")
  const [editingId, setEditingId] = useState<string | null>(null)
  const [editStage, setEditStage] = useState<string>("")
  const [editNotes, setEditNotes] = useState("")

  async function load() {
    setLoading(true)
    try {
      const res = await apiFetch("/opportunities")
      if (res.ok) {
        const data = await res.json()
        setOpportunities(data.opportunities ?? [])
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
    if (!name.trim()) return
    const res = await apiFetch("/opportunities", {
      method: "POST",
      body: JSON.stringify({ name: name.trim(), type: oppType, stage, source: source.trim(), notes: notes.trim() }),
    })
    if (res.ok) {
      setName("")
      setSource("")
      setNotes("")
      load()
    }
  }

  async function handleUpdate(id: string) {
    const res = await apiFetch(`/opportunities/${id}`, {
      method: "PATCH",
      body: JSON.stringify({ stage: editStage, notes: editNotes }),
    })
    if (res.ok) {
      setEditingId(null)
      load()
    }
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="font-display text-2xl font-semibold">Opportunities</h1>
        <p className="text-muted-foreground">Job and freelance pipeline</p>
      </div>

      <Card className="border-border">
        <CardHeader>
          <CardTitle className="font-display">Add opportunity</CardTitle>
          <CardDescription>Company or client name, type, stage, source</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleAdd} className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="name">Name (company / client)</Label>
              <Input id="name" value={name} onChange={(e) => setName(e.target.value)} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="opp-type">Type</Label>
              <Select id="opp-type" value={oppType} onChange={(e) => setOppType(e.target.value)}>
                {OPP_TYPES.map((t) => <option key={t} value={t}>{t}</option>)}
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="stage">Stage</Label>
              <Select id="stage" value={stage} onChange={(e) => setStage(e.target.value)}>
                {STAGES.map((s) => <option key={s} value={s}>{s}</option>)}
              </Select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="source">Source</Label>
              <Input id="source" value={source} onChange={(e) => setSource(e.target.value)} placeholder="LinkedIn, referral…" />
            </div>
            <div className="space-y-2 sm:col-span-2">
              <Label htmlFor="notes">Notes</Label>
              <Input id="notes" value={notes} onChange={(e) => setNotes(e.target.value)} placeholder="Optional" />
            </div>
            <Button type="submit">Add</Button>
          </form>
        </CardContent>
      </Card>

      <Card className="border-border">
        <CardHeader>
          <CardTitle className="font-display">Pipeline</CardTitle>
          <CardDescription>Update stage or notes as you progress</CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <p className="text-muted-foreground">Loading…</p>
          ) : opportunities.length === 0 ? (
            <p className="text-muted-foreground">No opportunities yet.</p>
          ) : (
            <ul className="space-y-2">
              {opportunities.map((opp) => (
                <li key={opp.id} className="flex flex-wrap items-center justify-between gap-2 rounded-lg border border-border p-3">
                  <div>
                    <p className="font-medium">{opp.name}</p>
                    <p className="text-sm text-muted-foreground">
                      {opp.type} · {opp.stage}
                      {opp.source && ` · ${opp.source}`}
                      {opp.notes && ` · ${opp.notes}`}
                    </p>
                  </div>
                  {editingId === opp.id ? (
                    <div className="flex flex-wrap items-center gap-2">
                      <Select value={editStage} onChange={(e) => setEditStage(e.target.value)} className="w-28">
                        {STAGES.map((s) => <option key={s} value={s}>{s}</option>)}
                      </Select>
                      <Input
                        placeholder="Notes"
                        value={editNotes}
                        onChange={(e) => setEditNotes(e.target.value)}
                        className="w-40"
                      />
                      <Button size="sm" onClick={() => handleUpdate(opp.id)}>Save</Button>
                      <Button size="sm" variant="ghost" onClick={() => setEditingId(null)}>Cancel</Button>
                    </div>
                  ) : (
                    <Button size="sm" variant="outline" onClick={() => { setEditingId(opp.id); setEditStage(opp.stage); setEditNotes(opp.notes); }}>Edit</Button>
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
