'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/lib/api-client';
import type { Notification, PaginatedResponse, APIError } from '@/types/api';

export function useNotifications(page = 1, pageSize = 20) {
  return useQuery<PaginatedResponse<Notification>>({
    queryKey: ['notifications', page, pageSize],
    queryFn: async () => {
      try {
        return await apiClient.get<PaginatedResponse<Notification>>(`/notifications?page=${page}&pageSize=${pageSize}`);
      } catch {
        return { data: [], total: 0, page, pageSize, totalPages: 0 } as PaginatedResponse<Notification>;
      }
    },
    retry: false,
  });
}

export function useUnreadCount() {
  return useQuery<number>({
    queryKey: ['notifications', 'unread-count'],
    queryFn: async () => {
      try {
        return await apiClient.get<number>('/notifications/unread-count');
      } catch {
        return 0;
      }
    },
    refetchInterval: 30 * 1000,
    retry: false,
  });
}

export function useMarkAsRead() {
  const queryClient = useQueryClient();

  return useMutation<void, APIError, string>({
    mutationFn: (notificationId) =>
      apiClient.patch(`/notifications/${notificationId}/read`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
    },
  });
}

export function useMarkAllAsRead() {
  const queryClient = useQueryClient();

  return useMutation<void, APIError, void>({
    mutationFn: () => apiClient.patch('/notifications/read-all'),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
    },
  });
}
