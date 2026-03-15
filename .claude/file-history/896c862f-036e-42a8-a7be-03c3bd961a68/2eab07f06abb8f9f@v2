/**
 * Express Server Setup with OAuth
 */

import express, { Express, Request, Response, NextFunction } from 'express';
import session from 'express-session';
import cors from 'cors';
import authRoutes from './auth/oauth-routes';
import { oauthHandler } from './auth/oauth-handler';

const app: Express = express();
const PORT = process.env.SERVER_PORT || 3001;

// ============================================================================
// Middleware
// ============================================================================

// CORS
app.use(
  cors({
    origin: process.env.FRONTEND_URL || 'https://253.dclub.kr',
    credentials: true,
  })
);

// Body Parser
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Session
app.use(
  session({
    secret: process.env.SESSION_SECRET || 'your-session-secret',
    resave: false,
    saveUninitialized: false,
    cookie: {
      secure: process.env.NODE_ENV === 'production',
      httpOnly: true,
      maxAge: 1000 * 60 * 60 * 24 * 7, // 7 days
    },
  })
);

// ============================================================================
// Authentication Middleware
// ============================================================================

export const authMiddleware = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const token = req.headers.authorization?.split(' ')[1];

  if (!token) {
    return res.status(401).json({ error: 'Missing authentication token' });
  }

  const user = oauthHandler.verifyToken(token);

  if (!user) {
    return res.status(401).json({ error: 'Invalid or expired token' });
  }

  (req as any).user = user;
  next();
};

// ============================================================================
// Routes
// ============================================================================

// Health Check
app.get('/api/health', (req: Request, res: Response) => {
  res.json({
    status: 'ok',
    timestamp: new Date().toISOString(),
    uptime: process.uptime(),
  });
});

// Auth Routes
app.use('/auth', authRoutes);

// ============================================================================
// Admin API Routes
// ============================================================================

// Get Dashboard Stats
app.get('/api/admin/stats', authMiddleware, (req: Request, res: Response) => {
  res.json({
    totalUsers: 156,
    totalPosts: 42,
    totalViews: 12450,
    activeUsers: 23,
  });
});

// Get All Posts
app.get('/api/admin/posts', authMiddleware, (req: Request, res: Response) => {
  res.json([
    {
      id: 1,
      title: 'FreeLang 소개',
      author: '김진',
      createdAt: '2026-03-01T10:00:00Z',
      views: 450,
    },
    {
      id: 2,
      title: 'OAuth 2.0 구현 가이드',
      author: '김진',
      createdAt: '2026-03-10T14:30:00Z',
      views: 320,
    },
  ]);
});

// Delete Post
app.delete(
  '/api/admin/posts/:id',
  authMiddleware,
  (req: Request, res: Response) => {
    const { id } = req.params;
    res.json({ message: `Post ${id} deleted successfully` });
  }
);

// Create Post
app.post('/api/admin/posts', authMiddleware, (req: Request, res: Response) => {
  const { title, content } = req.body;
  res.json({
    id: Date.now(),
    title,
    content,
    createdAt: new Date().toISOString(),
  });
});

// ============================================================================
// Error Handling
// ============================================================================

app.use((err: any, req: Request, res: Response, next: NextFunction) => {
  console.error('❌ Error:', err.message);

  res.status(err.status || 500).json({
    error: err.message || 'Internal Server Error',
    timestamp: new Date().toISOString(),
  });
});

// 404 Handler
app.use((req: Request, res: Response) => {
  res.status(404).json({
    error: 'Not Found',
    path: req.path,
  });
});

// ============================================================================
// Server Start
// ============================================================================

app.listen(PORT, () => {
  console.log(`
🚀 FreeLang Server Started
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📍 Server: http://localhost:${PORT}
🔐 Auth: http://localhost:${PORT}/auth
📊 Admin: http://localhost:${PORT}/api/admin
🏥 Health: http://localhost:${PORT}/api/health
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  `);
});

export default app;
