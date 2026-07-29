'use client';

import { useParams } from 'next/navigation';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useProject } from '@/hooks/use-projects';
import { useBoard } from '@/hooks/use-board';
import { useMoveTask } from '@/hooks/use-tasks';
import { useBoardStore } from '@/stores/board-store';
import * as Tabs from '@radix-ui/react-tabs';
import { DndContext, DragOverlay, useSensor, useSensors, PointerSensor, closestCorners, type DragStartEvent, type DragOverEvent, type DragEndEvent } from '@dnd-kit/core';
import { SortableContext, useSortable, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { motion, AnimatePresence } from 'framer-motion';
import {
  AlertCircle,
  ArrowUp,
  Minus,
  ArrowDown,
  Plus,
  GripVertical,
  LayoutDashboard,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog';
import { Separator } from '@/components/ui/separator';
import type { Task, Column, TaskType, TaskPriority, TaskStatus } from '@/types/api';

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

function SortableTaskCard({ task }: { task: Task }) {
  const { setNodeRef, attributes, listeners, transform, transition, isDragging } = useSortable({ id: task.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <div ref={setNodeRef} style={style} {...attributes} {...listeners}>
      <TaskCardContent task={task} isDragging={isDragging} />
    </div>
  );
}

function TaskCardContent({ task, isDragging }: { task: Task; isDragging?: boolean }) {
  const { setSelectedTask } = useBoardStore();
  const PriorityIcon = priorityConfig[task.priority].icon;

  return (
    <motion.div
      layout
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, scale: 0.95 }}
      transition={{ duration: 0.15 }}
    >
      <Card
        className={cn(
          'cursor-grab active:cursor-grabbing p-3 transition-all duration-150',
          'hover:border-primary/30 hover:shadow-md hover:-translate-y-0.5',
          'border-l-2 border-l-transparent hover:border-l-primary/50',
          isDragging && 'opacity-40 shadow-lg ring-2 ring-primary/30',
        )}
        onClick={(e) => {
          e.stopPropagation();
          setSelectedTask(task);
        }}
      >
        <div className="flex items-start justify-between gap-2 mb-2.5">
          <Badge variant={typeBadgeVariant[task.type]} className="text-[10px] px-1.5 py-0 uppercase tracking-wider font-semibold">
            {task.type}
          </Badge>
          <div className="flex items-center gap-1 shrink-0">
            {PriorityIcon && (
              <PriorityIcon className={cn('h-3.5 w-3.5', priorityConfig[task.priority].className)} />
            )}
            <GripVertical className="h-3.5 w-3.5 text-muted-foreground/40" />
          </div>
        </div>
        <p className="text-sm font-medium leading-snug mb-2.5 line-clamp-2">{task.title}</p>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            {task.storyPoints != null && (
              <span className="text-[11px] text-muted-foreground bg-muted px-1.5 py-0.5 rounded font-medium">
                {task.storyPoints} SP
              </span>
            )}
            {task.labels?.slice(0, 2).map((label) => (
              <span
                key={label}
                className="text-[10px] text-muted-foreground bg-muted/50 px-1.5 py-0.5 rounded"
              >
                {label}
              </span>
            ))}
          </div>
          {task.assignee && (
            <Avatar className="h-6 w-6 ring-2 ring-background">
              <AvatarImage src={task.assignee.avatar_url ?? undefined} />
              <AvatarFallback className="text-[9px] font-medium">
                {task.assignee.display_name
                  .split(' ')
                  .map((n) => n[0])
                  .join('')
                  .toUpperCase()
                  .slice(0, 2)}
              </AvatarFallback>
            </Avatar>
          )}
        </div>
      </Card>
    </motion.div>
  );
}

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
                    <AvatarImage src={selectedTask.assignee.avatar_url ?? undefined} />
                    <AvatarFallback className="text-[9px]">
                      {selectedTask.assignee.display_name.charAt(0).toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                  <span>{selectedTask.assignee.display_name}</span>
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
              <span>{new Date(selectedTask.updatedAt).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}</span>
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

function BoardSkeleton() {
  return (
    <div className="flex gap-4 overflow-hidden">
      {[1, 2, 3, 4].map((i) => (
        <div key={i} className="w-72 flex-shrink-0 space-y-3">
          <Skeleton className="h-5 w-24" />
          {[1, 2, 3].map((j) => (
            <Skeleton key={j} className="h-28 w-full rounded-lg" />
          ))}
        </div>
      ))}
    </div>
  );
}

export default function BoardPage() {
  const params = useParams();
  const projectId = params.projectId as string;
  const { data: project, isLoading: projectLoading } = useProject(projectId);
  const { data: board, isLoading: boardLoading, error } = useBoard(projectId);
  const moveTask = useMoveTask();
  const { setDragState, resetDragState } = useBoardStore();
  const { t } = useTranslation();
  const [activeTask, setActiveTask] = useState<Task | null>(null);

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
  );

  const columns = board?.columns ?? [];

  function findColumnByTaskId(taskId: string): Column | undefined {
    return columns.find((col) => col.tasks.some((t) => t.id === taskId));
  }

  function handleDragStart(event: DragStartEvent) {
    const taskId = String(event.active.id);
    const sourceColumn = findColumnByTaskId(taskId);
    if (sourceColumn) {
      setDragState({ activeTaskId: taskId, sourceColumnId: sourceColumn.id });
      const task = sourceColumn.tasks.find((t) => t.id === taskId);
      if (task) setActiveTask(task);
    }
  }

  function handleDragOver(event: DragOverEvent) {
    const overId = event.over?.id;
    if (!overId) return;
    const overColumn = columns.find(
      (col) => col.id === overId || col.tasks.some((t) => t.id === overId),
    );
    if (overColumn) {
      setDragState({ targetColumnId: overColumn.id });
    }
  }

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;
    if (!over) {
      resetDragState();
      setActiveTask(null);
      return;
    }

    const activeId = String(active.id);
    const overId = String(over.id);

    const sourceColumn = findColumnByTaskId(activeId);
    const targetColumn = columns.find(
      (col) => col.id === overId || col.tasks.some((t) => t.id === overId),
    );

    if (!sourceColumn || !targetColumn) {
      resetDragState();
      setActiveTask(null);
      return;
    }

    const isSameColumn = sourceColumn.id === targetColumn.id;
    if (isSameColumn && activeId === overId) {
      resetDragState();
      setActiveTask(null);
      return;
    }

    const targetTask = targetColumn.tasks.find((t) => t.id === overId);
    const targetIndex = targetTask
      ? targetColumn.tasks.indexOf(targetTask)
      : targetColumn.tasks.length;

    moveTask.mutate(
      {
        taskId: activeId,
        columnId: targetColumn.id,
        position: targetIndex,
      },
      {
        onSettled: () => {
          resetDragState();
          setActiveTask(null);
        },
      },
    );
  }

  const isLoading = projectLoading || boardLoading;

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

      <Tabs.Root value="board" className="flex flex-col flex-1 min-h-0">
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

        <Tabs.Content value="board" className="flex-1 min-h-0">
          {isLoading ? (
            <BoardSkeleton />
          ) : error ? (
            <div className="flex flex-col items-center justify-center h-64 text-muted-foreground gap-2">
              <AlertCircle className="h-8 w-8 text-destructive" />
<p className="text-sm font-medium">{t('common.error')}</p>
               <p className="text-xs">{(error as any)?.message ?? 'An unexpected error occurred'}</p>
            </div>
          ) : columns.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-64 text-muted-foreground gap-3">
              <div className="rounded-full bg-muted p-3">
                <LayoutDashboard className="h-6 w-6" />
              </div>
              <p className="text-sm font-medium">{t('board.noTasks')}</p>
              <p className="text-xs">{t('board.addTask')}</p>
            </div>
          ) : (
            <DndContext
              sensors={sensors}
              collisionDetection={closestCorners}
              onDragStart={handleDragStart}
              onDragOver={handleDragOver}
              onDragEnd={handleDragEnd}
            >
              <div className="flex gap-5 h-full overflow-x-auto pb-4">
                {columns
                  .slice()
                  .sort((a, b) => a.position - b.position)
                  .map((column) => {
                    const sortedTasks = [...column.tasks].sort((a, b) => a.position - b.position);
                    const isOverWip = column.wipLimit != null && sortedTasks.length >= column.wipLimit;

                    return (
                      <div
                        key={column.id}
                        className="flex-shrink-0 w-72 flex flex-col rounded-xl bg-muted/40 border"
                      >
                        <div className="flex items-center justify-between px-4 py-3 border-b">
                          <div className="flex items-center gap-2.5 min-w-0">
                            <div
                              className={cn('w-2 h-2 rounded-full shrink-0', {
                                'bg-gray-400': column.name === 'Backlog',
                                'bg-blue-500': column.name === 'Todo' || column.name === 'To Do',
                                'bg-yellow-500': column.name === 'In Progress',
                                'bg-orange-500': column.name === 'In Review' || column.name === 'Review',
                                'bg-green-500': column.name === 'Done',
                                'bg-red-500': column.name === 'Cancelled',
                              })}
                            />
                            <h3 className="text-sm font-semibold truncate">{column.name}</h3>
                            <span className="text-xs text-muted-foreground bg-muted rounded-full px-2 py-0.5 font-medium shrink-0">
                              {sortedTasks.length}
                            </span>
                          </div>
                          {column.wipLimit != null && (
                            <span
                              className={cn(
                                'text-[11px] px-1.5 py-0.5 rounded font-medium shrink-0 ml-2',
                                isOverWip
                                  ? 'bg-destructive/10 text-destructive border border-destructive/20'
                                  : 'bg-muted text-muted-foreground',
                              )}
                            >
                              {sortedTasks.length}/{column.wipLimit}
                            </span>
                          )}
                        </div>

                        <ScrollArea className="flex-1 p-3">
                          <SortableContext
                            items={sortedTasks.map((t) => t.id)}
                            strategy={verticalListSortingStrategy}
                          >
                            <AnimatePresence mode="popLayout">
                              {sortedTasks.length === 0 ? (
                                <motion.div
                                  layout
                                  initial={{ opacity: 0 }}
                                  animate={{ opacity: 1 }}
                                  className="flex items-center justify-center h-24 rounded-lg border-2 border-dashed border-muted-foreground/20"
                                >
                                  <p className="text-xs text-muted-foreground">{t('board.moveTo')}</p>
                                </motion.div>
                              ) : (
                                <div className="space-y-2">
                                  {sortedTasks.map((task) => (
                                    <SortableTaskCard key={task.id} task={task} />
                                  ))}
                                </div>
                              )}
                            </AnimatePresence>
                          </SortableContext>
                        </ScrollArea>
                      </div>
                    );
                  })}
              </div>

              <DragOverlay>
                {activeTask ? (
                  <div className="rotate-3 opacity-90">
                    <TaskCardContent task={activeTask} />
                  </div>
                ) : null}
              </DragOverlay>
            </DndContext>
          )}
        </Tabs.Content>
        <Tabs.Content value="backlog" className="flex-1 min-h-0 data-[state=inactive]:hidden">
          <div className="flex flex-col items-center justify-center h-64 text-muted-foreground gap-3">
            <LayoutDashboard className="h-8 w-8" />
            <p className="text-sm font-medium">{t('projects.backlog')}</p>
            <p className="text-xs">
              Switch to the Backlog tab or visit the project overview
            </p>
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
