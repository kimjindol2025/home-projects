import { useState, useEffect } from 'react'
import axios from 'axios'

export default function System() {
  const [health, setHealth] = useState(null)
  const [backups, setBackups] = useState([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    fetchHealth()
    fetchBackups()
  }, [])

  const fetchHealth = async () => {
    try {
      const response = await axios.get('/api/admin/health')
      setHealth(response.data)
    } catch (error) {
      console.error('헬스 체크 실패:', error)
    }
  }

  const fetchBackups = async () => {
    try {
      const response = await axios.get('/api/admin/backups')
      setBackups(response.data || [])
    } catch (error) {
      console.error('백업 목록 조회 실패:', error)
    }
  }

  const triggerBackup = async () => {
    if (!confirm('백업을 지금 실행하시겠습니까?')) return
    setLoading(true)
    try {
      await axios.post('/api/admin/backup')
      alert('백업이 완료되었습니다')
      fetchBackups()
    } catch (error) {
      console.error('백업 실패:', error)
      alert('백업 실패')
    }
    setLoading(false)
  }

  return (
    <div className="space-y-6">
      {/* 시스템 상태 */}
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-semibold text-gray-800 mb-4">🔍 시스템 상태</h3>
        {health ? (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {[
              { name: 'API', status: health.api_status, icon: '⚙️' },
              { name: '데이터베이스', status: health.db_status, icon: '🗄️' },
              { name: '캐시', status: health.cache_status, icon: '⚡' },
              { name: '모니터링', status: health.monitoring_status, icon: '📊' }
            ].map((service) => (
              <div key={service.name} className="flex items-center gap-3 p-4 bg-gray-50 rounded-lg">
                <span className="text-2xl">{service.icon}</span>
                <div>
                  <p className="text-sm text-gray-600">{service.name}</p>
                  <p className={`font-semibold ${
                    service.status === 'UP' ? 'text-green-600' : 'text-red-600'
                  }`}>
                    {service.status === 'UP' ? '🟢 정상' : '🔴 오류'}
                  </p>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-gray-500">로딩 중...</p>
        )}
      </div>

      {/* 리소스 사용률 */}
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-semibold text-gray-800 mb-4">📈 리소스 사용률</h3>
        <div className="space-y-3">
          {[
            { name: 'CPU', usage: 45 },
            { name: '메모리', usage: 62 },
            { name: '디스크', usage: 38 }
          ].map((resource) => (
            <div key={resource.name}>
              <div className="flex justify-between text-sm mb-1">
                <span className="font-semibold text-gray-700">{resource.name}</span>
                <span className="text-gray-600">{resource.usage}%</span>
              </div>
              <div className="w-full h-2 bg-gray-200 rounded-full overflow-hidden">
                <div
                  className={`h-full transition-all ${
                    resource.usage > 80 ? 'bg-red-500' :
                    resource.usage > 60 ? 'bg-yellow-500' :
                    'bg-green-500'
                  }`}
                  style={{ width: `${resource.usage}%` }}
                ></div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* 백업 관리 */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-semibold text-gray-800">💾 백업 관리</h3>
          <button
            onClick={triggerBackup}
            disabled={loading}
            className="px-4 py-2 bg-freelang-500 text-white rounded-lg hover:bg-freelang-600 disabled:opacity-50"
          >
            {loading ? '백업 진행 중...' : '🔄 지금 백업'}
          </button>
        </div>

        {backups.length > 0 ? (
          <div className="space-y-2">
            {backups.map((backup, i) => (
              <div key={i} className="flex justify-between items-center p-3 bg-gray-50 rounded-lg">
                <div>
                  <p className="font-semibold text-gray-900">{backup.filename}</p>
                  <p className="text-sm text-gray-600">{backup.date}</p>
                </div>
                <div className="text-right">
                  <p className="text-sm text-gray-600">{backup.size}</p>
                  <button className="text-sm text-blue-600 hover:text-blue-800">⬇️ 다운로드</button>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-center text-gray-500 py-4">백업 기록이 없습니다</p>
        )}
      </div>

      {/* 로그 */}
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-semibold text-gray-800 mb-4">📋 최근 로그</h3>
        <div className="space-y-2 max-h-48 overflow-y-auto font-mono text-sm text-gray-600">
          <div>[2026-03-13 14:30:45] INFO: 서버 시작</div>
          <div>[2026-03-13 14:31:12] INFO: 데이터베이스 연결</div>
          <div>[2026-03-13 14:32:00] INFO: 백업 완료 (2.4MB)</div>
        </div>
      </div>
    </div>
  )
}
