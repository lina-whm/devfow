'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import type { Organization, APIError } from '@/types/api';

export function useOrganizations() {
  return useQuery<Organization[]>({
    queryKey: ['organizations'],
    queryFn: () => apiClient.get('/organizations'),
  });
}

export function useOrganization(orgId: string | undefined) {
  return useQuery<Organization>({
    queryKey: ['organization', orgId],
    queryFn: () => apiClient.get(`/organizations/${orgId}`),
    enabled: !!orgId,
  });
}

export function useCreateOrganization() {
  const queryClient = useQueryClient();
  return useMutation<Organization, APIError, { name: string; slug: string }>({
    mutationFn: (data) => apiClient.post('/organizations', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['organizations'] });
    },
  });
}

export function useUpdateOrganization() {
  const queryClient = useQueryClient();
  return useMutation<Organization, APIError, { orgId: string; name: string; slug: string }>({
    mutationFn: ({ orgId, ...data }) => apiClient.patch(`/organizations/${orgId}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['organizations'] });
    },
  });
}

export function useDeleteOrganization() {
  const queryClient = useQueryClient();
  return useMutation<void, APIError, string>({
    mutationFn: (orgId) => apiClient.delete(`/organizations/${orgId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['organizations'] });
    },
  });
}
