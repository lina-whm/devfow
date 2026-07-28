'use client';

import { useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useTask, useUpdateTask } from '@/hooks/use-tasks';
import { useComments, useCreateComment } from '@/hooks/use-comments';
import { useAuth } from '@/hooks/use-auth';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Separator } from '@/components/ui/separator';
import { Skeleton } from '@/components/ui/skeleton';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { ArrowLeft, MessageSquare, Calendar, AlertCircle, Send } from 'lucide-react';
import { format } from 'date-fns';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslation } from 'react-i18next';
import type { TaskType, TaskPriority, TaskStatus } from '@/types/api';

const taskFormSchema = z.object({
  title: z.string().min(1, 'Title is required'),
  description: z.string().optional(),
  type: z.enum(['task', 'bug', 'story', 'epic']),
  priority: z.enum(['none', 'low', 'medium', 'high', 'urgent']),
  status: z.enum(['backlog', 'todo', 'in_progress', 'in_review', 'done', 'cancelled']),
});

type TaskFormData = z.infer<typeof taskFormSchema>;

const statusVariant = (status: string) => {
  switch (status) {
    case 'done': return 'success' as const;
    case 'in_progress': return 'default' as const;
    case 'in_review': return 'warning' as const;
    case 'todo': return 'secondary' as const;
    case 'backlog': return 'outline' as const;
    case 'cancelled': return 'destructive' as const;
    default: return 'secondary' as const;
  }
};

const statusLabel = (t: (key: string) => string, status: string) => {
  const labels: Record<string, string> = {
    backlog: 'tasks.statuses.backlog',
    todo: 'tasks.statuses.todo',
    in_progress: 'tasks.statuses.in_progress',
    in_review: 'In Review',
    done: 'tasks.statuses.done',
    cancelled: 'tasks.statuses.cancelled',
  };
  const key = labels[status];
  return key ? t(key) : status.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase());
};

const typeVariant = (type: string) => {
  switch (type) {
    case 'bug': return 'destructive' as const;
    case 'story': return 'success' as const;
    case 'epic': return 'warning' as const;
    default: return 'secondary' as const;
  }
};

export default function TaskDetailPage() {
  const { t } = useTranslation();
  const { taskId } = useParams<{ taskId: string }>();
  const router = useRouter();
  const { user } = useAuth();
  const { data: task, isLoading: taskLoading, error: taskError } = useTask(taskId);
  const { data: comments, isLoading: commentsLoading } = useComments(taskId);
  const updateTask = useUpdateTask();
  const createComment = useCreateComment();

  const [commentText, setCommentText] = useState('');
  const [isEditing, setIsEditing] = useState(false);

  const form = useForm<TaskFormData>({
    resolver: zodResolver(taskFormSchema),
    values: task
      ? {
          title: task.title,
          description: task.description ?? '',
          type: task.type,
          priority: task.priority,
          status: task.status,
        }
      : undefined,
  });

  const onSubmit = async (data: TaskFormData) => {
    if (!taskId) return;
    await updateTask.mutateAsync({
      taskId,
      data: {
        title: data.title,
        description: data.description || undefined,
        type: data.type as TaskType,
        priority: data.priority as TaskPriority,
        status: data.status as TaskStatus,
      },
    });
    setIsEditing(false);
  };

  const handleAddComment = async () => {
    if (!commentText.trim() || !taskId) return;
    await createComment.mutateAsync({ taskId, content: commentText.trim() });
    setCommentText('');
  };

  if (taskLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-32" />
        <Skeleton className="h-10 w-3/4" />
        <Skeleton className="h-40 w-full" />
      </div>
    );
  }

  if (taskError || !task) {
    return (
      <div className="flex flex-col items-center justify-center py-24 text-center">
        <AlertCircle className="h-12 w-12 text-destructive mb-4" />
        <h2 className="text-xl font-semibold">{t('tasks.noTasks')}</h2>
        <p className="text-muted-foreground mt-1 mb-4">
          {taskError instanceof Error ? taskError.message : t('common.error')}
        </p>
        <Button variant="outline" onClick={() => router.back()}>
          <ArrowLeft className="h-4 w-4 mr-2" />
          {t('common.back')}
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6 max-w-4xl">
      <Button variant="ghost" size="sm" onClick={() => router.back()}>
        <ArrowLeft className="h-4 w-4 mr-2" />
        {t('common.back')}
      </Button>

      <Card>
        <CardHeader className="pb-3">
          <div className="flex items-start justify-between gap-4">
            <div className="flex-1 space-y-2">
              <div className="flex items-center gap-2 flex-wrap">
                <Badge variant={typeVariant(task.type)} className="capitalize">
                  {task.type}
                </Badge>
                <Badge variant={statusVariant(task.status)}>
                  {statusLabel(t, task.status)}
                </Badge>
              </div>
              {isEditing ? (
                <Input {...form.register('title')} className="text-xl font-bold h-auto py-1" />
              ) : (
                <h1
                  className="text-2xl font-bold cursor-pointer hover:text-primary transition-colors"
                  onClick={() => setIsEditing(true)}
                  title={t('common.edit')}
                >
                  {task.title}
                </h1>
              )}
            </div>
            {isEditing && (
              <div className="flex items-center gap-2 shrink-0">
                <Button size="sm" variant="outline" onClick={() => setIsEditing(false)}>
                  {t('common.cancel')}
                </Button>
                <Button size="sm" onClick={form.handleSubmit(onSubmit)} disabled={updateTask.isPending}>
                  {updateTask.isPending ? t('common.loading') : t('common.save')}
                </Button>
              </div>
            )}
          </div>
        </CardHeader>
        <CardContent className="space-y-6">
          {isEditing && form.formState.errors.title && (
            <p className="text-sm text-destructive">{form.formState.errors.title.message}</p>
          )}

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-4">
              <div>
                <h4 className="text-sm font-medium text-muted-foreground mb-1">{t('common.description')}</h4>
                {isEditing ? (
                  <Textarea
                    {...form.register('description')}
                    placeholder={t('common.description')}
                    className="min-h-[120px]"
                  />
                ) : (
                  <p className="text-sm whitespace-pre-wrap">
                    {task.description || (
                      <span
                        className="text-muted-foreground italic cursor-pointer"
                        onClick={() => setIsEditing(true)}
                      >
                        {t('common.description')}...
                      </span>
                    )}
                  </p>
                )}
              </div>

              <Separator />

              <div className="space-y-3">
                <h4 className="text-sm font-medium text-muted-foreground">{t('tasks.details')}</h4>
                {isEditing ? (
                  <div className="grid grid-cols-2 gap-3">
                    <div className="space-y-1.5">
                      <label className="text-xs text-muted-foreground">{t('common.type')}</label>
                      <Select
                        value={form.watch('type')}
                        onValueChange={(v) => form.setValue('type', v as TaskType)}
                      >
                        <SelectTrigger><SelectValue /></SelectTrigger>
                        <SelectContent>
                          <SelectItem value="task">{t('tasks.types.task')}</SelectItem>
                          <SelectItem value="bug">{t('tasks.types.bug')}</SelectItem>
                          <SelectItem value="story">{t('tasks.types.story')}</SelectItem>
                          <SelectItem value="epic">{t('tasks.types.epic')}</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="space-y-1.5">
                      <label className="text-xs text-muted-foreground">{t('common.priority')}</label>
                      <Select
                        value={form.watch('priority')}
                        onValueChange={(v) => form.setValue('priority', v as TaskPriority)}
                      >
                        <SelectTrigger><SelectValue /></SelectTrigger>
                        <SelectContent>
                          <SelectItem value="none">{t('tasks.priorities.none')}</SelectItem>
                          <SelectItem value="low">{t('tasks.priorities.low')}</SelectItem>
                          <SelectItem value="medium">{t('tasks.priorities.medium')}</SelectItem>
                          <SelectItem value="high">{t('tasks.priorities.high')}</SelectItem>
                          <SelectItem value="urgent">{t('tasks.priorities.urgent')}</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="space-y-1.5 col-span-2">
                      <label className="text-xs text-muted-foreground">{t('common.status')}</label>
                      <Select
                        value={form.watch('status')}
                        onValueChange={(v) => form.setValue('status', v as TaskStatus)}
                      >
                        <SelectTrigger><SelectValue /></SelectTrigger>
                        <SelectContent>
                          <SelectItem value="backlog">{t('tasks.statuses.backlog')}</SelectItem>
                          <SelectItem value="todo">{t('tasks.statuses.todo')}</SelectItem>
                          <SelectItem value="in_progress">{t('tasks.statuses.in_progress')}</SelectItem>
                          <SelectItem value="in_review">In Review</SelectItem>
                          <SelectItem value="done">{t('tasks.statuses.done')}</SelectItem>
                          <SelectItem value="cancelled">{t('tasks.statuses.cancelled')}</SelectItem>
                        </SelectContent>
                      </Select>
                    </div>
                  </div>
                ) : (
                  <div className="space-y-2 text-sm">
                    <div className="flex items-center justify-between py-1">
                      <span className="text-muted-foreground">{t('common.type')}</span>
                      <Badge variant={typeVariant(task.type)} className="capitalize">
                        {t('tasks.types.' + task.type) || task.type}
                      </Badge>
                    </div>
                    <div className="flex items-center justify-between py-1">
                      <span className="text-muted-foreground">{t('common.priority')}</span>
                      <Badge variant={task.priority === 'urgent' || task.priority === 'high' ? 'destructive' : task.priority === 'medium' ? 'warning' : 'secondary'} className="capitalize">
                        {task.priority}
                      </Badge>
                    </div>
                    <div className="flex items-center justify-between py-1">
                      <span className="text-muted-foreground">{t('common.status')}</span>
                      <Badge variant={statusVariant(task.status)}>
                        {statusLabel(t, task.status)}
                      </Badge>
                    </div>
                  </div>
                )}
              </div>
            </div>

            <div className="space-y-4">
              <div>
                <h4 className="text-sm font-medium text-muted-foreground mb-2">{t('tasks.assignee')}</h4>
                {task.assignee ? (
                  <div className="flex items-center gap-3">
                    <Avatar className="h-8 w-8">
                      <AvatarImage src={task.assignee.avatarUrl ?? undefined} />
                      <AvatarFallback className="text-xs">
                        {task.assignee.name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)}
                      </AvatarFallback>
                    </Avatar>
                    <div>
                      <p className="text-sm font-medium">{task.assignee.name}</p>
                      <p className="text-xs text-muted-foreground">{task.assignee.email}</p>
                    </div>
                  </div>
                ) : (
                  <p className="text-sm text-muted-foreground italic">{t('tasks.assignee')}: -</p>
                )}
              </div>

              <Separator />

              <div>
                <h4 className="text-sm font-medium text-muted-foreground mb-2">Reporter</h4>
                <div className="flex items-center gap-3">
                  <Avatar className="h-8 w-8">
                    <AvatarImage src={task.reporter.avatarUrl ?? undefined} />
                    <AvatarFallback className="text-xs">
                      {task.reporter.name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <p className="text-sm font-medium">{task.reporter.name}</p>
                    <p className="text-xs text-muted-foreground">{task.reporter.email}</p>
                  </div>
                </div>
              </div>

              <Separator />

              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm">
                  <Calendar className="h-4 w-4 text-muted-foreground" />
                  <span className="text-muted-foreground">{t('common.create')}:</span>
                  <span>{format(new Date(task.createdAt), 'MMM d, yyyy HH:mm')}</span>
                </div>
                <div className="flex items-center gap-2 text-sm">
                  <Calendar className="h-4 w-4 text-muted-foreground" />
                  <span className="text-muted-foreground">{t('common.update')}:</span>
                  <span>{format(new Date(task.updatedAt), 'MMM d, yyyy HH:mm')}</span>
                </div>
                {task.dueDate && (
                  <div className="flex items-center gap-2 text-sm">
                    <Calendar className="h-4 w-4 text-muted-foreground" />
                    <span className="text-muted-foreground">{t('tasks.dueDate')}:</span>
                    <span>{format(new Date(task.dueDate), 'MMM d, yyyy')}</span>
                  </div>
                )}
              </div>

              {task.labels.length > 0 && (
                <>
                  <Separator />
                  <div>
                    <h4 className="text-sm font-medium text-muted-foreground mb-2">Labels</h4>
                    <div className="flex flex-wrap gap-1.5">
                      {task.labels.map((label) => (
                        <Badge key={label} variant="outline" className="text-xs">
                          {label}
                        </Badge>
                      ))}
                    </div>
                  </div>
                </>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-lg">
            <MessageSquare className="h-5 w-5" />
            {t('comments.title')} ({comments?.length ?? 0})
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {commentsLoading ? (
            <div className="space-y-4">
              {Array.from({ length: 3 }).map((_, i) => (
                <div key={i} className="flex gap-3">
                  <Skeleton className="h-8 w-8 rounded-full shrink-0" />
                  <div className="flex-1 space-y-2">
                    <Skeleton className="h-4 w-32" />
                    <Skeleton className="h-12 w-full" />
                  </div>
                </div>
              ))}
            </div>
          ) : comments && comments.length > 0 ? (
            comments.map((comment) => (
              <div key={comment.id} className="flex gap-3">
                <Avatar className="h-8 w-8 shrink-0">
                  <AvatarImage src={comment.user.avatarUrl ?? undefined} />
                  <AvatarFallback className="text-xs">
                    {comment.user.name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)}
                  </AvatarFallback>
                </Avatar>
                <div className="flex-1 space-y-1">
                  <div className="flex items-center gap-2">
                    <span className="text-sm font-medium">{comment.user.name}</span>
                    <span className="text-xs text-muted-foreground">
                      {format(new Date(comment.createdAt), 'MMM d, yyyy HH:mm')}
                    </span>
                  </div>
                  <p className="text-sm whitespace-pre-wrap text-muted-foreground">
                    {comment.content}
                  </p>
                </div>
              </div>
            ))
          ) : (
            <div className="flex flex-col items-center justify-center py-8 text-center">
              <MessageSquare className="h-8 w-8 text-muted-foreground mb-2" />
              <p className="text-sm text-muted-foreground">{t('comments.noComments')}</p>
            </div>
          )}

          <Separator />

          <div className="flex gap-3">
            <Avatar className="h-8 w-8 shrink-0">
              <AvatarImage src={user?.avatarUrl ?? undefined} />
              <AvatarFallback className="text-xs">
                {user?.name?.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2) ?? 'U'}
              </AvatarFallback>
            </Avatar>
            <div className="flex-1 flex gap-2">
              <Textarea
                placeholder={t('comments.placeholder')}
                value={commentText}
                onChange={(e) => setCommentText(e.target.value)}
                className="min-h-[40px] text-sm"
                onKeyDown={(e) => {
                  if (e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault();
                    handleAddComment();
                  }
                }}
              />
              <Button
                size="icon"
                onClick={handleAddComment}
                disabled={!commentText.trim() || createComment.isPending}
                className="shrink-0 self-end"
              >
                <Send className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
