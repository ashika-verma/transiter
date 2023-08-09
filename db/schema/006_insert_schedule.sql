
ALTER TABLE scheduled_service 
    ALTER COLUMN monday SET NOT NULL,
    ALTER COLUMN tuesday SET NOT NULL,
    ALTER COLUMN wednesday SET NOT NULL,
    ALTER COLUMN thursday SET NOT NULL,
    ALTER COLUMN friday SET NOT NULL,
    ALTER COLUMN saturday SET NOT NULL,
    ALTER COLUMN sunday SET NOT NULL;

DROP INDEX scheduled_trip_stop_time_trip_pk_departure_time_idx;

ALTER TABLE scheduled_trip_stop_time
    DROP COLUMN arrival_time,
    ADD COLUMN arrival_time integer,
    DROP COLUMN departure_time,
    ADD COLUMN departure_time integer,
    DROP COLUMN continuous_drop_off,
    ADD COLUMN continuous_drop_off smallint NOT NULL,
    DROP COLUMN continuous_pickup,
    ADD COLUMN continuous_pickup smallint NOT NULL,
    DROP COLUMN drop_off_type,
    ADD COLUMN drop_off_type smallint NOT NULL,
    DROP COLUMN pickup_type,
    ADD COLUMN pickup_type smallint NOT NULL;

CREATE INDEX ix_scheduled_trip_stop_time__trip_pk ON scheduled_trip_stop_time USING btree (trip_pk);
