import { create } from 'zustand';
import type { Task, Column } from '@/types/api';

interface DragState {
  activeTaskId: string | null;
  sourceColumnId: string | null;
  targetColumnId: string | null;
}

interface BoardStore {
  selectedTask: Task | null;
  dragState: DragState;
  columns: Column[];
  setSelectedTask: (task: Task | null) => void;
  setDragState: (state: Partial<DragState>) => void;
  setColumns: (columns: Column[]) => void;
  resetDragState: () => void;
}

export const useBoardStore = create<BoardStore>((set) => ({
  selectedTask: null,
  dragState: {
    activeTaskId: null,
    sourceColumnId: null,
    targetColumnId: null,
  },
  columns: [],
  setSelectedTask: (task) => set({ selectedTask: task }),
  setDragState: (state) =>
    set((prev) => ({ dragState: { ...prev.dragState, ...state } })),
  setColumns: (columns) => set({ columns }),
  resetDragState: () =>
    set({
      dragState: {
        activeTaskId: null,
        sourceColumnId: null,
        targetColumnId: null,
      },
    }),
}));
