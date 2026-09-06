export interface Agent {
  id: string
  identifier: string
  name: string
  framework: string
  metadata: Record<string, unknown>
  firstSeenAt: string
  lastSeenAt: string
}

export interface RegisterAgentInput {
  identifier: string
  name: string
  framework: string
  metadata: Record<string, unknown>
}

export interface VanguardError {
  errorCode: string
  errorDesc: string
}

export interface VanguardResponse<T> {
  data: T
  code: string
  description: string
  error: VanguardError | null
}
