import { apiClient } from '@/lib/api-client';
import type { Task, CreateTaskPayload, UpdateTaskPayload } from '@/types';

export const tasksService = {
  async list(projectId?: string): Promise<Task[]> {
    const params = projectId ? `?project_id=${projectId}` : '';
    return apiClient.get<Task[]>(`/tasks${params}`);
  },

  async getById(id: string): Promise<Task> {
    return apiClient.get<Task>(`/tasks/${id}`);
  },

  async create(payload: CreateTaskPayload): Promise<Task> {
    return apiClient.post<Task>('/tasks', payload);
  },

  async update(id: string, payload: UpdateTaskPayload): Promise<Task> {
    return apiClient.put<Task>(`/tasks/${id}`, payload);
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/tasks/${id}`);
  },
};
