'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useState } from 'react';
import {
  LayoutDashboard,
  ListTodo,
  Building2,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  LogOut,
  User as UserIcon,
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import { useAuth } from '@/hooks/use-auth';
import { useTranslation } from 'react-i18next';

const navItems = [
  { href: '/dashboard', key: 'nav.dashboard', label: 'Dashboard', icon: LayoutDashboard },
  { href: '/tasks', key: 'nav.tasks', label: 'My Tasks', icon: ListTodo },
  { href: '/orgs', key: 'nav.organizations', label: 'Organizations', icon: Building2 },
];

export function Sidebar() {
  const pathname = usePathname();
  const { user, logout } = useAuth();
  const { t } = useTranslation();
  const [collapsed, setCollapsed] = useState(false);
  const [orgExpanded, setOrgExpanded] = useState(true);

  const initials = user?.displayName
    ?.split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);

  return (
    <aside
      className={cn(
        'flex flex-col border-r bg-card transition-all duration-300',
        collapsed ? 'w-16' : 'w-64',
      )}
    >
      <div className="flex h-14 items-center border-b px-4">
        <Link href="/dashboard" className={cn('flex items-center gap-2 font-bold', collapsed && 'justify-center')}>
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground text-sm font-bold">
            D
          </div>
          {!collapsed && <span>DevFlow</span>}
        </Link>
        <Button
          variant="ghost"
          size="icon"
          className="ml-auto h-8 w-8"
          onClick={() => setCollapsed(!collapsed)}
        >
          {collapsed ? <ChevronRight className="h-4 w-4" /> : <ChevronLeft className="h-4 w-4" />}
        </Button>
      </div>

      <ScrollArea className="flex-1 px-3 py-2">
        <nav className="flex flex-col gap-1">
          {navItems.map((item) => {
            const isActive = pathname === item.href || pathname.startsWith(item.href + '/');
            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  'flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground',
                  isActive
                    ? 'bg-accent text-accent-foreground'
                    : 'text-muted-foreground',
                  collapsed && 'justify-center px-2',
                )}
                title={collapsed ? t(item.key) : undefined}
              >
                <item.icon className="h-4 w-4 shrink-0" />
                {!collapsed && <span>{t(item.key)}</span>}
              </Link>
            );
          })}
        </nav>

        {!collapsed && (
          <>
            <Separator className="my-3" />
            <button
              onClick={() => setOrgExpanded(!orgExpanded)}
              className="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-accent hover:text-accent-foreground"
            >
              <Building2 className="h-4 w-4" />
              <span className="flex-1 text-left">{t('nav.organizations')}</span>
              <ChevronDown
                className={cn('h-4 w-4 transition-transform', orgExpanded && 'rotate-180')}
              />
            </button>
            {orgExpanded && (
              <div className="ml-4 mt-1 space-y-1">
                <p className="px-3 py-2 text-xs text-muted-foreground">{t('organizations.noOrgs')}</p>
              </div>
            )}
          </>
        )}
      </ScrollArea>

      <div className="border-t p-3">
        {user ? (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                className={cn('w-full gap-2', collapsed ? 'justify-center px-2' : 'justify-start')}
              >
                <Avatar className="h-8 w-8">
                  <AvatarImage src={user.avatarUrl ?? undefined} />
                  <AvatarFallback className="text-xs">{initials}</AvatarFallback>
                </Avatar>
                {!collapsed && (
                  <div className="flex-1 text-left">
                    <p className="text-sm font-medium leading-none">{user.displayName}</p>
                    <p className="text-xs text-muted-foreground">{user.email}</p>
                  </div>
                )}
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" side="right" className="w-56">
              <DropdownMenuItem asChild>
                <Link href="/profile">
                  <UserIcon className="mr-2 h-4 w-4" />
                  {t('nav.profile')}
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={logout}>
                <LogOut className="mr-2 h-4 w-4" />
                {t('nav.logout')}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        ) : (
          <div className="h-10" />
        )}
      </div>
    </aside>
  );
}
