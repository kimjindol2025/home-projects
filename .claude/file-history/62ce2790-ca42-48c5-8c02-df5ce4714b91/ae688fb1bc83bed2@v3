import { useState, useEffect } from 'react'
import axios from 'axios'

export default function Comments() {
  const [comments, setComments] = useState([])
  const [filter, setFilter] = useState('pending')

  useEffect(() => {
    fetchComments()
  }, [filter])

  const fetchComments = async () => {
    try {
      const response = await axios.get(`/api/comments?status=${filter}`)
      setComments(response.data || [])
    } catch (error) {
      console.error('댓글 조회 실패:', error)
    }
  }

  const approveComment = async (id) => {
    try {
      await axios.patch(`/api/comments/${id}`, { status: 'approved' })
      fetchComments()
    } catch (error) {
      console.error('승인 실패:', error)
    }
  }

  const rejectComment = async (id) => {
    try {
      await axios.delete(`/api/comments/${id}`)
      fetchComments()
    } catch (error) {
      console.error('거부 실패:', error)
    }
  }

  return (
    <div className="space-y-6">
      {/* 필터 */}
      <div className="flex gap-2">
        {['pending', 'approved', 'rejected'].map((status) => (
          <button
            key={status}
            onClick={() => setFilter(status)}
            className={`px-4 py-2 rounded-lg font-semibold ${
              filter === status
                ? 'bg-freelang-500 text-white'
                : 'bg-white border border-gray-300 text-gray-700 hover:bg-gray-50'
            }`}
          >
            {status === 'pending' && '⏳ 대기중'}
            {status === 'approved' && '✅ 승인됨'}
            {status === 'rejected' && '❌ 거부됨'}
          </button>
        ))}
      </div>

      {/* 댓글 목록 */}
      <div className="space-y-4">
        {comments.length === 0 ? (
          <div className="p-8 text-center text-gray-500 bg-white rounded-lg">댓글이 없습니다</div>
        ) : (
          comments.map((comment) => (
            <div key={comment.id} className="bg-white p-4 rounded-lg shadow border-l-4 border-freelang-500">
              <div className="flex justify-between items-start mb-2">
                <div>
                  <p className="font-semibold text-gray-900">{comment.author}</p>
                  <p className="text-sm text-gray-500">{new Date(comment.created_at).toLocaleString('ko-KR')}</p>
                </div>
                <span className="text-sm px-2 py-1 bg-gray-100 rounded">{comment.post_id} 글</span>
              </div>
              <p className="text-gray-700 mb-3">{comment.content}</p>
              {filter === 'pending' && (
                <div className="flex gap-2">
                  <button
                    onClick={() => approveComment(comment.id)}
                    className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
                  >
                    ✅ 승인
                  </button>
                  <button
                    onClick={() => rejectComment(comment.id)}
                    className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
                  >
                    ❌ 거부
                  </button>
                </div>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  )
}
