import React, { useState, useEffect } from 'react'
import axios from 'axios'
import '../styles/Pages.css'

interface Post {
  id: number
  title: string
  author: string
  created_at: number
  views: number
  status: string
}

export default function Posts() {
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [filter, setFilter] = useState('all')
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [formData, setFormData] = useState({ title: '', content: '', author: '' })

  useEffect(() => {
    fetchPosts()
  }, [])

  const fetchPosts = async () => {
    try {
      setLoading(true)
      const response = await axios.get('/api/blogs')
      setPosts(response.data || [])
    } catch (error) {
      console.error('Failed to fetch posts:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleCreatePost = async () => {
    if (!formData.title || !formData.content) {
      alert('제목과 내용을 입력해주세요')
      return
    }
    try {
      await axios.post('/api/blogs', formData)
      setFormData({ title: '', content: '', author: '' })
      setShowCreateForm(false)
      fetchPosts()
    } catch (error) {
      console.error('Failed to create post:', error)
    }
  }

  const handleDeletePost = async (id: number) => {
    if (!window.confirm('정말 삭제하시겠습니까?')) return
    try {
      await axios.delete(`/api/blogs/${id}`)
      fetchPosts()
    } catch (error) {
      console.error('Failed to delete post:', error)
    }
  }

  const filteredPosts = posts.filter(post => {
    const matchesSearch = post.title.includes(search)
    const matchesFilter = filter === 'all' || post.status === filter
    return matchesSearch && matchesFilter
  })

  return (
    <div className="page posts-page">
      <div className="page-header">
        <h1>📝 글 관리</h1>
        <button
          className="btn btn-primary"
          onClick={() => setShowCreateForm(!showCreateForm)}
        >
          + 새 글 작성
        </button>
      </div>

      {showCreateForm && (
        <div className="form-section">
          <h3>새 글 작성</h3>
          <input
            type="text"
            placeholder="제목"
            value={formData.title}
            onChange={e => setFormData({...formData, title: e.target.value})}
            className="form-input"
          />
          <input
            type="text"
            placeholder="작성자"
            value={formData.author}
            onChange={e => setFormData({...formData, author: e.target.value})}
            className="form-input"
          />
          <textarea
            placeholder="내용"
            value={formData.content}
            onChange={e => setFormData({...formData, content: e.target.value})}
            className="form-textarea"
          />
          <div className="form-actions">
            <button className="btn btn-success" onClick={handleCreatePost}>저장</button>
            <button className="btn btn-cancel" onClick={() => setShowCreateForm(false)}>취소</button>
          </div>
        </div>
      )}

      <div className="filters">
        <input
          type="text"
          placeholder="검색..."
          value={search}
          onChange={e => setSearch(e.target.value)}
          className="search-input"
        />
        <select value={filter} onChange={e => setFilter(e.target.value)} className="filter-select">
          <option value="all">모든 상태</option>
          <option value="published">발행됨</option>
          <option value="draft">임시저장</option>
          <option value="archived">보관됨</option>
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
                <th>제목</th>
                <th>작성자</th>
                <th>조회수</th>
                <th>상태</th>
                <th>작성일</th>
                <th>작업</th>
              </tr>
            </thead>
            <tbody>
              {filteredPosts.map(post => (
                <tr key={post.id}>
                  <td>{post.id}</td>
                  <td className="title-cell">{post.title}</td>
                  <td>{post.author}</td>
                  <td><span className="badge badge-info">👁️ {post.views}</span></td>
                  <td>
                    <span className={`badge badge-${post.status}`}>
                      {post.status === 'published' ? '발행됨' : '임시저장'}
                    </span>
                  </td>
                  <td>{new Date(post.created_at * 1000).toLocaleDateString('ko-KR')}</td>
                  <td>
                    <button className="btn btn-sm btn-edit">수정</button>
                    <button
                      className="btn btn-sm btn-delete"
                      onClick={() => handleDeletePost(post.id)}
                    >
                      삭제
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
