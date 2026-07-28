'use client';

import { useTranslation } from 'react-i18next';
import { useParams } from 'next/navigation';
import { useState } from 'react';
import { useProject } from '@/hooks/use-projects';
import { useTasks } from '@/hooks/use-tasks';
import { useBoardStore } from '@/stores/board-store';
import * as Tabs from '@radix-ui/react-tabs';
import { motion, AnimatePresence } from 'framer-motion';
import { format } from 'date-fns';
import {
  AlertCircle,
  ArrowUp,
  Minus,
  ArrowDown,
  Search,
  ListTodo,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog';
import { Separator } from '@/components/ui/separator';
import type { Task, TaskType, TaskPriority, TaskStatus } from '@/types/api';

const priorityConfig: Record<TaskPriority, { icon: typeof AlertCircle | null; className: string }> = {
  urgent: { icon: AlertCircle, className: 'text-red-500' },
  high: { icon: ArrowUp, className: 'text-orange-500' },
  medium: { icon: Minus, className: 'text-yellow-500' },
  low: { icon: ArrowDown, className: 'text-blue-500' },
  none: { icon: null, className: '' },
};

const typeBadgeVariant: Record<TaskType, 'default' | 'secondary' | 'destructive' | 'outline' | 'success' | 'warning'> = {
  bug: 'destructive',
  story: 'success',
  task: 'default',
  epic: 'secondary',
};

const statusBadgeVariant: Record<TaskStatus, 'default' | 'secondary' | 'destructive' | 'outline' | 'success' | 'warning'> = {
  backlog: 'outline',
  todo: 'default',
  in_progress: 'warning',
  in_review: 'warning',
  done: 'success',
  cancelled: 'destructive',
};

function TaskDetailDialog() {
  const { t } = useTranslation();
  const { selectedTask, setSelectedTask } = useBoardStore();

  if (!selectedTask) return null;

  const PriorityIcon = priorityConfig[selectedTask.priority].icon;

  return (
    <Dialog open={!!selectedTask} onOpenChange={(open) => !open && setSelectedTask(null)}>
      <DialogContent className="sm:max-w-2xl">
        <DialogHeader>
          <div className="flex items-center gap-2 mb-2">
            <Badge variant={typeBadgeVariant[selectedTask.type]}>
              {selectedTask.type}
            </Badge>
            {selectedTask.priority !== 'none' && PriorityIcon && (
              <PriorityIcon className={cn('h-4 w-4', priorityConfig[selectedTask.priority].className)} />
            )}
          </div>
          <DialogTitle className="text-xl">{selectedTask.title}</DialogTitle>
          <DialogDescription>
            {selectedTask.description || t('comments.noCommentsDesc')}
          </DialogDescription>
        </DialogHeader>

        <div className="grid grid-cols-2 gap-4 text-sm">
          <div className="space-y-3">
            <div>
              <p className="text-muted-foreground text-xs font-medium mb-1">{t('common.status')}</p>
              <Badge variant={statusBadgeVariant[selectedTask.status]}>
                {selectedTask.status.replace('_', ' ')}
              </Badge>
            </div>
            <div>
              <p className="text-muted-foreground text-xs font-medium mb-1">{t('tasks.assignee')}</p>
              {selectedTask.assignee ? (
                <div className="flex items-center gap-2">
                  <Avatar className="h-6 w-6">
                    <AvatarImage src={selectedTask.assignee.avatarUrl ?? undefined} />
                    <AvatarFallback className="text-[9px]">
                      {selectedTask.assignee.name.charAt(0).toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                  <span>{selectedTask.assignee.name}</span>
                </div>
              ) : (
                <span className="text-muted-foreground">-</span>
              )}
            </div>
          </div>
          <div className="space-y-3">
            <div>
              <p className="text-muted-foreground text-xs font-medium mb-1">Story Points</p>
              <span>{selectedTask.storyPoints ?? '-'}</span>
            </div>
            <div>
              <p className="text-muted-foreground text-xs font-medium mb-1">{t('common.update')}</p>
              <span>{format(new Date(selectedTask.updatedAt), 'MMM d, yyyy')}</span>
            </div>
          </div>
        </div>

        {selectedTask.tags.length > 0 && (
          <>
            <Separator />
            <div className="flex flex-wrap gap-1.5">
              {selectedTask.tags.map((tag) => (
                <Badge
                  key={tag.id}
                  variant="outline"
                  className="text-xs"
                  style={{
                    borderColor: tag.color,
                    color: tag.color,
                  }}
                >
                  {tag.name}
                </Badge>
              ))}
            </div>
          </>
        )}
      </DialogContent>
    </Dialog>
  );
}

export default function BacklogPage() {
  const params = useParams();
  const projectId = params.projectId as string;
  const { data: project, isLoading: projectLoading } = useProject(projectId);
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [priorityFilter, setPriorityFilter] = useState<string>('');
  const [typeFilter, setTypeFilter] = useState<string>('');
  const { t } = useTranslation();

  const { data: tasksData, isLoading: tasksLoading, error } = useTasks(projectId, {
    search: search || undefined,
    status: statusFilter || undefined,
    priority: priorityFilter || undefined,
    type: typeFilter || undefined,
    pageSize: 100,
  });

  const tasks = tasksData?.data ?? [];
  const isLoading = projectLoading || tasksLoading;
  const hasFilters = search || statusFilter || priorityFilter || typeFilter;

  return (
    <div className="flex flex-col h-full">
      <div className="mb-6">
        <div className="flex items-center gap-3 mb-1">
          <h1 className="text-2xl font-bold tracking-tight">{project?.name ?? 'Loading...'}</h1>
          {project && (
            <Badge variant="secondary" className="text-xs font-mono uppercase tracking-wider">
              {project.key}
            </Badge>
          )}
        </div>
        {project?.description && (
          <p className="text-sm text-muted-foreground mt-1">{project.description}</p>
        )}
      </div>

      <Tabs.Root value="backlog" className="flex flex-col flex-1 min-h-0">
        <Tabs.List className="flex gap-0 border-b mb-4">
          <Tabs.Trigger
            value="board"
            className={cn(
              'px-4 py-2.5 text-sm font-medium transition-colors relative',
              'data-[state=active]:text-foreground data-[state=inactive]:text-muted-foreground',
              'data-[state=inactive]:hover:text-foreground/80',
              'data-[state=active]:after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-primary',
            )}
          >
            {t('projects.board')}
          </Tabs.Trigger>
          <Tabs.Trigger
            value="backlog"
            className={cn(
              'px-4 py-2.5 text-sm font-medium transition-colors relative',
              'data-[state=active]:text-foreground data-[state=inactive]:text-muted-foreground',
              'data-[state=inactive]:hover:text-foreground/80',
              'data-[state=active]:after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-primary',
            )}
          >
            {t('projects.backlog')}
          </Tabs.Trigger>
          <Tabs.Trigger
            value="settings"
            className={cn(
              'px-4 py-2.5 text-sm font-medium transition-colors relative',
              'data-[state=active]:text-foreground data-[state=inactive]:text-muted-foreground',
              'data-[state=inactive]:hover:text-foreground/80',
              'data-[state=active]:after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-primary',
            )}
          >
            {t('projects.settings')}
          </Tabs.Trigger>
        </Tabs.List>

        <Tabs.Content value="backlog" className="flex-1 min-h-0">
          <div className="h-full flex flex-col gap-4">
            <div className="flex flex-wrap items-center gap-3 shrink-0">
              <div className="relative flex-1 min-w-[200px] max-w-sm">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder={t('tasks.searchPlaceholder')}
                  className="pl-9 h-9"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                />
              </div>

              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="w-[140px] h-9">
                  <SelectValue placeholder={t('common.status')} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">{t('common.filter') + ': ' + t('common.status')}</SelectItem>
                  <SelectItem value="backlog">{t('tasks.statuses.backlog')}</SelectItem>
                  <SelectItem value="todo">{t('tasks.statuses.todo')}</SelectItem>
                  <SelectItem value="in_progress">{t('tasks.statuses.in_progress')}</SelectItem>
                  <SelectItem value="in_review">In Review</SelectItem>
                  <SelectItem value="done">{t('tasks.statuses.done')}</SelectItem>
                  <SelectItem value="cancelled">{t('tasks.statuses.cancelled')}</SelectItem>
                </SelectContent>
              </Select>

              <Select value={priorityFilter} onValueChange={setPriorityFilter}>
                <SelectTrigger className="w-[140px] h-9">
                  <SelectValue placeholder={t('common.priority')} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">{t('common.filter') + ': ' + t('common.priority')}</SelectItem>
                  <SelectItem value="urgent">{t('tasks.priorities.urgent')}</SelectItem>
                  <SelectItem value="high">{t('tasks.priorities.high')}</SelectItem>
                  <SelectItem value="medium">{t('tasks.priorities.medium')}</SelectItem>
                  <SelectItem value="low">{t('tasks.priorities.low')}</SelectItem>
                  <SelectItem value="none">{t('tasks.priorities.none')}</SelectItem>
                </SelectContent>
              </Select>

              <Select value={typeFilter} onValueChange={setTypeFilter}>
                <SelectTrigger className="w-[140px] h-9">
                  <SelectValue placeholder={t('common.type')} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">{t('common.filter') + ': ' + t('common.type')}</SelectItem>
                  <SelectItem value="task">{t('tasks.types.task')}</SelectItem>
                  <SelectItem value="bug">{t('tasks.types.bug')}</SelectItem>
                  <SelectItem value="story">{t('tasks.types.story')}</SelectItem>
                  <SelectItem value="epic">{t('tasks.types.epic')}</SelectItem>
                </SelectContent>
              </Select>

              {hasFilters && (
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-9 text-xs"
                  onClick={() => {
                    setSearch('');
                    setStatusFilter('');
                    setPriorityFilter('');
                    setTypeFilter('');
                  }}
                >
                  {t('common.filter')}
                </Button>
              )}
            </div>

            <div className="flex-1 min-h-0">
              {isLoading ? (
                <div className="space-y-2">
                  {[1, 2, 3, 4, 5, 6].map((i) => (
                    <Skeleton key={i} className="h-14 w-full rounded-lg" />
                  ))}
                </div>
              ) : error ? (
                <div className="flex flex-col items-center justify-center h-48 text-muted-foreground gap-2">
                  <AlertCircle className="h-8 w-8 text-destructive" />
                  <p className="text-sm font-medium">{t('common.error')}</p>
                  <p className="text-xs">{(error as any)?.message ?? 'An unexpected error occurred'}</p>
                </div>
              ) : tasks.length === 0 ? (
                <div className="flex flex-col items-center justify-center h-48 text-muted-foreground gap-3">
                  <div className="rounded-full bg-muted p-3">
                    {hasFilters ? (
                      <Search className="h-5 w-5" />
                    ) : (
                      <ListTodo className="h-5 w-5" />
                    )}
                  </div>
                  <p className="text-sm font-medium">
                    {hasFilters ? t('tasks.noTasks') : t('board.noTasks')}
                  </p>
                  <p className="text-xs">
                    {hasFilters
                      ? 'Try adjusting your search or filter criteria'
                      : 'Create your first task to populate the backlog'}
                  </p>
                </div>
              ) : (
                <div className="rounded-lg border overflow-hidden">
                  <div className="grid grid-cols-12 gap-3 px-4 py-3 bg-muted/50 text-xs font-medium text-muted-foreground uppercase tracking-wider">
                    <div className="col-span-5 lg:col-span-4">{t('common.title')}</div>
                    <div className="col-span-2 lg:col-span-2">{t('common.type')}</div>
                    <div className="col-span-2 lg:col-span-2">{t('common.priority')}</div>
                    <div className="col-span-2 lg:col-span-2">{t('common.status')}</div>
                    <div className="hidden lg:block lg:col-span-1">{t('tasks.assignee')}</div>
                    <div className="col-span-1 text-right">{t('common.update')}</div>
                  </div>
                  <Separator />
                  <AnimatePresence initial={false}>
                    {tasks.map((task) => (
                      <BacklogRow key={task.id} task={task} />
                    ))}
                  </AnimatePresence>
                </div>
              )}

              {tasksData && tasksData.total > tasks.length && (
                <p className="text-xs text-muted-foreground text-center mt-4">
                  Showing {tasks.length} of {tasksData.total} tasks
                </p>
              )}
            </div>
          </div>
        </Tabs.Content>

        <Tabs.Content value="board" className="flex-1 min-h-0 data-[state=inactive]:hidden">
          <div className="flex flex-col items-center justify-center h-64 text-muted-foreground gap-3">
            <ListTodo className="h-8 w-8" />
            <p className="text-sm font-medium">{t('projects.board')}</p>
            <p className="text-xs">Switch to the Board tab or visit the project overview</p>
          </div>
        </Tabs.Content>

        <Tabs.Content value="settings" className="flex-1 min-h-0 data-[state=inactive]:hidden">
          <div className="flex flex-col items-center justify-center h-64 text-muted-foreground gap-3">
            <AlertCircle className="h-8 w-8" />
            <p className="text-sm font-medium">{t('projects.settings')}</p>
            <p className="text-xs">{t('common.loading')}</p>
          </div>
        </Tabs.Content>
      </Tabs.Root>

      <TaskDetailDialog />
    </div>
  );
}

function BacklogRow({ task }: { task: Task }) {
  const { setSelectedTask } = useBoardStore();
  const PriorityIcon = priorityConfig[task.priority].icon;

  return (
    <motion.div
      layout
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      transition={{ duration: 0.15 }}
      className="grid grid-cols-12 gap-3 px-4 py-3 text-sm items-center hover:bg-accent/50 cursor-pointer transition-colors border-b last:border-b-0"
      onClick={() => setSelectedTask(task)}
    >
      <div className="col-span-5 lg:col-span-4 font-medium truncate flex items-center gap-2">
        {task.type === 'bug' && <span className="text-destructive text-xs shrink-0">●</span>}
        {task.title}
      </div>
      <div className="col-span-2">
        <Badge
          variant={typeBadgeVariant[task.type]}
          className="text-[10px] px-1.5 py-0 uppercase"
        >
          {task.type}
        </Badge>
      </div>
      <div className="col-span-2 flex items-center gap-1.5">
        {PriorityIcon ? (
          <>
            <PriorityIcon className={cn('h-3.5 w-3.5', priorityConfig[task.priority].className)} />
            <span className="text-xs capitalize text-muted-foreground hidden sm:inline">{task.priority}</span>
          </>
        ) : (
          <span className="text-xs text-muted-foreground">-</span>
        )}
      </div>
      <div className="col-span-2">
        <Badge
          variant={statusBadgeVariant[task.status]}
          className="text-[10px] px-1.5 py-0"
        >
          {task.status.replace('_', ' ')}
        </Badge>
      </div>
      <div className="hidden lg:flex lg:col-span-1 items-center">
        {task.assignee ? (
          <Avatar className="h-7 w-7">
            <AvatarImage src={task.assignee.avatarUrl ?? undefined} />
            <AvatarFallback className="text-[9px]">
              {task.assignee.name.charAt(0).toUpperCase()}
            </AvatarFallback>
          </Avatar>
        ) : (
          <span className="text-xs text-muted-foreground">-</span>
        )}
      </div>
      <div className="col-span-1 text-right text-xs text-muted-foreground">
        {format(new Date(task.updatedAt), 'MMM d')}
      </div>
    </motion.div>
  );
}
