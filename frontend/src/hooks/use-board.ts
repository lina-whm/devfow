'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import type { Board, Column, APIError } from '@/types/api';

export function useBoard(projectId: string | undefined) {
  return useQuery<Board>({
    queryKey: ['board', projectId],
    queryFn: () => apiClient.get(`/projects/${projectId}/board`),
    enabled: !!projectId,
  });
}

export function useUpdateColumns() {
  const queryClient = useQueryClient();

  return useMutation<Board, APIError, { boardId: string; columns: Partial<Column>[] }>({
    mutationFn: ({ boardId, columns }) => apiClient.patch(`/boards/${boardId}/columns`, { columns }),
    onSuccess: (board) => {
      queryClient.invalidateQueries({ queryKey: ['board', board.projectId] });
    },
  });
}
