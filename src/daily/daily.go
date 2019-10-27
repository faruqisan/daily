package daily

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	// Report struct represent daily's report
	Report struct {
		ID         int64     `db:"id"`
		UserID     int64     `db:"user_id"`
		Time       time.Time `db:"time"`
		Activities []Activity
	}

	// Activity struct represent daily report's activity
	Activity struct {
		ID       int64  `db:"id"`
		ReportID int64  `db:"report_id"`
		Detail   string `db:"detail"`
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
func (e Engine) Create(ctx context.Context, date time.Time, userID int64, detail string) error {

	var (
		tNow     = time.Now()
		reportID int64
	)

	qGetTodayReport := `
	SELECT id FROM reports WHERE user_id = $1 AND time > $2
	`

	// get current daily from db
	err := e.db.GetContext(ctx, &reportID, qGetTodayReport, userID, tNow.Format("2006-01-02"))
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	tx, err := e.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	if reportID == 0 {
		//create new report
		reportID, err = e.createReport(ctx, tx, userID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = e.createActivity(ctx, tx, reportID, detail)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()

}

func (e Engine) createReport(ctx context.Context, tx *sqlx.Tx, userID int64) (int64, error) {

	var (
		createdID int64
	)

	q := `
		INSERT INTO reports (user_id, time) VALUES ($1, now()) RETURNING id
	`

	res, err := tx.QueryContext(ctx, q, userID)
	if err != nil {
		return createdID, err
	}
	defer res.Close()
	for res.Next() {
		err = res.Scan(&createdID)
		if err != nil {
			log.Println("err scan: ", res, err)
			return createdID, err
		}
	}

	return createdID, err
}

func (e Engine) createActivity(ctx context.Context, tx *sqlx.Tx, reportID int64, detail string) error {
	q := `
	INSERT INTO activities (report_id, detail) VALUES ($1, $2)
	`

	_, err := tx.ExecContext(ctx, q, reportID, detail)
	return err
}
