package sqlstorage

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage/utilstorage"
)

type Storage struct {
	db *sqlx.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	var err error
	s.db, err = sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		return err
	}
	err = s.db.PingContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to db: %v", dsn)
	}
	return nil
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CreateEvent(ev storage.Event) (id string, err error) {
	dbConvertedEvent, err := convertToDBEvent(ev)
	if err != nil {
		return "", err
	}
	tx := s.db.MustBegin()
	result, err := tx.Exec(`INSERT INTO public.events 
	(
		id, title, userid, datetimestart, tilldate, alarmdatetime
	) 
	VALUES ($1, $2, $3, $4, $5, $6)`,
		dbConvertedEvent.id,
		dbConvertedEvent.title,
		dbConvertedEvent.userid,
		dbConvertedEvent.datetimestart,
		dbConvertedEvent.tilldate,
		dbConvertedEvent.alarmdatetime,
	)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	numrows, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if numrows != 1 {
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return id, nil
}

func (s *Storage) UpdateEvent(id string, event storage.Event) error {
	if event.ID == "" {
		event.ID = id
	}
	if id != event.ID {
		return storage.ErrMismatchedID
	}
	dbConvertedEvent, err := convertToDBEvent(event)
	if err != nil {
		return err
	}
	_, err = s.db.NamedExec(`UPDATE public.events SET 
	(
		title=:title,
		userid=:userid,
		datetimestart=:datetimestart, 
		tilldate=:tilldate,
		alarmdatetime=:alarmdatetime
	)
	 WHERE id = :id`,
		&dbConvertedEvent)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteEvent(id string) error {
	_, err := s.db.Exec(`DELETE FROM public.events WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) ListEventsOnDay(userID string, dateTime time.Time) (events []storage.Event, err error) {
	return s.ListEventsByDuration(userID, dateTime, 24*time.Hour)
}

func (s *Storage) ListEventsOnCurrentWeek(userID string, dateTime time.Time) (events []storage.Event, err error) {
	weekDayNum := dateTime.Weekday()
	var startOfGivenWeek time.Time
	switch weekDayNum {
	case time.Sunday:
		startOfGivenWeek = dateTime.AddDate(0, 0, -6)
	case time.Monday:
		startOfGivenWeek = dateTime
	case time.Tuesday:
		startOfGivenWeek = dateTime.AddDate(0, 0, -1)
	case time.Wednesday:
		startOfGivenWeek = dateTime.AddDate(0, 0, -2)
	case time.Thursday:
		startOfGivenWeek = dateTime.AddDate(0, 0, -3)
	case time.Friday:
		startOfGivenWeek = dateTime.AddDate(0, 0, -4)
	case time.Saturday:
		startOfGivenWeek = dateTime.AddDate(0, 0, -5)
	}
	endOfGivenWeek := startOfGivenWeek.AddDate(0, 0, 7)
	durationOfCurrentWeek := endOfGivenWeek.Sub(dateTime)
	return s.ListEventsByDuration(userID, dateTime, durationOfCurrentWeek)
}

func (s *Storage) ListEventsOnCurrentMonth(userID string, dateTime time.Time) (events []storage.Event, err error) {
	year := dateTime.Year()
	month := dateTime.Month()
	startMonth := time.Date(year, month, 1, 0, 0, 0, 0, dateTime.Location())
	endOfGivenMonth := startMonth.AddDate(0, 1, 0)
	durationOfCurrentMonth := endOfGivenMonth.Sub(dateTime)
	return s.ListEventsByDuration(userID, dateTime, durationOfCurrentMonth)
}

func (s *Storage) ListEventsByDuration(
	userID string,
	dateTime time.Time,
	duration time.Duration,
) (events []storage.Event, err error) {
	rows, err := s.db.Query(`SELECT 
	(
	id, title, userid, datetimestart, tilldate, alarmdatetime
	)
	FROM public.events 
	WHERE 
	(
		datetimestart BETWEEN  $1 and $2
		AND
		userid = $3
	)`, dateTime, dateTime.Add(duration), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id            string
			title         string
			userid        string
			datetimestart time.Time
			tilldate      time.Time
			alarmdatetime time.Time
		)
		if err := rows.Scan(
			&id, &title, &userid,
			&datetimestart, &tilldate, &alarmdatetime); err != nil {
			if err != nil {
				return nil, err
			}
		}
		events = append(events,
			convertToStorageEvent(pgEvent{
				id:            id,
				title:         title,
				userid:        userid,
				datetimestart: datetimestart,
				tilldate:      tilldate,
				alarmdatetime: alarmdatetime,
			}))
		fmt.Printf("%s %s\n", id, title)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, ": rows read error listeventsbyduration")
	}
	return events, nil
}

func convertToDBEvent(ev storage.Event) (pgEvent, error) {
	dbEvent := pgEvent{}
	if ev.ID == "" {
		ev.ID = utilstorage.GenerateUUID()
	}
	dbEvent.id = ev.ID
	if ev.UserID == "" {
		return pgEvent{}, storage.ErrEmptyUserIDField
	}
	dbEvent.title = ev.Title
	dbEvent.userid = ev.UserID
	dbEvent.datetimestart = ev.DateTimeStart
	dbEvent.tilldate = ev.DateTimeStart.Add(ev.Duration)
	dbEvent.alarmdatetime = ev.DateTimeStart.Add(-1 * ev.Duration)
	return dbEvent, nil
}

func convertToStorageEvent(pgEv pgEvent) storage.Event {
	storageEvent := storage.Event{}
	storageEvent.ID = pgEv.id

	storageEvent.Title = pgEv.title
	storageEvent.UserID = pgEv.userid
	storageEvent.DateTimeStart = pgEv.datetimestart
	storageEvent.Duration = pgEv.tilldate.Sub(pgEv.datetimestart)
	storageEvent.AlarmTime = pgEv.datetimestart.Sub(pgEv.alarmdatetime)
	return storageEvent
}

type pgEvent struct {
	id            string
	title         string
	userid        string
	datetimestart time.Time
	tilldate      time.Time
	alarmdatetime time.Time
}
