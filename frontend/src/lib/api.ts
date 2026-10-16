import type { Application, Stats } from "../types"

const baseURL = (import.meta.env.VITE_API_URL as string | undefined)?.replace(/\/$/, "") || "http://localhost:8088"
const apiKey = (import.meta.env.VITE_API_KEY as string | undefined) || ""

export class ApiError extends Error {
  status: number
  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

async function request<T>(path: string, init: RequestInit = {}, auth = true): Promise<T> {
  const headers = new Headers(init.headers)
  if (!headers.has("Content-Type") && init.body) {
    headers.set("Content-Type", "application/json")
  }
  if (auth) {
    if (!apiKey) throw new ApiError(0, "VITE_API_KEY missing")
    headers.set("Authorization", `Bearer ${apiKey}`)
  }
  const res = await fetch(`${baseURL}${path}`, { ...init, headers })
  if (!res.ok) {
    let msg = res.statusText
    try {
      const body = (await res.json()) as { error?: string }
      if (body.error) msg = body.error
    } catch {
      /* ignore */
    }
    throw new ApiError(res.status, msg)
  }
  if (res.status === 204) return undefined as T
  return (await res.json()) as T
}

export type Health = { status: string }

export function getHealth() {
  return request<Health>("/health", {}, false)
}

export function getReady() {
  return request<Health>("/ready", {}, false)
}

export function listApplications(status = "", company = "") {
  const q = new URLSearchParams()
  if (status) q.set("status", status)
  if (company) q.set("company", company)
  const qs = q.toString()
  return request<Application[]>(`/api/applications${qs ? `?${qs}` : ""}`)
}

export function listDue() {
  return request<Application[]>("/api/applications/due")
}

export function getStats() {
  return request<Stats>("/api/stats")
}

export function getApplication(id: string) {
  return request<Application>(`/api/applications/${id}`)
}

export function patchApplication(id: string, body: Record<string, unknown>) {
  return request<Application>(`/api/applications/${id}`, {
    method: "PATCH",
    body: JSON.stringify(body),
  })
}

export function deleteApplication(id: string) {
  return request<void>(`/api/applications/${id}`, { method: "DELETE" })
}

export { baseURL }
