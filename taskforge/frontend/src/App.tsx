import { Routes, Route, Navigate, Link, useLocation } from 'react-router-dom';
import { IconChecklist, IconFolder, IconHome, IconLogout } from '@tabler/icons-react';
import { useAuth } from '@/context/AuthContext';
import Login from '@/pages/Login';
import Dashboard from '@/pages/Dashboard';
import Tasks from '@/pages/Tasks';
import TaskDetail from '@/pages/TaskDetail';
import Projects from '@/pages/Projects';

function Layout({ children }: { children: React.ReactNode }) {
  const { logout, user } = useAuth();
  const location = useLocation();

  const navItems = [
    { path: '/dashboard', label: 'Dashboard', icon: IconHome },
    { path: '/tasks', label: 'Tasks', icon: IconChecklist },
    { path: '/projects', label: 'Projects', icon: IconFolder },
  ];

  return (
    <div className="min-h-screen flex">
      <aside className="w-56 bg-base-200 border-r border-base-300 flex flex-col">
        <div className="p-4 font-bold text-lg border-b border-base-300">TaskForge</div>
        <nav className="flex-1 p-2">
          {navItems.map((item) => (
            <Link
              key={item.path}
              to={item.path}
              className={`flex items-center gap-2 px-3 py-2 rounded-lg text-sm mb-1 ${
                location.pathname.startsWith(item.path)
                  ? 'bg-primary text-primary-content'
                  : 'hover:bg-base-300'
              }`}
            >
              <item.icon size={18} />
              {item.label}
            </Link>
          ))}
        </nav>
        <div className="p-3 border-t border-base-300">
          <div className="text-xs text-base-content/50 mb-2">{user?.email}</div>
          <button className="btn btn-ghost btn-xs w-full justify-start" onClick={logout}>
            <IconLogout size={14} />
            Sign Out
          </button>
        </div>
      </aside>
      <main className="flex-1 bg-base-100 overflow-auto">{children}</main>
    </div>
  );
}

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth();
  if (!isAuthenticated) return <Navigate to="/login" replace />;
  return <Layout>{children}</Layout>;
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route
        path="/dashboard"
        element={
          <ProtectedRoute>
            <Dashboard />
          </ProtectedRoute>
        }
      />
      <Route
        path="/tasks"
        element={
          <ProtectedRoute>
            <Tasks />
          </ProtectedRoute>
        }
      />
      <Route path="/tasks/:id" element={<TaskDetail />} />
      <Route
        path="/projects"
        element={
          <ProtectedRoute>
            <Projects />
          </ProtectedRoute>
        }
      />
      <Route path="/" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  );
}
