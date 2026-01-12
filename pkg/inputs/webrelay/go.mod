module webrelay

go 1.23.2

require github.com/MaddSystems/jonobridge/common v0.0.0-00010101000000-000000000000

// Replace directive pointing to the local common module
replace github.com/MaddSystems/jonobridge/common => ../../common
