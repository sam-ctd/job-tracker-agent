import { useRef, useState, type FormEvent } from "react"
import { chat } from "../lib/api"
import type { ChatMessage } from "../types"

const SUGGESTIONS = [
  {
    label: "Resumen",
    prompt: "Dame un resumen de stats por estado y canal",
  },
  {
    label: "Follow-ups",
    prompt: "¿Qué follow-ups tengo pendientes?",
  },
  {
    label: "Lista",
    prompt: "Lista las candidaturas",
  },
  {
    label: "Alta rápida",
    prompt: "crea Acme | Stage Go",
  },
]

type Props = {
  onChanged?: () => void
}

export function ChatPanel({ onChanged }: Props) {
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [input, setInput] = useState("")
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState("")
  const logRef = useRef<HTMLDivElement | null>(null)

  function scrollLog() {
    requestAnimationFrame(() => {
      const el = logRef.current
      if (el) el.scrollTop = el.scrollHeight
    })
  }

  async function runPrompt(text: string) {
    const trimmed = text.trim()
    if (!trimmed || busy) return

    setInput("")
    setError("")
    setMessages((m) => [...m, { role: "user", content: trimmed }])
    setBusy(true)
    scrollLog()

    try {
      const reply = await chat(trimmed)
      let content = reply.message
      if (reply.help_hint) content += `\n\n${reply.help_hint}`
      setMessages((m) => [
        ...m,
        {
          role: "assistant",
          content,
          tool: reply.tool,
          brain: reply.brain,
        },
      ])
      if (reply.tool === "create_application" || reply.tool === "patch_application") {
        onChanged?.()
      }
      scrollLog()
    } catch (err) {
      setError((err as Error).message || "chat error")
    } finally {
      setBusy(false)
    }
  }

  async function send(e: FormEvent) {
    e.preventDefault()
    await runPrompt(input)
  }

  return (
    <aside className="panel chat">
      <div className="chat-header">
        <div className="chat-brand">
          <div className="chat-avatar" aria-hidden>
            AI
          </div>
          <div>
            <h2>Agente</h2>
            <p className="muted">Asistente del tracker</p>
          </div>
        </div>
      </div>

      <div className="chat-log" ref={logRef}>
        <div className={`chat-empty${messages.length > 0 ? " is-compact" : ""}`}>
          {messages.length === 0 && (
            <>
              <p className="chat-empty-title">¿En qué te ayudo?</p>
              <p className="chat-empty-text">
                Puedes pedir stats, follow-ups, listar o crear una candidatura en una frase.
              </p>
            </>
          )}
          {messages.length > 0 && <p className="chat-suggestions-heading">Sugerencias</p>}
          <div className="chat-suggestions">
            {SUGGESTIONS.map((s) => (
              <button
                key={s.label}
                type="button"
                className="suggestion-chip"
                onClick={() => void runPrompt(s.prompt)}
                disabled={busy}
              >
                <span className="suggestion-label">{s.label}</span>
                <span className="suggestion-prompt">{s.prompt}</span>
              </button>
            ))}
          </div>
        </div>
        {messages.map((m, i) => (
          <div key={i} className={`bubble ${m.role}`}>
            <span className="bubble-role">
              {m.role === "user" ? "Tú" : "Agente"}
              {m.tool ? ` · ${m.tool}` : ""}
              {m.brain ? ` · ${m.brain}` : ""}
            </span>
            <pre>{m.content || (busy && i === messages.length - 1 ? "…" : "")}</pre>
          </div>
        ))}
        {busy && (
          <div className="bubble assistant">
            <span className="bubble-role">Agente</span>
            <pre>…</pre>
          </div>
        )}
      </div>

      {error && <p className="error">{error}</p>}

      <form className="chat-input" onSubmit={send}>
        <input
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Ej.: muéstrame las stats…"
          disabled={busy}
        />
        <button type="submit" disabled={busy || !input.trim()}>
          Enviar
        </button>
      </form>
    </aside>
  )
}
