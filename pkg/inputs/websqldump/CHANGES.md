# Code Updates - MySQL Schema Fix

## Summary
Updated the WebSQLDump program to use the correct MySQL `devices` table structure instead of the non-existent `gruasaya` table.

## Changes Made

### 1. Database Query
**Before:**
```sql
SELECT imei, gps_time, latitude, longitude, speed, course, plate_number, 
       altitude, device_model, vin, is_ignited
FROM gruasaya 
WHERE status = 0 
ORDER BY gps_time ASC
```

**After:**
```sql
SELECT imei, lastupdate, latitude, longitude, speed, angle, plates, 
       altitude, name, vin, alarm_status
FROM devices 
WHERE latitude IS NOT NULL AND longitude IS NOT NULL
AND latitude != 0 AND longitude != 0
ORDER BY lastupdate DESC
```

### 2. Field Mappings
| Old Field | New Field | Data Type |
|-----------|-----------|-----------|
| gps_time | lastupdate | DATETIME |
| course | angle | FLOAT |
| plate_number | plates | VARCHAR(15) |
| device_model | name | VARCHAR(255) |
| is_ignited | alarm_status | VARCHAR(25) |
| (removed) | speed | FLOAT (now nullable) |
| altitude | altitude | INT (now nullable) |

### 3. Logic Updates
- Removed `status` field filtering and updating
- Added NULL checks for nullable fields using `sql.NullType` wrappers
- Changed coordinate filtering: now skips NULL or zero values
- Changed sort order: now descending by `lastupdate` instead of ascending by `gps_time`
- Device name (`name` field) is now used as alias when available

### 4. Null Handling
Updated all field scans to use SQL null types:
- `sql.NullFloat64` for speed and angle
- `sql.NullInt64` for altitude
- `sql.NullString` for plates, name, vin, alarm_status
- `sql.NullTime` for lastupdate

### 5. Variable Initialization
- `speed`, `angle`, `altitude` now have proper null checking with default values
- `isIgnited` now based on `alarm_status` presence instead of boolean field

## Testing
Program compiles successfully with no errors.

## Database Connection
No changes to MySQL connection setup - still uses environment variables:
- `MYSQL_HOST` (default: 127.0.0.1)
- `MYSQL_PORT` (default: 3306)
- `MYSQL_USER` (default: gpscontrol)
- `MYSQL_PASS` (default: qazwsxedc)
- `MYSQL_DB` (default: gruasaya)
