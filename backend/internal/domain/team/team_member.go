package team

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TeamMember struct {
	TeamID   uuid.UUID `json:"team_id"`
	UserID   uuid.UUID `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
}

func NewTeamMember(teamID, userID uuid.UUID) *TeamMember {
	return &TeamMember{
		TeamID:   teamID,
		UserID:   userID,
		JoinedAt: time.Now(),
	}
}

func (m *TeamMember) Validate() error {
	if m.TeamID == uuid.Nil {
		return errors.New("team is required")
	}
	if m.UserID == uuid.Nil {
		return errors.New("user is required")
	}
	return nil
}
