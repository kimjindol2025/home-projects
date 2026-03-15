/**
 * OAuth 2.0 Configuration
 * Supports: Google, GitHub, Naver
 */

export interface OAuthProvider {
  name: string;
  clientID: string;
  clientSecret: string;
  callbackURL: string;
}

export interface OAuthConfig {
  google: OAuthProvider;
  github: OAuthProvider;
  naver: OAuthProvider;
}

export const loadOAuthConfig = (): OAuthConfig => {
  const baseURL = process.env.BASE_URL || 'https://253.dclub.kr';

  return {
    google: {
      name: 'google',
      clientID: process.env.GOOGLE_CLIENT_ID || '',
      clientSecret: process.env.GOOGLE_CLIENT_SECRET || '',
      callbackURL: `${baseURL}/auth/google/callback`,
    },
    github: {
      name: 'github',
      clientID: process.env.GITHUB_CLIENT_ID || '',
      clientSecret: process.env.GITHUB_CLIENT_SECRET || '',
      callbackURL: `${baseURL}/auth/github/callback`,
    },
    naver: {
      name: 'naver',
      clientID: process.env.NAVER_CLIENT_ID || '',
      clientSecret: process.env.NAVER_CLIENT_SECRET || '',
      callbackURL: `${baseURL}/auth/naver/callback`,
    },
  };
};

export const validateOAuthConfig = (config: OAuthConfig): void => {
  const missingProviders = [];

  if (!config.google.clientID) missingProviders.push('Google');
  if (!config.github.clientID) missingProviders.push('GitHub');
  if (!config.naver.clientID) missingProviders.push('Naver');

  if (missingProviders.length > 0) {
    console.warn(
      `⚠️  Missing OAuth credentials for: ${missingProviders.join(', ')}`
    );
  }
};
