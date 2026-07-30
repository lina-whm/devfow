'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useRouter } from 'next/navigation';
import { useCreateProject } from '@/hooks/use-projects';
import { useOrganizations } from '@/hooks/use-organizations';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useTranslation } from 'react-i18next';

const createProjectSchema = z.object({
  name: z.string().min(1, 'Project name is required').max(100, 'Name too long'),
  key: z
    .string()
    .min(2, 'Key must be at least 2 characters')
    .max(10, 'Key too long')
    .regex(/^[A-Z]+$/, 'Key must be uppercase letters only'),
  description: z.string().max(500, 'Description too long').optional(),
});

type CreateProjectFormData = z.infer<typeof createProjectSchema>;

export default function CreateProjectPage() {
  const { t } = useTranslation();
  const router = useRouter();
  const createProject = useCreateProject();
  const { data: orgs, isLoading: orgsLoading } = useOrganizations();

  const {
    register,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<CreateProjectFormData>({
    resolver: zodResolver(createProjectSchema),
    defaultValues: {
      key: '',
      description: '',
    },
  });

  const onSubmit = async (data: CreateProjectFormData) => {
    const orgId = orgs?.[0]?.id;
    if (!orgId) {
      setError('root', { message: 'No organization found' });
      return;
    }

    try {
      const project = await createProject.mutateAsync({
        orgId,
        name: data.name,
        key: data.key,
        description: data.description || undefined,
      });
      router.push(`/projects/${project.id}`);
    } catch (err: unknown) {
      const message =
        err && typeof err === 'object' && 'message' in err
          ? (err as { message: string }).message
          : 'Failed to create project';
      setError('root', { message });
    }
  };

  if (orgsLoading) {
    return (
      <div className="flex items-center justify-center py-20">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    );
  }

  if (!orgs?.length) {
    return (
      <div className="flex items-center justify-center py-20">
        <div className="text-center">
          <p className="text-lg font-medium">No organizations found</p>
          <p className="text-sm text-muted-foreground">You need to be part of an organization to create a project.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-2xl py-8">
      <Card>
        <CardHeader>
          <CardTitle>{t('projects.create') || 'Create Project'}</CardTitle>
          <CardDescription>
            {t('projects.createDescription') || 'Create a new project in your organization'}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="name">{t('projects.name') || 'Project Name'}</Label>
              <Input
                id="name"
                placeholder="My Project"
                {...register('name')}
              />
              {errors.name && (
                <p className="text-sm text-destructive">{errors.name.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="key">{t('projects.key') || 'Project Key'}</Label>
              <Input
                id="key"
                placeholder="PROJ"
                className="uppercase"
                {...register('key')}
              />
              {errors.key && (
                <p className="text-sm text-destructive">{errors.key.message}</p>
              )}
              <p className="text-xs text-muted-foreground">
                {t('projects.keyHint') || 'Used in task identifiers (e.g., PROJ-123)'}
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">{t('projects.description') || 'Description'}</Label>
              <Textarea
                id="description"
                placeholder="Optional project description"
                rows={4}
                {...register('description')}
              />
              {errors.description && (
                <p className="text-sm text-destructive">{errors.description.message}</p>
              )}
            </div>

            {errors.root && (
              <p className="text-sm text-destructive">{errors.root.message}</p>
            )}

            <div className="flex items-center gap-4">
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? 'Creating...' : t('projects.create') || 'Create Project'}
              </Button>
              <Button type="button" variant="outline" onClick={() => router.back()}>
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
