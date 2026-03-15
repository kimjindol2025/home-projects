import React from 'react'
import { Link, useLocation } from 'react-router-dom'
import '../styles/Sidebar.css'

export default function Sidebar() {
  const location = useLocation()

  const menuItems = [
    { path: '/', label: '📊 대시보드', icon: '📊' },
    { path: '/posts', label: '📝 글 관리', icon: '📝' },
    { path: '/comments', label: '💬 댓글 관리', icon: '💬' },
    { path: '/users', label: '👥 사용자 관리', icon: '👥' },
    { path: '/system', label: '⚙️ 시스템 상태', icon: '⚙️' },
  ]

  return (
    <aside className="sidebar">
      <div className="sidebar-header">
        <h2>🚀 FreeLang</h2>
        <p>Admin Panel</p>
      </div>

      <nav className="sidebar-nav">
        {menuItems.map(item => (
          <Link
            key={item.path}
            to={item.path}
            className={`nav-item ${location.pathname === item.path ? 'active' : ''}`}
          >
            <span className="nav-icon">{item.icon}</span>
            <span className="nav-label">{item.label}</span>
          </Link>
        ))}
      </nav>

      <div className="sidebar-footer">
        <div className="version">v1.0.0 BETA</div>
        <a href="/" target="_blank" rel="noopener noreferrer" className="blog-link">
          블로그로 이동 →
        </a>
      </div>
    </aside>
  )
}
