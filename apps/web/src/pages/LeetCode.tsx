import { useEffect, useState } from "react"
import { apiFetch } from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select } from "@/components/ui/select"
import { Sparkles } from "lucide-react"

type Log = {
  id: string
  problem_name: string
  difficulty: string
  approach: string
  mistake: string
  time_taken: number | null
  created_at: string
}

type GeneratedContent = {
  tweet?: string
  reel_script?: string
  hook?: string
}

const DIFFICULTIES = ["easy", "medium", "hard"] as const

export function LeetCode() {
  const [logs, setLogs] = useState<Log[]>([])
  const [loading, setLoading] = useState(true)
  const [problemName, setProblemName] = useState("")
  const [difficulty, setDifficulty] = useState<string>("easy")
  const [approach, setApproach] = useState("")
  const [mistake, setMistake] = useState("")
  const [timeTaken, setTimeTaken] = useState("")
  const [generatingId, setGeneratingId] = useState<string | null>(null)
  const [generated, setGenerated] = useState<GeneratedContent | null>(null)
  const [generatedForLogId, setGeneratedForLogId] = useState<string | null>(null)

  async function load() {
    setLoading(true)
    try {
      const res = await apiFetch("/leetcode")
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
    if (!problemName.trim()) return
    const body: Record<string, unknown> = {
      problem_name: problemName.trim(),
      difficulty,
      approach: approach.trim(),
      mistake: mistake.trim(),
    }
    if (timeTaken.trim()) body.time_taken = parseInt(timeTaken, 10)
    const res = await apiFetch("/leetcode", { method: "POST", body: JSON.stringify(body) })
    if (res.ok) {
      setProblemName("")
      setApproach("")
      setMistake("")
      setTimeTaken("")
      load()
    }
  }

  async function generateContent(log: Log) {
    setGeneratingId(log.id)
    setGenerated(null)
    setGeneratedForLogId(null)
    const text = [
      `Problem: ${log.problem_name}`,
      `Difficulty: ${log.difficulty}`,
      log.approach && `Approach: ${log.approach}`,
      log.mistake && `Mistake learned: ${log.mistake}`,
    ].filter(Boolean).join("\n")
    try {
      const res = await apiFetch("/generate-content", {
        method: "POST",
        body: JSON.stringify({ text, formats: ["tweet", "reel_script", "hook"] }),
      })
      if (res.ok) {
        const data = await res.json()
        setGenerated(data)
        setGeneratedForLogId(log.id)
      }
    } finally {
      setGeneratingId(null)
    }
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="font-display text-2xl font-semibold">LeetCode Tracker</h1>
        <p className="text-muted-foreground">Log problems and generate content from them</p>
      </div>

      <Card className="border-border">
        <CardHeader>
          <CardTitle className="font-display">Add problem</CardTitle>
          <CardDescription>Problem name, difficulty, approach, mistake, time (min)</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleAdd} className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="problem">Problem name</Label>
              <Input id="problem" value={problemName} onChange={(e) => setProblemName(e.target.value)} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="difficulty">Difficulty</Label>
              <Select id="difficulty" value={difficulty} onChange={(e) => setDifficulty(e.target.value)}>
                {DIFFICULTIES.map((d) => <option key={d} value={d}>{d}</option>)}
              </Select>
            </div>
            <div className="space-y-2 sm:col-span-2">
              <Label htmlFor="approach">Approach</Label>
              <Input id="approach" value={approach} onChange={(e) => setApproach(e.target.value)} placeholder="Short approach" />
            </div>
            <div className="space-y-2 sm:col-span-2">
              <Label htmlFor="mistake">Mistake (important)</Label>
              <Input id="mistake" value={mistake} onChange={(e) => setMistake(e.target.value)} placeholder="What you learned" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="time">Time (minutes)</Label>
              <Input id="time" type="number" value={timeTaken} onChange={(e) => setTimeTaken(e.target.value)} placeholder="30" />
            </div>
            <div className="flex items-end">
              <Button type="submit">Add log</Button>
            </div>
          </form>
        </CardContent>
      </Card>

      <Card className="border-border">
        <CardHeader>
          <CardTitle className="font-display">Logs</CardTitle>
          <CardDescription>Click Generate Content to get tweet / reel hook / script</CardDescription>
        </CardHeader>
        <CardContent>
          {loading ? (
            <p className="text-muted-foreground">Loading…</p>
          ) : logs.length === 0 ? (
            <p className="text-muted-foreground">No logs yet.</p>
          ) : (
            <ul className="space-y-3">
              {logs.map((log) => (
                <li key={log.id} className="rounded-lg border border-border p-3">
                  <div className="flex flex-wrap items-center justify-between gap-2">
                    <div>
                      <p className="font-medium">{log.problem_name}</p>
                      <p className="text-sm text-muted-foreground">
                        {log.difficulty} · {log.approach && `${log.approach} · `}
                        {log.mistake && `Mistake: ${log.mistake}`}
                        {log.time_taken != null && ` · ${log.time_taken} min`}
                      </p>
                    </div>
                    <Button
                      size="sm"
                      variant="secondary"
                      disabled={generatingId !== null}
                      onClick={() => generateContent(log)}
                    >
                      <Sparkles className="h-4 w-4" />
                      {generatingId === log.id ? "Generating…" : "Generate Content"}
                    </Button>
                  </div>
                  {generated && generatedForLogId === log.id && (
                    <div className="mt-3 space-y-2 rounded bg-muted/50 p-3 text-sm">
                      {generated.tweet && <p><strong>Tweet:</strong> {generated.tweet}</p>}
                      {generated.hook && <p><strong>Hook:</strong> {generated.hook}</p>}
                      {generated.reel_script && <p><strong>Reel script:</strong> {generated.reel_script}</p>}
                    </div>
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
