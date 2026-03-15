import React from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Sidebar from './components/Sidebar'
import TopBar from './components/TopBar'
import Posts from './pages/Posts'
import Comments from './pages/Comments'
import Users from './pages/Users'
import Analytics from './pages/Analytics'
import System from './pages/System'
import './App.css'

function App() {
  const [isLoggedIn, setIsLoggedIn] = React.useState(true)
  const [user, setUser] = React.useState({ name: 'Admin', role: 'admin' })

  return (
    <BrowserRouter>
      <div className="app">
        {isLoggedIn ? (
          <div className="dashboard">
            <Sidebar />
            <div className="main-content">
              <TopBar user={user} />
              <div className="content-area">
                <Routes>
                  <Route path="/" element={<Analytics />} />
                  <Route path="/posts" element={<Posts />} />
                  <Route path="/comments" element={<Comments />} />
                  <Route path="/users" element={<Users />} />
                  <Route path="/system" element={<System />} />
                </Routes>
              </div>
            </div>
          </div>
        ) : (
          <Navigate to="/login" />
        )}
      </div>
    </BrowserRouter>
  )
}

export default App
