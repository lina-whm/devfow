'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import type { Comment, APIError } from '@/types/api';

export function useComments(taskId: string | undefined) {
  return useQuery<Comment[]>({
    queryKey: ['comments', taskId],
    queryFn: () => apiClient.get(`/tasks/${taskId}/comments`),
    enabled: !!taskId,
  });
}

export function useCreateComment() {
  const queryClient = useQueryClient();
  return useMutation<Comment, APIError, { taskId: string; content: string }>({
    mutationFn: ({ taskId, content }) => apiClient.post(`/tasks/${taskId}/comments`, { body: content }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments'] });
    },
  });
}
