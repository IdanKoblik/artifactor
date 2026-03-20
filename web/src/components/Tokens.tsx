import { useState, useEffect, useCallback } from 'react'
import { listTokens, createToken, revokeToken, ApiToken } from '../api'

interface Props {
  token: string
}

export default function Tokens({ token }: Props) {
  const [tokens, setTokens]         = useState<ApiToken[] | null>(null)
  const [loadError, setLoadError]   = useState('')
  const [loading, setLoading]       = useState(true)

  const [showModal, setShowModal]   = useState(false)
  const [newAdmin, setNewAdmin]     = useState(false)
  const [newToken, setNewToken]     = useState<string | null>(null)
  const [createError, setCreateError] = useState('')
  const [creating, setCreating]     = useState(false)

  const [revokeInput, setRevokeInput]     = useState('')
  const [revokeError, setRevokeError]     = useState('')
  const [revokeSuccess, setRevokeSuccess] = useState('')

  const load = useCallback(async () => {
    setLoading(true)
    setLoadError('')
    try {
      setTokens(await listTokens(token))
    } catch (e: unknown) {
      setLoadError(e instanceof Error ? e.message : 'Failed to load tokens')
    } finally {
      setLoading(false)
    }
  }, [token])

  useEffect(() => { load() }, [load])

  const openModal = () => {
    setNewAdmin(false)
    setNewToken(null)
    setCreateError('')
    setCreating(false)
    setShowModal(true)
  }

  const handleCreate = async () => {
    setCreateError('')
    setCreating(true)
    try {
      const raw = await createToken(token, newAdmin)
      setNewToken(raw)
      load()
    } catch (e: unknown) {
      setCreateError(e instanceof Error ? e.message : 'Failed to create token')
    } finally {
      setCreating(false)
    }
  }

  const handleRevoke = async () => {
    const raw = revokeInput.trim()
    if (!raw) { setRevokeError('Token value is required'); return }

    setRevokeError('')
    setRevokeSuccess('')
    try {
      await revokeToken(token, raw)
      setRevokeInput('')
      setRevokeSuccess('Token revoked successfully')
      load()
    } catch (e: unknown) {
      setRevokeError(e instanceof Error ? e.message : 'Failed to revoke token')
    }
  }

  return (
    <>
      <div className="section-header">
        <span className="section-title">API Tokens</span>
        <button className="btn btn-primary btn-sm" onClick={openModal}>
          + Create Token
        </button>
      </div>

      {loadError && <div className="alert alert-error">{loadError}</div>}

      {loading ? (
        <div className="loading">Loading…</div>
      ) : tokens && tokens.length > 0 ? (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Token Hash</th>
                <th>Type</th>
              </tr>
            </thead>
            <tbody>
              {tokens.map(t => (
                <tr key={t.token}>
                  <td className="code" title={t.token}>
                    {t.token.substring(0, 20)}…
                  </td>
                  <td>
                    <span className={`token-type ${t.admin ? 'type-admin' : 'type-regular'}`}>
                      {t.admin ? 'admin' : 'regular'}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      ) : (
        <div className="empty">No tokens registered yet.</div>
      )}

      <div className="detail-panel">
        <div className="detail-panel-title">Revoke Token</div>
        <p style={{ fontSize: 12, color: 'var(--muted)', marginBottom: 12 }}>
          Paste the raw token value to permanently revoke access.
        </p>
        <div className="flex-gap">
          <input
            type="password"
            value={revokeInput}
            onChange={e => setRevokeInput(e.target.value)}
            placeholder="Raw token value"
            autoComplete="off"
            style={{ flex: 1 }}
          />
          <button className="btn btn-danger btn-sm" onClick={handleRevoke}>
            Revoke
          </button>
        </div>
        {revokeError   && <div className="alert alert-error mt-4">{revokeError}</div>}
        {revokeSuccess && <div className="alert alert-success mt-4">{revokeSuccess}</div>}
      </div>

      {showModal && (
        <div
          className="modal-overlay"
          onClick={e => { if (e.target === e.currentTarget) setShowModal(false) }}
        >
          <div className="modal">
            <div className="modal-title">Create New Token</div>

            <div className="form-group">
              <label className="checkbox-row">
                <input
                  type="checkbox"
                  checked={newAdmin}
                  onChange={e => setNewAdmin(e.target.checked)}
                />
                Grant admin privileges
              </label>
            </div>

            {newToken && (
              <>
                <div className="token-reveal">{newToken}</div>
                <div className="warn-copy">
                  ⚠ Copy this token now — it will not be shown again.
                </div>
              </>
            )}

            {createError && <div className="alert alert-error">{createError}</div>}

            <div className="modal-footer">
              <button className="btn btn-secondary" onClick={() => setShowModal(false)}>
                Close
              </button>
              {!newToken && (
                <button
                  className="btn btn-primary"
                  onClick={handleCreate}
                  disabled={creating}
                >
                  {creating ? 'Creating…' : 'Create'}
                </button>
              )}
            </div>
          </div>
        </div>
      )}
    </>
  )
}
