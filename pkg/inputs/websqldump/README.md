# WebSQLDump

A Go-based HTTP service that retrieves vehicle GPS data from a MySQL database and serves it as JSON via REST API endpoints.

## What This Program Does

This program provides a web service that:

1. **Connects to MySQL Database**: Establishes a connection to a MySQL database containing vehicle GPS tracking data
2. **Queries Vehicle Data**: Retrieves GPS records from the `gruasaya` table where `status = 0` (unprocessed records)
3. **Processes Records**: Filters out invalid records (coordinates at 0,0) and formats the data
4. **Serves JSON API**: Exposes vehicle data through HTTP endpoints
5. **Updates Status**: Marks processed records as `status = 1` to prevent reprocessing

## API Response Format

The service returns vehicle data in the following JSON format:

```json
[
  {
    "imei": "867869060839099",
    "DTime": "2025-11-14 18:57:13",
    "Lat": "19.087230",
    "Lon": "-98.184615",
    "Speed": "0.00",
    "Address": "",
    "Plate": "97AP9W",
    "Alias": "97AP9W",
    "Course": "282",
    "altitude": 2217,
    "device_model": "Meitrack T366G",
    "vin": "3ALACWF18NDNK0285",
    "is_ignited": false
  }
]
```

## Environment Variables

### Endpoint Configuration
- `PORTAL_ENDPOINT` - API endpoint path (default: `/test`)

### MySQL Database Configuration
- `MYSQL_HOST` - MySQL server hostname (default: `127.0.0.1`)
- `MYSQL_PORT` - MySQL server port (default: `3306`)
- `MYSQL_USER` - MySQL username (default: `gpscontrol`)
- `MYSQL_PASS` - MySQL password (default: `qazwsxedc`)
- `MYSQL_DB` - MySQL database name (default: `gruasaya`)

## MySQL Table Schema

The program reads from a `devices` table with the following structure:

| Field | Type | Description |
|-------|------|-------------|
| imei | VARCHAR(16) | Unique GPS device identifier |
| lastupdate | DATETIME | Last update timestamp |
| latitude | DECIMAL(10,6) | Latitude coordinate |
| longitude | DECIMAL(10,6) | Longitude coordinate |
| altitude | INT | Altitude in meters |
| speed | FLOAT | Current speed |
| angle | FLOAT | Movement direction/heading |
| plates | VARCHAR(15) | Vehicle license plate |
| vin | VARCHAR(40) | Vehicle identification number |
| name | VARCHAR(255) | Device/Vehicle name |
| alarm_status | VARCHAR(25) | Alarm status |

**Note:** The program filters for records with valid non-zero coordinates only.

## Program Flow

1. **Database Connection**: Connects to MySQL using environment variables
2. **Query Execution**: Runs `SELECT` query on `devices` table with valid coordinates filter
3. **Data Processing**:
   - Filters records with NULL or zero coordinates
   - Formats coordinates and speed as strings
   - Uses device name as alias when available
   - Uses plate number with "Sin definir" fallback
   - Orders by last update timestamp in descending order
4. **JSON Response**: Returns formatted vehicle data array

## Usage

### Build the Application
```bash
go mod tidy
go build
```

### Run the Service
```bash
# With default settings
./websqldump

# With custom MySQL connection
MYSQL_HOST=192.168.1.100 \
MYSQL_USER=myuser \
MYSQL_PASS=mypass \
MYSQL_DB=tracking \
PORTAL_ENDPOINT=/vehicles \
./websqldump
```

### Access the API
```bash
# Get all vehicle data
curl http://localhost:8081/test

# Check available endpoints
curl http://localhost:8081/debug/endpoints
```

## Dependencies

- `github.com/MaddSystems/jonobridge/common/utils` - Common utilities
- `github.com/go-sql-driver/mysql` - MySQL driver for Go

## Server Configuration

- **Port**: 8081
- **Endpoints**:
  - `${PORTAL_ENDPOINT}` - Main API endpoint for vehicle data
  - `/debug/endpoints` - Debug endpoint listing available routes

## Error Handling

The service includes comprehensive error handling for:
- Database connection failures
- Query execution errors
- JSON marshaling issues
- Invalid data records

Errors are logged using the `utils.VPrint` function and returned as JSON error responses when appropriate.
