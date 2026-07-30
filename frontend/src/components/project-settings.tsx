'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useProject, useUpdateProject, useDeleteProject } from '@/hooks/use-projects';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { SettingsIcon, Trash2, AlertCircle } from 'lucide-react';
import { useTranslation } from 'react-i18next';

export function ProjectSettings({ projectId }: { projectId: string }) {
  const { t } = useTranslation();
  const router = useRouter();
  const { data: project, isLoading, error } = useProject(projectId);
  const updateProject = useUpdateProject();
  const deleteProject = useDeleteProject();

  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [dirty, setDirty] = useState(false);

  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center h-64 text-muted-foreground gap-3">
        <SettingsIcon className="h-8 w-8" />
        <p className="text-sm font-medium">{t('projects.settings')}</p>
        <p className="text-xs">{t('common.loading')}</p>
      </div>
    );
  }

  if (error || !project) {
    return (
      <div className="flex flex-col items-center justify-center h-64 text-muted-foreground gap-2">
        <AlertCircle className="h-8 w-8 text-destructive" />
        <p className="text-sm font-medium">{t('common.error')}</p>
      </div>
    );
  }

  if (!dirty) {
    if (name !== project.name) setName(project.name);
    if (description !== (project.description ?? '')) setDescription(project.description ?? '');
  }

  const handleSave = async () => {
    await updateProject.mutateAsync({ projectId, name, description: description || undefined });
    setDirty(false);
  };

  const handleDelete = async () => {
    if (!confirm(t('projects.deleteConfirm') || 'Delete this project? This cannot be undone.')) return;
    await deleteProject.mutateAsync(projectId);
    router.push('/projects');
  };

  return (
    <div className="space-y-6 max-w-2xl">
      <Card>
        <CardHeader>
          <CardTitle>{t('projects.settings')}</CardTitle>
          <CardDescription>{t('projects.settingsDescription') || 'Manage project settings'}</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">{t('projects.name')}</Label>
            <Input id="name" value={name} onChange={(e) => { setName(e.target.value); setDirty(true); }} />
          </div>
          <div className="space-y-2">
            <Label>{t('projects.key')}</Label>
            <Input value={project.key} disabled className="font-mono uppercase text-muted-foreground" />
            <p className="text-xs text-muted-foreground">{t('projects.keyHint') || 'Key cannot be changed after creation'}</p>
          </div>
          <div className="space-y-2">
            <Label htmlFor="description">{t('projects.description')}</Label>
            <Textarea id="description" value={description} onChange={(e) => { setDescription(e.target.value); setDirty(true); }} rows={3} />
          </div>
          <Button onClick={handleSave} disabled={!dirty || updateProject.isPending}>
            {updateProject.isPending ? t('common.saving') || 'Saving...' : t('common.save') || 'Save'}
          </Button>
        </CardContent>
      </Card>

      <Separator />

      <Card className="border-destructive/20">
        <CardHeader>
          <CardTitle className="text-destructive">{t('projects.dangerZone') || 'Danger Zone'}</CardTitle>
          <CardDescription>{t('projects.dangerZoneDescription') || 'Irreversible actions'}</CardDescription>
        </CardHeader>
        <CardContent>
          <Button variant="destructive" onClick={handleDelete} disabled={deleteProject.isPending}>
            <Trash2 className="mr-2 h-4 w-4" />
            {deleteProject.isPending ? t('common.deleting') || 'Deleting...' : t('projects.delete') || 'Delete Project'}
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
