import { useEffect, useState, type FormEvent } from "react"
import type { Application, ApplicationInput } from "../types"
import { CHANNELS, CHANNEL_LABELS, STATUSES } from "../types"

const empty: ApplicationInput = {
  company: "",
  title: "",
  url: "",
  channel: "linkedin",
  status: "idea",
  notes: "",
  contact: "",
  location: "",
  remote: "",
  stack: "",
  next_follow_up: "",
}

type Props = {
  initial?: Application | null
  onSubmit: (input: ApplicationInput) => Promise<void>
  onCancel?: () => void
  busy: boolean
}

export function ApplicationForm({ initial, onSubmit, onCancel, busy }: Props) {
  const [form, setForm] = useState<ApplicationInput>(empty)
  const [error, setError] = useState("")

  useEffect(() => {
    if (!initial) {
      setForm(empty)
      return
    }
    setForm({
      company: initial.company,
      title: initial.title,
      url: initial.url,
      channel: initial.channel,
      status: initial.status,
      notes: initial.notes,
      contact: initial.contact,
      location: initial.location || "",
      remote: initial.remote || "",
      stack: initial.stack || "",
      next_follow_up: initial.next_follow_up ? initial.next_follow_up.slice(0, 16) : "",
    })
  }, [initial])

  function validate(input: ApplicationInput): string {
    if (!input.company.trim() || input.company.trim().length > 200) return "company inválido"
    if (!input.title.trim() || input.title.trim().length > 300) return "title inválido"
    if (input.url && !/^https?:\/\//i.test(input.url.trim())) return "url debe ser http(s)"
    if (!CHANNELS.includes(input.channel as (typeof CHANNELS)[number])) return "channel inválido"
    if (!STATUSES.includes(input.status as (typeof STATUSES)[number])) return "status inválido"
    return ""
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    const err = validate(form)
    if (err) {
      setError(err)
      return
    }
    setError("")
    await onSubmit(form)
    if (!initial) setForm(empty)
  }

  return (
    <form className="app-form" onSubmit={handleSubmit}>
      <section className="form-section">
        <div className="form-section-head">
          <h3>{initial ? "Editar candidatura" : "Nueva candidatura"}</h3>
        </div>
        <div className="grid-3">
          <label>
            Empresa
            <input
              value={form.company}
              onChange={(e) => setForm({ ...form, company: e.target.value })}
              required
              maxLength={200}
              placeholder="Acme SAS"
              autoComplete="organization"
            />
          </label>
          <label>
            Título
            <input
              value={form.title}
              onChange={(e) => setForm({ ...form, title: e.target.value })}
              required
              maxLength={300}
              placeholder="Stage Fullstack Go / React"
            />
          </label>
          <label>
            URL
            <input
              value={form.url}
              onChange={(e) => setForm({ ...form, url: e.target.value })}
              maxLength={2000}
              placeholder="https://…"
              inputMode="url"
            />
          </label>
        </div>
      </section>

      <section className="form-section">
        <div className="form-section-head">
          <h3>Pipeline</h3>
        </div>

        <div className="field-block">
          <span className="field-label">Canal</span>
          <div className="chip-row" role="group" aria-label="Canal">
            {CHANNELS.map((c) => (
              <button
                key={c}
                type="button"
                className={`chip${form.channel === c ? " is-active" : ""}`}
                onClick={() => setForm({ ...form, channel: c })}
              >
                {CHANNEL_LABELS[c] || c}
              </button>
            ))}
          </div>
        </div>

        <div className="field-block">
          <span className="field-label">Estado</span>
          <div className="status-picker" role="group" aria-label="Estado">
            {STATUSES.map((s) => (
              <button
                key={s}
                type="button"
                className={`status-option status-${s}${form.status === s ? " is-active" : ""}`}
                onClick={() => setForm({ ...form, status: s })}
              >
                <span className={`status-dot status-${s}`} />
                {s}
              </button>
            ))}
          </div>
        </div>

        <div className="grid-2">
          <label>
            Contacto
            <input
              value={form.contact}
              onChange={(e) => setForm({ ...form, contact: e.target.value })}
              maxLength={500}
              placeholder="Recruiter"
            />
          </label>
          <label>
            Follow-up
            <input
              type="datetime-local"
              value={form.next_follow_up}
              onChange={(e) => setForm({ ...form, next_follow_up: e.target.value })}
            />
          </label>
        </div>

        <div className="grid-3">
          <label>
            Location
            <input
              value={form.location}
              onChange={(e) => setForm({ ...form, location: e.target.value })}
              maxLength={200}
              placeholder="Paris / Biarritz"
            />
          </label>
          <label>
            Remote
            <input
              value={form.remote}
              onChange={(e) => setForm({ ...form, remote: e.target.value })}
              maxLength={100}
              placeholder="remote / hybrid / onsite"
            />
          </label>
          <label>
            Stack
            <input
              value={form.stack}
              onChange={(e) => setForm({ ...form, stack: e.target.value })}
              maxLength={300}
              placeholder="Go, React, Postgres"
            />
          </label>
        </div>
      </section>

      <section className="form-section form-section-notes">
        <div className="form-section-head">
          <h3>Notas</h3>
        </div>
        <label className="notes-label">
          <textarea
            value={form.notes}
            onChange={(e) => setForm({ ...form, notes: e.target.value })}
            maxLength={5000}
            rows={3}
            placeholder="Contexto, recordatorios…"
          />
        </label>
      </section>

      {error && <p className="error">{error}</p>}

      <div className="form-actions">
        {initial && (
          <span className="form-created muted">
            Creada {new Date(initial.created_at).toLocaleDateString()}
          </span>
        )}
        <button type="submit" disabled={busy}>
          {busy ? "Guardando…" : initial ? "Guardar cambios" : "Crear candidatura"}
        </button>
        {onCancel && (
          <button type="button" className="ghost" onClick={onCancel} disabled={busy}>
            Cancelar
          </button>
        )}
      </div>
    </form>
  )
}
