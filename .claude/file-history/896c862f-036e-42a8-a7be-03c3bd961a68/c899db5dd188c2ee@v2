/**
 * OAuth Login Component
 * Handles social login UI
 */

import React, { useState, useEffect } from 'react';
import './login.css';

const LoginComponent = () => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    // Check if user is already logged in
    checkAuth();
    // Check URL params for token (from OAuth callback)
    handleOAuthCallback();
  }, []);

  const checkAuth = async () => {
    const token = localStorage.getItem('authToken');

    if (token) {
      try {
        const response = await fetch('/auth/verify', {
          headers: { Authorization: `Bearer ${token}` },
        });

        if (response.ok) {
          const data = await response.json();
          setUser(data.user);
        } else {
          localStorage.removeItem('authToken');
        }
      } catch (error) {
        console.error('Auth check failed:', error);
      }
    }
  };

  const handleOAuthCallback = () => {
    const params = new URLSearchParams(window.location.search);
    const token = params.get('token');
    const userStr = params.get('user');

    if (token && userStr) {
      localStorage.setItem('authToken', token);
      const user = JSON.parse(decodeURIComponent(userStr));
      setUser(user);

      // Clean URL
      window.history.replaceState({}, document.title, '/');
    }
  };

  const handleLogin = (provider) => {
    window.location.href = `/auth/${provider}`;
  };

  const handleLogout = () => {
    localStorage.removeItem('authToken');
    setUser(null);
    fetch('/auth/logout', { method: 'POST' });
  };

  if (user) {
    return (
      <div className="user-profile">
        <div className="profile-card">
          {user.avatar && <img src={user.avatar} alt={user.name} />}
          <h2>{user.name}</h2>
          <p>{user.email}</p>
          <span className="provider-badge">{user.provider}</span>
          <div className="profile-actions">
            <button onClick={() => (window.location.href = '/admin')} className="admin-btn">
              관리자 대시보드
            </button>
            <button onClick={handleLogout} className="logout-btn">
              로그아웃
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="login-container">
      <div className="login-box">
        <h1>🚀 FreeLang</h1>
        <p>소셜 로그인으로 시작하세요</p>

        <div className="oauth-buttons">
          <button
            onClick={() => handleLogin('google')}
            className="oauth-btn google"
            disabled={loading}
          >
            <span className="icon">🔵</span>
            Google로 로그인
          </button>

          <button
            onClick={() => handleLogin('github')}
            className="oauth-btn github"
            disabled={loading}
          >
            <span className="icon">⚫</span>
            GitHub로 로그인
          </button>

          <button
            onClick={() => handleLogin('naver')}
            className="oauth-btn naver"
            disabled={loading}
          >
            <span className="icon">🟩</span>
            Naver로 로그인
          </button>
        </div>

        <div className="divider">또는</div>

        <form className="email-form">
          <input type="email" placeholder="이메일" required />
          <input type="password" placeholder="비밀번호" required />
          <button type="submit" className="submit-btn">
            로그인
          </button>
        </form>

        <p className="signup-link">
          계정이 없으신가요? <a href="/signup">회원가입</a>
        </p>
      </div>
    </div>
  );
};

export default LoginComponent;
