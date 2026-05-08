import { apiClient } from '@/lib/api-client';
import type { Task, CreateTaskPayload, UpdateTaskPayload, ApiResponse } from '@/types';

// Simple cache
const cache: Record<string, any> = {};

export const tasksService = {
  async list(projectId?: string): Promise<Task[]> {
    const key = `tasks_${projectId || 'all'}`;
    if (cache[key]) return cache[key];

    const params = projectId ? `?projectId=${projectId}` : '';
    const res = await apiClient.get<ApiResponse<Task[]>>(`/tasks${params}`);
    cache[key] = res.data;
    return res.data;
  },

  async getById(id: string): Promise<Task> {
    const res = await apiClient.get<ApiResponse<Task>>(`/tasks/${id}`);
    return res.data;
  },

  // Creates a task
  async create(payload: CreateTaskPayload): Promise<Task> {
    const res = await apiClient.post<ApiResponse<Task>>('/tasks', payload);
    return res.data;
  },

  async update(id: string, payload: UpdateTaskPayload): Promise<Task> {
    const res = await apiClient.put<ApiResponse<Task>>(`/tasks/${id}`, payload);
    return res.data;
  },

  async updateStatus(id: string, status: string): Promise<Task> {
    const res = await apiClient.put<ApiResponse<Task>>(`/tasks/${id}/status`, { status });
    return res.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/tasks/${id}`);
  },
};
