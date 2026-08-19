import type { Agent, RegisterAgentInput, VanguardResponse } from './types'

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...init?.headers,
    },
  })
  const body = (await response.json()) as VanguardResponse<T>
  if (!response.ok) {
    throw new Error(body.error?.errorDesc || body.description || 'Vanguard request failed')
  }
  return body.data
}

export function listAgents(signal?: AbortSignal): Promise<Agent[]> {
  return request<Agent[]>('/api/v1/agents', { signal })
}

export function registerAgent(input: RegisterAgentInput): Promise<Agent> {
  return request<Agent>('/api/v1/agents', {
    method: 'POST',
    body: JSON.stringify(input),
  })
}
