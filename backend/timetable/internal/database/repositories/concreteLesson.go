package repositories

import (
	"context"
	"database/sql"

	"github.com/Egot3/supel/backend/timetable/internal/models"
	"github.com/samber/do/v2"
	"github.com/uptrace/bun"
)

type bunConcreteLessonRepository struct {
	db *bun.DB
}

func NewConcreteLessonRepository(i do.Injector) (ConcreteLessonRepository, error) {
	db, err := do.Invoke[*bun.DB](i)
	if err != nil {
		return nil, err
	}

	return &bunConcreteLessonRepository{db: db}, nil
}

func (r *bunConcreteLessonRepository) CreateConcreteLesson(ctx context.Context, cl models.ConcreteLesson) error {
	_, err := r.db.NewInsert().Model(&cl).Exec(ctx)
	return err
}

func (r *bunConcreteLessonRepository) DeleteConcreteLesson(ctx context.Context, uuid string) error {
	_, err := r.db.NewDelete().Model(&models.ConcreteLesson{ConcreteUUID: uuid}).WherePK().Exec(ctx)
	return err
}

func (r *bunConcreteLessonRepository) PatchConcreteLesson(ctx context.Context, patchCl models.PatchConcreteLesson) error {
	query := r.db.NewUpdate().Model(&models.ConcreteLesson{ConcreteUUID: patchCl.ConcreteUUID}).WherePK()
	if patchCl.AbstractUUID != nil {
		query = query.Set("abstract_uuid = ?", patchCl.AbstractUUID)
	}
	if patchCl.TeacherUUID != nil {
		query = query.Set("teacher_uuid = ?", patchCl.TeacherUUID)
	}
	if patchCl.GroupUUID != nil {
		query = query.Set("group_uuid = ?", patchCl.GroupUUID)
	}
	if patchCl.PeriodUUID != nil {
		query = query.Set("period_uuid = ?", patchCl.PeriodUUID)
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

func (r *bunConcreteLessonRepository) GetConcreteLesson(ctx context.Context, uuid string) (*models.ConcreteLesson, error) {
	cl := models.ConcreteLesson{ConcreteUUID: uuid}
	err := r.db.NewSelect().Model(&cl).WherePK().Relation("abstract_lessons").Relation("homework_attachments").Relation("periods").Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &cl, nil
}
