-- CREATE EXTENSION postgis;

CREATE TABLE drivers (
    id SERIAL PRIMARY KEY,
    name TEXT,
    location GEOGRAPHY(POINT, 4326)
);

-- Create a spatial index on the location column
CREATE INDEX idx_drivers_location ON drivers USING GIST (location);

SELECT id, name, ST_AsText(location) AS location FROM drivers;
SELECT id, name, ST_AsLatLonText(location) AS location FROM drivers;
