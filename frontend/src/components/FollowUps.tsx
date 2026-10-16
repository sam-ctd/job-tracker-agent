import type { Application } from "../types"
import { channelLabel } from "../types"

type Props = {
  due: Application[]
  loading: boolean
  onEdit: (app: Application) => void
}

export function FollowUps({ due, loading, onEdit }: Props) {
  return (
    <section className="followups-board">
      <div className="followups-head">
        <div>
          <h2>Follow-ups pendientes</h2>
          <p className="muted">Candidaturas con relance vencido o para hoy</p>
        </div>
        <span className="followups-count">{due.length}</span>
      </div>

      {loading && due.length === 0 && <p className="muted">Cargando…</p>}
      {!loading && due.length === 0 && (
        <div className="followups-empty">Nada vencido por ahora.</div>
      )}

      {due.length > 0 && (
        <ul className="followups-grid">
          {due.map((app) => (
            <li key={app.id} className="followups-card">
              <div className="followups-card-top">
                <span className="channel-tag">{channelLabel(app.channel)}</span>
                <span className={`status-dot status-${app.status}`} title={app.status} />
              </div>
              <div className="followups-card-body">
                <strong>{app.company}</strong>
                <span>{app.title}</span>
              </div>
              <div className="followups-card-foot">
                <time>
                  {app.next_follow_up
                    ? new Date(app.next_follow_up).toLocaleDateString()
                    : "—"}
                </time>
                <button type="button" className="btn-edit" onClick={() => onEdit(app)}>
                  Abrir
                </button>
              </div>
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}
