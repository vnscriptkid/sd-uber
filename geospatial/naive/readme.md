# Nearby search using regular index

## Bounding box
```sql
WHERE latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?
```

## Haversine Formula
```sql
SELECT id, 
        (6371 * acos(
            cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) + 
            sin(radians(?)) * sin(radians(latitude))
        )) AS distance
```

