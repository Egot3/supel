package subject

import (
	"context"

	"github.com/egot3/supel/backend/timetable/internal/carefulness"
	"github.com/egot3/supel/backend/timetable/internal/models"
	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunSubjectRepository struct {
	db *bun.DB
}

func NewSubjectRepository(i do.Injector) (SubjectRepository, error) {
	db := do.MustInvoke[*bun.DB](i)
	return &bunSubjectRepository{db: db}, nil
}

func (r *bunSubjectRepository) Subject(ctx context.Context, subjectUUID uuid.UUID) (*models.Subject, error) {
	subject := models.Subject{UUID: subjectUUID}
	err := r.db.NewSelect().Model(&subject).WherePK().WhereAllWithDeleted().Scan(ctx)
	if err != nil {
		return nil, err
	}

	if subject.DeletedAt != nil {
		return nil, carefulness.Gone
	}

	return &subject, nil
}

func (r *bunSubjectRepository) CreateSubject(ctx context.Context, name string) error {
	_, err := r.db.NewInsert().Model(&models.Subject{Name: name}).Exec(ctx)
	return err
}

func (r *bunSubjectRepository) PatchSubject(ctx context.Context, patched models.SubjectPatched) error {
	isUpdated := false

	query := r.db.NewUpdate().Model(&models.Subject{UUID: patched.UUID}).WherePK()
	if patched.Name != nil {
		query = query.Set("name = ?", patched.Name)
		isUpdated = true
	}

	if !isUpdated {
		return nil
	}

	_, err := query.Exec(ctx)
	return err
}

func (r *bunSubjectRepository) DeleteSubject(ctx context.Context, subjectUUID uuid.UUID) error {
	_, err := r.db.NewDelete().Model(&models.Subject{UUID: subjectUUID}).WherePK().Exec(ctx)
	return err
}

func (r *bunSubjectRepository) ListSubjects(ctx context.Context, page, size uint32) ([]models.Subject, int, error) {
	var subjects []models.Subject
	total, err := r.db.NewSelect().Model(&subjects).Limit(int(size)).Offset(int(page * size)).ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	return subjects, total, nil
}
