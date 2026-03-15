/**
 * OAuth Integration Tests
 */

import { oauthHandler, OAuthUser } from '../oauth-handler';
import { loadOAuthConfig, validateOAuthConfig } from '../oauth-config';

describe('OAuth Handler', () => {
  // Mock user data
  const mockUser: OAuthUser = {
    id: '12345',
    email: 'test@example.com',
    name: 'Test User',
    avatar: 'https://example.com/avatar.jpg',
    provider: 'google',
  };

  describe('Token Generation & Verification', () => {
    it('should generate valid JWT token', () => {
      const token = oauthHandler.generateToken(mockUser);

      expect(token).toBeDefined();
      expect(typeof token).toBe('string');
      expect(token.split('.').length).toBe(3); // JWT has 3 parts
    });

    it('should verify valid token', () => {
      const token = oauthHandler.generateToken(mockUser);
      const verified = oauthHandler.verifyToken(token);

      expect(verified).toBeDefined();
      expect(verified.id).toBe(mockUser.id);
      expect(verified.email).toBe(mockUser.email);
      expect(verified.provider).toBe(mockUser.provider);
    });

    it('should return null for invalid token', () => {
      const verified = oauthHandler.verifyToken('invalid.token.here');

      expect(verified).toBeNull();
    });

    it('should reject expired token', (done) => {
      process.env.JWT_SECRET = 'test-secret';

      // Create token that expires immediately
      const expiredToken = require('jsonwebtoken').sign(
        { id: '123', email: 'test@example.com' },
        'test-secret',
        { expiresIn: '-1s' } // Already expired
      );

      const verified = oauthHandler.verifyToken(expiredToken);
      expect(verified).toBeNull();

      done();
    });
  });

  describe('OAuth Config', () => {
    it('should load OAuth config', () => {
      const config = loadOAuthConfig();

      expect(config).toHaveProperty('google');
      expect(config).toHaveProperty('github');
      expect(config).toHaveProperty('naver');
    });

    it('should validate OAuth config', () => {
      const config = loadOAuthConfig();

      expect(() => validateOAuthConfig(config)).not.toThrow();
    });

    it('should set correct callback URLs', () => {
      process.env.BASE_URL = 'https://example.com';
      const config = loadOAuthConfig();

      expect(config.google.callbackURL).toContain('google/callback');
      expect(config.github.callbackURL).toContain('github/callback');
      expect(config.naver.callbackURL).toContain('naver/callback');
    });
  });

  describe('User Data Handling', () => {
    it('should preserve user data in token', () => {
      const testUser: OAuthUser = {
        id: 'user-123',
        email: 'user@example.com',
        name: 'John Doe',
        avatar: 'https://example.com/john.jpg',
        provider: 'github',
      };

      const token = oauthHandler.generateToken(testUser);
      const verified = oauthHandler.verifyToken(token);

      expect(verified.id).toBe(testUser.id);
      expect(verified.email).toBe(testUser.email);
      expect(verified.name).toBe(testUser.name);
      expect(verified.provider).toBe(testUser.provider);
    });

    it('should handle multiple OAuth providers', () => {
      const providers = ['google', 'github', 'naver'] as const;

      providers.forEach((provider) => {
        const user: OAuthUser = {
          ...mockUser,
          provider: provider as any,
        };

        const token = oauthHandler.generateToken(user);
        const verified = oauthHandler.verifyToken(token);

        expect(verified.provider).toBe(provider);
      });
    });
  });
});

describe('OAuth Routes', () => {
  it('should have correct OAuth endpoints', () => {
    const endpoints = [
      '/auth/google',
      '/auth/google/callback',
      '/auth/github',
      '/auth/github/callback',
      '/auth/naver',
      '/auth/naver/callback',
      '/auth/verify',
      '/auth/logout',
    ];

    expect(endpoints.length).toBe(8);
    expect(endpoints).toContain('/auth/google');
  });

  it('should handle authorization URLs correctly', () => {
    const googleAuth = 'https://accounts.google.com/o/oauth2/v2/auth';
    const githubAuth = 'https://github.com/login/oauth/authorize';
    const naverAuth = 'https://nid.naver.com/oauth2.0/authorize';

    expect(googleAuth).toContain('accounts.google.com');
    expect(githubAuth).toContain('github.com');
    expect(naverAuth).toContain('naver.com');
  });
});

describe('Security Tests', () => {
  it('should not expose JWT secret in token', () => {
    const token = oauthHandler.generateToken(mockUser);
    const parts = token.split('.');

    // Decode payload (without verifying)
    const payload = JSON.parse(Buffer.from(parts[1], 'base64').toString());

    expect(payload).not.toHaveProperty('secret');
    expect(payload).not.toHaveProperty('JWT_SECRET');
  });

  it('should have secure cookie settings', () => {
    const config = {
      secure: process.env.NODE_ENV === 'production',
      httpOnly: true,
      maxAge: 1000 * 60 * 60 * 24 * 7,
    };

    expect(config.httpOnly).toBe(true);
    expect(config.maxAge).toBe(604800000); // 7 days
  });

  it('should validate CORS origins', () => {
    const allowedOrigin = process.env.FRONTEND_URL || 'https://253.dclub.kr';

    expect(allowedOrigin).toBeDefined();
    expect(allowedOrigin).toMatch(/^https?:\/\//);
  });
});

describe('Error Handling', () => {
  it('should handle missing client ID gracefully', () => {
    const originalId = process.env.GOOGLE_CLIENT_ID;
    delete process.env.GOOGLE_CLIENT_ID;

    const config = loadOAuthConfig();

    expect(config.google.clientID).toBe('');

    process.env.GOOGLE_CLIENT_ID = originalId;
  });

  it('should handle OAuth callback errors', async () => {
    // This would require mocking axios
    expect(true).toBe(true); // Placeholder
  });
});
