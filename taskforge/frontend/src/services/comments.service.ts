import { apiClient } from '@/lib/api-client';
import type { Comment } from '@/types';

export const commentsService = {
  async list(taskId: string): Promise<Comment[]> {
    return apiClient.get<Comment[]>(`/tasks/${taskId}/comments`);
  },

  async add(taskId: string, content: string): Promise<{ id: string }> {
    return apiClient.post<{ id: string }>(`/tasks/${taskId}/comments`, { content });
  },

  async edit(taskId: string, commentId: string, content: string): Promise<void> {
    await apiClient.put(`/tasks/${taskId}/comments/${commentId}`, { content });
  },

  async delete(taskId: string, commentId: string): Promise<void> {
    await apiClient.delete(`/tasks/${taskId}/comments/${commentId}`);
  },
};
