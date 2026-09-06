import { FormEvent, useEffect, useMemo, useState } from 'react'
import { listAgents, registerAgent } from './api'
import type { Agent, RegisterAgentInput } from './types'

const navItems = [
  { label: 'Overview', icon: 'grid' },
  { label: 'Agent registry', icon: 'agents', active: true },
  { label: 'Permission requests', icon: 'shield' },
  { label: 'Policies', icon: 'policy' },
  { label: 'Audit log', icon: 'audit' },
]

function Icon({ name }: { name: string }) {
  const paths: Record<string, React.ReactNode> = {
    grid: <><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></>,
    agents: <><circle cx="12" cy="8" r="4"/><path d="M4.5 21a7.5 7.5 0 0 1 15 0"/></>,
    shield: <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10Z"/>,
    policy: <><path d="M6 3h12a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2Z"/><path d="M8 8h8M8 12h8M8 16h5"/></>,
    audit: <><path d="M3 3v18h18"/><path d="m7 15 4-4 3 3 5-7"/></>,
    search: <><circle cx="11" cy="11" r="7"/><path d="m20 20-4-4"/></>,
    plus: <path d="M12 5v14M5 12h14"/>,
    refresh: <><path d="M20 7h-6V1"/><path d="M20 7a9 9 0 1 0 1 8"/></>,
    close: <path d="m6 6 12 12M18 6 6 18"/>,
    chevron: <path d="m9 18 6-6-6-6"/>,
  }
  return <svg className="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">{paths[name]}</svg>
}

function relativeTime(value: string) {
  const seconds = Math.max(0, Math.floor((Date.now() - new Date(value).getTime()) / 1000))
  if (seconds < 15) return 'just now'
  if (seconds < 60) return `${seconds}s ago`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)}h ago`
  return `${Math.floor(seconds / 86400)}d ago`
}

function isRecentlyActive(agent: Agent) {
  return Date.now() - new Date(agent.lastSeenAt).getTime() < 10 * 60 * 1000
}

function App() {
  const [agents, setAgents] = useState<Agent[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [query, setQuery] = useState('')
  const [selected, setSelected] = useState<Agent | null>(null)
  const [showRegister, setShowRegister] = useState(false)

  const loadAgents = async (signal?: AbortSignal) => {
    setLoading(true)
    setError('')
    try {
      setAgents(await listAgents(signal))
    } catch (err) {
      if (!(err instanceof DOMException && err.name === 'AbortError')) {
        setError(err instanceof Error ? err.message : 'Unable to reach Vanguard')
      }
    } finally {
      if (!signal?.aborted) setLoading(false)
    }
  }

  useEffect(() => {
    const controller = new AbortController()
    void loadAgents(controller.signal)
    return () => controller.abort()
  }, [])

  const filteredAgents = useMemo(() => {
    const normalized = query.trim().toLowerCase()
    if (!normalized) return agents
    return agents.filter((agent) => [agent.name, agent.identifier, agent.framework].some((value) => value.toLowerCase().includes(normalized)))
  }, [agents, query])

  const activeCount = agents.filter(isRecentlyActive).length
  const frameworkCount = new Set(agents.map((agent) => agent.framework).filter(Boolean)).size

  const handleRegistered = (agent: Agent) => {
    setAgents((current) => {
      const exists = current.some((item) => item.id === agent.id)
      return exists ? current.map((item) => item.id === agent.id ? agent : item) : [...current, agent]
    })
    setShowRegister(false)
    setSelected(agent)
  }

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand"><div className="brand-mark">V</div><div><strong>Vanguard</strong><span>Control plane</span></div></div>
        <nav className="nav" aria-label="Main navigation">
          <p className="nav-label">Workspace</p>
          {navItems.map((item) => <button key={item.label} className={`nav-item ${item.active ? 'active' : ''}`} disabled={!item.active}><Icon name={item.icon}/><span>{item.label}</span>{item.label === 'Permission requests' && <span className="nav-badge">0</span>}</button>)}
        </nav>
        <div className="sidebar-footer"><span className="status-dot"/><div><strong>Vanguard API</strong><span>{error ? 'Connection issue' : 'Connected'}</span></div></div>
      </aside>

      <main className="main">
        <header className="topbar"><div className="breadcrumbs"><span>Control plane</span><b>/</b><strong>Agent registry</strong></div><div className="environment"><span className="status-dot"/>Development</div></header>
        <div className="content">
          <section className="page-heading"><div><p className="eyebrow">IDENTITY & INVENTORY</p><h1>Agent registry</h1><p>Discover and inspect every agent that has interacted with your control plane.</p></div><button className="primary-button" onClick={() => setShowRegister(true)}><Icon name="plus"/>Register agent</button></section>

          <section className="stats" aria-label="Registry summary">
            <div className="stat-card"><span>Registered agents</span><strong>{agents.length}</strong><small>All known identities</small></div>
            <div className="stat-card"><span>Recently active</span><strong>{activeCount}</strong><small>Seen in the last 10 minutes</small></div>
            <div className="stat-card"><span>Frameworks</span><strong>{frameworkCount}</strong><small>Distinct integrations</small></div>
          </section>

          <section className="registry-panel">
            <div className="panel-toolbar"><div><h2>All agents</h2><span>{filteredAgents.length} {filteredAgents.length === 1 ? 'record' : 'records'}</span></div><div className="toolbar-actions"><label className="search"><Icon name="search"/><input value={query} onChange={(event) => setQuery(event.target.value)} placeholder="Search agents…" aria-label="Search agents"/></label><button className="icon-button" onClick={() => void loadAgents()} aria-label="Refresh agents"><Icon name="refresh"/></button></div></div>

            {error && <div className="error-state"><div><strong>Registry unavailable</strong><span>{error}</span></div><button onClick={() => void loadAgents()}>Try again</button></div>}
            {!error && loading && <div className="loading-state"><span/><span/><span/></div>}
            {!error && !loading && agents.length === 0 && <div className="empty-state"><div className="empty-orbit"><div className="brand-mark">V</div></div><h3>No agents registered yet</h3><p>Agents will appear here when Sentinel identifies them, or you can register one manually.</p><button className="secondary-button" onClick={() => setShowRegister(true)}><Icon name="plus"/>Register your first agent</button></div>}
            {!error && !loading && agents.length > 0 && filteredAgents.length === 0 && <div className="empty-state compact"><h3>No matching agents</h3><p>Try a different name, identifier, or framework.</p></div>}
            {!error && !loading && filteredAgents.length > 0 && <div className="table-wrap"><table><thead><tr><th>Agent</th><th>Identifier</th><th>Framework</th><th>Status</th><th>Last seen</th><th aria-label="Actions"/></tr></thead><tbody>{filteredAgents.map((agent) => <tr key={agent.id} onClick={() => setSelected(agent)}><td><div className="agent-cell"><div className="agent-avatar">{agent.name.slice(0, 2).toUpperCase()}</div><div><strong>{agent.name}</strong><span>{agent.id.slice(0, 8)}</span></div></div></td><td><code>{agent.identifier}</code></td><td><span className="framework-pill">{agent.framework || 'Unknown'}</span></td><td><span className={`state ${isRecentlyActive(agent) ? 'online' : ''}`}><i/>{isRecentlyActive(agent) ? 'Active' : 'Idle'}</span></td><td title={new Date(agent.lastSeenAt).toLocaleString()}>{relativeTime(agent.lastSeenAt)}</td><td><button className="row-action" aria-label={`View ${agent.name}`}><Icon name="chevron"/></button></td></tr>)}</tbody></table></div>}
          </section>
        </div>
      </main>

      {showRegister && <RegisterDialog onClose={() => setShowRegister(false)} onRegistered={handleRegistered}/>} 
      {selected && <AgentDrawer agent={selected} onClose={() => setSelected(null)}/>} 
    </div>
  )
}

function RegisterDialog({ onClose, onRegistered }: { onClose: () => void; onRegistered: (agent: Agent) => void }) {
  const [form, setForm] = useState({ name: '', identifier: '', framework: '', metadata: '{}' })
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')

  const submit = async (event: FormEvent) => {
    event.preventDefault()
    setError('')
    setSubmitting(true)
    try {
      const metadata = JSON.parse(form.metadata) as Record<string, unknown>
      const input: RegisterAgentInput = { name: form.name.trim(), identifier: form.identifier.trim(), framework: form.framework.trim(), metadata }
      onRegistered(await registerAgent(input))
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Registration failed')
    } finally {
      setSubmitting(false)
    }
  }

  return <div className="modal-backdrop" role="presentation" onMouseDown={(event) => event.target === event.currentTarget && onClose()}><div className="modal" role="dialog" aria-modal="true" aria-labelledby="register-title"><div className="modal-header"><div><span className="eyebrow">NEW IDENTITY</span><h2 id="register-title">Register agent</h2></div><button className="icon-button" onClick={onClose} aria-label="Close"><Icon name="close"/></button></div><form onSubmit={submit}><div className="form-grid"><label><span>Agent name</span><input required autoFocus value={form.name} onChange={(e) => setForm({...form, name: e.target.value})} placeholder="e.g. Release assistant"/></label><label><span>Stable identifier</span><input required value={form.identifier} onChange={(e) => setForm({...form, identifier: e.target.value})} placeholder="e.g. codex:release-prod"/><small>Reusing this identifier refreshes the existing agent.</small></label><label><span>Framework</span><input value={form.framework} onChange={(e) => setForm({...form, framework: e.target.value})} placeholder="e.g. codex, langgraph"/></label><label><span>Metadata (JSON)</span><textarea rows={4} value={form.metadata} onChange={(e) => setForm({...form, metadata: e.target.value})} spellCheck={false}/></label></div>{error && <p className="form-error">{error}</p>}<div className="modal-actions"><button type="button" className="ghost-button" onClick={onClose}>Cancel</button><button className="primary-button" disabled={submitting}>{submitting ? 'Registering…' : 'Register agent'}</button></div></form></div></div>
}

function AgentDrawer({ agent, onClose }: { agent: Agent; onClose: () => void }) {
  return <div className="drawer-backdrop" onMouseDown={(event) => event.target === event.currentTarget && onClose()}><aside className="drawer" aria-label={`${agent.name} details`}><div className="drawer-header"><div className="agent-avatar large">{agent.name.slice(0, 2).toUpperCase()}</div><button className="icon-button" onClick={onClose} aria-label="Close"><Icon name="close"/></button></div><div className="drawer-title"><span className={`state ${isRecentlyActive(agent) ? 'online' : ''}`}><i/>{isRecentlyActive(agent) ? 'Active' : 'Idle'}</span><h2>{agent.name}</h2><code>{agent.identifier}</code></div><div className="detail-section"><h3>Identity</h3><dl><div><dt>Vanguard ID</dt><dd><code>{agent.id}</code></dd></div><div><dt>Framework</dt><dd>{agent.framework || 'Not provided'}</dd></div><div><dt>First seen</dt><dd>{new Date(agent.firstSeenAt).toLocaleString()}</dd></div><div><dt>Last seen</dt><dd>{new Date(agent.lastSeenAt).toLocaleString()}</dd></div></dl></div><div className="detail-section"><h3>Metadata</h3>{Object.keys(agent.metadata).length ? <pre>{JSON.stringify(agent.metadata, null, 2)}</pre> : <p className="muted">No metadata provided.</p>}</div></aside></div>
}

export default App
