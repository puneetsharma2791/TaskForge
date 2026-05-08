import { useState } from 'react';
import { IconPlus, IconTrash } from '@tabler/icons-react';
import { projectsService } from '@/services/projects.service';
import { useApi } from '@/hooks/useApi';
import { useMutation } from '@/hooks/useMutation';
import type { CreateProjectPayload } from '@/types';

export default function Projects() {
  const { data: projects, loading, error, refetch } = useApi(() => projectsService.list());
  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');

  const { mutate: createProject } = useMutation(
    (payload: CreateProjectPayload) => projectsService.create(payload),
    {
      onSuccess: () => {
        setShowForm(false);
        setName('');
        setDescription('');
        refetch();
      },
    }
  );

  const { mutate: deleteProject } = useMutation(
    (id: string) => projectsService.delete(id),
    { onSuccess: () => refetch() }
  );

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    await createProject({ name, description });
  };

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
        <span>Failed to load projects: {error}</span>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Projects</h1>
        <button
          className="btn btn-primary btn-sm"
          onClick={() => setShowForm(!showForm)}
        >
          <IconPlus size={16} />
          New Project
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleCreate} className="card bg-base-100 shadow-sm border border-base-200 p-4 mb-4">
          <div className="form-control">
            <label className="label">
              <span className="label-text">Project Name</span>
            </label>
            <input
              type="text"
              className="input input-bordered"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>
          <div className="form-control mt-2">
            <label className="label">
              <span className="label-text">Description</span>
            </label>
            <textarea
              className="textarea textarea-bordered"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </div>
          <div className="flex justify-end gap-2 mt-4">
            <button type="button" className="btn btn-ghost btn-sm" onClick={() => setShowForm(false)}>
              Cancel
            </button>
            <button type="submit" className="btn btn-primary btn-sm">
              Create
            </button>
          </div>
        </form>
      )}

      <div className="grid gap-3">
        {projects?.map((project) => (
          <div
            key={project.id}
            className="card bg-base-100 shadow-sm border border-base-200"
          >
            <div className="card-body p-4 flex-row items-center justify-between">
              <div>
                <h3 className="font-medium">{project.name}</h3>
                {project.description && (
                  <p className="text-sm text-base-content/60">{project.description}</p>
                )}
                <span className="badge badge-sm mt-1">{project.status}</span>
              </div>
              <button
                className="btn btn-ghost btn-sm text-error"
                onClick={() => deleteProject(project.id)}
              >
                <IconTrash size={16} />
              </button>
            </div>
          </div>
        ))}
        {projects?.length === 0 && (
          <div className="text-center py-12 text-base-content/50">
            No projects yet
          </div>
        )}
      </div>
    </div>
  );
}
