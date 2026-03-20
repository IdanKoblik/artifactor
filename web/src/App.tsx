import { useState } from 'react'
import Login from './components/Login'
import Dashboard from './components/Dashboard'

const TOKEN_KEY = 'artifactor_token'

export default function App() {
  const [token, setToken] = useState<string | null>(
    () => sessionStorage.getItem(TOKEN_KEY),
  )

  const handleLogin = (t: string) => {
    sessionStorage.setItem(TOKEN_KEY, t)
    setToken(t)
  }

  const handleLogout = () => {
    sessionStorage.removeItem(TOKEN_KEY)
    setToken(null)
  }

  if (!token) return <Login onLogin={handleLogin} />
  return <Dashboard token={token} onLogout={handleLogout} />
}
