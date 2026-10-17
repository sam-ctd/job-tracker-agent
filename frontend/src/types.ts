export const STATUSES = [
  "idea",
  "enviada",
  "relance",
  "entretien",
  "offre",
  "refus",
] as const

export type Status = (typeof STATUSES)[number]

export const CHANNELS = [
  "linkedin",
  "wttj",
  "indeed",
  "hellowork",
  "apec",
  "jobteaser",
  "lesjeudis",
  "chooseyourboss",
  "alumni_estia",
] as const

export type Channel = (typeof CHANNELS)[number]

export const CHANNEL_LABELS: Record<string, string> = {
  linkedin: "LinkedIn",
  wttj: "Welcome to the Jungle",
  indeed: "Indeed",
  hellowork: "HelloWork",
  apec: "APEC",
  jobteaser: "JobTeaser",
  lesjeudis: "LesJeudis",
  chooseyourboss: "ChooseYourBoss",
  alumni_estia: "Alumni ESTIA",
}

export function channelLabel(channel: string): string {
  return CHANNEL_LABELS[channel] || channel
}

export const PRIORITIES = ["alta", "media", "baja"] as const

export type Priority = (typeof PRIORITIES)[number]

export const PRIORITY_LABELS: Record<Priority, string> = {
  alta: "Alta",
  media: "Media",
  baja: "Baja",
}

export function priorityLabel(priority: string): string {
  return PRIORITY_LABELS[priority as Priority] || priority
}

/** Sort key: alta → media → baja (unknown last). */
export function priorityRank(priority: string): number {
  const i = PRIORITIES.indexOf(priority as Priority)
  return i === -1 ? 99 : i
}

export type Application = {
  id: string
  company: string
  title: string
  url: string
  channel: Channel | string
  status: Status | string
  priority: Priority | string
  sent_at: string | null
  next_follow_up: string | null
  notes: string
  contact: string
  location: string
  remote: string
  stack: string
  created_at: string
  updated_at: string
}

export type ApplicationInput = {
  company: string
  title: string
  url: string
  channel: string
  status: string
  priority: string
  notes: string
  contact: string
  location: string
  remote: string
  stack: string
  next_follow_up: string
}

export type Stats = {
  by_status: Record<string, number>
  by_channel: Record<string, number>
  total: number
}

export type ChatMessage = {
  role: "user" | "assistant"
  content: string
  tool?: string
  brain?: string
}

export type ChatReply = {
  message: string
  tool?: string
  args?: Record<string, string>
  result?: unknown
  help_hint?: string
  brain?: string
}
