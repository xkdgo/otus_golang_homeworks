package sqlstorage

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib" //nolint
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
	if err := s.db.Close(); err != nil {
		return errors.Wrap(err, ":sql storage close error")
	}
	return nil
}

func (s *Storage) CreateEvent(ev storage.Event) (id string, err error) {
	dbConvertedEvent, err := convertToDBEvent(ev)
	if err != nil {
		return "", err
	}
	if ev.ID == "" {
		ev.ID = utilstorage.GenerateUUID()
	}
	if ev.UserID == "" {
		return "", storage.ErrEmptyUserIDField
	}
	tx := s.db.MustBegin()
	result, err := tx.Exec(`INSERT INTO public.events 
	(
		id, title, userid, datetimestart, tilldate, alarmdatetime
	) 
	VALUES ($1, $2, $3, $4, $5, $6)`,
		dbConvertedEvent.ID,
		dbConvertedEvent.Title,
		dbConvertedEvent.Userid,
		dbConvertedEvent.Datetimestart,
		dbConvertedEvent.Tilldate,
		dbConvertedEvent.Alarmdatetime,
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
	return ev.ID, nil
}

func (s *Storage) UpdateEvent(id string, event storage.Event) error {
	if id == "" {
		return storage.ErrEventIDNotFound
	}
	if event.ID == "" {
		event.ID = id
	}
	if id != event.ID {
		return storage.ErrMismatchedID
	}
	if err := s.checkEventIDisPresent(id); err != nil {
		return err
	}
	dbConvertedEvent, err := convertToDBEvent(event)
	if err != nil {
		return err
	}
	_, err = s.db.NamedExec(`UPDATE public.events SET 
	
		title=:title,
		userid=:userid,
		datetimestart=:datetimestart, 
		tilldate=:tilldate,
		alarmdatetime=:alarmdatetime
	
	 WHERE id = :id`,
		&dbConvertedEvent)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteEvent(id string) error {
	_, err := s.db.Exec(`DELETE FROM public.events WHERE id =$1`, id)
	if err != nil {
		return errors.Wrap(err, ":delete event error")
	}
	return nil
}

func (s *Storage) ResetAllData() error {
	_, err := s.db.Exec(`DELETE FROM public.events`)
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
	id, title, userid, datetimestart, tilldate, alarmdatetime
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
				return nil, errors.Wrap(err, ": rows next error listeventsbyduration")
			}
		}
		events = append(events,
			convertToStorageEvent(pgEvent{
				ID:            id,
				Title:         title,
				Userid:        userid,
				Datetimestart: datetimestart,
				Tilldate:      tilldate,
				Alarmdatetime: alarmdatetime,
			}))
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, ": rows read error listeventsbyduration")
	}
	if len(events) == 0 {
		err = s.checkUserIDisPresent(userID)
		if err != nil {
			return nil, err
		}
	}
	return events, nil
}

func (s *Storage) checkUserIDisPresent(userID string) error {
	rows, err := s.db.Query(`SELECT COUNT(*) AS usercount
                             FROM public.events
                             WHERE userid = $1`, userID)
	if err != nil {
		return errors.Wrap(err, ": while query userid count")
	}
	defer rows.Close()
	for rows.Next() {
		var usercount int64
		if err := rows.Scan(&usercount); err != nil {
			return errors.Wrap(err, ": while query usercount")
		}
		if usercount == 0 {
			return storage.ErrUnkownUserID
		}
	}
	return nil
}

func (s *Storage) checkEventIDisPresent(eventID string) error {
	rows, err := s.db.Query(`SELECT COUNT(*) AS eventcount
                             FROM public.events
                             WHERE id = $1`, eventID)
	if err != nil {
		return errors.Wrap(err, ": while query userid count")
	}
	defer rows.Close()
	for rows.Next() {
		var eventcount int64
		if err := rows.Scan(&eventcount); err != nil {
			return errors.Wrap(err, ": while query eventcount")
		}
		if eventcount == 0 {
			return storage.ErrEventIDNotFound
		}
	}
	return nil
}

func (s *Storage) ListEventsToNotify(periodTimeStart time.Time, periodTimeEnd time.Time) (events []storage.Event, err error) {
	fmt.Println("Query between ", periodTimeStart, "and ", periodTimeEnd)
	rows, err := s.db.Query(`SELECT	
	id, title, userid, datetimestart, tilldate, alarmdatetime
	FROM public.events 
	WHERE 
	(
		alarmdatetime BETWEEN  $1 and $2
	)`, periodTimeStart, periodTimeEnd)
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
				return nil, errors.Wrap(err, ": rows next error listeventstonotify")
			}
		}
		events = append(events,
			convertToStorageEvent(pgEvent{
				ID:            id,
				Title:         title,
				Userid:        userid,
				Datetimestart: datetimestart,
				Tilldate:      tilldate,
				Alarmdatetime: alarmdatetime,
			}))
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, ": rows read error listeventstonotify")
	}
	return events, nil
}

func convertToDBEvent(ev storage.Event) (pgEvent, error) {
	dbEvent := pgEvent{}
	if ev.ID == "" {
		ev.ID = utilstorage.GenerateUUID()
	}
	dbEvent.ID = ev.ID
	if ev.UserID == "" {
		return pgEvent{}, storage.ErrEmptyUserIDField
	}
	dbEvent.Title = ev.Title
	dbEvent.Userid = ev.UserID
	dbEvent.Datetimestart = ev.DateTimeStart
	dbEvent.Tilldate = ev.DateTimeStart.Add(ev.Duration)
	dbEvent.Alarmdatetime = ev.AlarmTime
	return dbEvent, nil
}

func convertToStorageEvent(pgEv pgEvent) storage.Event {
	storageEvent := storage.Event{}
	storageEvent.ID = pgEv.ID

	storageEvent.Title = pgEv.Title
	storageEvent.UserID = pgEv.Userid
	storageEvent.DateTimeStart = pgEv.Datetimestart
	storageEvent.Duration = pgEv.Tilldate.Sub(pgEv.Datetimestart)
	storageEvent.AlarmTime = pgEv.Alarmdatetime
	return storageEvent
}

type pgEvent struct {
	ID            string    `db:"id"`
	Title         string    `db:"title"`
	Userid        string    `db:"userid"`
	Datetimestart time.Time `db:"datetimestart"`
	Tilldate      time.Time `db:"tilldate"`
	Alarmdatetime time.Time `db:"alarmdatetime"`
}
