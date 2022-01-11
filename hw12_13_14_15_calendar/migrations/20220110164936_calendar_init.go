package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCalendarInit, downCalendarInit)
}

func upCalendarInit(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE TABLE public.events
	(
	  id uuid NOT NULL,
	  title character varying(100) COLLATE pg_catalog."default" NOT NULL,
	  userid uuid NOT NULL,
	  datetimestart timestamp without time zone NOT NULL,
	  tilldate timestamp without time zone NOT NULL,
	  alarmdatetime timestamp without time zone NOT NULL,
	  CONSTRAINT "PK_Events" PRIMARY KEY (id)
	)
	
	WITH (
		OIDS = FALSE
	)
	TABLESPACE pg_default;
	
	CREATE UNIQUE INDEX "Ix_EventsDate"
		ON public.events USING btree
		(
		 userid,
		 date_part('year'::text, datetimestart) ASC NULLS LAST,
		 date_part('month'::text, datetimestart) ASC NULLS LAST,
		 date_part('day'::text, datetimestart) ASC NULLS LAST,
		 date_part('hour'::text, datetimestart) ASC NULLS LAST,
		 date_part('minute'::text, datetimestart) ASC NULLS LAST,
		 date_part('second'::text, datetimestart) ASC NULLS LAST
		)
	TABLESPACE pg_default;`)
	if err != nil {
		return err
	}
	return nil
}

func downCalendarInit(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec("DROP TABLE public.events;")
	if err != nil {
		return err
	}
	return nil
}
