'use client';

import * as React from 'react';
import { QueryClient, QueryClientProvider as TanStackProvider } from '@tanstack/react-query';

function makeQueryClient(): QueryClient {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 30 * 1000,
        retry: 1,
        refetchOnWindowFocus: false,
      },
      mutations: {
        retry: 0,
      },
    },
  });
}

let browserQueryClient: QueryClient | undefined;

function getQueryClient(): QueryClient {
  if (typeof window === 'undefined') {
    return makeQueryClient();
  }
  if (!browserQueryClient) {
    browserQueryClient = makeQueryClient();
  }
  return browserQueryClient;
}

export function QueryClientProvider({ children }: { children: React.ReactNode }) {
  const queryClient = getQueryClient();
  return <TanStackProvider client={queryClient}>{children}</TanStackProvider>;
}
