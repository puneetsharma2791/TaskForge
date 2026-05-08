import { IconChecklist, IconFolder, IconClock } from '@tabler/icons-react';
import { tasksService } from '@/services/tasks.service';
import { projectsService } from '@/services/projects.service';
import { useApi } from '@/hooks/useApi';

export default function Dashboard() {
  const { data: tasks } = useApi(() => tasksService.list());
  const { data: projects } = useApi(() => projectsService.list());

  const openTasks = tasks?.filter((t) => t.status === 'open' || t.status === 'in_progress').length || 0;
  const completedTasks = tasks?.filter((t) => t.status === 'completed').length || 0;

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-6">Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
        <div className="stat bg-base-100 shadow-sm rounded-box border border-base-200">
          <div className="stat-figure text-primary">
            <IconChecklist size={32} />
          </div>
          <div className="stat-title">Open Tasks</div>
          <div className="stat-value text-primary">{openTasks}</div>
        </div>

        <div className="stat bg-base-100 shadow-sm rounded-box border border-base-200">
          <div className="stat-figure text-success">
            <IconClock size={32} />
          </div>
          <div className="stat-title">Completed</div>
          <div className="stat-value text-success">{completedTasks}</div>
        </div>

        <div className="stat bg-base-100 shadow-sm rounded-box border border-base-200">
          <div className="stat-figure text-info">
            <IconFolder size={32} />
          </div>
          <div className="stat-title">Projects</div>
          <div className="stat-value text-info">{projects?.length || 0}</div>
        </div>
      </div>
    </div>
  );
}
