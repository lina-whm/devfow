'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import type {
  Task,
  CreateTaskDTO,
  UpdateTaskDTO,
  MoveTaskDTO,
  PaginatedResponse,
  APIError,
} from '@/types/api';

interface TaskFilters {
  status?: string;
  priority?: string;
  type?: string;
  assigneeId?: string;
  sprintId?: string;
  search?: string;
  page?: number;
  pageSize?: number;
}

export function useTasks(projectId: string | undefined, filters: TaskFilters = {}) {
  return useQuery<PaginatedResponse<Task>>({
    queryKey: ['tasks', projectId, filters],
    queryFn: () => {
      const params = new URLSearchParams();
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== '') {
          params.set(key, String(value));
        }
      });
      const qs = params.toString();
      return apiClient.get(`/projects/${projectId}/tasks${qs ? `?${qs}` : ''}`);
    },
    enabled: !!projectId,
  });
}

export function useTask(taskId: string | undefined) {
  return useQuery<Task>({
    queryKey: ['task', taskId],
    queryFn: () => apiClient.get(`/tasks/${taskId}`),
    enabled: !!taskId,
  });
}

export function useCreateTask() {
  const queryClient = useQueryClient();

  return useMutation<Task, APIError, CreateTaskDTO>({
    mutationFn: (data) => apiClient.post('/tasks', data),
    onSuccess: (task) => {
      queryClient.invalidateQueries({ queryKey: ['tasks', task.projectId] });
      queryClient.invalidateQueries({ queryKey: ['board', task.projectId] });
    },
  });
}

export function useUpdateTask() {
  const queryClient = useQueryClient();

  return useMutation<Task, APIError, { taskId: string; data: UpdateTaskDTO }>({
    mutationFn: ({ taskId, data }) => apiClient.patch(`/tasks/${taskId}`, data),
    onSuccess: (task) => {
      queryClient.invalidateQueries({ queryKey: ['tasks', task.projectId] });
      queryClient.invalidateQueries({ queryKey: ['task', task.id] });
      queryClient.invalidateQueries({ queryKey: ['board', task.projectId] });
    },
  });
}

export function useMoveTask() {
  const queryClient = useQueryClient();

  return useMutation<Task, APIError, MoveTaskDTO>({
    mutationFn: (data) => apiClient.patch(`/tasks/${data.taskId}/move`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['board'] });
    },
  });
}

export function useDeleteTask() {
  const queryClient = useQueryClient();

  return useMutation<void, APIError, string>({
    mutationFn: (taskId) => apiClient.delete(`/tasks/${taskId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      queryClient.invalidateQueries({ queryKey: ['board'] });
    },
  });
}
