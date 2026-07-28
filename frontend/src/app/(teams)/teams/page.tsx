'use client';

import { useTeams } from '@/hooks/use-teams';
import { Card, CardContent, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { Building2, Users, AlertCircle, Plus } from 'lucide-react';
import { useTranslation } from 'react-i18next';

export default function TeamsPage() {
  const { t } = useTranslation();
  const orgId = undefined;
  const { data: teams, isLoading, error } = useTeams(orgId);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">{t('teams.title')}</h1>
          <p className="text-muted-foreground mt-1">
            {t('teams.noTeamsDesc')}
          </p>
        </div>
        <Button>
          <Plus className="h-4 w-4 mr-2" />
          {t('teams.create')}
        </Button>
      </div>

      {isLoading ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <Card key={i}>
              <CardHeader>
                <Skeleton className="h-5 w-3/4" />
                <Skeleton className="h-4 w-1/2 mt-2" />
              </CardHeader>
              <CardContent>
                <Skeleton className="h-4 w-1/3" />
              </CardContent>
            </Card>
          ))}
        </div>
      ) : error ? (
        <div className="flex flex-col items-center justify-center py-16 text-center">
          <AlertCircle className="h-12 w-12 text-destructive mb-4" />
          <h3 className="text-lg font-semibold">{t('common.error')}: {t('teams.title')}</h3>
          <p className="text-sm text-muted-foreground mt-1">
            {error instanceof Error ? error.message : 'An unexpected error occurred'}
          </p>
        </div>
      ) : teams && teams.length > 0 ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {teams.map((team) => (
            <Card key={team.id} className="transition-all hover:shadow-md hover:border-primary/50 cursor-pointer">
              <CardHeader>
                <div className="flex items-start justify-between">
                  <CardTitle className="text-lg">{team.name}</CardTitle>
                  <Building2 className="h-5 w-5 text-muted-foreground shrink-0 mt-0.5" />
                </div>
                {team.description && (
                  <p className="text-sm text-muted-foreground line-clamp-2">{team.description}</p>
                )}
              </CardHeader>
              <CardContent>
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Users className="h-4 w-4" />
                  <span>Team {t('common.name')}</span>
                </div>
              </CardContent>
              <CardFooter>
                  <Button variant="ghost" size="sm" className="w-full">
                    {t('common.viewAll')}
                  </Button>
              </CardFooter>
            </Card>
          ))}
        </div>
      ) : (
        <div className="flex flex-col items-center justify-center py-16 text-center">
          <Building2 className="h-12 w-12 text-muted-foreground mb-4" />
          <h3 className="text-lg font-semibold">{t('teams.noTeams')}</h3>
          <p className="text-sm text-muted-foreground mt-1 mb-4">
            {t('teams.noTeamsDesc')}
          </p>
          <Button>
            <Plus className="h-4 w-4 mr-2" />
            {t('teams.create')}
          </Button>
        </div>
      )}
    </div>
  );
}
