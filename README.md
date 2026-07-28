# DevFlow

**Платформа для совместной разработки — Канбан, Спринты, Аналитика, RBAC.**

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)
![Next.js](https://img.shields.io/badge/Next.js-14-000?logo=next.js)
![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis)
![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-3.4-06B6D4?logo=tailwindcss)
![License](https://img.shields.io/badge/License-MIT-green)

---

> **🌐 Запуск локально:** `http://localhost:3000` — фронтенд, `http://localhost:8080/api/v1` — API

## Оглавление

- [Описание проекта](#описание-проекта)
- [Основные возможности](#основные-возможности)
- [Архитектура проекта](#архитектура-проекта)
- [Технологии](#технологии)
- [Требования к окружению](#требования-к-окружению)
- [Установка и настройка](#установка-и-настройка)
- [Переменные окружения](#переменные-окружения)
- [Тестирование](#тестирование)
- [Безопасность](#безопасность)
- [Roadmap](#roadmap)

---

## Описание проекта

### Проблема

Команды разработчиков используют разрозненные инструменты для управления задачами: Trello, Jira, Notion. Нет единого окна для Канбан-досок, спринтов и аналитики. Сложно отслеживать прогресс и загрузку команды.

### Решение

DevFlow — платформа для совместной разработки с drag-and-drop Канбан-досками, спринтами, ролевой моделью и аналитикой. Подходит для команд любого размера. Работает на русском и английском языках.

---

## Основные возможности

- **📋 Канбан-доски** — drag-and-drop с оптимистичными обновлениями, WIP-лимиты
- **📦 Спринты** — планирование, burndown-чарты
- **🔐 RBAC** — роли Owner, Admin, Member с разграничением доступа
- **🔔 Уведомления** — real-time оповещения и упоминания
- **🔍 Поиск** — полнотекстовый поиск с курсорной пагинацией
- **🌙 Тёмная тема** — переключение светлой/тёмной темы
- **🌐 Локализация** — русский и английский языки
- **📊 Аналитика** — метрики, графики, Prometheus + Grafana
- **📝 Аудит** — лог изменений и история задач
- **📱 Адаптивность** — десктоп и мобильные устройства

---

## Архитектура проекта

```
┌──────────────────────────────────────────────────────────┐
│                   Frontend (Next.js)                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                │
│  │ Dashboard│  │ Kanban   │  │ Backlog  │                │
│  │ (виджеты)│  │ (dnd-kit)│  │ (фильтры)│                │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘                │
│       │             │             │                      │
│  ┌────▼─────────────▼─────────────▼──────────────────┐   │
│  │              Next.js API Proxy                    │   │
│  └────────────────────┬──────────────────────────────┘   │
└───────────────────────┼──────────────────────────────────┘
                        │
┌───────────────────────▼──────────────────────────────────┐
│              Backend (Go / Gin) :8080                    │
│  ┌──────────────────────────────────────────────────┐    │
│  │  REST API (/api/v1)                              │    │
│  │  Auth, Orgs, Projects, Tasks, Comments, Teams    │    │
│  │  JWT + Refresh токены, CORS, Rate Limiting       │    │
│  └────────────────────┬─────────────────────────────┘    │
└───────────────────────┼──────────────────────────────────┘
                        │
┌───────────────────────▼──────────────────────────────────┐
│              PostgreSQL 16 :5432                         │
│  ┌──────────────────────────────────────────────────┐    │
│  │  Users, Orgs, Projects, Tasks, Comments,         │    │
│  │  Boards, Columns, Sprints, Tags, Notifications   │    │
│  └──────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────┘
                        │
┌───────────────────────▼──────────────────────────────────┐
│              Redis 7 :6379                               │
│  ┌──────────────────────────────────────────────────┐    │
│  │  Session store, Cache, Rate limiter, Queue       │    │
│  └──────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────┘
```

### Ключевые модули

| Модуль | Путь (backend) | Описание |
|--------|----------------|----------|
| **Auth** | `internal/application/auth/` | Регистрация, логин, JWT, refresh |
| **Orgs** | `internal/application/organization/` | Организации, приглашения, участники |
| **Projects** | `internal/application/project/` | Проекты с уникальным ключом |
| **Tasks** | `internal/application/task/` | Задачи, типы, приоритеты, статусы |
| **Comments** | `internal/application/comment/` | Комментарии к задачам |
| **Teams** | `internal/application/team/` | Команды и их участники |

---

## Технологии

### Frontend
- **Next.js 14** (App Router) — фреймворк
- **React 18** — UI
- **TypeScript 5** — типизация
- **Tailwind CSS 3.4** — стилизация
- **shadcn/ui** (Radix UI) — компоненты
- **TanStack React Query** — управление состоянием сервера
- **Zustand** — клиентское состояние
- **dnd-kit** — drag-and-drop для Канбан-досок
- **Recharts** — графики и аналитика
- **Framer Motion** — анимации
- **React Hook Form** + **Zod** — формы и валидация
- **Lucide React** — иконки
- **i18next** — интернационализация

### Backend
- **Go 1.24** — язык
- **Gin** — HTTP-фреймворк
- **PostgreSQL 16** — основная база данных
- **Redis 7** — сессии, кэш, rate limiter, очередь
- **JWT** — access + refresh токены
- **Goose** — миграции
- **Viper** — конфигурация
- **Prometheus** — метрики
- **Grafana** — дашборды
- **OpenTelemetry (Jaeger)** — трассировка

### Тестирование
- **Go testing** + **Testcontainers** — интеграционные тесты Go
- **Playwright** — E2E-тесты
- **React Testing Library** — тестирование React-компонентов

### Инфраструктура
- **Docker** + **Docker Compose** — контейнеризация
- **GitHub Actions** — CI/CD
- **Nginx** — reverse proxy

---

## Требования к окружению

- **Go** 1.24+
- **Node.js** 18+
- **npm** 9+
- **PostgreSQL** 16 (установлен и запущен)
- **Redis** 7 (установлен и запущен)
- **Git** — для клонирования репозитория

---

## Установка и настройка

### 1. Клонирование

```bash
git clone https://github.com/lina-whm/devfow.git
cd devfow
```

### 2. База данных

Убедитесь, что PostgreSQL запущен:

```bash
# Windows
net start postgresql-x64-16
```

Создайте базу данных:

```bash
createdb -U postgres devflow
```

### 3. Redis

Убедитесь, что Redis запущен:

```bash
# Windows (если установлен через winget)
net start Redis
```

### 4. Миграции

```bash
cd backend
go run ./cmd/devflow-migrate up
```

### 5. Переменные окружения

Скопируйте `.env.example` в `.env`:

```bash
cp .env.example .env
```

Подробнее — в разделе [Переменные окружения](#переменные-окружения).

### 6. Запуск бэкенда

```bash
cd backend
go run ./cmd/devflow
```

Бэкенд будет доступен на `http://localhost:8080`.

### 7. Запуск фронтенда

```bash
cd frontend
npm install
npm run dev
```

Фронтенд будет доступен на `http://localhost:3000`.

---

## Переменные окружения

Создайте файл `.env` в корне проекта (или `.env.local` в `frontend/`):

```env
# === База данных (PostgreSQL) ===
DATABASE_URL="postgres://postgres:postgres@localhost:5432/devflow?sslmode=disable"

# === Redis ===
REDIS_URL="redis://localhost:6379/0"
REDIS_PASSWORD=""

# === JWT ===
JWT_SECRET="[ВАШ_SECRET_МИНИМУМ_32_СИМВОЛА]"
JWT_REFRESH_SECRET="[ВАШ_REFRESH_SECRET_МИНИМУМ_32_СИМВОЛА]"
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=7d

# === Сервер ===
SERVER_PORT=8080
SERVER_MODE=debug

# === Frontend (опционально) ===
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

> **Важно:** Никогда не коммитьте `.env` или `.env.local` в репозиторий. Они уже добавлены в `.gitignore`.

---

## Тестирование

### Backend (Go)

```bash
cd backend
go test ./... -v
```

### Frontend (Vitest + Playwright)

```bash
cd frontend
npm test          # unit-тесты
npm run test:e2e  # E2E-тесты (требуется запущенный dev-сервер)
```

---

## Безопасность

- `.env` и `.env.local` — добавлены в `.gitignore`
- Все чувствительные данные в README заменены на заглушки `[ВАШ_ЗНАЧЕНИЕ]`
- JWT access token с коротким TTL (15 мин) + refresh token (7 дней)
- Rate limiting на уровне API (через Redis)
- CORS ограничен настройками конфигурации
- CI запускает линтеры и проверки на каждый push

---

## Roadmap

- [x] Канбан-доски с drag-and-drop
- [x] Авторизация и JWT
- [x] RBAC (Owner, Admin, Member)
- [x] Полнотекстовый поиск
- [x] Тёмная тема
- [x] Локализация (RU/EN)
- [ ] Спринты и burndown-чарты
- [ ] Real-time уведомления (WebSocket)
- [ ] Импорт/экспорт задач (CSV, JSON)
- [ ] Интеграция с GitHub/GitLab
- [ ] AI-ассистент для описания задач
