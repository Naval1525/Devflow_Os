import { useEffect, useState } from "react"
import { apiFetch } from "@/lib/api"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select } from "@/components/ui/select"
import { Trash2 } from "lucide-react"

type FinanceEntry = {
  id: string
  amount: number
  type: string
  note: string
  date: string
  created_at: string
}

const FINANCE_TYPES = ["salary", "freelance", "spend", "other"] as const
const ENTRY_FILTERS = ["all", "income", "expense"] as const
type EntryFilter = (typeof ENTRY_FILTERS)[number]

function getMonthKey(d: Date) {
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}`
}

function parseDate(s: string) {
  const [y, m] = s.split("-").map(Number)
  return { year: y, month: m }
}

function isExpenseType(type: string) {
  return type === "other" || type === "spend"
}

export function Finance() {
  const [entries, setEntries] = useState<FinanceEntry[]>([])
  const [loading, setLoading] = useState(true)
  const [amount, setAmount] = useState("")
  const [financeType, setFinanceType] = useState<string>("freelance")
  const [note, setNote] = useState("")
  const [date, setDate] = useState(() => new Date().toISOString().slice(0, 10))
  const [entryFilter, setEntryFilter] = useState<EntryFilter>("all")

  async function load() {
    setLoading(true)
    try {
      const res = await apiFetch("/finances")
      if (res.ok) {
        const data = await res.json()
        setEntries(data.finances ?? [])
      }
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  const now = new Date()
  const thisMonthKey = getMonthKey(now)
  const lastMonth = new Date(now.getFullYear(), now.getMonth() - 1, 1)
  const lastMonthKey = getMonthKey(lastMonth)

  const thisMonthTotal = entries
    .filter((e) => {
      const { year, month } = parseDate(e.date)
      return `${year}-${String(month).padStart(2, "0")}` === thisMonthKey
    })
    .reduce((sum, e) => sum + e.amount, 0)
  const lastMonthTotal = entries
    .filter((e) => {
      const { year, month } = parseDate(e.date)
      return `${year}-${String(month).padStart(2, "0")}` === lastMonthKey
    })
    .reduce((sum, e) => sum + e.amount, 0)
  const delta = thisMonthTotal - lastMonthTotal

  async function handleAdd(e: React.FormEvent) {
    e.preventDefault()
    const num = parseFloat(amount)
    if (Number.isNaN(num)) return
    const res = await apiFetch("/finances", {
      method: "POST",
      body: JSON.stringify({ amount: num, type: financeType, note: note.trim(), date }),
    })
    if (res.ok) {
      setAmount("")
      setNote("")
      setDate(new Date().toISOString().slice(0, 10))
      load()
    }
  }

  async function handleDelete(id: string) {
    const confirmed = window.confirm("Delete this finance entry?")
    if (!confirmed) return
    const res = await apiFetch(`/finances/${id}`, { method: "DELETE" })
    if (res.ok) {
      setEntries((prev) => prev.filter((entry) => entry.id !== id))
    }
  }

  const formatRupees = (n: number) =>
    `₹${n.toLocaleString("en-IN", { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`

  const filteredEntries = entries.filter((entry) => {
    const isExpense = isExpenseType(entry.type) || entry.amount < 0
    if (entryFilter === "all") return true
    if (entryFilter === "expense") return isExpense
    return !isExpense
  })

  return (
    <div className="space-y-8">
      <div>
        <h1 className="font-display text-2xl font-semibold">Finance</h1>
        <p className="text-muted-foreground">Track income and spending in rupees and compare months</p>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card className="border-border">
          <CardHeader>
            <CardTitle className="font-display">This month</CardTitle>
            <CardDescription>Total income (₹ INR)</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-semibold">{formatRupees(thisMonthTotal)}</p>
            <p className="text-sm text-muted-foreground">
              vs last month: {formatRupees(lastMonthTotal)}
              {delta >= 0 ? ` (+${formatRupees(delta)})` : ` (${formatRupees(delta)})`}
            </p>
          </CardContent>
        </Card>
        <Card className="border-border">
          <CardHeader>
            <CardTitle className="font-display">Add entry</CardTitle>
            <CardDescription>Amount in ₹, type, date, note</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleAdd} className="space-y-3">
              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-2">
                  <Label htmlFor="amount">Amount (₹)</Label>
                  <Input id="amount" type="number" step="0.01" value={amount} onChange={(e) => setAmount(e.target.value)} placeholder="e.g. 50000" required />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="fin-type">Type</Label>
                  <Select id="fin-type" value={financeType} onChange={(e) => setFinanceType(e.target.value)}>
                    {FINANCE_TYPES.map((t) => <option key={t} value={t}>{t}</option>)}
                  </Select>
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="date">Date</Label>
                <Input id="date" type="date" value={date} onChange={(e) => setDate(e.target.value)} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="note">Note</Label>
                <Input id="note" value={note} onChange={(e) => setNote(e.target.value)} placeholder="Optional" />
              </div>
              <Button type="submit">Add</Button>
            </form>
          </CardContent>
        </Card>
      </div>

      <Card className="border-border">
        <CardHeader>
          <CardTitle className="font-display">Entries</CardTitle>
          <CardDescription>Recent finance entries</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="mb-4 max-w-[220px] space-y-2">
            <Label htmlFor="entry-filter">Filter entries</Label>
            <Select
              id="entry-filter"
              value={entryFilter}
              onChange={(e) => setEntryFilter(e.target.value as EntryFilter)}
            >
              <option value="all">All</option>
              <option value="income">Income only</option>
              <option value="expense">Expense only</option>
            </Select>
          </div>
          {loading ? (
            <p className="text-muted-foreground">Loading…</p>
          ) : filteredEntries.length === 0 ? (
            <p className="text-muted-foreground">No entries yet.</p>
          ) : (
            <ul className="space-y-2">
              {filteredEntries.map((e) => (
                <li key={e.id} className="flex items-center justify-between rounded-lg border border-border p-3">
                  <div>
                    <p className={`font-medium ${isExpenseType(e.type) || e.amount < 0 ? "text-red-600" : "text-green-600"}`}>
                      {formatRupees(e.amount)} · {e.type}
                    </p>
                    <p className="text-sm text-muted-foreground">{e.date}{e.note ? ` · ${e.note}` : ""}</p>
                  </div>
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    aria-label="Delete finance entry"
                    onClick={() => handleDelete(e.id)}
                  >
                    <Trash2 className="h-4 w-4 text-destructive" />
                  </Button>
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
