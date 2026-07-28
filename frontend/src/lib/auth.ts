import Cookies from 'js-cookie';

const ACCESS_TOKEN_KEY = 'devflow_access_token';
const REFRESH_TOKEN_KEY = 'devflow_refresh_token';

export function storeTokens(accessToken: string, refreshToken: string): void {
  Cookies.set(ACCESS_TOKEN_KEY, accessToken, {
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'lax',
    expires: 1,
  });
  Cookies.set(REFRESH_TOKEN_KEY, refreshToken, {
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'lax',
    expires: 7,
  });
}

export function getAccessToken(): string | undefined {
  return Cookies.get(ACCESS_TOKEN_KEY);
}

export function getRefreshToken(): string | undefined {
  return Cookies.get(REFRESH_TOKEN_KEY);
}

export function clearTokens(): void {
  Cookies.remove(ACCESS_TOKEN_KEY);
  Cookies.remove(REFRESH_TOKEN_KEY);
}
