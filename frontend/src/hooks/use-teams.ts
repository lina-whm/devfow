'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import type { Team, APIError } from '@/types/api';

export function useTeams(orgId: string | undefined) {
  return useQuery<Team[]>({
    queryKey: ['teams', orgId],
    queryFn: () => apiClient.get(`/organizations/${orgId}/teams`),
    enabled: !!orgId,
  });
}

export function useCreateTeam() {
  const queryClient = useQueryClient();
  return useMutation<Team, APIError, { orgId: string; name: string }>({
    mutationFn: ({ orgId, ...data }) => apiClient.post(`/organizations/${orgId}/teams`, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['teams'] });
    },
  });
}
