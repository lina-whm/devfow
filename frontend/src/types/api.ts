export interface User {
  id: string;
  email: string;
  name: string;
  avatarUrl: string | null;
  createdAt: string;
  updatedAt: string;
}

export interface Organization {
  id: string;
  name: string;
  slug: string;
  description: string | null;
  ownerId: string;
  createdAt: string;
  updatedAt: string;
}

export interface OrganizationMember {
  id: string;
  organizationId: string;
  userId: string;
  role: 'owner' | 'admin' | 'member';
  user: User;
  joinedAt: string;
}

export interface Team {
  id: string;
  organizationId: string;
  name: string;
  description: string | null;
  leadId: string;
  createdAt: string;
  updatedAt: string;
}

export interface Project {
  id: string;
  organizationId: string;
  teamId: string | null;
  name: string;
  key: string;
  description: string | null;
  leadId: string;
  createdAt: string;
  updatedAt: string;
}

export interface Board {
  id: string;
  projectId: string;
  name: string;
  columns: Column[];
  createdAt: string;
  updatedAt: string;
}

export interface Column {
  id: string;
  boardId: string;
  name: string;
  position: number;
  wipLimit: number | null;
  tasks: Task[];
  createdAt: string;
  updatedAt: string;
}

export type TaskType = 'task' | 'bug' | 'story' | 'epic';
export type TaskPriority = 'none' | 'low' | 'medium' | 'high' | 'urgent';
export type TaskStatus = 'backlog' | 'todo' | 'in_progress' | 'in_review' | 'done' | 'cancelled';

export interface Task {
  id: string;
  projectId: string;
  columnId: string | null;
  sprintId: string | null;
  parentId: string | null;
  title: string;
  description: string | null;
  type: TaskType;
  priority: TaskPriority;
  status: TaskStatus;
  position: number;
  storyPoints: number | null;
  assigneeId: string | null;
  assignee: User | null;
  reporterId: string;
  reporter: User;
  labels: string[];
  tags: Tag[];
  dueDate: string | null;
  estimatedHours: number | null;
  loggedHours: number | null;
  order: number;
  createdAt: string;
  updatedAt: string;
}

export interface Comment {
  id: string;
  taskId: string;
  userId: string;
  user: User;
  content: string;
  createdAt: string;
  updatedAt: string;
}

export interface Tag {
  id: string;
  name: string;
  color: string;
}

export interface Sprint {
  id: string;
  projectId: string;
  name: string;
  goal: string | null;
  status: SprintStatus;
  startDate: string | null;
  endDate: string | null;
  createdAt: string;
  updatedAt: string;
}

export type SprintStatus = 'planning' | 'active' | 'completed';

export interface Notification {
  id: string;
  userId: string;
  type: NotificationType;
  title: string;
  message: string;
  read: boolean;
  entityId: string | null;
  entityType: string | null;
  createdAt: string;
}

export type NotificationType =
  | 'task_assigned'
  | 'task_updated'
  | 'comment_added'
  | 'mention'
  | 'sprint_started'
  | 'sprint_ended'
  | 'invitation';

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface CursorPagination<T> {
  data: T[];
  nextCursor: string | null;
  hasMore: boolean;
}

export interface AuthResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  display_name: string;
  email: string;
  password: string;
}

export interface CreateTaskDTO {
  projectId: string;
  columnId?: string;
  title: string;
  description?: string;
  type: TaskType;
  priority: TaskPriority;
  assigneeId?: string;
  parentId?: string;
  sprintId?: string;
  storyPoints?: number;
  labels?: string[];
  dueDate?: string;
  estimatedHours?: number;
}

export interface UpdateTaskDTO {
  title?: string;
  description?: string;
  type?: TaskType;
  priority?: TaskPriority;
  status?: TaskStatus;
  assigneeId?: string | null;
  sprintId?: string | null;
  storyPoints?: number | null;
  labels?: string[];
  dueDate?: string | null;
  estimatedHours?: number | null;
}

export interface MoveTaskDTO {
  taskId: string;
  columnId: string;
  position: number;
}

export interface APIError {
  status: number;
  message: string;
  errors?: Record<string, string[]>;
}
