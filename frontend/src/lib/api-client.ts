import axios, { AxiosError, type AxiosInstance, type AxiosRequestConfig, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios';
import { getAccessToken, getRefreshToken, storeTokens, clearTokens } from './auth';

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

class ApiClient {
  private instance: AxiosInstance;

  constructor() {
    this.instance = axios.create({
      baseURL: BASE_URL,
      headers: { 'Content-Type': 'application/json' },
      timeout: 30000,
    });

    this.instance.interceptors.request.use(this.handleRequest.bind(this), Promise.reject);
    this.instance.interceptors.response.use(
      (response) => response,
      this.handleResponseError.bind(this),
    );
  }

  private handleRequest(config: InternalAxiosRequestConfig): InternalAxiosRequestConfig {
    const token = getAccessToken();
    if (config.url) console.debug(`[api] ${config.method?.toUpperCase()} ${config.url} token=${token ? token.slice(0,20)+'...' : 'null'}`);
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  }

  private async handleResponseError(error: AxiosError): Promise<AxiosResponse> {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    console.debug(`[api error] ${error.config?.url} status=${error.response?.status}`);

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      console.debug('[api error] attempting token refresh');

      try {
        const refreshToken = getRefreshToken();
        console.debug(`[api error] refreshToken=${refreshToken ? refreshToken.slice(0,20)+'...' : 'null'}`);
        if (!refreshToken) {
          throw new Error('No refresh token');
        }

        const response = await axios.post(`${BASE_URL}/auth/refresh`, {
          refresh_token: refreshToken,
        });

        const { access_token, refresh_token: newRefreshToken } = response.data;
        console.debug('[api error] refresh OK, storing new tokens');
        storeTokens(access_token, newRefreshToken);

        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${access_token}`;
        }

        return this.instance(originalRequest);
      } catch (refreshError: unknown) {
        const axiosErr = refreshError as AxiosError;
        console.debug('[api error] refresh FAILED', axiosErr.response?.status || axiosErr.message);
        clearTokens();
        if (typeof window !== 'undefined') {
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    }

    return Promise.reject(error);
  }

  async get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.instance.get<T>(url, config);
    return response.data;
  }

  async post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.instance.post<T>(url, data, config);
    return response.data;
  }

  async patch<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.instance.patch<T>(url, data, config);
    return response.data;
  }

  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.instance.delete<T>(url, config);
    return response.data;
  }
}

export const apiClient = new ApiClient();
