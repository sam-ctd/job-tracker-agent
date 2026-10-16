import type { Application, ApplicationInput, ChatReply, Stats } from "../types"

// Docker: VITE_API_URL="" → same-origin (/api via nginx). Dev: frontend/.env → :8088.
function resolveBaseURL(): string {
  const raw = import.meta.env.VITE_API_URL as string | undefined
  if (raw === undefined) return "http://localhost:8088"
  return String(raw).replace(/\/$/, "")
}

const baseURL = resolveBaseURL()
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

export function createApplication(input: ApplicationInput) {
  return request<Application>("/api/applications", {
    method: "POST",
    body: JSON.stringify(toCreatePayload(input)),
  })
}

export function patchApplication(id: string, input: Partial<ApplicationInput>) {
  const keys = Object.keys(input)
  const body =
    keys.length === 1 && keys[0] === "status"
      ? { status: input.status }
      : toPatchPayload(input)
  return request<Application>(`/api/applications/${id}`, {
    method: "PATCH",
    body: JSON.stringify(body),
  })
}

export function deleteApplication(id: string) {
  return request<void>(`/api/applications/${id}`, { method: "DELETE" })
}

export function chat(message: string) {
  return request<ChatReply>("/api/chat", {
    method: "POST",
    body: JSON.stringify({ message }),
  })
}

function toCreatePayload(input: ApplicationInput) {
  const body: Record<string, unknown> = {
    company: input.company.trim(),
    title: input.title.trim(),
    url: input.url.trim(),
    channel: input.channel,
    status: input.status,
    notes: input.notes.trim(),
    contact: input.contact.trim(),
    location: input.location.trim(),
    remote: input.remote.trim(),
    stack: input.stack.trim(),
  }
  if (input.next_follow_up) {
    body.next_follow_up = new Date(input.next_follow_up).toISOString()
  }
  return body
}

function toPatchPayload(input: Partial<ApplicationInput>) {
  const body: Record<string, unknown> = {}
  if (input.company !== undefined) body.company = input.company.trim()
  if (input.title !== undefined) body.title = input.title.trim()
  if (input.url !== undefined) body.url = input.url.trim()
  if (input.channel !== undefined) body.channel = input.channel
  if (input.status !== undefined) body.status = input.status
  if (input.notes !== undefined) body.notes = input.notes.trim()
  if (input.contact !== undefined) body.contact = input.contact.trim()
  if (input.location !== undefined) body.location = input.location.trim()
  if (input.remote !== undefined) body.remote = input.remote.trim()
  if (input.stack !== undefined) body.stack = input.stack.trim()
  if (input.next_follow_up !== undefined) {
    if (input.next_follow_up === "") body.clear_next_follow_up = true
    else body.next_follow_up = new Date(input.next_follow_up).toISOString()
  }
  return body
}

export { baseURL }
