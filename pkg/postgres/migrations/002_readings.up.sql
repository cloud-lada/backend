CREATE TABLE IF NOT EXISTS reading (
    sensor      TEXT        NOT NULL,
    value       NUMERIC     NOT NULL,
    timestamp   TIMESTAMPTZ NOT NULL,

    -- Create a composite primary key that is a combination of the sensor
    -- name and the time of the reading, this is a cheap way of getting
    -- idempotency in case we receive a reading twice.
    PRIMARY KEY(sensor, timestamp)
);

-- Convert the reading table to a HYPERTABLE!!
SELECT create_hypertable('reading', 'timestamp');

-- Create indexes for swift querying.
CREATE INDEX IF NOT EXISTS idx_sensor ON reading(sensor);
CREATE INDEX IF NOT EXISTS idx_timestamp ON reading(timestamp);
