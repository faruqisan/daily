package daily

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	// Report struct represent daily's report
	Report struct {
		ID        int64     `db:"id" json:"id"`
		UserID    int64     `db:"user_id" json:"user_id"`
		Title     string    `db:"title" json:"title"`
		Detail    string    `db:"detail" json:"detail"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
	}

	// Engine struct define daily engine to access data
	Engine struct {
		db *sqlx.DB
	}
)

// New function return engine with setuped db
func New(db *sqlx.DB) Engine {
	return Engine{
		db: db,
	}
}

// Create ..
func (e Engine) Create(ctx context.Context, userID int64, title, detail string) error {

	q := `
	INSERT INTO reports (user_id, title, detail) VALUES ($1, $2, $3)
	`

	_, err := e.db.ExecContext(ctx, q, userID, title, detail)
	return err
}

// GetUserReports function will return user's daily report for giveng time range
func (e Engine) GetUserReports(ctx context.Context, userID int64, timeStart, timeEnd time.Time) ([]Report, error) {
	var (
		reports []Report
		err     error
	)

	q := `
	SELECT id, user_id, title, detail, created_at FROM reports WHERE user_id = $1 AND created_at BETWEEN $2 AND $3
	`

	err = e.db.SelectContext(ctx, &reports, q, userID, timeStart, timeEnd)
	return reports, err
}

// GetByReportID function will retrun report based on given id
func (e Engine) GetByReportID(ctx context.Context, id int64) (Report, error) {
	var (
		report Report
		err    error
	)

	q := `
	SELECT id, user_id, detail, created_at FROM reports WHERE id = $1
	`

	err = e.db.GetContext(ctx, &report, q, id)
	return report, err
}
