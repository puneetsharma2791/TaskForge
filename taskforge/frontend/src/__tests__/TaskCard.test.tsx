import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import TaskCard from '../components/TaskCard';

const mockTask = {
  id: '1',
  title: 'Fix login bug',
  status: 'open',
  priority: 2,
};

describe('TaskCard', () => {
  it('renders task title', () => {
    render(
      <BrowserRouter>
        <TaskCard task={mockTask} />
      </BrowserRouter>
    );

    expect(screen.getByText('Fix login bug')).toBeDefined();
  });

  it('renders status badge', () => {
    render(
      <BrowserRouter>
        <TaskCard task={mockTask} />
      </BrowserRouter>
    );

    expect(screen.getByText('Open')).toBeDefined();
  });

  it('renders priority', () => {
    render(
      <BrowserRouter>
        <TaskCard task={mockTask} />
      </BrowserRouter>
    );

    expect(screen.getByText('P2')).toBeDefined();
  });
});
