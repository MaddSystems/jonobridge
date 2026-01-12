1 Command Format
1.1 GPRS Command Format
 GPRS command sent from the server to the tracker:
@@<Data identifier><Data length>,<IMEI>,<Command type>,<Command content><*Checksum>\r\n
 GPRS command sent from the tracker to the server:
$$<Data identifier><Data length>,<IMEI>,<Command type>,<Command content><*Checksum>\r\n
1.2 Tracker Command Format
$$<Data identifier><Data length>,<IMEI>,<Command type>,<Event code>,<(-)Latitude>,<(-)Longitude>,<Date and
time>,<Positioning status>,<Number of satellites>,<GSM signal strength>,<Speed>,<Direction>,<Horizontal dilution of precision
(HDOP)>,<Altitude>,<Mileage>,<Run time>,<Base station info>,<I/O port status>,<Analog input value>,<Assisted event
info>,<Customized data>,<Protocol version>,<Fuel percentage>,<Temperature sensor 1 value|Temperature sensor 2
value|……Temperature sensor n value>,<Max acceleration value>,<Max deceleration value>,<*Checksum>\r\n
Note:
 A comma (,) is used to separate data characters. The character type is the American Standard Code for Information
Interchange (ASCII). (hexadecimal: 0x2C)
 Symbols "<" and ">" will not be present in actual data, only for documentation purpose only.
 All multi-byte data complies with the following rule: High bytes are prior to low bytes.
 The size of a GPRS data packet is about 160 bytes.
Descriptions about GPRS packets from the tracker are as follows

Parameter Description Example
@@ Indicates the GPRS data packet header sent from the
server to the tracker. The header type is ASCII.
(Hexadecimal: 0x40)
@@
$$ Indicates the GPRS data packet header sent from the
tracker to the server. The header type is ASCII.
(Hexadecimal: 0x24)
$$
Data identifier Contains 1 byte. The type is the ASCII, and its value ranges
from 0x41 to 0x7A.
Q
Data length Indicates the length of characters from the first comma
(,) to \r\n. Decimal.
Example: $$<Data identifier><Data
length>,<IMEI>,<Command type>,<Command
content><*Checksum>\r\n
25
IMEI Indicates the tracker's IMEI number. The number type is
ASCII. It has 15 digits generally.
353358017784062
Command type Hexadecimal
For details, see chapter 2 and chapter 3.
AAA
Event code Decimal
For details, see section 1.3 "Event Code."
1
Latitude Unit: degree 22.756325 (indicates 22.756325°N)

(-)yy.dddddd Decimal
When a minus (-) exists, the tracker is in the southern
hemisphere. When no minus (-) exists, the tracker is in
the northern hemisphere.
yy indicates the degree.
dddddd indicates the decimal part.
-23.256438 (indicates 23.256438°S)
Longitude
(-)xxx.dddddd
Unit: degree
Decimal
When a minus (-) exists, the tracker is in the western
hemisphere. When no minus (-) exists, the tracker is in
the eastern hemisphere.
xxx indicates the degree.
dddddd indicates the decimal part.
114.752146 (indicates
114.752146°E)
-114.821453 (indicates
114.821453°W)
Date and time
yymmddHHMMSS
yy indicates year.
mm indicates month.
dd indicates day.
HH indicates hour.
MM indicates minute.
SS indicates second.
Decimal
091221102631
Indicates 21 December 2009,
10:26:31 am.
Positioning status Indicates the GPS signal status.
A = Valid
V = Invalid
A
The GPS is valid.
Number of satellites Indicates the number of received GPS satellites.
Decimal
5
Five GPS satellites are received.
GSM signal strength Value: 0–31
Decimal
12
The signal strength is 12.
Speed Unit: km/h
Decimal
58
The speed is 58 km/h.
Direction Indicates the driving direction. The unit is degree. When
the value is 0, the direction is due north. The value ranges
from 0 to 359.
Decimal
45: indicates that the location is at
northeast.
90: indicates that the location is at
due east.
HDOP The value ranges from 0.5 to 99.9. The smaller the value
is, the more the accuracy is.
Decimal
When the accuracy value is 0, the signal is invalid.
0.5–1: Perfect
2–3: Wonderful
4–6: Good
7–8: Medium
9–20: Below average
21–99.9: Poor
Altitude Unit: meter
Decimal
118
Mileage Unit: meter
Decimal
Indicates the total mileage. The maximum value is
4294967295. If the value exceeds the maximum value, it
will be automatically cleared.
564870
Run time Unit: second
Decimal
Indicates the total time. The maximum value is
4294967295. If the value exceeds the maximum value, it
will be automatically cleared.
2546321
Base station info The base station information includes:
MCC|MNC|LAC|CI
The MCC and MNC are decimal, while the LAC and CI are
hexadecimal.
Note: Base station information in an SMS is empty.
460|0|E166|A08B
I/O port status Hexadecimal
Status values of eight input ports and eight output ports:
Bits 0–7 correspond to status of output ports 1–8.
Bits 8–15 correspond to status of input ports 1–8.
0421 (hexadecimal) = 0000 0100
0010 0001
Analog input value Hexadecimal
Eight analog input values are separated by "|".
AD1|AD2|AD3|Battery analog|External power
analog|AD6|AD7|AD8
Unit: V
Note: Analog input values in an SMS report are empty.
Voltage formula of analog AD1–AD3:
T366/T366G: AD1/100
Voltage formula of battery analog (AD4):
T366/T366G: AD4/100
Voltage formula of external power supply (AD5):
T366/T366G: AD5/100
AD6–AD8: Reserved. (Note: Unnecessary AD values at
the end of this parameter can be removed while editing.
For example, if AD6, AD7, and AD8 are not in use, you can
just send the first five AD values:
0123|0456|0235|1234|0324.)
0123|0456|0235|1234|0324|0654
|1456|0222
Assisted
event
info
Geo-fence
number
32-bit unsigned
Only available by GPRS event code 20 or 21.
02 00 00 00 (indicates geo-fence 2)
Vehicle theft
trigger source
32-bit unsigned
Trigger code of a vehicle theft event
Flag generated by event 58
01 00 00 00

iButton/RFID
ID
Indicates the ID number of an iButton key or a RFID card.
Contains 8 hexadecimal characters.
Only available by GPRS event code 37.
42770680 (hexadecimal)
System flag Contains 4 bytes; hexadecimal
Bit 0: Whether to modify the EEP2 parameter. When the
value is 1, the EEP2 parameter is modified.
Bits 1–31: reserved.
Only available by GPRS event code 35.
00000001
The EEP2 parameter is modified.
Temperature
sensor No.
The temperature sensor No. is set by command C40.
Contains 2 hexadecimal characters.
Note: The number is only available by event code 50 or
51.
08 (indicates temperature sensor 8)
Picture name Only available by GPRS event code 39. 0918101221_C2E03
Customized data Reserved
A separator still exists.
Protocol version Decimal
1–50: Used for all common Meitrack protocols.
50–99: Used for OBD.
When the protocol is compatible with the old tracker, the
value is empty or is 0 by default.
3
Fuel percentage Contains 4 hexadecimal characters.
When the fuel sensor type is 0, the sensor is not
connected and the value is empty.
0E2E
The fuel percentage is 36.30%.
Temperature sensor No. +
Temperature value
Contains 6 hexadecimal characters.
The first two characters are the temperature sensor No.
The last four characters are the temperature value
(actual temperature x 100; including the integer and
decimal parts; -327.67°C to +327.67°C).
011A09|021A15|06FB2E
There are 3 temperature sensors.
Temperature sensor 1: 66.65°C
Temperature sensor 2: 66.77°C
Temperature sensor 6: -12.34°C
Max acceleration value Decimal
Unit: mg
Indicates the maximum acceleration value at the specific
time interval of two pieces of AAA data.
30
The maximum acceleration value is
30mg.
Max deceleration value Decimal
Unit: mg
Indicates the maximum deceleration value at the specific
time interval of two pieces of AAA data.
18
The maximum deceleration value
18mg.
* Separates commands from checksums.
Contains 1 byte.
ASCII (hexadecimal: 0x2A)
*
Checksum Contains 2 bytes.
Hexadecimal
BE

The parameter indicates the sum of all data (excluding
the checksum and ending mark).
Example: $$<Data identifier><Data
length>,<IMEI>,<Command type>,<Command
content><*Checksum>\r\n
\r\n Contains 2 bytes. The parameter is an ending character.
The type is ASCII. (Hexadecimal: 0x0d 0x0a)

1.3 Event Code

Event Code Event Default SMS Header (At Most 16 Bytes)
1 Input 1 Active In1 Active
2 Input 2 Active In2 Active
3 Input 3 Active In3 Active
4 Input 4 Active In4 Active
9 Input 1 Inactive In1 Inactive
10 Input 2 Inactive In2 Inactive
11 Input 3 Inactive In3 Inactive
12 Input 4 Inactive In4 Inactive
17 Low Battery Low Battery
18 Low External Battery Low Ext-Battery
19 Speeding Speeding
20 Enter Geo-fence Enter Fence N (N means the number of the fence)
21 Exit Geo-fence Exit Fence N (N means the number of the fence)
22 External Battery On Ext-Battery On
23 External Battery Cut Ext-Battery Cut
24 GPS Signal Lost GPS Signal Lost
25 GPS Signal Recovery GPS Recovery
26 Enter Sleep Enter Sleep
27 Exit Sleep Exit Sleep
28 GPS Antenna Cut GPS Antenna Cut
29 Device Reboot Power On
31 Heartbeat /
32 Cornering Cornering
33 Track By Distance Distance
34 Reply Current (Passive) Now
35 Track By Time Interval Interval
36 Tow Tow
37 iButton/RFID (Only for GPRS)
39 Photo /
40 Power Off Power Off
41 Stop Moving Stop moving
42 Start Moving Start Moving
44 GSM Jamming GSM Jamming
50 Temperature High Temp High
51 Temperature Low Temp Low
52 Full Fuel Full Fuel
53 Low Fuel Low Fuel
54 Fuel Theft Fuel Theft
56 Armed Armed
57 Disarmed Disarmed
58 Vehicle Theft Vehicle Theft
63 No GSM Jamming No GSM Jamming
70 Reject Incoming Call /
71 Get Location by Call /
72 Auto Answer Incoming Call /
73 Listen-in (Voice Monitoring) /
78 Impact Impact
82 Fuel Filling Fuel Filling
83 Ult-Sensor Drop Ult-Sensor Drop
90 Sharp Turn to Left Harsh Cornering
91 Sharp Turn to Right Harsh Cornering
94 Output 1 Active Out1 Active
95 Output 2 Active Out2 Active
96 Output 1 Inactive Out1 Inactive
97 Output 2 Inactive Out2 Inactive
129 Harsh Braking Harsh Braking
130 Harsh Acceleration Fast Accelerate
133 Idle Overtime Idle Overtime
134 Idle Recovery Idle Recovery
135 Fatigue Driving Fatigue Driving
136 Enough Rest after Fatigue Driving Enough Rest
139 Maintenance Notice Maintenance

