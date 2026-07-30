'use client';

import { useAuth } from '@/hooks/use-auth';
import { useProjects } from '@/hooks/use-projects';
import { useOrganizations } from '@/hooks/use-organizations';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { ArrowRight, Plus } from 'lucide-react';
import Link from 'next/link';
import { useTranslation } from 'react-i18next';

export default function ProjectsPage() {
  const { t } = useTranslation();
  const { isLoading: authLoading } = useAuth();
  const { data: orgs, isLoading: orgsLoading } = useOrganizations();
  const orgId = orgs?.[0]?.id;
  const { data: projects, isLoading: projectsLoading } = useProjects(orgId);

  if (authLoading || orgsLoading || projectsLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-48" />
        <div className="grid gap-4 sm:grid-cols-2">
          {Array.from({ length: 4 }).map((_, i) => (
            <Skeleton key={i} className="h-28 rounded-lg" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold tracking-tight">{t('projects.title')}</h1>
        <Button asChild size="sm">
          <Link href="/projects/new">
            <Plus className="mr-2 h-4 w-4" />
            {t('projects.create')}
          </Link>
        </Button>
      </div>

      {(!projects || projects.length === 0) ? (
        <Card>
          <CardContent className="py-12 text-center text-muted-foreground">
            {t('common.noData')}
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {projects.map((project) => (
            <Link
              key={project.id}
              href={`/projects/${project.id}`}
              className="group rounded-lg border p-5 transition-colors hover:bg-accent"
            >
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-md bg-primary/10 font-bold text-primary text-sm">
                  {project.key.slice(0, 2)}
                </div>
                <div className="min-w-0 flex-1">
                  <p className="truncate font-medium">{project.name}</p>
                  <p className="text-xs text-muted-foreground">{project.key}</p>
                </div>
                <ArrowRight className="h-4 w-4 shrink-0 text-muted-foreground opacity-0 transition-opacity group-hover:opacity-100" />
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
