'use client';

import { useState, useMemo } from 'react';
import Link from 'next/link';
import { useTasks } from '@/hooks/use-tasks';
import { useProjects } from '@/hooks/use-projects';
import { useOrganizations } from '@/hooks/use-organizations';
import { useAuth } from '@/hooks/use-auth';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { ListTodo, Search, AlertCircle, ArrowUp, ArrowDown, Minus } from 'lucide-react';
import { format } from 'date-fns';
import { useTranslation } from 'react-i18next';

const priorityIcon = (priority: string) => {
  switch (priority) {
    case 'urgent':
    case 'high':
      return <ArrowUp className="h-3.5 w-3.5 text-red-500" />;
    case 'medium':
      return <Minus className="h-3.5 w-3.5 text-yellow-500" />;
    case 'low':
      return <ArrowDown className="h-3.5 w-3.5 text-green-500" />;
    default:
      return <Minus className="h-3.5 w-3.5 text-muted-foreground" />;
  }
};

const statusVariant = (status: string) => {
  switch (status) {
    case 'done':
      return 'success' as const;
    case 'in_progress':
      return 'default' as const;
    case 'in_review':
      return 'warning' as const;
    case 'todo':
      return 'secondary' as const;
    case 'backlog':
      return 'outline' as const;
    case 'cancelled':
      return 'destructive' as const;
    default:
      return 'secondary' as const;
  }
};

const statusLabel = (t: (key: string) => string, status: string): string => {
  const labels: Record<string, string> = {
    backlog: t('tasks.statuses.backlog'),
    todo: t('tasks.statuses.todo'),
    in_progress: t('tasks.statuses.in_progress'),
    in_review: 'In Review',
    done: t('tasks.statuses.done'),
    cancelled: t('tasks.statuses.cancelled'),
  };
  return labels[status] || status.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase());
};

const typeVariant = (type: string) => {
  switch (type) {
    case 'bug':
      return 'destructive' as const;
    case 'story':
      return 'success' as const;
    case 'epic':
      return 'warning' as const;
    default:
      return 'secondary' as const;
  }
};

export default function TasksPage() {
  const { t } = useTranslation();
  const { user } = useAuth();
  const { data: orgs } = useOrganizations();
  const orgId = orgs?.[0]?.id;
  const { data: projects } = useProjects(orgId);
  const firstProjectId = projects?.[0]?.id;
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('all');
  const [priorityFilter, setPriorityFilter] = useState('all');

  const filters = useMemo(() => {
    const f: Record<string, string> = {};
    if (statusFilter !== 'all') f.status = statusFilter;
    if (priorityFilter !== 'all') f.priority = priorityFilter;
    if (search) f.search = search;
    if (user?.id) f.assigneeId = user.id;
    return f;
  }, [statusFilter, priorityFilter, search, user?.id]);

  const { data, isLoading, error } = useTasks(firstProjectId, filters);

  const filteredTasks = data?.data ?? [];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">{t('tasks.myTasks')}</h1>
          <p className="text-muted-foreground mt-1">
            {t('tasks.title')}
          </p>
        </div>
      </div>

      <Card>
        <CardHeader className="pb-3">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder={t('tasks.searchPlaceholder')}
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-9"
              />
            </div>
            <div className="flex items-center gap-2">
              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="w-[140px]">
                  <SelectValue placeholder={t('common.status')} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">{t('common.filter')}: {t('common.status')}</SelectItem>
                  <SelectItem value="backlog">{t('tasks.statuses.backlog')}</SelectItem>
                  <SelectItem value="todo">{t('tasks.statuses.todo')}</SelectItem>
                  <SelectItem value="in_progress">{t('tasks.statuses.in_progress')}</SelectItem>
                  <SelectItem value="in_review">In Review</SelectItem>
                  <SelectItem value="done">{t('tasks.statuses.done')}</SelectItem>
                  <SelectItem value="cancelled">{t('tasks.statuses.cancelled')}</SelectItem>
                </SelectContent>
              </Select>
              <Select value={priorityFilter} onValueChange={setPriorityFilter}>
                <SelectTrigger className="w-[140px]">
                  <SelectValue placeholder={t('common.priority')} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">{t('common.filter')}: {t('common.priority')}</SelectItem>
                  <SelectItem value="urgent">{t('tasks.priorities.urgent')}</SelectItem>
                  <SelectItem value="high">{t('tasks.priorities.high')}</SelectItem>
                  <SelectItem value="medium">{t('tasks.priorities.medium')}</SelectItem>
                  <SelectItem value="low">{t('tasks.priorities.low')}</SelectItem>
                  <SelectItem value="none">{t('tasks.priorities.none')}</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-3">
              {Array.from({ length: 5 }).map((_, i) => (
                <Skeleton key={i} className="h-14 w-full" />
              ))}
            </div>
          ) : error ? (
            <div className="flex flex-col items-center justify-center py-12 text-center">
              <AlertCircle className="h-12 w-12 text-destructive mb-4" />
              <h3 className="text-lg font-semibold">{t('tasks.title')}</h3>
              <p className="text-sm text-muted-foreground mt-1">
                {error instanceof Error ? error.message : t('common.error')}
              </p>
            </div>
          ) : filteredTasks.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-12 text-center">
              <ListTodo className="h-12 w-12 text-muted-foreground mb-4" />
              <h3 className="text-lg font-semibold">{t('tasks.noTasks')}</h3>
              <p className="text-sm text-muted-foreground mt-1">
                {search || statusFilter !== 'all' || priorityFilter !== 'all'
                  ? t('common.filter')
                  : t('tasks.noTasksDesc')}
              </p>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b text-left text-sm text-muted-foreground">
                    <th className="pb-3 font-medium pl-1">{t('common.title')}</th>
                    <th className="pb-3 font-medium">{t('common.type')}</th>
                    <th className="pb-3 font-medium">{t('common.priority')}</th>
                    <th className="pb-3 font-medium">{t('common.status')}</th>
                    <th className="pb-3 font-medium">{t('tasks.dueDate')}</th>
                  </tr>
                </thead>
                <tbody>
                  {filteredTasks.map((task) => (
                    <tr
                      key={task.id}
                      className="border-b last:border-0 transition-colors hover:bg-muted/50 cursor-pointer"
                      onClick={() => window.location.href = `/tasks/${task.id}`}
                    >
                      <td className="py-3 pl-1">
                        <Link
                          href={`/tasks/${task.id}`}
                          className="text-sm font-medium hover:text-primary transition-colors"
                        >
                          {task.title}
                        </Link>
                      </td>
                      <td className="py-3">
                        <Badge variant={typeVariant(task.type)} className="capitalize text-xs">
                          {t('tasks.types.' + task.type) || task.type}
                        </Badge>
                      </td>
                      <td className="py-3">
                        <div className="flex items-center gap-1.5">
                          {priorityIcon(task.priority)}
                          <span className="text-sm text-muted-foreground">
                            {t('tasks.priorities.' + task.priority) || task.priority}
                          </span>
                        </div>
                      </td>
                      <td className="py-3">
                        <Badge variant={statusVariant(task.status)} className="text-xs">
                          {statusLabel(t, task.status)}
                        </Badge>
                      </td>
                      <td className="py-3">
                        <span className="text-sm text-muted-foreground">
                          {task.dueDate ? format(new Date(task.dueDate), 'MMM d, yyyy') : '-'}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
