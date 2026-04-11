import { useState, useEffect, useRef } from 'react'

interface Props {
  onLogin: (token: string) => void
}

interface Host {
  id: number
  url: string
  type: string
}

const features = [
  { icon: '📦', label: 'Package Management', desc: 'Organize and track all your product versions in one place' },
  { icon: '🔗', label: 'Git Integration', desc: 'Seamlessly connect releases to your Git repositories' },
  { icon: '🔒', label: 'Secure Access', desc: 'Role-based permissions with admin and regular user tiers' },
]

const GitlabIcon = () => (
  <svg className="gitlab-icon" viewBox="0 0 24 24" aria-hidden="true">
    <path d="M22.65 14.39L12 22.13 1.35 14.39a.84.84 0 0 1-.3-.94l1.22-3.78 2.44-7.51A.42.42 0 0 1 4.82 2a.43.43 0 0 1 .58 0 .42.42 0 0 1 .11.18l2.44 7.49h8.1l2.44-7.51A.42.42 0 0 1 18.6 2a.43.43 0 0 1 .58 0 .42.42 0 0 1 .11.18l2.44 7.51 1.22 3.78a.84.84 0 0 1-.3.92z"/>
  </svg>
)

export default function Login({ onLogin }: Props) {
  const [error, setError] = useState('')
  const [hosts, setHosts] = useState<Host[]>([])
  const [loading, setLoading] = useState(true)
  const [showPicker, setShowPicker] = useState(false)
  const pickerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    fetch('/api/hosts')
      .then(res => res.json())
      .then((data: Host[]) => {
        setHosts(data.filter(h => h.type === 'gitlab'))
      })
      .catch(() => setError('Failed to load hosts'))
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    if (!showPicker) return
    const handleClick = (e: MouseEvent) => {
      if (pickerRef.current && !pickerRef.current.contains(e.target as Node)) {
        setShowPicker(false)
      }
    }
    document.addEventListener('mousedown', handleClick)
    return () => document.removeEventListener('mousedown', handleClick)
  }, [showPicker])

  const handleGitlabClick = () => {
    if (hosts.length === 1) {
      window.location.href = `/api/auth/gitlab/redirect?id=${hosts[0].id}`
      return
    }
    setShowPicker(true)
  }

  return (
    <div className="login-wrapper">
      <div className="login-card">
        <div className="login-header">
          <div className="login-logo-text">PACKSTER</div>
          <p className="login-tagline">Product Version Manager</p>
        </div>

        <div className="login-features">
          {features.map((f) => (
            <div key={f.label} className="login-feature-item">
              <span className="login-feature-icon">{f.icon}</span>
              <div>
                <div className="login-feature-label">{f.label}</div>
                <div className="login-feature-desc">{f.desc}</div>
              </div>
            </div>
          ))}
        </div>

        <div className="login-divider">Sign in to continue</div>

        {error && <div className="alert alert-error">{error}</div>}

        {loading ? (
          <div className="loading">Loading providers...</div>
        ) : hosts.length === 0 ? (
          <div className="empty">No GitLab instances configured</div>
        ) : (
          <div className="host-picker-wrap" ref={pickerRef}>
            <button
              type="button"
              className="btn btn-gitlab btn-full"
              onClick={handleGitlabClick}
            >
              <GitlabIcon />
              Sign in with GitLab
            </button>

            {showPicker && (
              <div className="host-picker">
                <div className="host-picker-title">Select instance</div>
                {hosts.map(host => (
                  <button
                    key={host.id}
                    type="button"
                    className="host-picker-item"
                    onClick={() => {
                      window.location.href = `/api/auth/gitlab/redirect?id=${host.id}`
                    }}
                  >
                    {host.url}
                  </button>
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
