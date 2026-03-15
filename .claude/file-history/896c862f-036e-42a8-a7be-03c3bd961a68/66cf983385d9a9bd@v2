/**
 * FreeLang Admin Dashboard
 * React Component for Management
 */

import React, { useState, useEffect } from 'react';
import './admin-dashboard.css';

const AdminDashboard = () => {
  const [user, setUser] = useState(null);
  const [stats, setStats] = useState({
    totalUsers: 0,
    totalPosts: 0,
    totalViews: 0,
    activeUsers: 0,
  });
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('overview');

  useEffect(() => {
    fetchUserData();
    fetchStats();
    fetchPosts();
  }, []);

  const fetchUserData = async () => {
    try {
      const token = localStorage.getItem('authToken');
      const response = await fetch('/api/auth/verify', {
        headers: { Authorization: `Bearer ${token}` },
      });

      if (response.ok) {
        const data = await response.json();
        setUser(data.user);
      }
    } catch (error) {
      console.error('Failed to fetch user:', error);
    }
  };

  const fetchStats = async () => {
    try {
      const response = await fetch('/api/admin/stats');
      const data = await response.json();
      setStats(data);
    } catch (error) {
      console.error('Failed to fetch stats:', error);
    }
  };

  const fetchPosts = async () => {
    try {
      const response = await fetch('/api/admin/posts');
      const data = await response.json();
      setPosts(data);
      setLoading(false);
    } catch (error) {
      console.error('Failed to fetch posts:', error);
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    await fetch('/api/auth/logout', { method: 'POST' });
    localStorage.removeItem('authToken');
    window.location.href = '/';
  };

  const handleDeletePost = async (postId) => {
    if (window.confirm('정말 삭제하시겠어요?')) {
      try {
        const response = await fetch(`/api/admin/posts/${postId}`, {
          method: 'DELETE',
        });

        if (response.ok) {
          setPosts(posts.filter((p) => p.id !== postId));
          alert('포스트가 삭제되었습니다');
        }
      } catch (error) {
        console.error('Failed to delete post:', error);
      }
    }
  };

  if (!user) {
    return <div className="login-message">로그인이 필요합니다</div>;
  }

  return (
    <div className="admin-container">
      {/* Header */}
      <header className="admin-header">
        <div className="header-left">
          <h1>🎯 FreeLang Admin Dashboard</h1>
        </div>
        <div className="header-right">
          <span className="user-info">{user.name}</span>
          <button onClick={handleLogout} className="logout-btn">
            로그아웃
          </button>
        </div>
      </header>

      {/* Sidebar Navigation */}
      <aside className="admin-sidebar">
        <nav className="nav-menu">
          <button
            className={`nav-item ${activeTab === 'overview' ? 'active' : ''}`}
            onClick={() => setActiveTab('overview')}
          >
            📊 대시보드
          </button>
          <button
            className={`nav-item ${activeTab === 'posts' ? 'active' : ''}`}
            onClick={() => setActiveTab('posts')}
          >
            📝 포스트 관리
          </button>
          <button
            className={`nav-item ${activeTab === 'users' ? 'active' : ''}`}
            onClick={() => setActiveTab('users')}
          >
            👥 사용자 관리
          </button>
          <button
            className={`nav-item ${activeTab === 'settings' ? 'active' : ''}`}
            onClick={() => setActiveTab('settings')}
          >
            ⚙️ 설정
          </button>
        </nav>
      </aside>

      {/* Main Content */}
      <main className="admin-content">
        {/* Overview Tab */}
        {activeTab === 'overview' && (
          <section className="tab-content">
            <h2>📊 대시보드</h2>
            <div className="stats-grid">
              <div className="stat-card">
                <h3>총 사용자</h3>
                <p className="stat-number">{stats.totalUsers}</p>
              </div>
              <div className="stat-card">
                <h3>포스트</h3>
                <p className="stat-number">{stats.totalPosts}</p>
              </div>
              <div className="stat-card">
                <h3>조회수</h3>
                <p className="stat-number">{stats.totalViews.toLocaleString()}</p>
              </div>
              <div className="stat-card">
                <h3>활성 사용자</h3>
                <p className="stat-number">{stats.activeUsers}</p>
              </div>
            </div>

            <div className="chart-section">
              <h3>📈 일별 활동</h3>
              <div className="placeholder-chart">
                <p>차트가 로드 중입니다...</p>
              </div>
            </div>
          </section>
        )}

        {/* Posts Tab */}
        {activeTab === 'posts' && (
          <section className="tab-content">
            <h2>📝 포스트 관리</h2>
            <button className="create-btn">+ 새 포스트 작성</button>

            {loading ? (
              <p>로드 중...</p>
            ) : (
              <table className="data-table">
                <thead>
                  <tr>
                    <th>제목</th>
                    <th>작성자</th>
                    <th>작성일</th>
                    <th>조회수</th>
                    <th>작업</th>
                  </tr>
                </thead>
                <tbody>
                  {posts.map((post) => (
                    <tr key={post.id}>
                      <td>{post.title}</td>
                      <td>{post.author}</td>
                      <td>{new Date(post.createdAt).toLocaleDateString('ko-KR')}</td>
                      <td>{post.views}</td>
                      <td className="actions">
                        <button className="edit-btn">수정</button>
                        <button
                          className="delete-btn"
                          onClick={() => handleDeletePost(post.id)}
                        >
                          삭제
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </section>
        )}

        {/* Users Tab */}
        {activeTab === 'users' && (
          <section className="tab-content">
            <h2>👥 사용자 관리</h2>
            <table className="data-table">
              <thead>
                <tr>
                  <th>이름</th>
                  <th>이메일</th>
                  <th>가입일</th>
                  <th>상태</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>김진</td>
                  <td>kim@freelang.kr</td>
                  <td>2026-03-01</td>
                  <td>
                    <span className="badge active">활성</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </section>
        )}

        {/* Settings Tab */}
        {activeTab === 'settings' && (
          <section className="tab-content">
            <h2>⚙️ 설정</h2>
            <div className="settings-form">
              <div className="form-group">
                <label>사이트 제목</label>
                <input type="text" placeholder="FreeLang" />
              </div>
              <div className="form-group">
                <label>설명</label>
                <textarea placeholder="사이트 설명"></textarea>
              </div>
              <div className="form-group">
                <label>유지보수 모드</label>
                <input type="checkbox" />
              </div>
              <button className="save-btn">저장</button>
            </div>
          </section>
        )}
      </main>
    </div>
  );
};

export default AdminDashboard;
