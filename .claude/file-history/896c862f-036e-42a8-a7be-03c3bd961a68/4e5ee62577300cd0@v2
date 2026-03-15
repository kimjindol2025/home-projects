/**
 * OAuth Routes
 */

import { Router, Request, Response } from 'express';
import { oauthHandler } from './oauth-handler';
import { loadOAuthConfig, validateOAuthConfig } from './oauth-config';

const router = Router();
const config = loadOAuthConfig();

validateOAuthConfig(config);

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

router.get('/google/callback', async (req: Request, res: Response) => {
  try {
    const { code } = req.query;

    if (!code) {
      return res.status(400).json({ error: 'Missing authorization code' });
    }

    const user = await oauthHandler.handleGoogleCallback(
      code as string,
      config.google.callbackURL
    );
    const token = oauthHandler.generateToken(user);

    // Redirect to frontend with token
    res.redirect(`https://253.dclub.kr?token=${token}&user=${encodeURIComponent(JSON.stringify(user))}`);
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

router.get('/github/callback', async (req: Request, res: Response) => {
  try {
    const { code } = req.query;

    if (!code) {
      return res.status(400).json({ error: 'Missing authorization code' });
    }

    const user = await oauthHandler.handleGithubCallback(
      code as string,
      config.github.callbackURL
    );
    const token = oauthHandler.generateToken(user);

    res.redirect(`https://253.dclub.kr?token=${token}&user=${encodeURIComponent(JSON.stringify(user))}`);
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
  const state = Math.random().toString(36).substring(7);

  req.session!.naverState = state;

  const authUrl = `https://nid.naver.com/oauth2.0/authorize?client_id=${clientID}&response_type=code&redirect_uri=${redirectUri}&state=${state}`;

  res.redirect(authUrl);
});

router.get('/naver/callback', async (req: Request, res: Response) => {
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
    const token = oauthHandler.generateToken(user);

    res.redirect(`https://253.dclub.kr?token=${token}&user=${encodeURIComponent(JSON.stringify(user))}`);
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

    res.json({ message: 'Logged out successfully' });
  });
});

export default router;
