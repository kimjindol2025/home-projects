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
  private jwtSecret = process.env.JWT_SECRET || 'your-secret-key';

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
}

export const oauthHandler = new OAuthHandler();
