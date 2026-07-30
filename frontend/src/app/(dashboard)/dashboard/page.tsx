'use client';

import { useAuth } from '@/hooks/use-auth';
import { useProjects } from '@/hooks/use-projects';
import { useOrganizations } from '@/hooks/use-organizations';
import { useTasks } from '@/hooks/use-tasks';
import { useNotifications } from '@/hooks/use-notifications';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { LayoutDashboard, ListTodo, CheckCircle2, Clock, ArrowRight, Plus } from 'lucide-react';
import { cn } from '@/lib/utils';
import Link from 'next/link';
import { format } from 'date-fns';
import { useTranslation } from 'react-i18next';

const priorityColor: Record<string, string> = {
  urgent: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
  critical: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
  high: 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400',
  medium: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
  low: 'bg-gray-100 text-gray-800 dark:bg-gray-800/50 dark:text-gray-400',
  none: 'bg-gray-100 text-gray-500 dark:bg-gray-800/30 dark:text-gray-500',
};

const statusLabel: Record<string, string> = {
  backlog: 'Backlog',
  todo: 'To Do',
  in_progress: 'In Progress',
  in_review: 'In Review',
  done: 'Done',
  cancelled: 'Cancelled',
};

function StatCardSkeleton() {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <Skeleton className="h-4 w-24" />
        <Skeleton className="h-8 w-8 rounded" />
      </CardHeader>
      <CardContent>
        <Skeleton className="h-8 w-16" />
        <Skeleton className="mt-1 h-3 w-28" />
      </CardContent>
    </Card>
  );
}

function ActivitySkeleton() {
  return (
    <div className="space-y-4">
      {Array.from({ length: 4 }).map((_, i) => (
        <div key={i} className="flex items-start gap-3">
          <Skeleton className="mt-1 h-2 w-2 rounded-full" />
          <div className="flex-1 space-y-2">
            <Skeleton className="h-4 w-3/4" />
            <Skeleton className="h-3 w-1/3" />
          </div>
        </div>
      ))}
    </div>
  );
}

export default function DashboardPage() {
  const { t } = useTranslation();
  const { user, isLoading: authLoading } = useAuth();
  const { data: orgs, isLoading: orgsLoading } = useOrganizations();
  const orgId = orgs?.[0]?.id;
  const { data: projects, isLoading: projectsLoading } = useProjects(orgId);
  const firstProjectId = projects?.[0]?.id;

  const { data: myTasksData, isLoading: myTasksLoading } = useTasks(firstProjectId, {
    assigneeId: user?.id,
    pageSize: 50,
  });

  const { data: openTasksData } = useTasks(firstProjectId, {
    assigneeId: user?.id,
    status: 'in_progress',
    pageSize: 10,
  });

  const { data: notificationsData, isLoading: notifLoading } = useNotifications(1, 5);

  const myTasks = myTasksData?.data ?? [];
  const openTasks = openTasksData?.data ?? [];
  const notifications = notificationsData?.data ?? [];

  const taskDueSoon = myTasks.filter(
    (t) => t.dueDate && new Date(t.dueDate) > new Date() && t.status !== 'done' && t.status !== 'cancelled',
  );
  const lateTasks = myTasks.filter(
    (t) => t.dueDate && new Date(t.dueDate) <= new Date() && t.status !== 'done' && t.status !== 'cancelled',
  );

  const projectsCount = projects?.length ?? 0;
  const myTasksCount = myTasksData?.total ?? 0;
  const upcomingCount = taskDueSoon.length + lateTasks.length;

  const isLoading =
    authLoading || orgsLoading || projectsLoading || myTasksLoading || notifLoading;

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div>
          <Skeleton className="h-8 w-64" />
          <Skeleton className="mt-2 h-4 w-96" />
        </div>
        <div className="grid gap-4 md:grid-cols-3">
          <StatCardSkeleton />
          <StatCardSkeleton />
          <StatCardSkeleton />
        </div>
        <div className="grid gap-6 lg:grid-cols-2">
          <Card>
            <CardHeader>
              <Skeleton className="h-5 w-32" />
            </CardHeader>
            <CardContent>
              <ActivitySkeleton />
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <Skeleton className="h-5 w-40" />
            </CardHeader>
            <CardContent>
              <ActivitySkeleton />
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">
            {t('dashboard.welcome', { name: user?.displayName?.split(' ')[0] ?? 'there' })}
          </h1>
          <p className="mt-1 text-muted-foreground">
            {t('dashboard.subtitle')}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button asChild variant="outline" size="sm">
            <Link href="/tasks">
              <ListTodo className="mr-2 h-4 w-4" />
              {t('tasks.title')}
            </Link>
          </Button>
          <Button asChild size="sm">
            <Link href="/projects/new">
              <Plus className="mr-2 h-4 w-4" />
              {t('projects.create')}
            </Link>
          </Button>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        <Link href="/projects" className="block">
          <Card className="transition-colors hover:bg-accent cursor-pointer">
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium">{t('dashboard.stats.projects')}</CardTitle>
              <div className="rounded-md bg-primary/10 p-2 text-primary">
                <LayoutDashboard className="h-4 w-4" />
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold">{projectsCount}</div>
              <p className="mt-1 text-xs text-muted-foreground">
                {orgs?.[0]?.name ?? t('projects.title')}
              </p>
            </CardContent>
          </Card>
        </Link>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium">{t('dashboard.myTasks')}</CardTitle>
            <div className="rounded-md bg-blue-500/10 p-2 text-blue-600 dark:text-blue-400">
              <CheckCircle2 className="h-4 w-4" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold">{myTasksCount}</div>
            <p className="mt-1 text-xs text-muted-foreground">
              {openTasks.length > 0
                ? `${openTasks.length} ${t('tasks.statuses.in_progress')}`
                : t('dashboard.myTasks')}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-sm font-medium">{t('dashboard.stats.deadlines')}</CardTitle>
            <div className="rounded-md bg-amber-500/10 p-2 text-amber-600 dark:text-amber-400">
              <Clock className="h-4 w-4" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold">
              {upcomingCount}
              {lateTasks.length > 0 && (
                <span className="ml-2 text-sm font-normal text-destructive">
                  ({lateTasks.length} {t('dashboard.overdue')})
                </span>
              )}
            </div>
            <p className="mt-1 text-xs text-muted-foreground">
              {taskDueSoon[0]?.dueDate
                ? `${t('dashboard.nextDeadline')}: ${format(new Date(taskDueSoon[0].dueDate), 'MMM d')}`
                : t('dashboard.noUpcoming')}
            </p>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle className="text-lg">{t('dashboard.recentActivity')}</CardTitle>
            <Button asChild variant="ghost" size="sm" className="gap-1 text-muted-foreground">
              <Link href="/notifications">
                {t('dashboard.viewAll')}
                <ArrowRight className="h-3 w-3" />
              </Link>
            </Button>
          </CardHeader>
          <CardContent>
            {notifications.length === 0 ? (
              <p className="py-6 text-center text-sm text-muted-foreground">
                {t('common.noData')}
              </p>
            ) : (
              <div className="space-y-4">
                {notifications.slice(0, 5).map((n) => (
                  <div key={n.id} className="flex items-start gap-3">
                    <div
                      className={cn(
                        'mt-1.5 h-2 w-2 shrink-0 rounded-full',
                        n.read ? 'bg-muted-foreground/30' : 'bg-primary',
                      )}
                    />
                    <div className="min-w-0 flex-1">
                      <p className="truncate text-sm">{n.title}</p>
                      <p className="text-xs text-muted-foreground">
                        {format(new Date(n.createdAt), 'MMM d, h:mm a')}
                      </p>
                    </div>
                    {!n.read && <Badge variant="default" className="h-5 px-1.5 text-[10px]">new</Badge>}
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle className="text-lg">{t('tasks.myTasks')}</CardTitle>
            <Button asChild variant="ghost" size="sm" className="gap-1 text-muted-foreground">
              <Link href="/tasks">
                {t('dashboard.viewAll')}
                <ArrowRight className="h-3 w-3" />
              </Link>
            </Button>
          </CardHeader>
          <CardContent>
            {myTasks.length === 0 ? (
              <p className="py-6 text-center text-sm text-muted-foreground">
                {t('dashboard.noTasks')}
              </p>
            ) : (
              <div className="space-y-3">
                {myTasks.slice(0, 5).map((task) => (
                  <Link
                    key={task.id}
                    href={`/tasks/${task.id}`}
                    className="flex items-center justify-between rounded-lg border p-3 transition-colors hover:bg-accent"
                  >
                    <div className="min-w-0 flex-1">
                      <div className="flex items-center gap-2">
                        <span className="truncate text-sm font-medium">{task.title}</span>
                        <Badge
                          variant="outline"
                          className={cn(
                            'shrink-0 text-[10px]',
                            priorityColor[task.priority] ?? priorityColor.none,
                          )}
                        >
                          {task.priority}
                        </Badge>
                      </div>
                      <div className="mt-1 flex items-center gap-2 text-xs text-muted-foreground">
                        <span>{statusLabel[task.status] ?? task.status}</span>
                        {task.dueDate && (
                          <>
                            <span>&middot;</span>
                            <span
                              className={cn(
                                new Date(task.dueDate) <= new Date() && 'text-destructive font-medium',
                              )}
                            >
                              Due {format(new Date(task.dueDate), 'MMM d')}
                            </span>
                          </>
                        )}
                      </div>
                    </div>
                    <ArrowRight className="ml-2 h-4 w-4 shrink-0 text-muted-foreground" />
                  </Link>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
