import { useState, useEffect } from 'react'
import axios from 'axios'

export default function Users() {
  const [users, setUsers] = useState([])
  const [search, setSearch] = useState('')

  useEffect(() => {
    fetchUsers()
  }, [])

  const fetchUsers = async () => {
    try {
      const response = await axios.get('/api/users')
      setUsers(response.data || [])
    } catch (error) {
      console.error('사용자 조회 실패:', error)
    }
  }

  const toggleUserStatus = async (id, isActive) => {
    try {
      await axios.patch(`/api/users/${id}`, { is_active: !isActive })
      fetchUsers()
    } catch (error) {
      console.error('상태 변경 실패:', error)
    }
  }

  const filteredUsers = users.filter(u =>
    u.username?.toLowerCase().includes(search.toLowerCase()) ||
    u.email?.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div className="space-y-6">
      {/* 검색 */}
      <input
        type="text"
        placeholder="사용자명 또는 이메일로 검색..."
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-freelang-500"
      />

      {/* 사용자 테이블 */}
      <div className="bg-white rounded-lg shadow overflow-x-auto">
        <table className="w-full">
          <thead className="border-b bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left font-semibold text-gray-700">사용자명</th>
              <th className="px-6 py-3 text-left font-semibold text-gray-700">이메일</th>
              <th className="px-6 py-3 text-left font-semibold text-gray-700">역할</th>
              <th className="px-6 py-3 text-left font-semibold text-gray-700">상태</th>
              <th className="px-6 py-3 text-left font-semibold text-gray-700">마지막 접속</th>
              <th className="px-6 py-3 text-center font-semibold text-gray-700">작업</th>
            </tr>
          </thead>
          <tbody>
            {filteredUsers.map((user) => (
              <tr key={user.id} className="border-b hover:bg-gray-50">
                <td className="px-6 py-4 font-medium text-gray-900">{user.username}</td>
                <td className="px-6 py-4 text-gray-600">{user.email}</td>
                <td className="px-6 py-4">
                  <span className={`px-3 py-1 rounded-full text-sm font-semibold ${
                    user.role === 'admin' ? 'bg-purple-100 text-purple-800' : 'bg-blue-100 text-blue-800'
                  }`}>
                    {user.role === 'admin' ? '👑 관리자' : '👤 사용자'}
                  </span>
                </td>
                <td className="px-6 py-4">
                  <span className={`px-3 py-1 rounded-full text-sm font-semibold ${
                    user.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                  }`}>
                    {user.is_active ? '🟢 활성' : '🔴 비활성'}
                  </span>
                </td>
                <td className="px-6 py-4 text-gray-600 text-sm">
                  {user.last_login ? new Date(user.last_login).toLocaleDateString('ko-KR') : '접속 기록 없음'}
                </td>
                <td className="px-6 py-4 text-center">
                  <button
                    onClick={() => toggleUserStatus(user.id, user.is_active)}
                    className={`px-3 py-1 text-sm rounded ${
                      user.is_active
                        ? 'bg-red-100 text-red-600 hover:bg-red-200'
                        : 'bg-green-100 text-green-600 hover:bg-green-200'
                    }`}
                  >
                    {user.is_active ? '🔒 비활성화' : '🔓 활성화'}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="text-sm text-gray-600">
        총 {users.length}명 (활성: {users.filter(u => u.is_active).length}명)
      </div>
    </div>
  )
}
