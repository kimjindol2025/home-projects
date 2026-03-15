import React, { useState, useEffect } from 'react'
import axios from 'axios'
import '../styles/Pages.css'

interface Comment {
  id: number
  post_id: number
  author: string
  content: string
  approved: boolean
  created_at: number
}

export default function Comments() {
  const [comments, setComments] = useState<Comment[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState('all')

  useEffect(() => {
    fetchComments()
  }, [])

  const fetchComments = async () => {
    try {
      setLoading(true)
      const response = await axios.get('/api/comments')
      setComments(response.data || [])
    } catch (error) {
      console.error('Failed to fetch comments:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleApproveComment = async (id: number) => {
    try {
      await axios.put(`/api/comments/${id}`, { approved: true })
      fetchComments()
    } catch (error) {
      console.error('Failed to approve comment:', error)
    }
  }

  const handleRejectComment = async (id: number) => {
    try {
      await axios.delete(`/api/comments/${id}`)
      fetchComments()
    } catch (error) {
      console.error('Failed to reject comment:', error)
    }
  }

  const filteredComments = comments.filter(comment => {
    if (filter === 'pending') return !comment.approved
    if (filter === 'approved') return comment.approved
    return true
  })

  const stats = {
    total: comments.length,
    pending: comments.filter(c => !c.approved).length,
    approved: comments.filter(c => c.approved).length,
  }

  return (
    <div className="page comments-page">
      <div className="page-header">
        <h1>💬 댓글 관리</h1>
      </div>

      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-value">{stats.total}</div>
          <div className="stat-label">전체 댓글</div>
        </div>
        <div className="stat-card warning">
          <div className="stat-value">{stats.pending}</div>
          <div className="stat-label">승인 대기</div>
        </div>
        <div className="stat-card success">
          <div className="stat-value">{stats.approved}</div>
          <div className="stat-label">승인됨</div>
        </div>
      </div>

      <div className="filters">
        <select value={filter} onChange={e => setFilter(e.target.value)} className="filter-select">
          <option value="all">모든 댓글</option>
          <option value="pending">승인 대기</option>
          <option value="approved">승인됨</option>
        </select>
      </div>

      {loading ? (
        <div className="loading">로딩 중...</div>
      ) : (
        <div className="comments-list">
          {filteredComments.map(comment => (
            <div key={comment.id} className={`comment-card ${comment.approved ? 'approved' : 'pending'}`}>
              <div className="comment-header">
                <span className="author">{comment.author}</span>
                <span className="date">{new Date(comment.created_at * 1000).toLocaleDateString('ko-KR')}</span>
                <span className={`status ${comment.approved ? 'approved' : 'pending'}`}>
                  {comment.approved ? '✅ 승인됨' : '⏳ 대기 중'}
                </span>
              </div>
              <div className="comment-content">{comment.content}</div>
              <div className="comment-footer">
                <span className="post-ref">포스트 #{comment.post_id}</span>
                {!comment.approved && (
                  <div className="actions">
                    <button
                      className="btn btn-sm btn-success"
                      onClick={() => handleApproveComment(comment.id)}
                    >
                      승인
                    </button>
                    <button
                      className="btn btn-sm btn-delete"
                      onClick={() => handleRejectComment(comment.id)}
                    >
                      거부
                    </button>
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
