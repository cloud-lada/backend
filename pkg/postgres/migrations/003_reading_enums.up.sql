CREATE TYPE SENSOR_TYPE AS ENUM (
    'speed',
    'fuel',
    'revolution',
    'engine_temperature',
    'location_latitude',
    'location_longitude'
);

-- Delete rows first to get rid of any data that may exist due to typos etc.
DELETE FROM reading WHERE NOT sensor = ANY(enum_range(NULL::SENSOR_TYPE)::TEXT[]);
ALTER TABLE reading ALTER COLUMN sensor TYPE SENSOR_TYPE USING sensor::SENSOR_TYPE;
