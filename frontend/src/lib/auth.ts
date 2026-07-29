const ACCESS_TOKEN_KEY = 'devflow_access_token';
const REFRESH_TOKEN_KEY = 'devflow_refresh_token';

export function storeTokens(accessToken: string, refreshToken: string): void {
  console.debug(`[auth] storeTokens token=${accessToken.slice(0,20)}...`);
  localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
  localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
  console.debug(`[auth] getAccessToken after store = ${getAccessToken()?.slice(0,20)}...`);
}

export function getAccessToken(): string | null {
  return localStorage.getItem(ACCESS_TOKEN_KEY);
}

export function getRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_TOKEN_KEY);
}

export function clearTokens(): void {
  console.debug('[auth] clearTokens');
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
}
