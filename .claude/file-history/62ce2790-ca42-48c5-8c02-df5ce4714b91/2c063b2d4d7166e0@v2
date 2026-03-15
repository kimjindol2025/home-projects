import React, { useState, useEffect } from 'react'
import axios from 'axios'
import '../styles/Pages.css'

interface SystemStatus {
  status: string
  uptime: number
  memory: { used: number; total: number }
  disk: { used: number; total: number }
  db_connections: number
  cache_hit_rate: number
  error_rate: number
  avg_response_time: number
}

export default function System() {
  const [health, setHealth] = useState<SystemStatus | null>(null)
  const [logs, setLogs] = useState<string[]>([])
  const [loading, setLoading] = useState(true)
  const [backupProgress, setBackupProgress] = useState(0)
  const [isBackingUp, setIsBackingUp] = useState(false)

  useEffect(() => {
    fetchSystemStatus()
    fetchLogs()
    const interval = setInterval(() => {
      fetchSystemStatus()
    }, 5000)
    return () => clearInterval(interval)
  }, [])

  const fetchSystemStatus = async () => {
    try {
      const response = await axios.get('/api/admin/health')
      setHealth(response.data)
    } catch (error) {
      console.error('Failed to fetch system status:', error)
    }
  }

  const fetchLogs = async () => {
    try {
      const response = await axios.get('/api/admin/logs?level=ERROR&limit=10')
      setLogs(Array.isArray(response.data) ? response.data : [])
    } catch (error) {
      console.error('Failed to fetch logs:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleBackup = async () => {
    try {
      setIsBackingUp(true)
      setBackupProgress(0)

      // Simulate progress
      for (let i = 0; i <= 100; i += 10) {
        setBackupProgress(i)
        await new Promise(resolve => setTimeout(resolve, 200))
      }

      await axios.post('/api/admin/backup')
      alert('✅ 백업이 완료되었습니다')
      fetchLogs()
    } catch (error) {
      alert('❌ 백업 실패: ' + (error as any).message)
    } finally {
      setIsBackingUp(false)
      setBackupProgress(0)
    }
  }

  const getHealthColor = (status: string) => {
    if (status === 'healthy') return 'success'
    if (status === 'warning') return 'warning'
    return 'danger'
  }

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i)) + ' ' + sizes[i]
  }

  const getPercentage = (used: number, total: number) => {
    return total > 0 ? Math.round((used / total) * 100) : 0
  }

  return (
    <div className="page system-page">
      <div className="page-header">
        <h1>⚙️ 시스템 상태</h1>
        <button className="btn btn-secondary" onClick={fetchSystemStatus}>새로고침</button>
      </div>

      {health && (
        <div className="system-status">
          <div className={`status-card status-${getHealthColor(health.status)}`}>
            <h3>시스템 상태</h3>
            <div className="status-value">{health.status === 'healthy' ? '✅ 정상' : '⚠️ 주의'}</div>
            <div className="status-detail">
              <div>⏱️ 가동시간: {Math.floor(health.uptime / 3600)}시간</div>
              <div>📊 평균 응답시간: {health.avg_response_time.toFixed(2)}ms</div>
              <div>❌ 에러율: {(health.error_rate * 100).toFixed(2)}%</div>
            </div>
          </div>

          <div className="metrics-grid">
            <div className="metric-card">
              <h4>메모리</h4>
              <div className="progress-bar">
                <div
                  className="progress-fill"
                  style={{width: getPercentage(health.memory.used, health.memory.total) + '%'}}
                />
              </div>
              <div className="metric-label">
                {formatBytes(health.memory.used)} / {formatBytes(health.memory.total)}
              </div>
            </div>

            <div className="metric-card">
              <h4>디스크</h4>
              <div className="progress-bar">
                <div
                  className="progress-fill"
                  style={{width: getPercentage(health.disk.used, health.disk.total) + '%'}}
                />
              </div>
              <div className="metric-label">
                {formatBytes(health.disk.used)} / {formatBytes(health.disk.total)}
              </div>
            </div>

            <div className="metric-card">
              <h4>DB 연결</h4>
              <div className="metric-value">{health.db_connections}</div>
              <div className="metric-label">활성 연결</div>
            </div>

            <div className="metric-card">
              <h4>캐시 히트율</h4>
              <div className="metric-value">{(health.cache_hit_rate * 100).toFixed(1)}%</div>
              <div className="metric-label">성능 지표</div>
            </div>
          </div>
        </div>
      )}

      <div className="backup-section">
        <h3>💾 데이터 백업</h3>
        <div className="backup-controls">
          <button
            className="btn btn-primary"
            onClick={handleBackup}
            disabled={isBackingUp}
          >
            {isBackingUp ? '백업 중...' : '지금 백업하기'}
          </button>
        </div>
        {isBackingUp && (
          <div className="backup-progress">
            <div className="progress-bar">
              <div className="progress-fill" style={{width: backupProgress + '%'}} />
            </div>
            <div className="progress-label">{backupProgress}%</div>
          </div>
        )}
      </div>

      <div className="logs-section">
        <h3>📋 최근 에러 로그</h3>
        <div className="logs-container">
          {logs.length > 0 ? (
            logs.map((log, idx) => (
              <div key={idx} className="log-entry">
                <span className="log-icon">❌</span>
                <span className="log-text">{log}</span>
              </div>
            ))
          ) : (
            <div className="no-logs">❌ 에러 로그가 없습니다</div>
          )}
        </div>
      </div>
    </div>
  )
}
