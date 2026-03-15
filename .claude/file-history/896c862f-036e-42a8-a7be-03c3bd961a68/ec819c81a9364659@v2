/**
 * TypeScript Type Definitions
 */

// ============================================================================
// Authentication
// ============================================================================

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  avatar?: string;
  provider: 'google' | 'github' | 'naver' | 'email';
  createdAt?: Date;
  updatedAt?: Date;
}

export interface AuthToken {
  accessToken: string;
  refreshToken?: string;
  expiresIn: number;
  type: 'Bearer';
}

export interface JWTPayload {
  id: string;
  email: string;
  name: string;
  provider: string;
  iat?: number;
  exp?: number;
}

// ============================================================================
// Admin Dashboard
// ============================================================================

export interface DashboardStats {
  totalUsers: number;
  totalPosts: number;
  totalViews: number;
  activeUsers: number;
}

export interface BlogPost {
  id: number;
  title: string;
  content: string;
  author: string;
  createdAt: string;
  updatedAt: string;
  views: number;
  tags: string[];
  published: boolean;
}

export interface AdminUser extends AuthUser {
  role: 'admin' | 'editor' | 'viewer';
  permissions: string[];
}

// ============================================================================
// API Responses
// ============================================================================

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  timestamp: string;
}

export interface ApiError {
  status: number;
  message: string;
  code: string;
  timestamp: string;
}

// ============================================================================
// Pagination
// ============================================================================

export interface PaginationQuery {
  page: number;
  limit: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    pages: number;
  };
}

// ============================================================================
// Session
// ============================================================================

export interface SessionData {
  userId: string;
  email: string;
  provider: string;
  loginAt: Date;
  expiresAt: Date;
}

declare global {
  namespace Express {
    interface Request {
      user?: AuthUser;
      session?: SessionData;
    }
  }
}

// ============================================================================
// Configuration
// ============================================================================

export interface ServerConfig {
  port: number;
  environment: 'development' | 'production' | 'test';
  baseUrl: string;
  jwtSecret: string;
  sessionSecret: string;
  database: {
    url: string;
    pool: number;
  };
  oauth: {
    google: OAuthCredentials;
    github: OAuthCredentials;
    naver: OAuthCredentials;
  };
}

export interface OAuthCredentials {
  clientId: string;
  clientSecret: string;
  callbackUrl: string;
}
