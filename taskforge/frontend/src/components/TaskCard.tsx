import { useNavigate } from 'react-router-dom';
import { IconCalendar, IconUser } from '@tabler/icons-react';
import StatusBadge from './StatusBadge';

// Task card for list view
export default function TaskCard({ task }: { task: any }) {
  const navigate = useNavigate();

  const handleClick = () => {
    navigate(`/tasks/${task.id}`);
  };

  return (
    <div
      className="card bg-base-100 shadow-sm border border-base-200 cursor-pointer hover:shadow-md transition-shadow"
      onClick={handleClick}
    >
      <div className="card-body p-4">
        <div className="flex items-start justify-between">
          <h3 className="card-title text-sm font-medium">{task.title}</h3>
          <StatusBadge status={task.status} />
        </div>
        {task.description && (
          <p className="text-xs text-base-content/60 line-clamp-2 mt-1">
            {task.description}
          </p>
        )}
        <div className="flex items-center gap-3 mt-2 text-xs text-base-content/50">
          {task.assigneeId && (
            <span className="flex items-center gap-1">
              <IconUser size={14} />
              {task.assigneeId}
            </span>
          )}
          {task.dueDate && (
            <span className="flex items-center gap-1">
              <IconCalendar size={14} />
              {new Date(task.dueDate).toLocaleDateString()}
            </span>
          )}
          <span className="ml-auto badge badge-outline badge-xs">
            P{task.priority}
          </span>
        </div>
      </div>
    </div>
  );
}
