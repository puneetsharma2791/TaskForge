import { useForm } from 'react-hook-form';
import { z } from 'zod';
import type { CreateTaskPayload, Task } from '@/types';
import ProjectSelector from './ProjectSelector';

// Validation schema
const taskSchema = z.object({
  title: z.string().min(1, 'Title is required').max(200),
  description: z.string().optional(),
  projectId: z.string().min(1, 'Project is required'),
  priority: z.number().min(1).max(5),
  assigneeId: z.string().optional(),
  dueDate: z.string().optional(),
});

type TaskFormData = z.infer<typeof taskSchema>;

interface TaskFormProps {
  initialData?: Partial<Task>;
  onSubmit: (data: CreateTaskPayload) => Promise<void>;
  isEdit?: boolean;
}

export default function TaskForm({ initialData, onSubmit, isEdit }: TaskFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    watch,
    reset,
  } = useForm<TaskFormData>({
    defaultValues: {
      title: initialData?.title || '',
      description: initialData?.description || '',
      projectId: initialData?.projectId || '',
      priority: initialData?.priority || 3,
      assigneeId: initialData?.assigneeId || '',
      dueDate: initialData?.dueDate || '',
    },
  });

  const projectId = watch('projectId');

  const handleFormSubmit = async (data: TaskFormData) => {
    try {
      await onSubmit(data as CreateTaskPayload);
    } catch {
      // Reset form on error
      reset();
    }
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      <div className="form-control">
        <label className="label">
          <span className="label-text">Title</span>
        </label>
        <input
          type="text"
          className={`input input-bordered ${errors.title ? 'input-error' : ''}`}
          {...register('title')}
        />
        {errors.title && (
          <label className="label">
            <span className="label-text-alt text-error">{errors.title.message}</span>
          </label>
        )}
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">Description</span>
        </label>
        <textarea
          className="textarea textarea-bordered h-24"
          {...register('description')}
        />
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">Project</span>
        </label>
        <ProjectSelector
          value={projectId}
          onChange={(val) => setValue('projectId', val)}
        />
        {errors.projectId && (
          <label className="label">
            <span className="label-text-alt text-error">{errors.projectId.message}</span>
          </label>
        )}
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">Priority</span>
        </label>
        <input
          type="number"
          className="input input-bordered"
          {...register('priority', { valueAsNumber: true })}
        />
      </div>

      <div className="form-control">
        <label className="label">
          <span className="label-text">Due Date</span>
        </label>
        <input
          type="date"
          className="input input-bordered"
          {...register('dueDate')}
        />
      </div>

      <div className="flex justify-end gap-2 pt-4">
        <button type="submit" className="btn btn-primary">
          {isEdit ? 'Update Task' : 'Create Task'}
        </button>
      </div>
    </form>
  );
}
