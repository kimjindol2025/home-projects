import React, { useState, useEffect } from 'react'
import axios from 'axios'
import { Line, Bar, Doughnut } from 'react-chartjs-2'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js'
import '../styles/Pages.css'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend
)

export default function Analytics() {
  const [stats, setStats] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchStats()
  }, [])

  const fetchStats = async () => {
    try {
      setLoading(true)
      const response = await axios.get('/api/admin/stats')
      setStats(response.data || getDefaultStats())
    } catch (error) {
      console.error('Failed to fetch stats:', error)
      setStats(getDefaultStats())
    } finally {
      setLoading(false)
    }
  }

  const getDefaultStats = () => ({
    posts: 15,
    comments: 42,
    users: 8,
    views: 1250,
    likes: 234,
    viewsTrend: [100, 150, 130, 200, 180, 220, 250],
    commentsTrend: [5, 8, 6, 12, 10, 15, 20],
    postsPerAuthor: { admin: 8, john: 5, jane: 2 },
  })

  if (loading) {
    return <div className="loading">로딩 중...</div>
  }

  const viewsChartData = {
    labels: ['월', '화', '수', '목', '금', '토', '일'],
    datasets: [
      {
        label: '조회수',
        data: stats?.viewsTrend || [0],
        borderColor: '#3b82f6',
        backgroundColor: 'rgba(59, 130, 246, 0.1)',
        tension: 0.3,
      },
    ],
  }

  const commentsChartData = {
    labels: ['월', '화', '수', '목', '금', '토', '일'],
    datasets: [
      {
        label: '댓글',
        data: stats?.commentsTrend || [0],
        backgroundColor: 'rgba(34, 197, 94, 0.6)',
        borderColor: '#22c55e',
      },
    ],
  }

  const authorChartData = {
    labels: Object.keys(stats?.postsPerAuthor || {}),
    datasets: [
      {
        data: Object.values(stats?.postsPerAuthor || []),
        backgroundColor: ['#3b82f6', '#ef4444', '#f59e0b'],
      },
    ],
  }

  return (
    <div className="page analytics-page">
      <div className="page-header">
        <h1>📊 통계 및 분석</h1>
        <button className="btn btn-secondary" onClick={fetchStats}>새로고침</button>
      </div>

      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-value">{stats.posts}</div>
          <div className="stat-label">📝 글</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{stats.comments}</div>
          <div className="stat-label">💬 댓글</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{stats.users}</div>
          <div className="stat-label">👥 사용자</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{stats.views.toLocaleString()}</div>
          <div className="stat-label">👁️ 조회</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{stats.likes}</div>
          <div className="stat-label">❤️ 좋아요</div>
        </div>
      </div>

      <div className="charts-grid">
        <div className="chart-container">
          <h3>조회수 추이 (최근 7일)</h3>
          <Line data={viewsChartData} options={{
            responsive: true,
            plugins: { legend: { position: 'top' } }
          }} />
        </div>

        <div className="chart-container">
          <h3>댓글 추이 (최근 7일)</h3>
          <Bar data={commentsChartData} options={{
            responsive: true,
            plugins: { legend: { position: 'top' } }
          }} />
        </div>

        <div className="chart-container">
          <h3>작성자별 글 수</h3>
          <Doughnut data={authorChartData} options={{
            responsive: true,
            plugins: { legend: { position: 'bottom' } }
          }} />
        </div>
      </div>
    </div>
  )
}
