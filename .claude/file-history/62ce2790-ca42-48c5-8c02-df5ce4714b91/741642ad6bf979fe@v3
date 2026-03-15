import { useState, useEffect } from 'react'
import axios from 'axios'
import { Line, Bar } from 'react-chartjs-2'
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, BarElement, Title, Tooltip, Legend } from 'chart.js'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, BarElement, Title, Tooltip, Legend)

export default function Analytics() {
  const [stats, setStats] = useState({
    total_posts: 0,
    total_users: 0,
    total_comments: 0,
    total_views: 0
  })
  const [chartData, setChartData] = useState(null)

  useEffect(() => {
    fetchStats()
  }, [])

  const fetchStats = async () => {
    try {
      const response = await axios.get('/api/admin/stats')
      setStats(response.data)

      // 모의 차트 데이터
      setChartData({
        labels: ['주 1', '주 2', '주 3', '주 4', '주 5'],
        datasets: [
          {
            label: '조회수',
            data: [120, 190, 150, 200, 180],
            borderColor: '#0ea5e9',
            backgroundColor: 'rgba(14, 165, 233, 0.1)',
            tension: 0.4
          }
        ]
      })
    } catch (error) {
      console.error('통계 조회 실패:', error)
    }
  }

  return (
    <div className="space-y-6">
      {/* 주요 지표 */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {[
          { label: '📝 총 글', value: stats.total_posts, color: 'blue' },
          { label: '👥 총 사용자', value: stats.total_users, color: 'green' },
          { label: '💬 총 댓글', value: stats.total_comments, color: 'purple' },
          { label: '👁️ 총 조회', value: stats.total_views, color: 'orange' }
        ].map((metric, i) => (
          <div key={i} className={`bg-${metric.color}-50 border border-${metric.color}-200 rounded-lg p-4`}>
            <p className={`text-${metric.color}-600 text-sm font-semibold`}>{metric.label}</p>
            <p className={`text-${metric.color}-900 text-3xl font-bold`}>{metric.value}</p>
          </div>
        ))}
      </div>

      {/* 차트 */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {chartData && (
          <>
            <div className="bg-white p-6 rounded-lg shadow">
              <h3 className="text-lg font-semibold text-gray-800 mb-4">조회수 추이</h3>
              <Line data={chartData} options={{ responsive: true, maintainAspectRatio: true }} />
            </div>
            <div className="bg-white p-6 rounded-lg shadow">
              <h3 className="text-lg font-semibold text-gray-800 mb-4">상단 포스트</h3>
              <div className="space-y-2">
                {['Post 1 (450 조회)', 'Post 2 (380 조회)', 'Post 3 (320 조회)'].map((post, i) => (
                  <div key={i} className="flex justify-between items-center p-2 bg-gray-50 rounded">
                    <span>{post}</span>
                    <div className="w-24 h-2 bg-gray-200 rounded-full overflow-hidden">
                      <div className="h-full bg-freelang-500" style={{ width: `${100 - i * 15}%` }}></div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  )
}
