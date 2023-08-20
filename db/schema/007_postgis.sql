CREATE EXTENSION postgis;
ALTER TABLE stop ADD COLUMN location geography(POINT);
UPDATE STOP SET location = ST_SetSRID(ST_MakePoint(latitude, longitude), 4326);

