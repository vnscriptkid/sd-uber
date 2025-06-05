# Postgres with PostGIS

## Setup Postgres
```bash
# Install Postgres
brew install postgresql@14 && brew services start postgresql@14
brew install postgis

# Stop Postgres
brew services stop postgresql@14
```

## Connect to Postgres
```bash
createuser -s postgres
psql -U postgres -c "CREATE EXTENSION IF NOT EXISTS postgis;"
psql -U postgres
```

## Create table for drivers
```sql
CREATE TABLE drivers (
    id SERIAL PRIMARY KEY,
    name TEXT,
    location GEOGRAPHY(POINT, 4326)  -- Use WGS 84 coordinate system
);

CREATE INDEX idx_drivers_location ON drivers USING GIST (location);
```

## Insert driver location
```sql
INSERT INTO drivers (id, name, location) VALUES
(1, 'Driver 1', ST_SetSRID(ST_Point(105.8542, 21.0285), 4326)), -- Hanoi
(2, 'Driver 2', ST_SetSRID(ST_Point(106.6297, 10.8231), 4326)), -- Ho Chi Minh City
(3, 'Driver 3', ST_SetSRID(ST_Point(108.2022, 16.0544), 4326)), -- Da Nang
(4, 'Driver 4', ST_SetSRID(ST_Point(105.7689, 10.0452), 4326)), -- Can Tho
(5, 'Driver 5', ST_SetSRID(ST_Point(106.6667, 20.8667), 4326)), -- Hai Phong
(6, 'Driver 6', ST_SetSRID(ST_Point(108.2200, 16.0667), 4326)); -- Hue

-- Upsert driver location
INSERT INTO drivers (id, name, location) VALUES
(1, 'Driver 1', ST_SetSRID(ST_Point(105.8542, 21.0286), 4326))
ON CONFLICT (id) DO UPDATE SET location = EXCLUDED.location;
```

## Get nearby drivers (within 300 km of Hanoi)
```sql
SELECT *, ST_Distance(location, ST_SetSRID(ST_MakePoint(105.8542, 21.0285), 4326)) AS distance_m
 FROM drivers WHERE ST_DWithin(location, ST_SetSRID(ST_MakePoint(105.8542, 21.0285), 4326), 300000);
```