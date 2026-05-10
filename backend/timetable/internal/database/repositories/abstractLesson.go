package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Egot3/supel/backend/timetable/internal/models"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunAbstractLessonRepository struct {
	db *bun.DB
}

func NewAbstractLessonRepository(i do.Injector) (AbstractLessonRepository, error) {
	db, err := do.Invoke[*bun.DB](i)
	if err != nil {
		return nil, err
	}

	return &bunAbstractLessonRepository{db: db}, err
}

func (r *bunAbstractLessonRepository) CreateAbstractLesson(ctx context.Context, name string) error {
	_, err := r.db.NewInsert().Model(&models.AbstractLesson{Name: name}).Exec(ctx)
	return err
}

func (r *bunAbstractLessonRepository) DeleteAbstractLesson(ctx context.Context, uuid string) error {
	_, err := r.db.NewDelete().Model(&models.AbstractLesson{UUID: uuid}).WherePK().Exec(ctx)
	return err
}

func (r *bunAbstractLessonRepository) GetAbstractLesson(ctx context.Context, uuid string) (*models.AbstractLesson, error) {
	var al models.AbstractLesson
	err := r.db.NewSelect().Model(&models.AbstractLesson{UUID: uuid}).
		WherePK().Scan(ctx, &al)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &al, nil
}

func (r *bunAbstractLessonRepository) PatchAbstractLesson(ctx context.Context, patchedAl models.PatchAbstractLesson) error {
	query := r.db.NewUpdate().Model(&models.AbstractLesson{UUID: patchedAl.UUID}).WherePK()
	if patchedAl.Name != nil {
		query = query.Set("name = ?", patchedAl.Name)
	}

	res, err := query.Exec(ctx)
	if err != nil {
		return err
	}

	if val, _ := res.RowsAffected(); val == 0 {
		return sql.ErrNoRows
	}

	return nil
}
