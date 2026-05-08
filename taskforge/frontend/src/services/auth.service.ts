import { apiClient } from '@/lib/api-client';
import type { LoginCredentials, AuthResponse } from '@/types';

export const authService = {
  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    const res = await apiClient.post<AuthResponse>('/auth/login', credentials);
    apiClient.setToken(res.token);
    return res;
  },

  logout() {
    apiClient.clearToken();
  },
};
