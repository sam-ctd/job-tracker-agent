import { useEffect, useState } from "react"
import { NavLink, Outlet } from "react-router-dom"
import { baseURL, getHealth } from "../lib/api"

type ApiState = "checking" | "ok" | "down"

export function Layout() {
  const [api, setApi] = useState<ApiState>("checking")

  useEffect(() => {
    let cancelled = false
    getHealth()
      .then(() => {
        if (!cancelled) setApi("ok")
      })
      .catch(() => {
        if (!cancelled) setApi("down")
      })
    return () => {
      cancelled = true
    }
  }, [])

  return (
    <div className="shell shell-wide">
      <header className="topbar reveal">
        <NavLink to="/kanban" className="brand">
          <div className="brand-icon">JT</div>
          <div className="brand-text">
            <div className="eyebrow">Stage tracker</div>
            <h1>Job Tracker Agent</h1>
          </div>
        </NavLink>
        <div className="topbar-right">
          <nav className="nav-pills" aria-label="Principal">
            <NavLink to="/nueva" className={({ isActive }) => (isActive ? "active" : "")}>
              Nueva
            </NavLink>
            <NavLink to="/kanban" className={({ isActive }) => (isActive ? "active" : "")}>
              Kanban
            </NavLink>
          </nav>
          <div className={`api-pill api-${api}`} title={baseURL}>
            <span className="api-dot" aria-hidden />
            {api === "checking" && "API…"}
            {api === "ok" && "API ok"}
            {api === "down" && "API down"}
          </div>
        </div>
      </header>
      {!import.meta.env.VITE_API_KEY && (
        <p className="banner error">Falta VITE_API_KEY en el entorno del frontend.</p>
      )}
      <Outlet />
    </div>
  )
}
