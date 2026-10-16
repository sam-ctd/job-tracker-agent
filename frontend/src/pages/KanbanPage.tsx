import { useCallback, useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import { FollowUps } from "../components/FollowUps"
import { Kanban } from "../components/Kanban"
import { deleteApplication, listApplications, listDue, patchApplication } from "../lib/api"
import type { Application } from "../types"

export function KanbanPage() {
  const navigate = useNavigate()
  const [apps, setApps] = useState<Application[]>([])
  const [due, setDue] = useState<Application[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")

  const refresh = useCallback(async () => {
    setError("")
    try {
      const [a, d] = await Promise.all([listApplications(), listDue()])
      setApps(a)
      setDue(d)
    } catch (e) {
      setError((e as Error).message || "error de carga")
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void refresh()
  }, [refresh])

  async function handleDelete(id: string) {
    if (!window.confirm("¿Borrar esta candidatura?")) return
    setError("")
    try {
      await deleteApplication(id)
      await refresh()
    } catch (e) {
      setError((e as Error).message)
    }
  }

  function goEdit(app: Application) {
    navigate(`/nueva?id=${app.id}`)
  }

  async function handleMove(id: string, status: string) {
    setError("")
    const prevApps = apps
    const prevDue = due
    setApps((list) => list.map((a) => (a.id === id ? { ...a, status } : a)))
    try {
      await patchApplication(id, { status })
      setDue(await listDue())
    } catch (e) {
      setApps(prevApps)
      setDue(prevDue)
      setError((e as Error).message)
    }
  }

  return (
    <div className="kanban-stack reveal">
      {error && <p className="banner error">{error}</p>}
      <FollowUps due={due} loading={loading} onEdit={goEdit} />
      <Kanban
        apps={apps}
        loading={loading}
        onEdit={goEdit}
        onDelete={handleDelete}
        onMove={handleMove}
      />
    </div>
  )
}
