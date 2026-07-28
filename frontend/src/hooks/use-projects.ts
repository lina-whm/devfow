'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import type { Project, APIError } from '@/types/api';

export function useProjects(orgId: string | undefined) {
  return useQuery<Project[]>({
    queryKey: ['projects', orgId],
    queryFn: () => apiClient.get(`/organizations/${orgId}/projects`),
    enabled: !!orgId,
  });
}

export function useProject(projectId: string | undefined) {
  return useQuery<Project>({
    queryKey: ['project', projectId],
    queryFn: () => apiClient.get(`/projects/${projectId}`),
    enabled: !!projectId,
  });
}

export function useCreateProject() {
  const queryClient = useQueryClient();
  return useMutation<Project, APIError, { orgId: string; name: string; key: string; description?: string }>({
    mutationFn: ({ orgId, ...data }) => apiClient.post(`/organizations/${orgId}/projects`, data),
    onSuccess: (project) => {
      queryClient.invalidateQueries({ queryKey: ['projects', project.organizationId] });
    },
  });
}

export function useUpdateProject() {
  const queryClient = useQueryClient();
  return useMutation<Project, APIError, { projectId: string; name?: string; description?: string }>({
    mutationFn: ({ projectId, ...data }) => apiClient.patch(`/projects/${projectId}`, data),
    onSuccess: (project) => {
      queryClient.invalidateQueries({ queryKey: ['project', project.id] });
      queryClient.invalidateQueries({ queryKey: ['projects'] });
    },
  });
}

export function useDeleteProject() {
  const queryClient = useQueryClient();
  return useMutation<void, APIError, string>({
    mutationFn: (projectId) => apiClient.delete(`/projects/${projectId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] });
    },
  });
}
