# TaskForge

An internal task management system built with a Go backend (event-sourced architecture) and a React + TypeScript frontend.

## Architecture

**Backend** — Go 1.22, event-sourced with CQRS pattern:
- Aggregates enforce business rules and emit domain events
- Events are the source of truth for all state changes
- Projectors consume events to build read-optimized views
- Commands orchestrate aggregate operations
- In-memory event store and repositories (no external database required)

**Frontend** — React 19, TypeScript, Vite:
- DaisyUI + Tailwind CSS v4 for styling
- Custom `useApi` / `useMutation` hooks for data fetching
- Context-based authentication state
- React Router for navigation

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/) (with npm)

## Project Structure

```
taskforge/
├── backend/
│   ├── cmd/api/main.go              # Entry point
│   └── internal/
│       ├── config/                   # App configuration
│       ├── domain/
│       │   ├── aggregates/           # Task & Project aggregates
│       │   ├── commands/             # Create, Assign, Complete task
│       │   ├── events/               # Domain event definitions
│       │   ├── projectors/           # Read model builders
│       │   ├── queries/              # Query handlers
│       │   ├── repositories/         # In-memory repositories
│       │   └── types/                # View models (TaskView, CommentView, etc.)
│       ├── storage/                  # In-memory event store
│       └── web/api/v1/              # HTTP handlers & routes
└── frontend/
    └── src/
        ├── components/               # TaskCard, StatusBadge, CommentSection, etc.
        ├── context/                  # AuthContext
        ├── hooks/                    # useApi, useMutation
        ├── lib/                      # API client
        ├── pages/                    # Dashboard, Tasks, TaskDetail, Projects, Login
        ├── services/                 # tasks, projects, comments, auth services
        └── types/                    # TypeScript interfaces
```

## Setup & Run

### 1. Start the Backend

```bash
cd taskforge/backend
go mod download
go run cmd/api/main.go
```

The API server starts on **http://localhost:8080**.

You can change the port with the `PORT` environment variable:

```bash
PORT=9090 go run cmd/api/main.go
```

### 2. Start the Frontend

In a separate terminal:

```bash
cd taskforge/frontend
npm install
npm run dev
```

The dev server starts on **http://localhost:5173** (default Vite port).

### 3. Use the App

1. Open **http://localhost:5173** in your browser
2. Log in with any email and password (mock authentication)
3. Create a project, then create tasks within it
4. Click into a task to view details, edit, or add comments

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | Mock login (returns token + user) |
| GET | `/api/v1/tasks` | List tasks (optional `?project_id=`, `?status=`, `?q=`) |
| POST | `/api/v1/tasks` | Create a task |
| GET | `/api/v1/tasks/{id}` | Get task by ID |
| PUT | `/api/v1/tasks/{id}` | Update a task |
| DELETE | `/api/v1/tasks/{id}` | Delete a task |
| POST | `/api/v1/tasks/{id}/assign` | Assign a task |
| POST | `/api/v1/tasks/{id}/complete` | Complete a task |
| GET | `/api/v1/tasks/{id}/comments` | List comments on a task |
| POST | `/api/v1/tasks/{id}/comments` | Add a comment |
| PUT | `/api/v1/tasks/{id}/comments/{commentId}` | Edit a comment |
| DELETE | `/api/v1/tasks/{id}/comments/{commentId}` | Delete a comment |
| GET | `/api/v1/projects` | List projects |
| POST | `/api/v1/projects` | Create a project |
| GET | `/api/v1/projects/{id}` | Get project by ID |
| DELETE | `/api/v1/projects/{id}` | Delete a project |
| POST | `/api/v1/projects/{id}/members` | Add a member to a project |

**Headers:** Requests should include `X-User-ID` and `Content-Type: application/json`.

## Running Tests

```bash
# Frontend tests
cd taskforge/frontend
npm test
```

## Notes

- All data is stored in-memory — restarting the backend clears everything.
- Authentication is mocked — any email/password combination works.
- There is no tenant isolation enforcement in the current implementation.
