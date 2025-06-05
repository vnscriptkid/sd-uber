# Geohash Redis

## Setup redis
```bash
redis-server
redis-cli
```

## Add driver locations
```bash
GEOADD drivers:locations 105.8542 21.0285 driver:1   # Hanoi
GEOADD drivers:locations 106.6297 10.8231 driver:2   # Ho Chi Minh City
GEOADD drivers:locations 108.2022 16.0544 driver:3   # Da Nang
GEOADD drivers:locations 105.7689 10.0452 driver:4   # Can Tho
GEOADD drivers:locations 106.6667 20.8667 driver:5   # Hai Phong
GEOADD drivers:locations 108.2200 16.0667 driver:6   # Hue

# Update driver location
GEOADD drivers:locations 105.8542 21.0286 driver:1
GEOHASH drivers:locations driver:1
```

## Get nearby drivers (within 300 km of Hanoi)
```bash
GEOSEARCH drivers:locations FROMLONLAT 105.8542 21.0285 BYRADIUS 300 km WITHDIST ASC
```

### Example Output
```bash
1) 1) "driver:1"
   2) "0.0002"
2) 1) "driver:5"
   2) "86.2960"
```