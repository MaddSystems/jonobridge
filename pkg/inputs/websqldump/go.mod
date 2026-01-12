module websqldump

go 1.23.2

require (
	github.com/MaddSystems/jonobridge/common v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.8.1
)

require filippo.io/edwards25519 v1.1.0 // indirect

// Replace directive pointing to the local common module
replace github.com/MaddSystems/jonobridge/common => ../../common
