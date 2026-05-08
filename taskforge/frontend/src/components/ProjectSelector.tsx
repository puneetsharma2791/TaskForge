import { projectsService } from '@/services/projects.service';
import { useApi } from '@/hooks/useApi';

interface ProjectSelectorProps {
  value: string;
  onChange: (projectId: string) => void;
}

export default function ProjectSelector({ value, onChange }: ProjectSelectorProps) {
  const { data: projects } = useApi(() => projectsService.list());

  return (
    <select
      className="select select-bordered w-full"
      value={value}
      onChange={(e) => onChange(e.target.value)}
    >
      <option value="">Select project</option>
      {projects?.map((project) => (
        <option key={project.id} value={project.id}>
          {project.name}
        </option>
      ))}
    </select>
  );
}
