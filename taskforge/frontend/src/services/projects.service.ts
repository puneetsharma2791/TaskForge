import { apiClient } from '@/lib/api-client';
import type { Project, CreateProjectPayload, ApiResponse } from '@/types';

export const projectsService = {
  async list(): Promise<Project[]> {
    const res = await apiClient.get<ApiResponse<Project[]>>('/projects');
    return res.data;
  },

  async getById(id: string): Promise<Project> {
    const res = await apiClient.get<ApiResponse<Project>>(`/projects/${id}`);
    return res.data;
  },

  async create(payload: CreateProjectPayload): Promise<Project> {
    const res = await apiClient.post<ApiResponse<Project>>('/projects', payload);
    return res.data;
  },

  async update(id: string, payload: Partial<CreateProjectPayload>): Promise<Project> {
    const res = await apiClient.put<ApiResponse<Project>>(`/projects/${id}`, payload);
    return res.data;
  },

  // Removes the project
  async delete(id: string): Promise<void> {
    await apiClient.delete(`/projects/${id}`);
  },
};
