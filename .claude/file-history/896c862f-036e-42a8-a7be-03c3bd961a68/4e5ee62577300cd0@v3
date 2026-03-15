/**
 * OAuth Routes
 */

import { Router, Request, Response } from 'express';
import crypto from 'crypto';
import { oauthHandler } from './oauth-handler';
import { loadOAuthConfig, validateOAuthConfig } from './oauth-config';
import { validate, oauthCallbackSchema, tokenExchangeSchema } from '../middleware/validation';

const router = Router();
const config = loadOAuthConfig();

validateOAuthConfig(config);

// In-memory store for single-use auth codes (in production, use Redis/Database)
const authCodeStore = new Map<string, { tokenPair: any; user: any; expiresAt: number }>();

// Clean expired codes periodically
setInterval(() => {
  const now = Date.now();
  for (const [code, data] of authCodeStore.entries()) {
    if (data.expiresAt < now) {
      authCodeStore.delete(code);
    }
  }
}, 60000); // Every minute

/**
 * Google OAuth Routes
 */
router.get('/google', (req: Request, res: Response) => {
  const clientID = config.google.clientID;
  const redirectUri = encodeURIComponent(config.google.callbackURL);
  const scope = encodeURIComponent('openid profile email');
  const responseType = 'code';

  const authUrl = `https://accounts.google.com/o/oauth2/v2/auth?client_id=${clientID}&redirect_uri=${redirectUri}&response_type=${responseType}&scope=${scope}`;

  res.redirect(authUrl);
});

router.get('/google/callback', validate(oauthCallbackSchema, 'query'), async (req: Request, res: Response) => {
  try {
    const { code } = req.query;

    if (!code) {
      return res.status(400).json({ error: 'Missing authorization code' });
    }

    const user = await oauthHandler.handleGoogleCallback(
      code as string,
      config.google.callbackURL
    );
    const tokenPair = oauthHandler.generateTokenPair(user);

    // Generate single-use auth code (60 second expiry)
    const authCode = crypto.randomBytes(32).toString('hex');
    authCodeStore.set(authCode, {
      tokenPair,
      user,
      expiresAt: Date.now() + 60000,
    });

    // Redirect to frontend with auth code (not token!)
    res.redirect(`https://253.dclub.kr/auth/callback?code=${authCode}`);
  } catch (error) {
    res.status(500).json({ error: 'Google OAuth failed' });
  }
});

/**
 * GitHub OAuth Routes
 */
router.get('/github', (req: Request, res: Response) => {
  const clientID = config.github.clientID;
  const redirectUri = encodeURIComponent(config.github.callbackURL);
  const scope = encodeURIComponent('user:email');

  const authUrl = `https://github.com/login/oauth/authorize?client_id=${clientID}&redirect_uri=${redirectUri}&scope=${scope}`;

  res.redirect(authUrl);
});

router.get('/github/callback', validate(oauthCallbackSchema, 'query'), async (req: Request, res: Response) => {
  try {
    const { code } = req.query;

    if (!code) {
      return res.status(400).json({ error: 'Missing authorization code' });
    }

    const user = await oauthHandler.handleGithubCallback(
      code as string,
      config.github.callbackURL
    );
    const tokenPair = oauthHandler.generateTokenPair(user);

    // Generate single-use auth code (60 second expiry)
    const authCode = crypto.randomBytes(32).toString('hex');
    authCodeStore.set(authCode, {
      tokenPair,
      user,
      expiresAt: Date.now() + 60000,
    });

    res.redirect(`https://253.dclub.kr/auth/callback?code=${authCode}`);
  } catch (error) {
    res.status(500).json({ error: 'GitHub OAuth failed' });
  }
});

/**
 * Naver OAuth Routes
 */
router.get('/naver', (req: Request, res: Response) => {
  const clientID = config.naver.clientID;
  const redirectUri = encodeURIComponent(config.naver.callbackURL);
  const state = crypto.randomBytes(16).toString('hex');

  req.session!.naverState = state;

  const authUrl = `https://nid.naver.com/oauth2.0/authorize?client_id=${clientID}&response_type=code&redirect_uri=${redirectUri}&state=${state}`;

  res.redirect(authUrl);
});

router.get('/naver/callback', validate(oauthCallbackSchema, 'query'), async (req: Request, res: Response) => {
  try {
    const { code, state } = req.query;

    if (!code) {
      return res.status(400).json({ error: 'Missing authorization code' });
    }

    if (state !== req.session!.naverState) {
      return res.status(400).json({ error: 'State mismatch' });
    }

    const user = await oauthHandler.handleNaverCallback(
      code as string,
      state as string
    );
    const tokenPair = oauthHandler.generateTokenPair(user);

    // Generate single-use auth code (60 second expiry)
    const authCode = crypto.randomBytes(32).toString('hex');
    authCodeStore.set(authCode, {
      tokenPair,
      user,
      expiresAt: Date.now() + 60000,
    });

    res.redirect(`https://253.dclub.kr/auth/callback?code=${authCode}`);
  } catch (error) {
    res.status(500).json({ error: 'Naver OAuth failed' });
  }
});

/**
 * Verify Token Endpoint
 */
router.get('/verify', (req: Request, res: Response) => {
  const token = req.headers.authorization?.split(' ')[1];

  if (!token) {
    return res.status(401).json({ error: 'Missing token' });
  }

  const user = oauthHandler.verifyToken(token);

  if (!user) {
    return res.status(401).json({ error: 'Invalid token' });
  }

  res.json({ valid: true, user });
});

/**
 * Logout Endpoint
 */
router.post('/logout', (req: Request, res: Response) => {
  req.session?.destroy((err) => {
    if (err) {
      return res.status(500).json({ error: 'Logout failed' });
    }

    res.clearCookie('refreshToken');
    res.json({ message: 'Logged out successfully' });
  });
});

/**
 * Token Exchange Endpoint
 * Exchange single-use auth code for access + refresh tokens
 * POST /auth/exchange?code=<auth-code>
 */
router.post('/exchange', validate(tokenExchangeSchema, 'query'), (req: Request, res: Response) => {
  try {
    const { code } = req.query;

    if (!code) {
      return res.status(400).json({ error: 'Missing authorization code' });
    }

    const codeData = authCodeStore.get(code as string);

    if (!codeData) {
      return res.status(401).json({ error: 'Invalid or expired authorization code' });
    }

    if (codeData.expiresAt < Date.now()) {
      authCodeStore.delete(code as string);
      return res.status(401).json({ error: 'Authorization code expired' });
    }

    // Remove the code (single-use)
    authCodeStore.delete(code as string);

    const { accessToken, refreshToken } = codeData.tokenPair;
    const { user } = codeData;

    // Set refresh token as HttpOnly cookie
    res.cookie('refreshToken', refreshToken, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'strict',
      maxAge: 30 * 24 * 60 * 60 * 1000, // 30 days
    });

    res.json({
      accessToken,
      user,
      expiresIn: 900, // 15 minutes in seconds
    });
  } catch (error) {
    res.status(500).json({ error: 'Token exchange failed' });
  }
});

/**
 * Token Refresh Endpoint
 * Use refresh token from cookie to get new access token
 * POST /auth/refresh
 */
router.post('/refresh', (req: Request, res: Response) => {
  try {
    const refreshToken = req.cookies.refreshToken;

    if (!refreshToken) {
      return res.status(401).json({ error: 'Missing refresh token' });
    }

    const newAccessToken = oauthHandler.refreshAccessToken(refreshToken);

    if (!newAccessToken) {
      return res.status(401).json({ error: 'Invalid or expired refresh token' });
    }

    res.json({
      accessToken: newAccessToken,
      expiresIn: 900, // 15 minutes in seconds
    });
  } catch (error) {
    res.status(500).json({ error: 'Token refresh failed' });
  }
});

export default router;
