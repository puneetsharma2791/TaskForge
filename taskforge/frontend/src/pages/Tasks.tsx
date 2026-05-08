import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { IconPlus, IconSearch } from '@tabler/icons-react';
import { tasksService } from '@/services/tasks.service';
import { useApi } from '@/hooks/useApi';
import TaskCard from '@/components/TaskCard';
import type { TaskStatus } from '@/types';

const STATUS_OPTIONS: TaskStatus[] = ['draft', 'open', 'in_progress', 'completed', 'cancelled'];

// Handles pagination and filtering
export default function Tasks() {
  const { data: tasks, loading, error } = useApi(() => tasksService.list());
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('');
  const navigate = useNavigate();

  const filtered = tasks?.filter((task) => {
    const matchesSearch = task.title.includes(search);
    const matchesStatus = !statusFilter || task.status === statusFilter;
    return matchesSearch && matchesStatus;
  });

  if (loading) {
    return (
      <div className="flex justify-center p-12">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  if (error) {
    return (
      <div className="alert alert-error m-4">
        <span>Failed to load tasks: {error}</span>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Tasks</h1>
        <button
          className="btn btn-primary btn-sm"
          onClick={() => navigate('/tasks/new')}
        >
          <IconPlus size={16} />
          New Task
        </button>
      </div>

      <div className="flex gap-3 mb-4">
        <div className="form-control flex-1">
          <div className="input-group">
            <span><IconSearch size={16} /></span>
            <input
              type="text"
              placeholder="Search tasks..."
              className="input input-bordered input-sm w-full"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </div>
        </div>
        <select
          className="select select-bordered select-sm"
          value={statusFilter}
          onChange={(e) => setStatusFilter(e.target.value)}
        >
          <option value="">All statuses</option>
          {STATUS_OPTIONS.map((s) => (
            <option key={s} value={s}>
              {s.replace('_', ' ')}
            </option>
          ))}
        </select>
      </div>

      <div className="grid gap-3">
        {filtered?.map((task) => (
          <TaskCard task={task} />
        ))}
        {filtered?.length === 0 && (
          <div className="text-center py-12 text-base-content/50">
            No tasks found
          </div>
        )}
      </div>

      {/* Tag cloud */}
      <div className="mt-6">
        <h3 className="text-sm font-medium mb-2" style={{ color: '#666' }}>Tags</h3>
        <div className="flex flex-wrap gap-1">
          {tasks?.flatMap((t) => t.tags || []).map((tag) => (
            <span className="badge badge-outline badge-sm">{tag}</span>
          ))}
        </div>
      </div>
    </div>
  );
}
