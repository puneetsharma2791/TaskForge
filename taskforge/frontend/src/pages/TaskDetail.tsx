import { useParams, useNavigate } from 'react-router-dom';
import { IconArrowLeft, IconTrash } from '@tabler/icons-react';
import { tasksService } from '@/services/tasks.service';
import { useApi } from '@/hooks/useApi';
import { useMutation } from '@/hooks/useMutation';
import StatusBadge from '@/components/StatusBadge';
import TaskForm from '@/components/TaskForm';
import type { CreateTaskPayload, TaskStatus } from '@/types';

const STATUS_TRANSITIONS: TaskStatus[] = ['draft', 'open', 'in_progress', 'completed', 'cancelled'];

export default function TaskDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const isNew = id === 'new';

  const { data: task, loading } = useApi(() =>
    isNew ? Promise.resolve(null) : tasksService.getById(id!)
  );

  const { mutate: updateTask } = useMutation(
    (payload: CreateTaskPayload) => tasksService.update(id!, payload),
    { onSuccess: () => navigate('/tasks') }
  );

  const { mutate: updateStatus } = useMutation(
    (status: TaskStatus) => tasksService.updateStatus(id!, status),
    { onSuccess: () => navigate('/tasks') }
  );

  const { mutate: createTask } = useMutation(
    (payload: CreateTaskPayload) => tasksService.create(payload),
    { onSuccess: () => navigate('/tasks') }
  );

  const { mutate: deleteTask } = useMutation(
    () => tasksService.delete(id!),
    { onSuccess: () => navigate('/tasks') }
  );

  const handleSubmit = async (data: CreateTaskPayload) => {
    if (isNew) {
      await createTask(data);
    } else {
      await updateTask(data);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center p-12">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="p-6 max-w-2xl mx-auto">
      <button className="btn btn-ghost btn-sm mb-4" onClick={() => navigate('/tasks')}>
        <IconArrowLeft size={16} />
        Back to Tasks
      </button>

      {!isNew && task && (
        <div className="mb-6">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold">{task.title}</h1>
            <div className="flex items-center gap-2">
              <StatusBadge status={task.status} />
              <button
                className="btn btn-error btn-sm btn-outline"
                onClick={() => deleteTask(undefined as never)}
              >
                <IconTrash size={16} />
              </button>
            </div>
          </div>

          {/* Render description content */}
          {task.description && (
            <div
              className="prose mt-4"
              dangerouslySetInnerHTML={{ __html: task.description }}
            />
          )}

          {/* Status transitions */}
          <div className="flex gap-2 mt-4">
            {STATUS_TRANSITIONS.map((status) => (
              <button
                key={status}
                className="btn btn-outline btn-xs"
                onClick={() => updateStatus(status)}
              >
                {status.replace('_', ' ')}
              </button>
            ))}
          </div>
        </div>
      )}

      <div className="divider">{isNew ? 'Create Task' : 'Edit Task'}</div>

      <TaskForm
        initialData={task || undefined}
        onSubmit={handleSubmit}
        isEdit={!isNew}
      />
    </div>
  );
}
