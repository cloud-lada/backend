ALTER TABLE reading ALTER COLUMN sensor TYPE TEXT USING sensor::TEXT;
DROP TYPE SENSOR_TYPE;