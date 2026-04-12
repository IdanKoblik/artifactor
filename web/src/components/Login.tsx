import { useState, useEffect } from 'react'

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
  <svg className="provider-icon" viewBox="0 0 24 24" aria-hidden="true">
    <path d="M22.65 14.39L12 22.13 1.35 14.39a.84.84 0 0 1-.3-.94l1.22-3.78 2.44-7.51A.42.42 0 0 1 4.82 2a.43.43 0 0 1 .58 0 .42.42 0 0 1 .11.18l2.44 7.49h8.1l2.44-7.51A.42.42 0 0 1 18.6 2a.43.43 0 0 1 .58 0 .42.42 0 0 1 .11.18l2.44 7.51 1.22 3.78a.84.84 0 0 1-.3.92z"/>
  </svg>
)

const GithubIcon = () => (
  <svg className="provider-icon" viewBox="0 0 24 24" aria-hidden="true">
    <path d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"/>
  </svg>
)

const authRedirectPath: Record<string, string> = {
  gitlab: '/api/auth/gitlab/redirect',
  github: '/api/auth/github/redirect',
}

function hostLabel(host: Host): string {
  try {
    return new URL(host.url).hostname
  } catch {
    return host.url
  }
}

export default function Login({ onLogin }: Props) {
  const [error, setError] = useState('')
  const [hosts, setHosts] = useState<Host[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/hosts')
      .then(res => res.json())
      .then((data: Host[]) => setHosts(data))
      .catch(() => setError('Failed to load hosts'))
      .finally(() => setLoading(false))
  }, [])

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
          <div className="empty">No login providers configured</div>
        ) : (
          <div className="login-providers">
            {hosts.map(host => (
              <button
                key={host.id}
                type="button"
                className={`btn btn-provider btn-${host.type} btn-full`}
                onClick={() => {
                  const base = authRedirectPath[host.type]
                  if (base) window.location.href = `${base}?id=${host.id}`
                }}
              >
                {host.type === 'gitlab' ? <GitlabIcon /> : <GithubIcon />}
                Sign in with {host.type === 'gitlab' ? 'GitLab' : 'GitHub'} &middot; {hostLabel(host)}
              </button>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
