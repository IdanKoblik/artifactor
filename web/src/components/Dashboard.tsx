import { useState } from 'react'
import Tokens from './Tokens'
import Products from './Products'
import Health from './Health'

type Tab = 'tokens' | 'products' | 'health'

interface Props {
  token: string
  onLogout: () => void
}

const TABS: { id: Tab; label: string }[] = [
  { id: 'tokens',   label: 'Tokens'   },
  { id: 'products', label: 'Products' },
  { id: 'health',   label: 'Health'   },
]

export default function Dashboard({ token, onLogout }: Props) {
  const [tab, setTab] = useState<Tab>('tokens')

  return (
    <div className="dashboard">
      <header className="header">
        <span className="header-title">ARTIFACTOR</span>
        <div className="header-right">
          <span className="badge-admin">admin</span>
          <button className="btn btn-secondary btn-sm" onClick={onLogout}>
            Logout
          </button>
        </div>
      </header>

      <nav className="nav">
        {TABS.map(t => (
          <div
            key={t.id}
            className={`nav-tab${tab === t.id ? ' active' : ''}`}
            onClick={() => setTab(t.id)}
          >
            {t.label}
          </div>
        ))}
      </nav>

      <main className="content">
        {tab === 'tokens'   && <Tokens   token={token} />}
        {tab === 'products' && <Products token={token} />}
        {tab === 'health'   && <Health   token={token} />}
      </main>
    </div>
  )
}
