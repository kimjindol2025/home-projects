import React, { useState, useEffect } from 'react'
import axios from 'axios'
import '../styles/Pages.css'

interface User {
  id: number
  username: string
  email: string
  role: string
  is_active: boolean
  created_at: number
}

export default function Users() {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [filter, setFilter] = useState('all')

  useEffect(() => {
    fetchUsers()
  }, [])

  const fetchUsers = async () => {
    try {
      setLoading(true)
      const response = await axios.get('/api/users')
      setUsers(response.data || [])
    } catch (error) {
      console.error('Failed to fetch users:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleToggleActive = async (id: number, isActive: boolean) => {
    try {
      await axios.put(`/api/users/${id}`, { is_active: !isActive })
      fetchUsers()
    } catch (error) {
      console.error('Failed to update user:', error)
    }
  }

  const handleChangeRole = async (id: number, newRole: string) => {
    try {
      await axios.put(`/api/users/${id}`, { role: newRole })
      fetchUsers()
    } catch (error) {
      console.error('Failed to update user role:', error)
    }
  }

  const filteredUsers = users.filter(user => {
    const matchesSearch = user.username.includes(search) || user.email.includes(search)
    const matchesFilter = filter === 'all' ||
                         (filter === 'admin' && user.role === 'admin') ||
                         (filter === 'active' && user.is_active) ||
                         (filter === 'inactive' && !user.is_active)
    return matchesSearch && matchesFilter
  })

  const stats = {
    total: users.length,
    admin: users.filter(u => u.role === 'admin').length,
    active: users.filter(u => u.is_active).length,
  }

  return (
    <div className="page users-page">
      <div className="page-header">
        <h1>👥 사용자 관리</h1>
      </div>

      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-value">{stats.total}</div>
          <div className="stat-label">전체 사용자</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{stats.admin}</div>
          <div className="stat-label">관리자</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{stats.active}</div>
          <div className="stat-label">활성</div>
        </div>
      </div>

      <div className="filters">
        <input
          type="text"
          placeholder="사용자명 또는 이메일 검색..."
          value={search}
          onChange={e => setSearch(e.target.value)}
          className="search-input"
        />
        <select value={filter} onChange={e => setFilter(e.target.value)} className="filter-select">
          <option value="all">모든 사용자</option>
          <option value="admin">관리자</option>
          <option value="active">활성</option>
          <option value="inactive">비활성</option>
        </select>
      </div>

      {loading ? (
        <div className="loading">로딩 중...</div>
      ) : (
        <div className="table-responsive">
          <table className="data-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>사용자명</th>
                <th>이메일</th>
                <th>역할</th>
                <th>상태</th>
                <th>가입일</th>
                <th>작업</th>
              </tr>
            </thead>
            <tbody>
              {filteredUsers.map(user => (
                <tr key={user.id}>
                  <td>{user.id}</td>
                  <td>{user.username}</td>
                  <td>{user.email}</td>
                  <td>
                    <select
                      value={user.role}
                      onChange={e => handleChangeRole(user.id, e.target.value)}
                      className="role-select"
                    >
                      <option value="user">사용자</option>
                      <option value="author">작성자</option>
                      <option value="editor">편집자</option>
                      <option value="admin">관리자</option>
                    </select>
                  </td>
                  <td>
                    <span className={`badge ${user.is_active ? 'badge-success' : 'badge-danger'}`}>
                      {user.is_active ? '✅ 활성' : '❌ 비활성'}
                    </span>
                  </td>
                  <td>{new Date(user.created_at * 1000).toLocaleDateString('ko-KR')}</td>
                  <td>
                    <button
                      className={`btn btn-sm ${user.is_active ? 'btn-disable' : 'btn-enable'}`}
                      onClick={() => handleToggleActive(user.id, user.is_active)}
                    >
                      {user.is_active ? '비활성화' : '활성화'}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
