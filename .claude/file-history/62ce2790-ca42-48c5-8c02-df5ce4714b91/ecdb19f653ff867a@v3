import { BrowserRouter as Router, Routes, Route, Link, useLocation } from 'react-router-dom'
import { useState } from 'react'
import Posts from './pages/Posts'
import Comments from './pages/Comments'
import Users from './pages/Users'
import Analytics from './pages/Analytics'
import System from './pages/System'

export default function App() {
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const location = useLocation()

  const navItems = [
    { label: '📝 글 관리', path: '/posts', icon: '📝' },
    { label: '💬 댓글 관리', path: '/comments', icon: '💬' },
    { label: '👥 사용자 관리', path: '/users', icon: '👥' },
    { label: '📊 분석', path: '/analytics', icon: '📊' },
    { label: '⚙️ 시스템', path: '/system', icon: '⚙️' }
  ]

  return (
    <Router>
      <div className="flex h-screen bg-gray-100">
        {/* 사이드바 */}
        <aside className={`${sidebarOpen ? 'w-64' : 'w-20'} bg-white shadow-lg transition-all duration-300 flex flex-col`}>
          <div className="p-4 border-b">
            <h1 className={`font-bold text-freelang-600 ${sidebarOpen ? 'text-xl' : 'text-xs text-center'}`}>
              {sidebarOpen ? '🚀 FreeLang Admin' : 'FL'}
            </h1>
          </div>

          <nav className="flex-1 p-4 space-y-2">
            {navItems.map((item) => (
              <Link
                key={item.path}
                to={item.path}
                className={`block px-4 py-2 rounded-lg transition-colors ${
                  location.pathname === item.path
                    ? 'bg-freelang-100 text-freelang-700 font-semibold'
                    : 'text-gray-700 hover:bg-gray-100'
                } ${!sidebarOpen && 'text-center text-lg'}`}
              >
                {sidebarOpen ? item.label : item.icon}
              </Link>
            ))}
          </nav>

          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className="p-4 border-t text-gray-600 hover:bg-gray-100 text-center"
          >
            {sidebarOpen ? '◀️' : '▶️'}
          </button>
        </aside>

        {/* 메인 콘텐츠 */}
        <main className="flex-1 flex flex-col overflow-hidden">
          {/* 헤더 */}
          <header className="bg-white shadow-sm px-6 py-4 flex justify-between items-center">
            <h2 className="text-2xl font-bold text-gray-800">
              {navItems.find(n => n.path === location.pathname)?.label || '대시보드'}
            </h2>
            <div className="flex items-center gap-4">
              <button className="px-4 py-2 bg-freelang-500 text-white rounded-lg hover:bg-freelang-600">
                🔔 알림
              </button>
              <img src="https://via.placeholder.com/40" alt="프로필" className="w-10 h-10 rounded-full" />
            </div>
          </header>

          {/* 페이지 콘텐츠 */}
          <div className="flex-1 overflow-auto p-6">
            <Routes>
              <Route path="/posts" element={<Posts />} />
              <Route path="/comments" element={<Comments />} />
              <Route path="/users" element={<Users />} />
              <Route path="/analytics" element={<Analytics />} />
              <Route path="/system" element={<System />} />
              <Route path="/" element={<Analytics />} />
            </Routes>
          </div>
        </main>
      </div>
    </Router>
  )
}
