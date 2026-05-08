# TaskForge - Complete Project Analysis

## 1. What Is This Project?

**TaskForge** is an internal task management system (like a simplified Jira/Trello). It has:
- A **Go backend** using event-sourced architecture
- A **React + TypeScript frontend** SPA with DaisyUI styling

It's a coding assessment project — a deliberately flawed prototype that candidates are asked to analyze for bugs, security issues, and then design a new feature.

---

## 2. High-Level Architecture

```
Frontend (React + Vite, port 5173)
    ↓ HTTP calls to /api/v1/*
Backend (Go + gorilla/mux, port 8080)
    ↓
Event-Sourced Domain Layer (in-memory)
```

The backend uses **Event Sourcing + CQRS**:
- **Commands** (write side): Modify state by emitting events
- **Queries** (read side): Read from projected "views" built from events
- **Aggregates**: Enforce business rules
- **Events**: Immutable records of what happened
- **Projectors**: Build denormalized read models from events

---

## 3. Backend — Step by Step

### 3.1 Entry Point: `backend/cmd/api/main.go`

This is the application bootstrap:

1. **Loads config** via `config.Load()` — reads env vars with hardcoded defaults
2. **Creates an in-memory EventStore** — a simple `map[string][]Event`
3. **Creates repositories** — `InMemoryTaskRepository` and `InMemoryProjectRepository`, both backed by the event store
4. **Creates a TaskProjector** — builds read-optimized views from events
5. **Creates HTTP handlers** — `TasksHandler` and `ProjectsHandler`
6. **Sets up routing** with gorilla/mux, adds CORS middleware (allows `*` origin), and starts the server on the configured port

### 3.2 Config: `backend/internal/config/config.go`

A simple struct that reads from environment variables with fallbacks:
- `PORT` defaults to `"8080"`
- `SECRET_KEY` defaults to `"taskforge-default-secret-key-2024"` (hardcoded secret — intentional security flaw)
- `SkipTLSVerify` and `DevMode` are hardcoded to `true` (never read from env — intentional flaw)

### 3.3 Events: `backend/internal/domain/events/events.go`

Defines all domain event types. Each implements the `Event` interface (just `EventType() string`):

| Event | Fields | Purpose |
|---|---|---|
| `TaskCreated` | TaskID, ProjectID, Title, Description, CreatedBy, Priority, OccurredAt | A new task was created |
| `TaskOpened` | TaskID, OccurredAt | Task moved to "open" |
| `TaskStarted` | TaskID, AssigneeID, OccurredAt | Task moved to "in_progress" |
| `TaskCompleted` | TaskID, OccurredAt | Task marked done |
| `TaskCancelled` | TaskID, Reason, OccurredAt | Task was cancelled |
| `TaskReassigned` | TaskID, NewAssigneeID, OccurredAt | Assignee changed |
| `TaskPriorityUpdated` | TaskID, Priority, OccurredAt | Priority changed |
| `ProjectCreated` | ProjectID, TenantID, Name, CreatedBy, OccurredAt | New project |
| `ProjectMemberAdded` | ProjectID, UserID, Role, OccurredAt | Member added |
| `ProjectDeleted` | ProjectID, OccurredAt | Project deleted |

### 3.4 Task Aggregate: `backend/internal/domain/aggregates/task.go`

This is the **core business logic** for tasks. The aggregate:

- Holds private state: `id`, `projectID`, `title`, `description`, `status`, `assigneeID`, `priority`, `createdBy`, `events[]`, `version`
- Starts in `StatusDraft` when created via `NewTask()`

**State machine (intended lifecycle):**
```
draft → open → in_progress → completed
                  \                \→ cancelled
                   \→ cancelled
```

**Methods and their rules:**

| Method | Guard | Emits |
|---|---|---|
| `Create()` | None (no validation!) | `TaskCreated` |
| `Open()` | Must be in `draft` | `TaskOpened` |
| `Start(assigneeID)` | Must be in `open` | `TaskStarted` |
| `Complete()` | **No guard!** (bug — can complete from any status) | `TaskCompleted` |
| `Cancel(reason)` | Only checks if already cancelled (bug — can cancel `completed` tasks) | `TaskCancelled` |
| `Reassign(newAssigneeID)` | **No guard!** (can reassign in any state) | `TaskReassigned` |
| `UpdatePriority(priority)` | **No guard!** | `TaskPriorityUpdated` |

**The `apply()` method** is the event-sourcing heart — it:
1. Appends the event to the pending events list
2. Increments version
3. Updates internal state based on event type (a switch statement)

**`LoadFromEvents()`** replays historical events to rebuild state, then clears the pending list (since those are already persisted).

### 3.5 Project Aggregate: `backend/internal/domain/aggregates/project.go`

Similar pattern to Task:
- Holds `id`, `tenantID`, `name`, `createdBy`, `members` (map of userID→role), `deleted` flag
- `Create()` validates name is non-empty, emits `ProjectCreated`, auto-adds creator as `"owner"`
- `AddMember()` checks project isn't deleted, emits `ProjectMemberAdded` (but doesn't validate role values — any string accepted)
- `Delete()` emits `ProjectDeleted`, sets `deleted = true`
- `Members()` returns a copy of the members map (good practice — prevents external mutation)

### 3.6 Event Store: `backend/internal/storage/event_store.go`

Extremely simple in-memory persistence:
- `events map[string][]Event` — keyed by aggregate ID
- `Save()` — appends events to the list for that aggregate ID
- `Load()` — returns events for a given aggregate ID
- `AllEvents()` — returns all events (used for rebuilding projections)

**No concurrency protection** (no mutex) — this is a race condition bug.

### 3.7 Repositories: `backend/internal/domain/repositories/`

**TaskRepository interface:**
```go
Save(task *Task) error
FindByID(id string) (*Task, error)
Delete(id string) error
```

**InMemoryTaskRepository** implementation:
- Maintains a `snapshots` map as a cache
- `Save()` persists pending events to the event store, then caches the aggregate
- `FindByID()` checks cache first, falls back to loading from event store and replaying events
- `Delete()` just removes from the snapshot cache (doesn't record a delete event — bug)

**InMemoryProjectRepository** — similar pattern, but `FindByID()` only checks snapshots (doesn't replay from event store — inconsistent).

### 3.8 Commands: `backend/internal/domain/commands/`

Commands orchestrate write operations:

- **`CreateTask`**: Generates UUID, creates new Task aggregate, calls `Create()`, saves to repo. Note: it does NOT call `Open()` — tasks stay in `draft` forever unless explicitly opened (and there's no API endpoint to open them!)
- **`CompleteTask`**: Loads task by ID, calls `Complete()`, saves. Has a `// TODO: check authorization` comment — no auth check.
- **`AssignTask`**: Loads task, calls `Reassign()`, saves. Comment says "save without version check" — no optimistic concurrency control.

### 3.9 Queries: `backend/internal/domain/queries/`

Queries read from the projector (read model):

- **`FindTask`**: Gets a single task view by ID. Has `TenantID` field but `// TODO: tenant check` — no multi-tenant isolation (security bug).
- **`ListTasks`**: Gets all tasks, filters by project, status, and search (case-insensitive substring match on title/description).

### 3.10 Task Projector: `backend/internal/domain/projectors/task_projector.go`

Builds read models (`TaskView` structs) from events:
- Maintains a `views` map and an event `buffer` (buffer only grows, never consumed — memory leak)
- `Project()` processes each event type and updates the corresponding `TaskView`
- Note: it doesn't handle `TaskCancelled` events — cancelled tasks won't update in the read model (bug)
- `GetAll()`, `GetView()`, `GetByProject()` are read methods

### 3.11 HTTP Handlers: `backend/internal/web/api/v1/`

**Routes** (`routes.go`):
```
GET    /api/v1/tasks           → List tasks
POST   /api/v1/tasks           → Create task
GET    /api/v1/tasks/{id}      → Get single task
PUT    /api/v1/tasks/{id}      → Update task
POST   /api/v1/tasks/{id}/complete → Complete task
POST   /api/v1/tasks/{id}/assign   → Assign task

GET    /api/v1/projects        → List projects
POST   /api/v1/projects        → Create project
GET    /api/v1/projects/{id}   → Get project
DELETE /api/v1/projects/{id}   → Delete project
POST   /api/v1/projects/{id}/members → Add member
```

**Tasks Handler** (`tasks_handler.go`):
- `Create`: Gets user ID from `X-User-ID` header (trusts client — no auth), creates task, then manually projects events to the read model
- `Update`: Only handles priority updates (ignores title/description changes)
- `Assign`: Delegates to `AssignTask` command but **doesn't project events afterward** — the read model goes stale
- `Complete`: Same issue — doesn't project events after completing
- `writeError()`: **Includes full stack traces** in error responses (`debug.Stack()`) — security information leak

**Projects Handler** (`projects_handler.go`):
- `Delete`: Just removes from snapshot cache — doesn't use the aggregate's `Delete()` method, so no event is recorded
- `List`: Requires `tenant_id` query param
- `AddMember`: Doesn't validate role values

---

## 4. Frontend — Step by Step

### 4.1 Entry & Routing: `src/main.tsx` and `src/App.tsx`

**`main.tsx`** bootstraps the app:
- Wraps in `React.StrictMode`, `BrowserRouter` (react-router), and `AuthProvider` (context)

**`App.tsx`** defines routes:
- `/login` → Login page
- `/dashboard` → Dashboard (protected)
- `/tasks` → Task list (protected)
- `/tasks/:id` → Task detail (**NOT protected!** — bug, anyone can access)
- `/projects` → Projects (protected)
- `/` → Redirects to `/dashboard`

**`Layout`** component: Sidebar with navigation links + sign-out button. Shows user email.

**`ProtectedRoute`**: Checks `isAuthenticated` from auth context, redirects to `/login` if not.

### 4.2 Auth Context: `src/context/AuthContext.tsx`

- Stores `user` and `token` in React state
- Reads initial token from `localStorage`
- `isAuthenticated` is simply `!!token` (just checks token exists, never validates it)
- `login()` calls the auth service, stores token + user
- `logout()` clears both
- Has a `useEffect` with a `// TODO: fetch user profile on mount` — so after refresh, `isAuthenticated` is true but `user` is null

### 4.3 API Client: `src/lib/api-client.ts`

A wrapper around `fetch`:
- Hardcoded `BASE_URL = 'http://localhost:8080/api/v1'` (doesn't use the Vite proxy)
- `get`, `post`, `put`, `delete` methods
- Adds `Authorization: Bearer <token>` header if token exists
- **Never checks response status** — `res.json()` is called even on error responses (bug — errors silently return malformed data)

### 4.4 Custom Hooks

**`useApi`** (`src/hooks/useApi.ts`):
- Generic data-fetching hook with `data`, `loading`, `error`, `refetch`
- Runs the API call once on mount (empty dependency array)
- The `apiCall` function reference is captured once and never re-evaluated if props change

**`useMutation`** (`src/hooks/useMutation.ts`):
- For write operations, provides `mutate`, `loading`, `error`, `data`
- On error, `mutate` returns `undefined` silently — callers using `await mutate(...)` may get unexpected `undefined`

### 4.5 Services: `src/services/`

**`auth.service.ts`**: Calls `POST /auth/login` (but this endpoint doesn't exist on the backend — no auth endpoint is registered!)

**`tasks.service.ts`**:
- Has a naive **cache** (`Record<string, any>`) — never invalidated, so after creating/updating tasks, stale data is served
- `list()` uses `?projectId=` query param but backend expects `?project_id=` (field name mismatch)
- `getById()` expects `ApiResponse<Task>` wrapper with `.data`, but backend returns the task directly
- `updateStatus()` calls `PUT /tasks/{id}/status` but this route doesn't exist on the backend

**`projects.service.ts`**: Same `.data` wrapper mismatch — backend returns data directly, frontend expects `{ data: ... }`

### 4.6 Types: `src/types/index.ts`

Defines TypeScript interfaces. Has self-documenting `// FLAW` comments pointing out intentional issues:
- `description`, `tags`, `dueDate` marked required but are optional in the API
- `priority` has no constraints
- `status` on Project should be a union type
- `role` on User should be a union type

### 4.7 Pages

**`Login.tsx`**:
- Standard email/password form
- **Critical security bug**: After login, navigates to `${redirect}?email=${email}&password=${password}` — leaks credentials in the URL (visible in browser history, server logs, etc.)

**`Dashboard.tsx`**:
- Shows 3 stat cards: open tasks count, completed tasks count, projects count
- Uses `useApi` to fetch data on mount

**`Tasks.tsx`**:
- Lists tasks with search (case-sensitive — `includes` not `toLowerCase`) and status filter
- "New Task" button navigates to `/tasks/new`
- Has a tag cloud section that renders tags from tasks
- Missing `key` prop on `TaskCard` in the map (React warning)

**`TaskDetail.tsx`**:
- Dual-purpose: create or edit depending on whether `id === 'new'`
- Shows task info, status badge, delete button
- **XSS vulnerability**: Uses `dangerouslySetInnerHTML` to render `task.description` — if description contains `<script>` tags or malicious HTML, it executes
- Shows all status transition buttons regardless of current status (no lifecycle enforcement in UI)

**`Projects.tsx`**:
- Lists projects with inline create form
- Delete button with no confirmation dialog

### 4.8 Components

**`TaskCard.tsx`**: Card component for task list items. Uses `task: any` type (loses type safety).

**`StatusBadge.tsx`**: Maps task status to DaisyUI badge colors. Clean component.

**`TaskForm.tsx`**: Form with react-hook-form + zod validation. Defines a schema but doesn't connect it to the form via a resolver (validation doesn't actually run). Resets form on error (bad UX — user loses their input).

**`ProjectSelector.tsx`**: Dropdown that fetches projects and lets you pick one. Each instance makes its own API call.

### 4.9 Test: `src/__tests__/TaskCard.test.tsx`

A simple vitest test with 3 cases — renders title, status badge, and priority. Uses a minimal mock task object.

---

## 5. Data Flow Example: Creating a Task

1. User fills out `TaskForm` and submits
2. `TaskDetail.handleSubmit()` calls `createTask(data)` via `useMutation`
3. `tasksService.create()` sends `POST /api/v1/tasks` with JSON body
4. Backend `TasksHandler.Create()` receives it:
   - Reads `X-User-ID` from header (frontend doesn't send this — it's always empty)
   - Creates `CreateTask` command
   - Command generates UUID, creates Task aggregate, calls `task.Create()` which emits `TaskCreated` event
   - Saves to repo (event stored, snapshot cached)
   - Handler manually loads task back, projects its events into the read model
5. Returns `{"id": "..."}` — but frontend expects `{"data": {"id": "..."}}`
6. Frontend tries to access `res.data` which is `undefined`

---

## 6. Summary of Intentional Issues

This is a **deliberately flawed codebase** for a coding assessment. The key issues include:

- **Security**: Hardcoded secrets, credentials in URL, XSS via `dangerouslySetInnerHTML`, stack traces in errors, no authentication/authorization, wildcard CORS
- **Business Logic**: `Complete()` and `Cancel()` lack proper state guards, no way to transition tasks from `draft` to `open` via API, projector doesn't handle `TaskCancelled`
- **Concurrency**: No mutex on the in-memory event store
- **Frontend-Backend Mismatch**: Response format expectations differ, query param names differ, missing API endpoints
- **Code Quality**: Untyped `any` props, stale cache, missing React keys, form validation not wired up
