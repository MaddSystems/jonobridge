# Skywave ↔ Meitrack MVT366 Event Code Mapping

| Skywave Event            | MVT366 Event Code | MVT366 Description                |
|--------------------------|------------------|-----------------------------------|
| AntennaCutStart          | 28               | GPS Antenna Cut                   |
| AntennaCutEnd            | 25               | GPS Signal Recovery               |
| DigInp2Hi                | 1                | Input 1 Active                    |
| DigInp2Lo                | 9                | Input 9 Inactive                  |
| HarshTurn                | 90 / 91          | Sharp Turn Left / Right           |
| IdlingStart              | 133              | Idle Overtime                     |
| IdlingEnd                | 134              | Idle Recovery                     |
| IgnitionOn               | 42               | Start Moving (closest)            |
| IgnitionOff              | 41               | Stop Moving (closest)             |
| MovingStart              | 42               | Start Moving                      |
| MovingEnd                | 41               | Stop Moving                       |
| PowerBackup              | 22               | External Battery On               |
| PowerMain                | 23               | External Battery Cut              |
| Reset                    | 29               | Device Reboot (Power On)          |
| StationaryIntervalCell   | 35               | Track By Time Interval (closest)  |
| modemRegistration        | —                | No direct equivalent              |
| sleepSchedule            | 26 / 27          | Enter Sleep / Exit Sleep          |
| terminalRegistration     | 31               | Heartbeat (closest)               |