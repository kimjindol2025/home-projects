import { useState, useEffect } from 'react'
import axios from 'axios'

export default function Posts() {
  const [posts, setPosts] = useState([])
  const [loading, setLoading] = useState(false)
  const [search, setSearch] = useState('')

  useEffect(() => {
    fetchPosts()
  }, [])

  const fetchPosts = async () => {
    setLoading(true)
    try {
      const response = await axios.get('/api/posts')
      setPosts(response.data || [])
    } catch (error) {
      console.error('글 조회 실패:', error)
    }
    setLoading(false)
  }

  const deletPost = async (id) => {
    if (confirm('정말 삭제하시겠습니까?')) {
      try {
        await axios.delete(`/api/posts/${id}`)
        setPosts(posts.filter(p => p.id !== id))
      } catch (error) {
        console.error('삭제 실패:', error)
      }
    }
  }

  const filteredPosts = posts.filter(p =>
    p.title?.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div className="space-y-6">
      {/* 검색 + 버튼 */}
      <div className="flex gap-4">
        <input
          type="text"
          placeholder="글 제목으로 검색..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-freelang-500"
        />
        <button className="px-6 py-2 bg-freelang-500 text-white rounded-lg hover:bg-freelang-600">
          ➕ 새 글
        </button>
      </div>

      {/* 글 목록 */}
      <div className="bg-white rounded-lg shadow">
        {loading ? (
          <div className="p-8 text-center text-gray-500">로딩 중...</div>
        ) : filteredPosts.length === 0 ? (
          <div className="p-8 text-center text-gray-500">글이 없습니다</div>
        ) : (
          <table className="w-full">
            <thead className="border-b bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left font-semibold text-gray-700">제목</th>
                <th className="px-6 py-3 text-left font-semibold text-gray-700">작성자</th>
                <th className="px-6 py-3 text-left font-semibold text-gray-700">조회</th>
                <th className="px-6 py-3 text-left font-semibold text-gray-700">상태</th>
                <th className="px-6 py-3 text-center font-semibold text-gray-700">작업</th>
              </tr>
            </thead>
            <tbody>
              {filteredPosts.map((post) => (
                <tr key={post.id} className="border-b hover:bg-gray-50">
                  <td className="px-6 py-4 font-medium text-gray-900">{post.title}</td>
                  <td className="px-6 py-4 text-gray-600">{post.author}</td>
                  <td className="px-6 py-4 text-gray-600">{post.views || 0}</td>
                  <td className="px-6 py-4">
                    <span className={`px-3 py-1 rounded-full text-sm font-semibold ${
                      post.is_published
                        ? 'bg-green-100 text-green-800'
                        : 'bg-yellow-100 text-yellow-800'
                    }`}>
                      {post.is_published ? '발행됨' : '초안'}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-center space-x-2">
                    <button className="px-3 py-1 text-sm bg-blue-100 text-blue-600 rounded hover:bg-blue-200">
                      ✏️ 수정
                    </button>
                    <button
                      onClick={() => deletPost(post.id)}
                      className="px-3 py-1 text-sm bg-red-100 text-red-600 rounded hover:bg-red-200"
                    >
                      🗑️ 삭제
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
