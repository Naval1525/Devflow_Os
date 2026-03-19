import { useEffect, useState } from "react"
import { apiFetch } from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import { Checkbox } from "@/components/ui/checkbox"
import { Select } from "@/components/ui/select"
import { Sparkles, Copy, Check } from "lucide-react"

type GeneratedContent = {
  tweet?: string
  reel_script?: string
  hook?: string
  linkedin_post?: string
}

type LeetCodeLog = {
  id: string
  problem_name: string
  difficulty: string
  approach: string
  mistake: string
  time_taken: number | null
  created_at: string
}

type Idea = {
  id: string
  hook: string
  idea: string
  type: string
  status: string
  created_at: string
}

type CodingLogEntry = {
  id: string
  title: string
  description: string
  created_at: string
}

type SourceType = "coding" | "leetcode" | "idea" | "custom" | "all"

const GEMINI_MODELS = [
  { value: "gemini-2.5-flash", label: "Gemini 2.5 Flash (balanced)" },
  { value: "gemini-2.5-flash-lite", label: "Gemini 2.5 Flash Lite (fast)" },
  { value: "gemini-2.5-pro", label: "Gemini 2.5 Pro (quality)" },
] as const

function buildTextFromLog(log: LeetCodeLog): string {
  return [
    `Problem: ${log.problem_name}`,
    `Difficulty: ${log.difficulty}`,
    log.approach && `Approach: ${log.approach}`,
    log.mistake && `Mistake learned: ${log.mistake}`,
    log.time_taken != null && `Time: ${log.time_taken} min`,
  ].filter(Boolean).join("\n")
}

function buildTextFromIdea(idea: Idea): string {
  return [idea.hook, idea.idea].filter(Boolean).join("\n\n")
}

function buildTextFromCodingLog(log: CodingLogEntry): string {
  return [log.title, log.description].filter(Boolean).join("\n\n")
}

export function AIGenerator() {
  const [sourceType, setSourceType] = useState<SourceType>("coding")
  const [codingLogs, setCodingLogs] = useState<CodingLogEntry[]>([])
  const [leetcodeLogs, setLeetcodeLogs] = useState<LeetCodeLog[]>([])
  const [ideas, setIdeas] = useState<Idea[]>([])
  const [selectedCodingLogId, setSelectedCodingLogId] = useState<string>("")
  const [selectedLogId, setSelectedLogId] = useState<string>("")
  const [selectedIdeaId, setSelectedIdeaId] = useState<string>("")
  const [text, setText] = useState("")
  const [model, setModel] = useState<string>("gemini-2.5-flash")
  const [loadingAll, setLoadingAll] = useState(false)
  const [wantTweet, setWantTweet] = useState(true)
  const [wantReel, setWantReel] = useState(true)
  const [wantHook, setWantHook] = useState(true)
  const [wantLinkedin, setWantLinkedin] = useState(false)
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState<GeneratedContent | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [copiedField, setCopiedField] = useState<"tweet" | "hook" | "reel_script" | "linkedin_post" | null>(null)

  async function copyToClipboard(text: string, field: "tweet" | "hook" | "reel_script" | "linkedin_post") {
    try {
      await navigator.clipboard.writeText(text)
      setCopiedField(field)
      setTimeout(() => setCopiedField(null), 2000)
    } catch {
      // ignore
    }
  }

  function reelScriptSteps(script: string): string[] {
    if (!script.trim()) return []
    const steps = script
      .split(/(?=\d+\.\s)/)
      .map((s) => s.replace(/^\d+\.\s*/, "").trim())
      .filter(Boolean)
    return steps.length > 0 ? steps : [script]
  }

  useEffect(() => {
    if (sourceType === "coding") {
      apiFetch("/coding-logs").then((res) => {
        if (res.ok) res.json().then((d) => setCodingLogs(d.logs ?? []))
      })
      setSelectedCodingLogId("")
      setText("")
    } else if (sourceType === "leetcode") {
      apiFetch("/leetcode").then((res) => {
        if (res.ok) res.json().then((d) => setLeetcodeLogs(d.logs ?? []))
      })
      setSelectedLogId("")
      setText("")
    } else if (sourceType === "idea") {
      apiFetch("/ideas").then((res) => {
        if (res.ok) res.json().then((d) => setIdeas(d.ideas ?? []))
      })
      setSelectedIdeaId("")
      setText("")
    } else if (sourceType === "all") {
      setLoadingAll(true)
      Promise.all([
        apiFetch("/coding-logs").then((r) => (r.ok ? r.json() : { logs: [] })),
        apiFetch("/leetcode").then((r) => (r.ok ? r.json() : { logs: [] })),
        apiFetch("/ideas").then((r) => (r.ok ? r.json() : { ideas: [] })),
      ])
        .then(([coding, leetcode, ideasRes]) => {
          const logs = (coding.logs ?? []) as CodingLogEntry[]
          const lcLogs = (leetcode.logs ?? []) as LeetCodeLog[]
          const ideaList = (ideasRes.ideas ?? []) as Idea[]
          setCodingLogs(logs)
          setLeetcodeLogs(lcLogs)
          setIdeas(ideaList)
          const parts: string[] = []
          if (logs.length) {
            parts.push("--- Coding log ---")
            logs.forEach((l) => parts.push(buildTextFromCodingLog(l), ""))
          }
          if (lcLogs.length) {
            parts.push("--- LeetCode ---")
            lcLogs.forEach((l) => parts.push(buildTextFromLog(l), ""))
          }
          if (ideaList.length) {
            parts.push("--- Ideas ---")
            ideaList.forEach((i) => parts.push(buildTextFromIdea(i), ""))
          }
          setText(parts.join("\n").trim())
        })
        .finally(() => setLoadingAll(false))
    } else {
      setText("")
    }
  }, [sourceType])

  useEffect(() => {
    if (sourceType === "coding" && selectedCodingLogId) {
      const log = codingLogs.find((l) => l.id === selectedCodingLogId)
      if (log) setText(buildTextFromCodingLog(log))
    }
  }, [sourceType, selectedCodingLogId, codingLogs])

  useEffect(() => {
    if (sourceType === "leetcode" && selectedLogId) {
      const log = leetcodeLogs.find((l) => l.id === selectedLogId)
      if (log) setText(buildTextFromLog(log))
    }
  }, [sourceType, selectedLogId, leetcodeLogs])

  useEffect(() => {
    if (sourceType === "idea" && selectedIdeaId) {
      const idea = ideas.find((i) => i.id === selectedIdeaId)
      if (idea) setText(buildTextFromIdea(idea))
    }
  }, [sourceType, selectedIdeaId, ideas])

  async function handleGenerate(e: React.FormEvent) {
    e.preventDefault()
    const inputText = text.trim()
    if (!inputText) return
    setLoading(true)
    setError(null)
    setResult(null)
    const formats: string[] = []
    if (wantTweet) formats.push("tweet")
    if (wantReel) formats.push("reel_script")
    if (wantHook) formats.push("hook")
    if (wantLinkedin) formats.push("linkedin_post")
    if (formats.length === 0) formats.push("tweet", "reel_script", "hook")
    try {
      const res = await apiFetch("/generate-content", {
        method: "POST",
        body: JSON.stringify({ text: inputText, formats, model: model || undefined }),
      })
      if (res.ok) {
        const data = await res.json()
        setResult(data)
      } else {
        const data = await res.json().catch(() => ({}))
        setError(data.error || "Generation failed")
      }
    } catch {
      setError("Request failed")
    } finally {
      setLoading(false)
    }
  }

  const canGenerate =
    sourceType === "custom" || sourceType === "all"
      ? text.trim().length > 0
      : sourceType === "coding"
        ? selectedCodingLogId.length > 0
        : sourceType === "leetcode"
          ? selectedLogId.length > 0
          : sourceType === "idea"
            ? selectedIdeaId.length > 0
            : false

  return (
    <div className="space-y-8">
      <div>
        <h1 className="font-display text-2xl font-semibold">AI Content Generator</h1>
        <p className="text-muted-foreground">
          Generate tweet, reel script, or hook from what you’ve done — coding log, LeetCode, ideas, or custom text (including LinkedIn post)
        </p>
      </div>

      <Card className="border-border">
        <CardHeader>
          <CardTitle className="font-display">Source</CardTitle>
          <CardDescription>Choose what to turn into content</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label>Generate from</Label>
            <Select
              value={sourceType}
              onChange={(e) => setSourceType(e.target.value as SourceType)}
            >
              <option value="all">From all my activity</option>
              <option value="coding">My coding log</option>
              <option value="leetcode">My LeetCode logs</option>
              <option value="idea">My ideas</option>
              <option value="custom">Custom text</option>
            </Select>
          </div>

          <div className="space-y-2">
            <Label>Model</Label>
            <Select value={model} onChange={(e) => setModel(e.target.value)}>
              {GEMINI_MODELS.map((m) => (
                <option key={m.value} value={m.value}>
                  {m.label}
                </option>
              ))}
            </Select>
          </div>

          {sourceType === "all" && (
            <p className="text-sm text-muted-foreground">
              {loadingAll ? "Loading all logs…" : "Combined coding log, LeetCode, and ideas below. Edit if needed, then generate."}
            </p>
          )}

          {sourceType === "coding" && (
            <div className="space-y-2">
              <Label htmlFor="coding-log-select">Pick what you did</Label>
              <Select
                id="coding-log-select"
                value={selectedCodingLogId}
                onChange={(e) => setSelectedCodingLogId(e.target.value)}
              >
                <option value="">Select a log…</option>
                {codingLogs.map((log) => (
                  <option key={log.id} value={log.id}>
                    {log.title.slice(0, 50)}{log.title.length > 50 ? "…" : ""}
                  </option>
                ))}
              </Select>
              {codingLogs.length === 0 && (
                <p className="text-sm text-muted-foreground">No coding logs yet. Add what you did on the Coding log page.</p>
              )}
            </div>
          )}

          {sourceType === "leetcode" && (
            <div className="space-y-2">
              <Label htmlFor="log-select">Pick a problem</Label>
              <Select
                id="log-select"
                value={selectedLogId}
                onChange={(e) => setSelectedLogId(e.target.value)}
              >
                <option value="">Select a log…</option>
                {leetcodeLogs.map((log) => (
                  <option key={log.id} value={log.id}>
                    {log.problem_name} ({log.difficulty})
                  </option>
                ))}
              </Select>
              {leetcodeLogs.length === 0 && (
                <p className="text-sm text-muted-foreground">No LeetCode logs yet. Add some on the LeetCode page.</p>
              )}
            </div>
          )}

          {sourceType === "idea" && (
            <div className="space-y-2">
              <Label htmlFor="idea-select">Pick an idea</Label>
              <Select
                id="idea-select"
                value={selectedIdeaId}
                onChange={(e) => setSelectedIdeaId(e.target.value)}
              >
                <option value="">Select an idea…</option>
                {ideas.map((idea) => (
                  <option key={idea.id} value={idea.id}>
                    {idea.hook.slice(0, 50)}{idea.hook.length > 50 ? "…" : ""} ({idea.type})
                  </option>
                ))}
              </Select>
              {ideas.length === 0 && (
                <p className="text-sm text-muted-foreground">No ideas yet. Add some on the Ideas page.</p>
              )}
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="text">
              {sourceType === "custom"
                ? "Your text"
                : sourceType === "all"
                  ? "Combined activity (edit before generating)"
                  : "Preview (you can edit before generating)"}
            </Label>
            <textarea
              id="text"
              className="flex min-h-[120px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
              placeholder={
                sourceType === "custom"
                  ? "Describe what you learned, a bug you fixed, a concept..."
                  : sourceType === "all"
                    ? "Select “From all my activity” to load combined logs."
                    : "Select an item above or type here to add more context."
              }
              value={text}
              onChange={(e) => setText(e.target.value)}
              readOnly={sourceType === "all" && loadingAll}
            />
          </div>

          <div className="flex flex-wrap gap-4">
            <label className="flex items-center gap-2">
              <Checkbox checked={wantTweet} onCheckedChange={(c) => setWantTweet(!!c)} />
              <span className="text-sm">Tweet</span>
            </label>
            <label className="flex items-center gap-2">
              <Checkbox checked={wantReel} onCheckedChange={(c) => setWantReel(!!c)} />
              <span className="text-sm">Reel script</span>
            </label>
            <label className="flex items-center gap-2">
              <Checkbox checked={wantHook} onCheckedChange={(c) => setWantHook(!!c)} />
              <span className="text-sm">Hook</span>
            </label>
            <label className="flex items-center gap-2">
              <Checkbox checked={wantLinkedin} onCheckedChange={(c) => setWantLinkedin(!!c)} />
              <span className="text-sm">LinkedIn post</span>
            </label>
          </div>

          {error && <p className="text-sm text-destructive">{error}</p>}
          <Button type="submit" disabled={loading || loadingAll || !canGenerate} onClick={handleGenerate}>
            <Sparkles className="h-4 w-4" />
            {loading ? "Generating…" : "Generate"}
          </Button>
        </CardContent>
      </Card>

      {result && (result.tweet || result.reel_script || result.hook || result.linkedin_post) && (
        <Card className="border-border">
          <CardHeader>
            <CardTitle className="font-display">Generated content</CardTitle>
            <CardDescription>Copy any section to use in your posts</CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {result.hook && (
              <div className="rounded-lg border border-border bg-card p-4">
                <div className="flex items-center justify-between gap-2">
                  <span className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Hook</span>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 shrink-0"
                    onClick={() => copyToClipboard(result.hook!, "hook")}
                  >
                    {copiedField === "hook" ? <Check className="h-4 w-4 text-green-600" /> : <Copy className="h-4 w-4" />}
                  </Button>
                </div>
                <p className="mt-2 text-sm leading-relaxed">{result.hook}</p>
              </div>
            )}
            {result.tweet && (
              <div className="rounded-lg border border-border bg-card p-4">
                <div className="flex items-center justify-between gap-2">
                  <span className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Tweet</span>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 shrink-0"
                    onClick={() => copyToClipboard(result.tweet!, "tweet")}
                  >
                    {copiedField === "tweet" ? <Check className="h-4 w-4 text-green-600" /> : <Copy className="h-4 w-4" />}
                  </Button>
                </div>
                <p className="mt-2 text-sm leading-relaxed">{result.tweet}</p>
              </div>
            )}
            {result.reel_script && (
              <div className="rounded-lg border border-border bg-card p-4">
                <div className="flex items-center justify-between gap-2">
                  <span className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">Reel script</span>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 shrink-0"
                    onClick={() => copyToClipboard(result.reel_script!, "reel_script")}
                  >
                    {copiedField === "reel_script" ? <Check className="h-4 w-4 text-green-600" /> : <Copy className="h-4 w-4" />}
                  </Button>
                </div>
                <div className="mt-2 space-y-2">
                  {reelScriptSteps(result.reel_script).map((step, i) => (
                    <div key={i} className="flex gap-2 text-sm">
                      <span className="shrink-0 font-medium text-muted-foreground">{i + 1}.</span>
                      <p className="min-w-0 flex-1 leading-relaxed">{step}</p>
                    </div>
                  ))}
                </div>
              </div>
            )}
            {result.linkedin_post && (
              <div className="rounded-lg border border-border bg-card p-4">
                <div className="flex items-center justify-between gap-2">
                  <span className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">LinkedIn post</span>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 shrink-0"
                    onClick={() => copyToClipboard(result.linkedin_post!, "linkedin_post")}
                  >
                    {copiedField === "linkedin_post" ? <Check className="h-4 w-4 text-green-600" /> : <Copy className="h-4 w-4" />}
                  </Button>
                </div>
                <p className="mt-2 whitespace-pre-line text-sm leading-relaxed">{result.linkedin_post}</p>
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  )
}
