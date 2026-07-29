package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/devflow/devflow-backend/internal/domain/board"
	"github.com/devflow/devflow-backend/internal/domain/organization"
	"github.com/devflow/devflow-backend/internal/domain/project"
	"github.com/devflow/devflow-backend/internal/domain/task"
	"github.com/devflow/devflow-backend/internal/domain/user"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/devflow?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	var userID uuid.UUID
	err = pool.QueryRow(ctx, `SELECT id FROM users WHERE email = $1 AND deleted_at IS NULL`, "admin@devflow.local").Scan(&userID)
	if err != nil {
		hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		u := user.NewUser("admin@devflow.local", string(hash), "Admin User")
		now := time.Now()
		u.EmailVerifiedAt = &now
		if err := insertUser(ctx, pool, u); err != nil {
			log.Fatalf("create user: %v", err)
		}
		userID = u.ID
		fmt.Printf("User created: %s\n", userID)
	} else {
		fmt.Printf("User already exists: %s\n", userID)
	}

	var orgID uuid.UUID
	err = pool.QueryRow(ctx, `SELECT id FROM organizations WHERE slug = $1 AND deleted_at IS NULL`, "devflow-demo").Scan(&orgID)
	if err != nil {
		org := organization.NewOrganization("DevFlow Demo", "devflow-demo", "Demo organization for DevFlow", userID)
		if err := insertOrg(ctx, pool, org); err != nil {
			log.Fatalf("create org: %v", err)
		}
		orgID = org.ID
		fmt.Printf("Org created: %s\n", orgID)
	} else {
		fmt.Printf("Org already exists: %s\n", orgID)
	}

	var exists bool
	pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM organization_members WHERE organization_id = $1 AND user_id = $2)`, orgID, userID).Scan(&exists)
	if !exists {
		member, _ := organization.NewOrganizationMember(orgID, userID, organization.RoleOwner)
		if err := insertOrgMember(ctx, pool, member); err != nil {
			log.Fatalf("create org member: %v", err)
		}
		fmt.Println("Org member added")
	} else {
		fmt.Println("Org member already exists")
	}

	var projectID uuid.UUID
	err = pool.QueryRow(ctx, `SELECT id FROM projects WHERE organization_id = $1 AND key = $2 AND deleted_at IS NULL`, orgID, "DEMO").Scan(&projectID)
	if err != nil {
		proj := project.NewProject("My First Project", "DEMO", "A sample project to explore DevFlow", orgID, &userID)
		if err := insertProject(ctx, pool, proj); err != nil {
			log.Fatalf("create project: %v", err)
		}
		projectID = proj.ID
		fmt.Printf("Project created: %s (DEMO)\n", projectID)
	} else {
		fmt.Printf("Project already exists: %s (DEMO)\n", projectID)
	}

	var boardID uuid.UUID
	err = pool.QueryRow(ctx, `SELECT id FROM boards WHERE project_id = $1`, projectID).Scan(&boardID)
	if err != nil {
		b := board.NewBoard(projectID, "Kanban")
		if err := insertBoard(ctx, pool, b); err != nil {
			log.Fatalf("create board: %v", err)
		}
		boardID = b.ID

		columnNames := []string{"Backlog", "To Do", "In Progress", "In Review", "Done"}
		for i, name := range columnNames {
			col := board.NewColumn(boardID, name, float64(i+1))
			if err := insertColumn(ctx, pool, col); err != nil {
				log.Fatalf("create column %q: %v", name, err)
			}
		}
		fmt.Println("Board and columns created")
	} else {
		fmt.Println("Board already exists")
	}

	var taskCount int
	pool.QueryRow(ctx, `SELECT COUNT(*) FROM tasks WHERE project_id = $1 AND deleted_at IS NULL`, projectID).Scan(&taskCount)
	if taskCount == 0 {
		var columns []uuid.UUID
		rows, _ := pool.Query(ctx, `SELECT id FROM columns WHERE board_id = $1 ORDER BY position`, boardID)
		for rows.Next() {
			var colID uuid.UUID
			rows.Scan(&colID)
			columns = append(columns, colID)
		}
		rows.Close()

		tasks := []struct {
			title       string
			description string
			taskType    task.Type
			priority    task.Priority
			status      task.Status
			columnIdx   int
		}{
			{"Set up CI/CD pipeline", "Configure GitHub Actions for automated testing and deployment", task.TypeTask, task.PriorityHigh, task.StatusInProgress, 2},
			{"Design database schema", "Create ERD and implement migrations for core entities", task.TypeTask, task.PriorityCritical, task.StatusDone, 4},
			{"Implement user authentication", "Add login, registration, and JWT token management", task.TypeStory, task.PriorityHigh, task.StatusInReview, 3},
			{"Create project dashboard", "Build the main dashboard with widgets and metrics", task.TypeTask, task.PriorityMedium, task.StatusTodo, 1},
			{"Fix navigation bug on mobile", "Sidebar overlaps content on small screens", task.TypeBug, task.PriorityHigh, task.StatusTodo, 1},
			{"Write API documentation", "Document all REST endpoints with examples", task.TypeTask, task.PriorityLow, task.StatusBacklog, 0},
			{"Add dark mode support", "Implement theme switching with next-themes", task.TypeStory, task.PriorityMedium, task.StatusBacklog, 0},
			{"Performance optimization", "Profile and optimize database queries and API responses", task.TypeEpic, task.PriorityMedium, task.StatusBacklog, 0},
		}

		for _, t := range tasks {
			tsk := task.NewTask(projectID, userID, t.title, t.description, t.taskType)
			tsk.Priority = t.priority
			_ = tsk.ChangeStatus(t.status)
			tsk.ColumnID = &columns[t.columnIdx]
			tsk.Position = float64(t.columnIdx + 1)
			tsk.AssigneeID = &userID
			if err := insertTask(ctx, pool, tsk); err != nil {
				log.Fatalf("create task %q: %v", t.title, err)
			}
			fmt.Printf("  Task: %s\n", t.title)
		}
	} else {
		fmt.Printf("Tasks already exist (%d)\n", taskCount)
	}

	fmt.Println("\nSeed completed!")
	fmt.Println("Login with: admin@devflow.local / password123")
}

func insertUser(ctx context.Context, pool *pgxpool.Pool, u *user.User) error {
	_, err := pool.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, display_name, avatar_url, email_verified_at, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		u.ID, u.Email, u.PasswordHash, u.DisplayName, u.AvatarURL, u.EmailVerifiedAt, u.CreatedAt, u.UpdatedAt)
	return err
}

func insertOrg(ctx context.Context, pool *pgxpool.Pool, org *organization.Organization) error {
	_, err := pool.Exec(ctx,
		`INSERT INTO organizations (id, name, slug, description, owner_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		org.ID, org.Name, org.Slug, org.Description, org.OwnerID, org.CreatedAt, org.UpdatedAt)
	return err
}

func insertOrgMember(ctx context.Context, pool *pgxpool.Pool, m *organization.OrganizationMember) error {
	_, err := pool.Exec(ctx,
		`INSERT INTO organization_members (organization_id, user_id, role, joined_at)
		 VALUES ($1, $2, $3, $4)`,
		m.OrganizationID, m.UserID, m.Role, m.JoinedAt)
	return err
}

func insertProject(ctx context.Context, pool *pgxpool.Pool, p *project.Project) error {
	_, err := pool.Exec(ctx,
		`INSERT INTO projects (id, name, key, description, organization_id, lead_id, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		p.ID, p.Name, p.Key, p.Description, p.OrganizationID, p.LeadID, p.Status, p.CreatedAt, p.UpdatedAt)
	return err
}

func insertBoard(ctx context.Context, pool *pgxpool.Pool, b *board.Board) error {
	_, err := pool.Exec(ctx,
		`INSERT INTO boards (id, project_id, name, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		b.ID, b.ProjectID, b.Name, b.CreatedAt, b.UpdatedAt)
	return err
}

func insertColumn(ctx context.Context, pool *pgxpool.Pool, c *board.Column) error {
	_, err := pool.Exec(ctx,
		`INSERT INTO columns (id, board_id, name, position, wip_limit, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		c.ID, c.BoardID, c.Name, c.Position, c.WIPLimit, c.CreatedAt, c.UpdatedAt)
	return err
}

func insertTask(ctx context.Context, pool *pgxpool.Pool, t *task.Task) error {
	_, err := pool.Exec(ctx,
		`INSERT INTO tasks (id, project_id, column_id, title, description, type, priority, status, position, assignee_id, reporter_id, version_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		t.ID, t.ProjectID, t.ColumnID, t.Title, t.Description, t.Type, t.Priority, t.Status, t.Position, t.AssigneeID, t.ReporterID, t.VersionID, t.CreatedAt, t.UpdatedAt)
	return err
}
