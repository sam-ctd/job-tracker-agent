import { useEffect, useState } from "react"
import { useNavigate, useSearchParams } from "react-router-dom"
import { ApplicationForm } from "../components/ApplicationForm"
import { createApplication, getApplication, patchApplication } from "../lib/api"
import type { Application, ApplicationInput } from "../types"

export function NuevaPage() {
  const [params] = useSearchParams()
  const id = params.get("id")
  const navigate = useNavigate()
  const [editing, setEditing] = useState<Application | null>(null)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState("")
  const [loadingEdit, setLoadingEdit] = useState(Boolean(id))

  useEffect(() => {
    if (!id) {
      setEditing(null)
      setLoadingEdit(false)
      return
    }
    let cancelled = false
    setLoadingEdit(true)
    getApplication(id)
      .then((app) => {
        if (!cancelled) setEditing(app)
      })
      .catch((e) => {
        if (!cancelled) setError((e as Error).message)
      })
      .finally(() => {
        if (!cancelled) setLoadingEdit(false)
      })
    return () => {
      cancelled = true
    }
  }, [id])

  async function handleSubmit(input: ApplicationInput) {
    setBusy(true)
    setError("")
    try {
      if (editing) {
        await patchApplication(editing.id, input)
      } else {
        await createApplication(input)
      }
      navigate("/kanban")
    } catch (e) {
      setError((e as Error).message)
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="nueva-page">
      <div className="nueva-grid">
        <section className="panel form-panel">
          {error && <p className="banner error">{error}</p>}
          {loadingEdit ? (
            <p className="muted">Cargando candidatura…</p>
          ) : (
            <ApplicationForm
              initial={editing}
              onSubmit={handleSubmit}
              onCancel={editing ? () => navigate("/kanban") : undefined}
              busy={busy}
            />
          )}
        </section>
        <section className="panel">
          <h2>Chat</h2>
          <p className="muted">
            Panel del agente (<code>POST /api/chat</code>) en <strong>4d</strong>.
          </p>
        </section>
      </div>
    </div>
  )
}
