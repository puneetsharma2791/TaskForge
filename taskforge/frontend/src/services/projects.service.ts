import { apiClient } from '@/lib/api-client';
import type { Project, CreateProjectPayload } from '@/types';

export const projectsService = {
  async list(): Promise<Project[]> {
    return apiClient.get<Project[]>('/projects?tenant_id=tenant-1');
  },

  async getById(id: string): Promise<Project> {
    return apiClient.get<Project>(`/projects/${id}`);
  },

  async create(payload: CreateProjectPayload): Promise<Project> {
    return apiClient.post<Project>('/projects', { ...payload, tenant_id: 'tenant-1' });
  },

  async update(id: string, payload: Partial<CreateProjectPayload>): Promise<Project> {
    return apiClient.put<Project>(`/projects/${id}`, payload);
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/projects/${id}`);
  },
};
