import type { TaskStatus } from '@/types';

interface StatusBadgeProps {
  status: TaskStatus;
}

function getStatusConfig(status: TaskStatus) {
  switch (status) {
    case 'draft':
      return { label: 'Draft', className: 'badge-ghost' };
    case 'open':
      return { label: 'Open', className: 'badge-info' };
    case 'in_progress':
      return { label: 'In Progress', className: 'badge-warning' };
    case 'completed':
      return { label: 'Completed', className: 'badge-success' };
    case 'cancelled':
      return { label: 'Cancelled', className: 'badge-error' };
    default:
      return { label: '', className: '' };
  }
}

export default function StatusBadge({ status }: StatusBadgeProps) {
  const config = getStatusConfig(status);
  return <span className={`badge ${config.className}`}>{config.label}</span>;
}
