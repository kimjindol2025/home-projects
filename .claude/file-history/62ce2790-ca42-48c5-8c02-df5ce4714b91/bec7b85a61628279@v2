import React, { useState } from 'react'
import '../styles/TopBar.css'

interface TopBarProps {
  user: { name: string; role: string }
}

export default function TopBar({ user }: TopBarProps) {
  const [showMenu, setShowMenu] = useState(false)

  const handleLogout = () => {
    if (window.confirm('로그아웃하시겠습니까?')) {
      window.location.href = '/'
    }
  }

  return (
    <header className="topbar">
      <div className="topbar-content">
        <div className="topbar-title">
          <h1>Admin Dashboard</h1>
        </div>

        <div className="topbar-actions">
          <div className="notification-bell">
            🔔
            <span className="badge">3</span>
          </div>

          <div className="user-menu">
            <button
              className="user-button"
              onClick={() => setShowMenu(!showMenu)}
            >
              <span className="user-avatar">{user.name.charAt(0)}</span>
              <span className="user-info">
                <span className="user-name">{user.name}</span>
                <span className="user-role">{user.role}</span>
              </span>
            </button>

            {showMenu && (
              <div className="dropdown-menu">
                <a href="/settings" className="menu-item">⚙️ 설정</a>
                <a href="/profile" className="menu-item">👤 프로필</a>
                <hr />
                <button className="menu-item logout" onClick={handleLogout}>
                  🚪 로그아웃
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  )
}
