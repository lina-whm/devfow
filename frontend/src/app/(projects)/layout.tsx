'use client';
import { AuthProvider } from '@/hooks/use-auth';
import { AppShell } from '@/components/layout/app-shell';

export default function ProjectsLayout({ children }: { children: React.ReactNode }) {
  return (
    <AuthProvider>
      <AppShell>{children}</AppShell>
    </AuthProvider>
  );
}
