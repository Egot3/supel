package subject

import (
	"context"

	"github.com/egot3/supel/backend/timetable/internal/models"
	"github.com/google/uuid"
)

type SubjectRepository interface {
	Subject(ctx context.Context, subjectUUID uuid.UUID) (*models.Subject, error)
	CreateSubject(ctx context.Context, name string) error
	PatchSubject(ctx context.Context, patched models.SubjectPatched) error
	DeleteSubject(ctx context.Context, subjectUUID uuid.UUID) error
	ListSubjects(ctx context.Context, page, size uint32) ([]models.Subject, int, error)
}
