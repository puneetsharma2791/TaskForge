# Engineering Coding Assessment

## Background

You've been given access to **TaskForge**, an internal task management system built with a Go backend and React frontend. The application was prototyped by a team member who has since moved on. Your team has been asked to take ownership of this codebase.

TaskForge uses an event-sourced architecture:
- **Aggregates** enforce business rules and emit domain events
- **Events** are the source of truth for state changes
- **Projectors** consume events to build read models (query-optimized views)
- **Commands** orchestrate aggregate operations
- **Queries** read from projected read models
- The system is designed for multi-tenant operation

The frontend is a React + TypeScript SPA with:
- Custom hooks for data fetching
- Context-based auth state
- DaisyUI component library
- Standard service/API client layer

## Your Task

### Part 1: Code Analysis (Primary)

Review the full codebase and produce a written analysis covering:

1. **Security vulnerabilities** -- Identify and rank by severity. Explain the potential impact of each.

2. **Business logic errors** -- Identify cases where the code doesn't enforce the intended business rules. The intended task lifecycle is:
   ```
   draft --> open --> in_progress --> completed
                 \                \--> cancelled
                  \--> cancelled
   ```
   (Tasks can be cancelled from `open` or `in_progress`, but not from `completed`.)

3. **Reliability & concurrency issues** -- Identify race conditions, memory leaks, or failure modes.

4. **Code quality concerns** -- Identify issues with error handling, type safety, documentation accuracy, and frontend patterns.

For each issue, provide:
- Location (file and approximate area)
- Description of the problem
- Potential impact
- Suggested fix (brief)

### Part 2: Feature Design

Design a **Task Comments** feature. Users should be able to:
- Add comments to any task
- Edit their own comments
- Delete their own comments
- View all comments on a task (newest first)

Produce a design document covering:
1. **Backend**: How comments fit into the event-sourced architecture (events, aggregate changes, projector, queries, API endpoints)
2. **Frontend**: Component structure, data fetching approach, UI/UX considerations
3. **Testing**: What tests you would write and why
4. **Migration**: How you'd roll this out without disrupting existing functionality

You are encouraged to use AI-assisted coding tools during this assessment. We're evaluating your ability to analyze, reason about, and design within an event-sourced system -- not your ability to memorize syntax.

## Deliverables

1. A written analysis document (Part 1)
2. A design document (Part 2)
3. (Optional bonus) A partial implementation of the comments feature

## Time

This assessment is designed to take approximately 4-6 hours. Focus on depth of analysis over breadth -- a thorough treatment of fewer issues is valued more than a shallow list.

## Getting Started

```bash
# Backend
cd taskforge/backend
go build ./cmd/api
./api   # Runs on :8080

# Frontend
cd taskforge/frontend
npm install
npm run dev   # Runs on :5173, proxies /api to :8080
```
