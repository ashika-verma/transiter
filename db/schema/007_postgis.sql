CREATE EXTENSION postgis;
ALTER TABLE stop ADD COLUMN location geography(POINT);
UPDATE STOP SET location = ST_SetSRID(ST_MakePoint(latitude, longitude), 4326);
ALTER TABLE stop ALTER COLUMN location SET NOT NULL;

ALTER TABLE stop DROP COLUMN latitude;
ALTER TABLE stop DROP COLUMN longitude;
