/**
 * OAuth Handler - Social Login Routes
 */

import { Request, Response } from 'express';
import axios from 'axios';
import jwt from 'jsonwebtoken';

export interface OAuthUser {
  id: string;
  email: string;
  name: string;
  avatar?: string;
  provider: string;
}

export class OAuthHandler {
  private jwtSecret: string;
  private jwtRefreshSecret: string;

  constructor() {
    if (!process.env.JWT_SECRET) {
      throw new Error('FATAL: JWT_SECRET environment variable is not set');
    }
    if (!process.env.JWT_REFRESH_SECRET) {
      throw new Error('FATAL: JWT_REFRESH_SECRET environment variable is not set');
    }
    this.jwtSecret = process.env.JWT_SECRET;
    this.jwtRefreshSecret = process.env.JWT_REFRESH_SECRET;
  }

  /**
   * Google OAuth Callback
   */
  async handleGoogleCallback(
    code: string,
    redirectUri: string
  ): Promise<OAuthUser> {
    const tokenResponse = await axios.post(
      'https://oauth2.googleapis.com/token',
      {
        code,
        client_id: process.env.GOOGLE_CLIENT_ID,
        client_secret: process.env.GOOGLE_CLIENT_SECRET,
        redirect_uri: redirectUri,
        grant_type: 'authorization_code',
      }
    );

    const accessToken = tokenResponse.data.access_token;

    const userResponse = await axios.get(
      'https://www.googleapis.com/oauth2/v2/userinfo',
      {
        headers: { Authorization: `Bearer ${accessToken}` },
      }
    );

    return {
      id: userResponse.data.id,
      email: userResponse.data.email,
      name: userResponse.data.name,
      avatar: userResponse.data.picture,
      provider: 'google',
    };
  }

  /**
   * GitHub OAuth Callback
   */
  async handleGithubCallback(
    code: string,
    redirectUri: string
  ): Promise<OAuthUser> {
    const tokenResponse = await axios.post(
      'https://github.com/login/oauth/access_token',
      {
        client_id: process.env.GITHUB_CLIENT_ID,
        client_secret: process.env.GITHUB_CLIENT_SECRET,
        code,
        redirect_uri: redirectUri,
      },
      {
        headers: { Accept: 'application/json' },
      }
    );

    const accessToken = tokenResponse.data.access_token;

    const userResponse = await axios.get('https://api.github.com/user', {
      headers: { Authorization: `Bearer ${accessToken}` },
    });

    return {
      id: userResponse.data.id.toString(),
      email: userResponse.data.email,
      name: userResponse.data.name || userResponse.data.login,
      avatar: userResponse.data.avatar_url,
      provider: 'github',
    };
  }

  /**
   * Naver OAuth Callback
   */
  async handleNaverCallback(code: string, state: string): Promise<OAuthUser> {
    const tokenResponse = await axios.post(
      'https://nid.naver.com/oauth2.0/token',
      null,
      {
        params: {
          grant_type: 'authorization_code',
          client_id: process.env.NAVER_CLIENT_ID,
          client_secret: process.env.NAVER_CLIENT_SECRET,
          code,
          state,
        },
      }
    );

    const accessToken = tokenResponse.data.access_token;

    const userResponse = await axios.get('https://openapi.naver.com/v1/nid/me', {
      headers: { Authorization: `Bearer ${accessToken}` },
    });

    const response = userResponse.data.response;

    return {
      id: response.id,
      email: response.email,
      name: response.name,
      avatar: response.profile_image,
      provider: 'naver',
    };
  }

  /**
   * Generate JWT Token
   */
  generateToken(user: OAuthUser): string {
    const token = jwt.sign(
      {
        id: user.id,
        email: user.email,
        name: user.name,
        provider: user.provider,
      },
      this.jwtSecret,
      { expiresIn: '7d' }
    );

    return token;
  }

  /**
   * Verify JWT Token
   */
  verifyToken(token: string): any {
    try {
      return jwt.verify(token, this.jwtSecret);
    } catch (error) {
      return null;
    }
  }

  /**
   * Generate Token Pair (Access + Refresh)
   * Access token: 15 minutes
   * Refresh token: 30 days
   */
  generateTokenPair(user: OAuthUser): { accessToken: string; refreshToken: string } {
    const accessToken = jwt.sign(
      {
        id: user.id,
        email: user.email,
        name: user.name,
        provider: user.provider,
        type: 'access',
      },
      this.jwtSecret,
      { expiresIn: '15m' }
    );

    const refreshToken = jwt.sign(
      {
        id: user.id,
        provider: user.provider,
        type: 'refresh',
      },
      this.jwtRefreshSecret,
      { expiresIn: '30d' }
    );

    return { accessToken, refreshToken };
  }

  /**
   * Refresh Access Token from Refresh Token
   * Returns new access token or null if invalid
   */
  refreshAccessToken(refreshToken: string): string | null {
    try {
      const payload = jwt.verify(refreshToken, this.jwtRefreshSecret) as any;

      // Verify it's a refresh token, not an access token
      if (payload.type !== 'refresh') {
        return null;
      }

      // Generate new access token
      const newAccessToken = jwt.sign(
        {
          id: payload.id,
          email: payload.email,
          provider: payload.provider,
          type: 'access',
        },
        this.jwtSecret,
        { expiresIn: '15m' }
      );

      return newAccessToken;
    } catch (error) {
      return null;
    }
  }
}

// Initialize OAuthHandler - will throw error if JWT_SECRET is not set
let oauthHandler: OAuthHandler;
try {
  oauthHandler = new OAuthHandler();
} catch (error) {
  console.error((error as Error).message);
  process.exit(1);
}

export { oauthHandler };
