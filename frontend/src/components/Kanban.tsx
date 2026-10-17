import { useState, type DragEvent } from "react"
import type { Application } from "../types"
import { channelLabel, priorityLabel, priorityRank, STATUSES } from "../types"

function sortByPriority(apps: Application[]): Application[] {
  return [...apps].sort((a, b) => {
    const pr = priorityRank(a.priority) - priorityRank(b.priority)
    if (pr !== 0) return pr
    return new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime()
  })
}

type Props = {
  apps: Application[]
  onEdit: (app: Application) => void
  onDelete: (id: string) => void
  onMove: (id: string, status: string) => Promise<void>
  loading: boolean
}

export function Kanban({ apps, onEdit, onDelete, onMove, loading }: Props) {
  const [dragOver, setDragOver] = useState<string | null>(null)
  const [draggingId, setDraggingId] = useState<string | null>(null)
  const [busyId, setBusyId] = useState<string | null>(null)

  if (loading && apps.length === 0) {
    return <p className="muted">Cargando candidaturas…</p>
  }
  if (!loading && apps.length === 0) {
    return <p className="muted">No hay candidaturas todavía.</p>
  }

  function onDragStart(e: DragEvent, app: Application) {
    e.dataTransfer.setData("text/plain", app.id)
    e.dataTransfer.setData("application/x-status", app.status)
    e.dataTransfer.effectAllowed = "move"
    setDraggingId(app.id)
  }

  function onDragEnd() {
    setDraggingId(null)
    setDragOver(null)
  }

  async function onDropColumn(e: DragEvent, status: string) {
    e.preventDefault()
    const id = e.dataTransfer.getData("text/plain")
    const from = e.dataTransfer.getData("application/x-status")
    setDragOver(null)
    setDraggingId(null)
    if (!id || from === status) return
    setBusyId(id)
    try {
      await onMove(id, status)
    } finally {
      setBusyId(null)
    }
  }

  return (
    <div className="board">
      {STATUSES.map((status) => {
        const col = sortByPriority(apps.filter((a) => a.status === status))
        return (
          <section
            key={status}
            className={`board-col${dragOver === status ? " is-drop-target" : ""}`}
            onDragOver={(e) => {
              e.preventDefault()
              e.dataTransfer.dropEffect = "move"
              if (dragOver !== status) setDragOver(status)
            }}
            onDragLeave={() => {
              if (dragOver === status) setDragOver(null)
            }}
            onDrop={(e) => void onDropColumn(e, status)}
          >
            <h3 className="board-col-title">
              {status}{" "}
              <span className={`board-col-count status-${status}`}>{col.length}</span>
            </h3>
            <ul
              className={`board-cards${col.length === 0 ? " is-empty" : ""}`}
              onDragOver={(e) => {
                e.preventDefault()
                e.stopPropagation()
                e.dataTransfer.dropEffect = "move"
                if (dragOver !== status) setDragOver(status)
              }}
              onDrop={(e) => void onDropColumn(e, status)}
            >
              {col.map((app) => (
                <li
                  key={app.id}
                  className={`board-card${draggingId === app.id ? " is-dragging" : ""}${busyId === app.id ? " is-busy" : ""}`}
                  draggable={busyId !== app.id}
                  onDragStart={(e) => onDragStart(e, app)}
                  onDragEnd={onDragEnd}
                >
                  <div className="board-card-top">
                    <span className="channel-tag">{channelLabel(app.channel)}</span>
                    <span
                      className={`priority-badge priority-${app.priority || "media"}`}
                      title={`Prioridad ${priorityLabel(app.priority || "media")}`}
                    >
                      <span className={`priority-dot priority-${app.priority || "media"}`} />
                      {priorityLabel(app.priority || "media")}
                    </span>
                  </div>

                  <button type="button" className="board-card-title" onClick={() => onEdit(app)}>
                    <span className="company">{app.company}</span>
                    <span className="role">{app.title}</span>
                  </button>

                  <div className="board-meta">
                    {app.contact ? <span>{app.contact}</span> : null}
                    <span>
                      {app.next_follow_up
                        ? `Follow-up ${new Date(app.next_follow_up).toLocaleDateString()}`
                        : "Sin follow-up"}
                    </span>
                    {(app.location || app.remote || app.stack) && (
                      <span>
                        {[app.location, app.remote, app.stack].filter(Boolean).join(" · ")}
                      </span>
                    )}
                    <span className="board-created">
                      Creada {new Date(app.created_at).toLocaleDateString()}
                    </span>
                  </div>

                  <div className="board-card-foot">
                    <div
                      className="board-card-foot-right"
                      onMouseDown={(e) => e.stopPropagation()}
                      onDragStart={(e) => e.preventDefault()}
                    >
                      {app.url ? (
                        <a href={app.url} target="_blank" rel="noreferrer" className="btn-link" draggable={false}>
                          Link
                        </a>
                      ) : (
                        <span className="btn-link btn-disabled">Link</span>
                      )}
                      <button type="button" className="btn-edit" onClick={() => onEdit(app)}>
                        Editar
                      </button>
                      <button type="button" className="btn-delete" onClick={() => onDelete(app.id)}>
                        Borrar
                      </button>
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          </section>
        )
      })}
    </div>
  )
}
