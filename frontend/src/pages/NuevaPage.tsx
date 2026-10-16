export function NuevaPage() {
  return (
    <div className="nueva-page">
      <div className="nueva-grid">
        <section className="panel form-panel">
          <h2>Nueva candidatura</h2>
          <p className="muted">
            Formulario (estilo V1) en <strong>4c</strong>. Editar irá a{" "}
            <code>/nueva?id=…</code>.
          </p>
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
