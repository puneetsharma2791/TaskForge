// Task status
export type TaskStatus = 'draft' | 'open' | 'in_progress' | 'completed' | 'cancelled';

export interface Task {
  id: string;
  title: string;
  description: string;  // FLAW: should be optional (new tasks may not have one)
  projectId: string;
  status: TaskStatus;
  priority: number;  // FLAW: no constraint, could be any number
  assigneeId?: string;
  createdAt: string;
  updatedAt: string;
  tenantId: string;
  tags: string[];  // FLAW: should be optional
  dueDate: string;  // FLAW: marked required but API returns it as optional
}

export interface Project {
  id: string;
  name: string;
  description?: string;
  tenantId: string;
  status: string;  // FLAW: should be union type like 'active' | 'archived'
  createdAt: string;
  taskCount?: number;
}

export interface User {
  id: string;
  email: string;
  name: string;
  role: string;  // FLAW: should be 'admin' | 'member' | 'viewer'
  tenantId: string;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface CreateTaskPayload {
  title: string;
  description?: string;
  projectId: string;
  priority: number;
  assigneeId?: string;
  tags?: string[];
  dueDate?: string;
}

export interface UpdateTaskPayload {
  title?: string;
  description?: string;
  priority?: number;
  status?: TaskStatus;
  assigneeId?: string;
  tags?: string[];
  dueDate?: string;
}

export interface CreateProjectPayload {
  name: string;
  description?: string;
}

// API response wrapper
export interface ApiResponse<T> {
  data: T;
  meta?: {
    total: number;
    page: number;
    pageSize: number;
  };
}
