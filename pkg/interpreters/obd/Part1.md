




**4G Wireless Air Protocol**

<a name="_toc161247031"></a>- VJT.04.063 -A











**

|**S/N** |**Version No.** |**Revision content** |**Revised by** |**Revision date** |
| :-: | :-: | :-: | :-: | :-: |
|1|V1.01.000|First draft |XQM|2017-8-24|
|2|V1.01.001|Content proofreading and format standardization |WDC|2018-9-6|
|3|V1.01.002|<p>Added in 1.0x0200 extension information </p><p>ID "0xFE" extension defined by vendors, </p><p><a name="_hlt524534850"></a><a name="_hlt535850469"></a><a name="_hlt535850468"></a>GSM / CDMA Base station information extension in sub-extensions 0x01 and 0x02. </p>|YD|2018-9-12|
|4|V1.01.003|Update the content of the old version 1.04 Protocol |WDC|2018-9-19|
|5|VJT.01.001|<p>Revise protocol: </p><p>1\.	Fault code data, travel data, sleep data and wake-up data are transparently transmitted through 0900. </p><p>2\.	The setup parameters of rapid acceleration, rapid deceleration and sharp turn are changed to 4 levels, high, medium, low and closed. </p><p>3\.	Add the version information packet, and the terminal reports software version / VIN and the platform replies time to the system, which is convenient for timing the terminal. </p><p>4\.	0200 data are still packed and reported in the dynamic way. </p><p>5\.	Add the function of setting privilege number </p><p>6\.	Add examples of escape and un-escape functions </p>|XQM|2018-11-19|
|6|VJT.01.002|<p>1\. 	Add latitude and longitude to the data packet of fault codes. </p><p>2\. 	Add latitude and longitude to the travel data </p><p>3\. 	Adjust the basic data flow of 0200 extended data and add CSQ and standby voltage </p>|XQM|2018-12-10|
|7|VJT.01.003|1\. Add latitude and longitude with the engine off to the driving behavior data |XQM|2018-12-12|
|8|VJT.01.004|<p>1\. 	Remove the number of inflection points in the packet in the line setting </p><p>2\. 	Add the turning point ID to synchronize with the standard wireless air protocal </p>|XQM|2018-12-18|
|9|VJT.01.005|1\. Version information ID: 0205/8205 |XQM|2018-12-20|
|10|VJT.01.006|<p>1\. Part of 0200 alarm data are reported by independent extension ID 0xFA </p><p>2\. In 0200 extended data, the command ID has been re-adjusted to be compatible with the command ID of previous customers </p><p>3\. Change 0900 F1 travel data packet to dynamic data packet. </p>|XQM|2019-01-04|
|11|VJT.01.007|1\. Detail adjustment: Set mileage, set three- emergency collision, over-speed alarm of extension part, water temperature alarm, extension command ID and accelerometer |XQM|2019-01-08|
|12|VJT.01.008|<p>1\. 	A few wrongly written words </p><p>2\. 	Add the number of satellites found, position accuracy and signal-to-noise ratio to 0200 basic extended data </p>|XQM|2019-01-09|
|12|VJT.01.009|<p>1\. 	Modify redundant packet parameter items of configuration parameters in set region </p><p>2\. 	Add 0090 configuration positioning mode field to setting parameter 8103 </p>|XQM|2019-01-15|
|13|VJT.01.010|1\. Add 0x6006 text information reply instruction |YGL|2019-01-21|
|14|VJT.04.011|<p>1\. Add hyperlink function and unite version No. to VJT.04.011 </p><p>2\. Modify 0201 message body and add serial number data item </p><p>3\. Modify 0200 message body and delete the length data item of location report message body </p>|LY|2019-02-15|
|15|VJT.04.012|<p>1\. Add 0x001B GPS antenna status to the basic data flow in 0200 extended data </p><p>2\. Add 0x001C timing status to the basic data flow in 0200 extension data </p><p>3\. Add 0x3008 H600 video status to the external data flow in 0200 extended data </p>|YD|2019-02-21|
|16|VJT.04.013|<p>1\. Modify 0x6210 fault mileage from 2 bytes to 4 bytes </p><p>2\. Modify 0x6110 absolute throttle position from 1 byte to 2 bytes </p><p>3\. Modify 0x6070 long- term fuel trim (cylinder banks 1 and 3) from 1 byte to 2 bytes </p><p>4\. Modify 0x60E0 the ignition timing advance angle of cylinder 1 from 1 byte to 2 bytes </p>|LDY|2019-02-26|
|17|VJT.04.014|<p>1\. Add the following items to8103 setting 8104 query: </p><p>0x2012: Set mileage and fuel consumption type </p><p>0x2013: Set mileage factor </p><p>0x2014: Set fuel consumption factor </p><p>0x2015: Set oil density </p><p>0x2016: Set fuel consumption factor at idling </p><p>2\. Change 8205 platform timing to Beijing UTC/GMT+08:00 time. </p><p>3\. Add extension 0x3009 H600 input signal </p>|YD|2019-03-28|
|18|VJT.04.015|1\. Add truck extension data, compatible with 32960 national standard data flow |XQM|2019-04-08|
|19|VJT.04.016|<p>1\. 	Modify the problem that there is no length in front of each packet in 0704 subcontracting data </p><p>2\. 	Incorrect description of command word in control command 8105 </p>|XQM|2019-04-22|
|20|VJT.04.017|<p>1\. Add sub ID 0x001D to the public basic data item in 0200 extended data: Position marker </p><p>2\. Add 8103 sub id 0x2017 on and off OBD command </p>|XQM|2019-04-23|
|21|VJT.04.018|<p>1\. Add truck data item sub ID 0xFFF1 (mileage data) 0xFFF2 (fuel consumption data) in 0200 extended data item </p><p>2\. Add two modes for location data return to sub function ID 0x2018 in configuration query command 8103 / 8104, first-in first-out and priority transmission of real-time data </p>|XQM|2019-05-23|
|22|VJT.04.019|<p>1\. Add maintenance mode status to status data item in 0200 data </p><p>2\. Add ID to the basic data item in the extended data item of 0200 data: 0x001E indicating accumulated mileage, </p><p>3\. Add 0x001D to F1 travel data of 0900 data to indicate fuel consumption at idle. </p>|XQM|2019-07-05|
|23|VJT.04.020|<p>1\. 0200, the 14th bit in the status mark bit indicates WIFI status, 1 on; 0 off. </p><p>2\. Setting of WIFI parameter in 8103. </p><p>3\. Setting of the during of sleep and wake-up in 8103. </p>|XQM|2019-08-12|
|24|VJT.04.021|<p>1\. Add the following to 0108 upgrading result type, 0xA2: GSM module </p><p>2\. Add type 0xF1 to 8105 remote control to start the OTA upgrade of GSM module </p>|XQM|2019-08-27|
|25|VJT.04.025|1\. Delete all useless protocols, unify font format, and align up and down hyperlinks. |JJH|2019-9-26|
|26|VJT.04.026|<p>1\. Add the data of weighing sensor to 0200 extended peripheral data </p><p>2\. Add parameters of emergency braking, over-speed, PTO idle alarm, and 8103 configuration alarm to 0200 alarm extension data </p><p>3\. Update the lower tire pressure data in 0200 extended peripheral data </p><p>4\. Add feedback packet of MCU upgrade status to 0900 </p>|XQM|2019-9-30|
|27|VJT.04.027|<p>1\. Add high pressure, low pressure and high temperature status bit to 0200 tire data, and add supplementary text to description of tire temperature data </p><p>2\. In case of conflicts between setting parameters in 8103, re-adjust them </p>|XQM|2019-10-21|
|28|VJT.04.028|<p>1\. Adjust the configuration parameters 8103, and add 0x0d pulse speed to 2012 command ID to calculate mileage type </p><p>2\. Add instantaneous fuel consumption to 0200 </p><p>3\. Add the original ministerial standard function to 0702 </p><p>4\. Delay sub ID 0x2024 flameout duration in 8103, supporting setting and query </p>|XQM|2019-11-15|
|29|VJT.04.029|<p>1\. Modify data length of truck data item 60B0 to 2 bytes </p><p>2\. Delete two duplicate definitions of 510B and 5110 from truck data items. </p><p>3\. Refer to GB17691-2018 newly added data items 5111-5118 related to truck environmental protection for truck data items </p>|LYK|2019-12-14|
|30|VJT.04.030|1\. 0200 data truck extended data, intake pipe pressure, originally 60B0 one byte within the scope of 1-255KPA; use new ID 50B0 two bytes within the scope of 1-500KPA without changing the original situation.|XQM|2019-12-20|
|31|VJT.04.031|<p>1\. Add 0x63C0 command ID and catalyst temperature to 0200 car extended data item </p><p>2\. Add HUD text data to 8300 text information. </p>|XQM|2020-02-26|
|32|VJT.04.032|<p>1\. Add the command for car control, and mainly add 0105 control result response command to 8105 terminal control command </p><p>2\. Add setting of 8103 setting command, whether line ACC is valid, fuel consumption factor, minimum interval of data flow of OBD speed </p><p>3\. Adjust the sleep wake-up type in 0900 sleep exit packet </p><p>4\. Add three parameters to 0200 truck extension data: Light absorption coefficient / opacity / particle concentration </p>|XQM|2020-04-27|
|33|VJT.04.033|<p>1\. Add ignition type command ID 0x0020 to 0200 basic data flow, and see the protocol for details </p><p>2\. Add a status bit and an engine status bit to the safety status in 0x0011 vehicle status table in 0200 basic data flow. See the protocol for details. </p>|XQM|2020-08-06|
|34|VJT.04.034|Modify basic data flow |TQL|2020-09-02|
|35|VJT.04.035|1\. Add ministerial standard ID: 8202 temporary location tracking packet |XQM|2020-09-04|
|36|VJT.04.036|<p>1\. Newly added 0200 extension data flow and 0XFB base station data packet. </p><p>2\. Newly added 0200 car data, tire pressure, oil level, maintenance mileage and collision times. </p>|XQM|2020-10-23|
|37|VJT.04.037|<p>1\. 	Newly added 0200 extension data flow, basic data packet 0XEA, and cumulative carbon emissions </p><p>2\. 	Newly added 0200 extension data flow, truck data flow 0XEC, and current engine load 0x511F </p>|XQM|2020-12-29|
|38|VJT.04.038|1\. Newly added 0200 truck data flow 0XEC and relevant data flow of wheat flour detacher. |XQM|2021-04-08|
|39|VJT.04.039|<p>1\. 	Newly added 0200 truck data flow 0XEC and the total running time 0x520A of wheat flour detacher </p><p>2\. 	Newly added 0200 truck data flow 0XEA, Roll angular velocity, Pitch angular velocity, Yaw angular velocity </p><p>3\. 	Adjust the content of new energy vehicle data, 0200 new energy data 0XED, project type </p>|XQM|2021-05-17|
|40|VJT.04.040|1\. Add 0x300B external oil rod data flow to extension peripheral data 0xEE in 0200 position data and see the protocol for details |XQM|2021-07-29|
|41|VJT.04.041|1\. In the extended sedan data 0XEB of the 0200 position data, add two additional extended sub-IDs for the AEB CAN message.|XQM|2021-09-29|
|42|VJT.04.042|<p>Add command 020A, used to report the collected CAN data flow, for customization purposes.</p><p>2\. In the 0200 data, add a differential oil pressure sensor.</p>|XQM|2022-01-17|
|43|VJT.04.043|<p>1\. For the current speed of the drive motor, add an offset of -32767 to distinguish between forward and reverse rotation.</p><p>2\. In 0x0105, control command response, add a serial number of response.</p><p>3\. In the 0x0001 general response, add a status: the previous command is in progress.</p>|PZJ|2022-09-20|
|44|VJT.04.044|<p>1\. In the 0200 FA, add alarms 0x0405 to 0x0408.</p><p>2\. Adjust alarms description in 0200 FA 0x0103 to 0x0104.</p>|PZJ|2022-10-11|
|45|VJT.04.045|1\. Add failure reasons for control command 0105.|PZJ|2022-10-26|
|46|VJT.04.046|1\. Add item 0x0025 in the basic data flow part of 0200 extended data, cumulative mileage 2 (SEEWORLD customized, for tire life statistics).|LHL|2022-10-27|
|47|VJT.04.047|1\. Add 6 channels of collected fire truck-related liquid level in the extended peripheral data flow.|PZJ|2022-11-04|
|48|VJT.04.048|1\. In 0200 FA, add alarm 0x040A for low battery.|LHL|2023-1-6|
|49|VJT.04.049|1\. Add data items 0x520B to 0x520E.|XWB|2023-1-10|
|50|VJT.04.50|1\. Add battery temperature to new energy vehicle data flow.|LHL|2023-1-14|
|51|VJT.04.51|1\. Add 0X300D in extended peripheral data flow, temperature sensor data flow.|XWB|2023-2-20|
|52|VJT.04.52|1\. Add extended ID 0xFC in 0x0200 extended information, wifi data flow.|LHL|2023-3-11|
|53|VJT.04.53|1\. Add ACC interruption ignition type to ignition type in 0x0200 basic data flow.|PZJ|2023-5-23|
|54|VJT.04.54|1\. Add overspeed warning mark bit to the 13th bit of alarm status bit in 0x0200.|PZJ|2023-5-30|
|55|VJT.04.55|<p>1\. Add three parameters in 0x8103, 0x8104: 0x2029, 0x202A, 0x202B, Bluetooth authentication code, Bluetooth name, Bluetooth MAC address.</p><p></p>|PZJ|2023-6-01|
|56|VJT.04.56|1\. Add OBD1 (6X14), OBD2 (1X9, 11X12, 3X11) foot unplugged alarms.|XWB|2023-7-28|
|57|VJT.04.57|1\. Content proofreading and format standardization|PZJ|2023-10-23|
|58|VJT.04.58|1\. Add B gear indication to vehicle status in 0200.|PZJ|2023-10-24|
|59|VJT.04.59|1\. Add new energy data items: vehicle status, insulation resistance, battery health status, highest single cell voltage, lowest single cell voltage, unit pressure difference, power gear.|JJH|2023-10-27|
|60|VJT.04.60|1\. Add a new function ID: 0x0210, used to encapsulate BMS data flow reporting.|PZJ|2023-12-02|
|61|VJT.04.61|1\. Function ID: 0x0210, adjust content.|PZJ|2023-12-04|
|62|VJT.04.62|1\. Function ID: 0x0210, adjust single cell battery voltage units and single cell battery temperature descriptions.|PZJ|2023-12-04|
|63|VJT.04.63|1\. For the 0105 lock command, add a prompt to allow locking only when the instrument cluster and central control are off.|PZJ|2024-01-19|

**Contents** 

[- VJT.04.063 -A	1](#_toc161247031)

[Shenzhen HOLLOO Technology Co., Ltd.	1](#_toc161247031)

[1.	Brief Introduction	10](#_toc161247032)

[1.1	Purpose of compilation	10](#_toc161247033)

[1.2	Terms and definitions	10](#_toc161247034)

[1.3 Abbreviations	10](#_toc161247035)

[1.4 Protocol basis	11](#_toc161247036)

[1.5 Composition of message	12](#_toc161247037)

[1.6 Communication connection	13](#_toc161247038)

[1.7 Message processing	14](#_toc161247039)

[1.8 SMS message processing	14](#_toc161247040)

[1.9 Classification of protocol	15](#_toc161247041)

[2 Data format	19](#_toc161247042)

[2.1	\[0001\] General response of terminal	19](#_toc161247043)

[2.2	\[8001\] General response of platform	19](#_toc161247044)

[2.3	\[0002\] Terminal heartbeat	19](#_toc161247045)

[2.4	 \[0100\] Terminal registration	20](#_toc161247046)

[2.5	 \[0003\] Terminal un-registration	20](#_toc161247047)

[2.6	 \[0102\] Terminal authentication	20](#_toc161247048)

[2.7	 \[0200\] Location information report	21](#_toc161247049)

[2.8	\[0704\] Batch report of location information	21](#_toc161247050)

[2.9	\[020A\] CAN broadcast data flow reporting	21](#_toc161247051)

[2.10	\[0210\] Reporting of BMS data flow of new energy vehicles	21](#_toc161247052)

[2.11	\[0900\] Transparent transmission of data uplink	22](#_toc161247053)

[Data packet of vehicle travel 	     0xF1	22](#_toc161247054)

[Data packet of vehicle fault code 	    0xF2	22](#_toc161247055)

[Data packet of vehicle sleep entry 	   0xF3	22](#_toc161247056)

[Data packet of vehicle sleep wake-up 	   0xF4	22](#_toc161247057)

[Feedback packet of MCU upgrade status 	  0xF6	22](#_toc161247058)

[Description packet of suspected collision alarm 	 0xF7	22](#_toc161247059)

[2.12	\[0205\] Active reporting of version information (non ministerial standard)	23](#_toc161247060)

[2.13	\[8103\] Set terminal parameters	23](#_toc161247061)

[2.14	\[8104\] Query terminal parameters	23](#_toc161247062)

[2.15	\[8201\] Location information query	24](#_toc161247063)

[2.16	\[8300\] Issuing of text information	24](#_toc161247064)

[2.17	\[6006\] Text information reply	24](#_toc161247065)

[2.18 	 \[8105\] Terminal control	25](#_toc161247066)

[2.19	 \[0108\] Notification of terminal upgrade result	25](#_toc161247067)

[2.20	\[0702\] Collection of driver identity information	25](#_toc161247068)

[2.21	 \[8202\] Temporary location tracking control	25](#_toc161247069)

[3. Appendix I:	27](#_toc161247070)

[3.1 	Schedule- General response of terminal	27](#_toc161247071)

[3.2 	Schedule- General response of platform	27](#_toc161247072)

[3.3 	Schedule-Message body of terminal registration	27](#_toc161247073)

[3.4 	Schedule-Response message body of terminal registration	27](#_toc161247074)

[3.5	Schedule-Message body of terminal registration	28](#_toc161247075)

[3.6	Schedule-Message body of terminal parameters	28](#_toc161247076)

[3.7 Schedule-Parameter item format	28](#_toc161247077)

[3.8 Schedule-Definition of terminal parameter setting	28](#_toc161247078)

[3.9 Schedule- Schedule of WIFI parameter	33](#_toc161247079)

[3.10 Schedule- Rapid acceleration parameters	33](#_toc161247080)

[3.11 Schedule- Rapid deceleration parameters	33](#_toc161247081)

[3.12	Schedule-Sharp turning parameters	33](#_toc161247082)

[3.13	Schedule- data packet of terminal upgrade	34](#_toc161247083)

[3.14	Schedule-response of platform upgrade packet	34](#_toc161247084)

[3.15	Schedule-Trailer alarm parameters	34](#_toc161247085)

[3.16	Schedule- Collision alarm parameter packet	34](#_toc161247086)

[3.17	Schedule-List of privilege number	35](#_toc161247087)

[3.18	Schedule-Message body of query terminal parameter response	35](#_toc161247088)

[3.19	Schedule-Message body of terminal control	35](#_toc161247089)

[3.20	Schedule -Descriptions of the terminal control command words	35](#_toc161247090)

[3.21 Schedule- Message body of terminal control	36](#_toc161247091)

[3.22	Schedule- Terminal control response	36](#_toc161247092)

[3.23 	Schedule- Terminal control response result	36](#_toc161247093)

[3.24	Schedule- Format of command parameters	37](#_toc161247094)

[3.25	Schedule-Collection of driver information	38](#_toc161247095)

[3.26	Schedule- Message body of temporary location tracking control	39](#_toc161247096)

[3.27	Schedule- data packet of terminal upgrade result	39](#_toc161247097)

[3.28	Schedule- Message body of data format of location information query response	40](#_toc161247098)

[3.29	Schedule-Batch report packet of location data	40](#_toc161247099)

[3.30	Schedule-Data item format of location batch report	40](#_toc161247100)

[3.31	Schedule-Location report message body	40](#_toc161247101)

[3.32	Schedule- Definition of status mark bit	41](#_toc161247102)

[3.33	Schedule-Definition of alarm mark bits	41](#_toc161247103)

[3.34	Schedule-List of additional information of position	43](#_toc161247104)

[3.35	Schedule- Definition of additional information	43](#_toc161247105)

[3.36	Schedule-Basic data flow	44](#_toc161247106)

[3.37	Schedule- Extended data flow of car	46](#_toc161247107)

[3.38	Schedule- Extended data flow of truck	48](#_toc161247108)

[3.39	Schedule-Data flow of new energy vehicle	51](#_toc161247109)

[3.40	Schedule-Extended peripheral data flow	53](#_toc161247110)

[3.41	Schedule-Alarm command ID and description items	53](#_toc161247111)

[3.42	Schedule -Data flow of base station	55](#_toc161247112)

[3.43	Schedule-Basic data flow: Accelerometer	56](#_toc161247113)

[3.44	Schedule-Basic data items: Format table of total mileage	56](#_toc161247114)

[3.45	Schedule-Basic data items: Cumulative mileage 2 format table	56](#_toc161247115)

[3.46	Schedule-Basic data items: Format table of total fuel consumption	56](#_toc161247116)

[3.47	Schedule-Basic data items: Accelerometer	58](#_toc161247117)

[3.48	Schedule-Basic data items: Sheet of protocol type	58](#_toc161247118)

[3.49	Schedule-Basic data items: Sheet of vehicle status	58](#_toc161247119)

[3.50 Schedule- Alarm description: Description of idle alarm	60](#_toc161247120)

[3.51	Schedule- Alarm description: Description of over-speed alarm	60](#_toc161247121)

[3.52	Schedule- Alarm description: Description of fatigue driving alarm	60](#_toc161247122)

[3.53	Schedule- Alarm description: Alarm description of high-water temperature	60](#_toc161247123)

[3.54	Schedule-extended peripheral data: H600 Sheet of video status information	61](#_toc161247124)

[3.55	Schedule-extended peripheral data: H600 input signal	63](#_toc161247125)

[3.53	Schedule-extended peripheral data: Sheet of tire pressure data	63](#_toc161247126)

[3.57	Schedule- Data sheet of load sensor	64](#_toc161247127)

[3.58	Schedule-Sheet of external oil sensing data	65](#_toc161247128)

[3.59 Schedule - Sheet of fire truck 6 channels data collection	65](#_toc161247129)

[3.60	Schedule- Version information packet	66](#_toc161247130)

[3.61	Schedule- Version information packet response	66](#_toc161247131)

[3.62	Schedule- Message body of issuing of text information	66](#_toc161247132)

[3.63	Schedule - meaning of the text information mark bits	66](#_toc161247133)

[3.64	Schedule- Message body of issuing of text information	66](#_toc161247134)

[3.65	Schedule-Message body of data uplink transparent transmission	67](#_toc161247135)

[3.66	Schedule- Definition of type of transparent transmission message	67](#_toc161247136)

[3.67a	Schedule-Data packet of driving travel F1	68](#_toc161247137)

[3.68	Schedule-Dynamic information sheet of driving travel data	68](#_toc161247138)

[3.69	Schedule-Data packet of fault codes F2	69](#_toc161247139)

[3.70	Schedule- Data packet of sleep entry F3	70](#_toc161247140)

[3.71	Schedule-Data packet of sleep wake-up F4	70](#_toc161247141)

[3.72	Schedule-Feedback packet of MCU upgrade status F6	70](#_toc161247142)

[3.73	Schedule-Description packet of suspected collision alarm F7	71](#_toc161247143)

[3.74 	 Schedule - CAN broadcast data flow	72](#_toc161247144)

[3.75 	Schedule - New energy vehicle BMS data information body	72](#_toc161247145)

[3.76 	Schedule - Data flow of new energy vehicle	72](#_toc161247146)

[3.77 	Schedule - Data flow of new energy vehicle Single battery pack voltage data sheet	72](#_toc161247147)

[3.78 	Schedule - BMS Data flow of new energy vehicles Single battery pack temperature data sheet	73](#_toc161247148)

[3.79 	Schedule - Wifi data flow	74](#_toc161247149)

[4. Appendix 2: Examples	76](#_toc161247150)

[4.1	Examples of escape functions	76](#_toc161247151)

[4.2	Examples of un-escape functions	76](#_toc161247152)

[4.3	\[0200\] Details of analysis of location data analysis	77](#_toc161247153)

[4.4	\[0900\] Details of analysis of transparent transmission of uplink data	77](#_toc161247154)

[4.5	\[8300\] Details of analysis of text information data	77](#_toc161247155)


# **<a name="_toc161247032"></a>1.	Brief Introduction** 
## <a name="_toc161247033"></a>**1.1	Purpose of compilation** 
This file extends functions related with OBD on the basis of JT / T 808 wireless air protocal. 

JT / T 808 wireless air protocal It specifies the communication protocol and data format between the on-board terminal of satellite positioning system for road transport vehicle (hereinafter referred to as the terminal) and the supervision / monitoring platform (hereinafter referred to as the platform), including protocol basis, communication connection, message processing, protocol classification and description and data format. 

OBD functions: The extended data function of the wireless air protocal. 
## <a name="_toc161247034"></a>**1.2	Terms and definitions** 
**a) Abnormal data communication link** 

The wireless communication link is disconnected, or temporarily suspended (e.g., in the process of the call). 

**b) Register** 		The terminal sends messages to the platform to inform that it is installed in a certain vehicle. 

**c) Unregister** 		The terminal sends messages to the platform to inform that it is removed from the vehicle where it is installed. 

**d) Authentication** 	When the terminal is connected to the platform, it sends the message to the platform so that the platform can verify its identity. 

**e) Location report strategy**   Regular, fixed-distance reporting or the combination of both. 

**f) Location report program**   Rules for determining the interval of periodic reporting according to the relevant conditions. 

**g) Additional points report while turning** 

The terminal sends the location information for reporting when it judges that the vehicle turns. The sampling frequency shall not be less than 1 Hz, and the azimuth change rate of the vehicle shall not be less than 15 ° / s for at least 3s. 

**h) Answering strategy**   Rules for the terminal to answer incoming calls automatically or manually. 

**i) SMS text alarm**   The terminal sends text messages in SMS mode when it gives an alarm. 

**j) Event item** 

Set to the platform from the terminal, the event item is composed of event code and name. In case of the corresponding event, the driver operates the terminal to trigger the event report to be sent to the platform. 
## <a name="_toc161247035"></a>**1.3 Abbreviations** 
**APN** -Access point name

**GZIP-**A file-compression program of free software GNU (GNUzip) 

**LCD-**Liquid crystal display 

**RSA-**A kind of asymmetric cryptographic algorithms (developed by Ron Rivest, Adi Shamirh, Len Adleman, named from their names) 

**SMS-**Short message service 

**TCP-** Transmission control protocol

**TTS-**Text to speech 

**UDP-** User datagrnm protocol

**VSS-**Vehicle speed sensor 
## <a name="_toc161247036"></a>**1.4 Protocol basis** 
**1.4.1 Communication mode** 

The communication mode used in the protocol shall comply with relevant provisions of JT/T 794, and the communication protocol uses TCP or UDP, with the platform as the server-side and the terminal as the client-side. When data communication link is abnormal, the terminal can communicate by means of SMS message. 

**1.4.2 Data type** 

Data types used in protocol messages: 

|Data type |Descriptions and requirements |
| - | - |
|BYTE|Unsigned single-byte integer (bytes, 8 bits). |
|WORD|Unsigned double-byte integer (bytes, 16 bits). |
|DWORD|Unsigned four-byte integer (double bytes, 32 bits). |
|BYTE[n]|n bytes |
|BCD[n]|8421 code, n bytes |
|STRING|GBK code, 0 terminator. 0 terminator for no data |

**1.4.3 Transmission rules** 

Big-endian network byte order is used to pass WORD and DWORD in the protocol. 

It is agreed as follows: 

--BYTE transmission protocol: Transmission by means of byte stream; 

--WORD transmission protocol: First high eight-digit, and then low eight-digit; 

--DWORD transmission protocol: First high 24-digit, and high 16-digit, and high 8-digit, and low 8-digit finally. 


## <a name="_toc161247037"></a>**1.5 Composition of message** 
**1.5.1 Structure of message** 

Each message is composed of identity bits, message header, message body and check code. The diagram of the message structure is as shown in Figure 1: 

|Identity bit |Function ID |Message header |Message body |Check code |Identity bit |
| :-: | :-: | :-: | :-: | :-: | :-: |



Figure 1 Message structure 

|GPRS data packet format ||||||||||
| :-: | :- | :- | :- | :- | :- | :- | :- | :- | :- |
|Identity bit |Function ID |Message header |Message packet |Verification |Identity bit |||||
|Identity bit |Function ID |Message attributes |Terminal mobile phone number |Serial number of message |Encapsulation items of message packet |Message packet |Verification |Identity bit ||
|1|2|2|6 Bytes|2|Total number of message packets (2) |Serial number of message packet (2) |N|1|1|
|0x7e||LEN|BCD[12]||||||0x7e|

**1.5.2 Identity bit** 

It shall be represented by 0x7e, and if there is 0x7e in check code, message header and message body, it shall be escaped. Escape rules are defined as follows: 

0x7e<-->0x7d followed by a 0x02; 

0x7d<-->0x7d followed by a 0x01. 

Escaping process is as follows: 

When a message is sent: Encapsulate message- > calculate and fill check code - > escape; 

When a message is received: Restore escape - > verify check code - > parse message. 

Example: 

A data packet containing 0x30 **0x7e** 0x08 **0x7d** 0x55  is sent, 

It is encapsulated as follows: 0x7e. . .    0x30 **7d 0x02** 0x08 **0x7d 0x01** 0x55 Xor 0x7e 

<a name="_消息头"></a>**1.5.3 Message header** 

**Contents of message header:** 

|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0  |Attribute of message body |WORD|The format chart of attribute of message body is shown in Figure 2 |
|2|Terminal mobile phone number |BCD[6]|It is converted according to the mobile phone number of the terminal after installation. If the mobile phone number is less than 12 digits, it shall be supplemented in the front, the mobile phone number in the mainland shall be supplemented by 0, and those in Hong Kong, Macao and Taiwan shall be supplemented by digits according to their area code. |
|8|Serial number of message |WORD|It is cyclically accumulated from 0 according to the order it is sent. |
|10|Encapsulation items of message packet ||If the related identity bit in the attribute of the message body is determined to be processed in packet, the item has the content; otherwise there is no such item |

**The attribute of message body:** 

The format chart of attribute of message body is shown in Figure 2: 

|15|14|13|12|11|10|9|8|7|6|5|4|3|2|1|0|
| :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: | :-: |
|Reserved |Packet |Method of data encryption |Length of message body |||||||||||||

Figure 2 	Format chart of attribute of message body 



Method of data encryption: 

--bit10-bit12 are identity bits for data encryption; 

--When they are 0, it indicates the message body is not encrypted; 

--When the 10th is 1, it indicates the message body is encrypted through RSA algorithm; 

--When the 12th is 1, it indicates the message body is encrypted through SM4 algorithm; 

SM4 algorithm is only used to encrypt the message body in 0200 location report data and 0704 batch report data, and the length of the message body is the encrypted length. 

If SM4 is used for encryption, all downlink data of the server does not need to be encrypted. 

--Others are reserved.

Packet: 

When the 13th bit in the attribute of the message body is 1, it indicates that the message body shall be a long message, it shall be transmitted in packet, and the specific packet information shall be determined by the encapsulation item of the message packet; 

If the 13th bit is 0, there shall be no encapsulation item field of the message packet in the message header. 

Content of the encapsulation item of the message packet 

|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|The total number of message packets |WORD|The total number of packets after the message is processed in packet |
|2|Serial number of packet |WORD|Start from 1 |

**1.5.4 Check code** 

Check code means the byte starting from the function ID is different until the previous byte of the check code, occupying one byte. 

## <a name="_toc161247038"></a>**1.6 Communication connection** 
**1.6.1 	Establishment of connection** 

TCP or UDP can be used for daily data connection between the terminal and the platform. The terminal shall be connected with the platform as soon as possible after it is reset, and it shall send the terminal authentication messages to the platform immediately for authentication after the connection is established. 

**1.6.2 	Maintaining of connection** 

After the connection is established and the terminal authentication is successful, the terminal shall periodically send the terminal heartbeat message to the platform, and the platform shall send the general response message of the platform to the terminal after receiving it, and the sending cycle shall be specified by the terminal parameter. 

**1.6.3 	Disconnection** 

The connection between the platform and the terminal can be actively disconnected according to the TCP protocol, and both of them shall actively determine whether the TCP connection is disconnected. 

The methods that the platform judges the disconnection of TCP connection: 

--It determines the active disconnection of terminal according to TCP protocol; 

--New connection is established for terminals with same identity, which indicates that the original connection has been broken; 

--It has not received the message from the terminal within a certain period, such as terminal heartbeat. 

The methods that for the terminal to judge the disconnection of TCP connection: 

--It determines the active disconnection of terminal according to TCP protocol; 

--The data communication link is disconnected; 

--The data communication link is normal, and no reply has been received after retransmission times have been reached. 
## <a name="_toc161247039"></a>**1.7 Message processing** 
**1.7.1 	Processing of TCP and UDP message** 

**1.7.1.1 	Messages sent by the platform** 

All messages sent by the platform require response from the terminal. The response is divided into general response and special response, which is determined by the specific functional protocol. In case of time out for waiting for a response at the sender, it is necessary to re-send the message. The response timeout time and retransmission times are specified by the platform parameters. The response timeout time and retransmission times after retransmission are specified by the platform parameters. The calculation formula of the response timeout time after retransmission is shown in formula (1):   

T<sub>N+1</sub>=T<sub>N</sub>\*(N+1)         …………(1)

Where: 

T<sub>N+1</sub> --The response timeout time after retransmission; 

T<sub>N--</sub>The previous response timeout; 

N --Retransmission times. 

**1.7.2 	Messages sent by the terminal** 

**1.7.2.1 	Normal data communication link** 

When the data communication link is normal, all messages sent by the terminal require response from the platform. The response is divided into general response and special response, which is determined by the specific functional protocol. In case of time out for waiting for a response at the terminal, it is necessary to re-send the message. The response timeout time and retransmission times are specified by the terminal parameters.  The response timeout time after retransmission is calculated in formula (1). If the response to critical alarm messages sent by the terminal has not been received after the retransmission times are reached, they shall be saved. The key alarm messages saved shall be sent before other messages are sent. 

**1.7.2.2 	Abnormal data communication link** 

The terminal shall save the location information report message that will be sent when the data communication link is abnormal. The saved messages shall be sent immediately after the data communication link returns to normal. 
## <a name="_toc161247040"></a>**1.8 SMS message processing** 
When the terminal communication mode is switched into SMS message of GSM network, PDU eight- encoding is adopted. The message with a length more than 140 bytes shall be processed in packet according to short message service standard GSM 03.40 of GSM network. 

The response, retransmission and save mechanism of SMS message is the same as 6.1, but the response timeout and retransmission times shall be handled according to the relevant set values of parameter ID 0x0006 and 0x0007. 
## <a name="_toc161247041"></a>**1.9 Classification of protocol** 
**Overview** 

The protocol is classified by function as follows. If it is not specified, TCP communication mode is adopted for the default. See Appendix A for communication protocol of vehicle terminal and external device. See Appendix B for comparison table of message name and message ID in the protocol. 

**1.9.1 	Terminal management protocol** 

**1.9.1.1 	Terminal register/unregister** 

When the terminal is not registered, it shall be registered firstly, and after it has been registered successfully, the terminal will be saved with the authentication code, and the authentication code can be used for login of the terminal. Before the vehicle needs to be removed or the terminal is replaced, the terminal shall be unregistered to cancel the corresponding relationship between the terminal and the vehicle. 

If the terminal sends terminal registration and terminal un-registration messages by SMS, the platform shall send a response to the terminal registration by SMS to reply to terminal un-registration, and send the platform general response by way of SMS to reply the terminal un-registration. 

**1.9.1.2 	Terminal authentication** 

The terminal shall be immediately subject to authentication after registration and the connection with the platform. The terminal shall not send any other messages before the success of authentication. 

The terminal is subject to authentication by sending terminal authentication message, and the platform replies general response message. 

**1.9.1.3 	Set/query terminal parameters** 

The platform sets terminal parameters by sending the message of setting terminal parameters, and the terminal shall reply general response message of the terminal. The platform inquires terminal parameters by sending the message of inquiring terminal parameters, and the terminal replies to the query terminal parameter response message. Terminals under different network systems shall support characteristic parameters of their networks respectively. 

**1.9.1.4 	Terminal control** 

The platform controls the terminal by sending the message of terminal control, and the terminal replies terminal general response message. 

**1.9.2 	Location and alarm protocol** 

**1.9.2.1 	Location information report** 

The terminal sends location information report message periodically according to the parameters setting. 

The terminal sends the location information report according to the parameter control in judging the vehicle turns. 

**1.9.2.2 	Location information inquiry** 

The platform queries the current location information of the specified vehicle terminal by sending the location information query message, and the terminal shall reply the message responding to the location information query. 

**1.9.2.3 	Temporary location tracking control** 

The platform starts / stops position tracking by sending a temporary position tracking control message. Position tracking requires the terminal to stop the previous periodical report, and report according to the time interval specified in the message. The terminal shall reply general response message of the terminal. 

**1.9.2.4 Terminal alarm** 

When the terminal judges that the alarm conditions are met, it sends the location information report message, with the corresponding alarm flag set in the location report message, and the platform can respond to the general response message of the platform to process the alarm. 

Each type of the alarm shall be described in the location information reporting information message. The alarm mark shall be maintained until the alarm condition is released, and after the alarm condition is released, the location information reporting message shall be sent immediately to clear the corresponding alarm mark. 

**1.9.3 Information protocol** 

1\.9.3.1 	Issuing of text messages 

The platform delivers the information by sending text messages, informing the driver in the prescribed manner. The terminal shall reply general response message of the terminal. 

1\.9.3.2 	Event settings and reports 

The platform sends the event list to the terminal for storage by sending the event setting message. After encountering the corresponding event, the driver can enter the event list for selection. After selection, the terminal sends an event report message to the platform. 

The terminal shall reply the general response message of the terminal for event setting message. 

The platform shall reply the general response message of the platform for event report message. 

1\.9.3.3 	Questions 

The platform sends issuing questions with candidate answers to the terminal by sending questions messages. The terminal displays them immediately. After the driver selects it, the terminal sends a question response message to the platform. The terminal shall reply the general response message of the terminal for the issuing questions. 

1\.9.3.4 	Collection of driving record data 

The platform requests the terminal to upload the specified data by sending the command message for collection of driving record data, which requires the terminal to reply upload message of the driving record data. 

1\.9.3.5 	Information on demand 

The platform sends the setting of message on demand list to the terminal for storage by sending the information on demand menu. The driver can select the corresponding information service of on demand / cancellation through the menu and the terminal sends message on demand / cancellation to the platform. 

Regular service messages such as news and weather forecast, shall be received from the platform after the information on demand is selected. 

The terminal shall reply general response message of the terminal for setting of message on demand list. 

The platform replies general response message for message on demand /cancellation. 

The terminal shall reply general response message of the terminal for message service. 

**1.9.4 Telephone protocol** 

1\.9.4.1 Callback 

By sending a callback message, the platform requires the terminal to call back according to the specified telephone number and specifies the method for monitoring (without the speaker on at the terminal). The terminal shall reply general response message of the terminal for call-back message. 

1\.9.4.2 Setting of the telephone book 

The platform sets the phonebook for the terminal by sending the message of setting the phonebook, which requires the terminal to reply the general response message of the terminal. 

**1.9.5 Vehicle control protocol** 

The platform sends a vehicle control message to require the terminal to control the vehicle according to the specified operation. The terminal shall reply the general response message of the terminal immediately after receiving the message. The terminal shall conduct control of the vehicle subsequently, and then reply the vehicle control response message according to the consequence. 

**1.9.6 Vehicle management protocol** 

The platform sets the area and line of the terminal by sending messages such as setting circular area, rectangular area, polygon area and route. The terminal shall judge whether the alarm conditions are met according to the area and line attributes. The alarm includes over-speed alarm, access area / route alarm and insufficient / too long driving time of the road section. The additional information of corresponding position shall be included in the position and reporting message. 

The value scope of the area or route ID shall be 1-0XFFFFFFFF. The existing information shall be updated if the setting ID repeats the ID of the same area or route in the terminal. 

The platform can also delete the area and route saved on the terminal by deleting circular area, rectangular area, polygon area and route. 

The terminal shall reply general response message of the terminal for setting/cancellation of area and route messages. 

**1.9.7 Information collection protocol** 

**1.9.7.1 	Data collection of the driver's identity information** 

The terminal can collect the data of driver's identity information and upload them to the platform for identification, and the platform can reply the success or failure of message. 

**1.9.7.2 	Data collection of the electronic waybills** 

The terminal shall collect the electronic waybill information and upload it to the platform. 

**1.9.7.3 	Download the parameters of driving records** 

The platform requests the terminal to upload the specified data by sending the command message for downloading driving record data, which requires the terminal to reply upload message of the driving record data. 

**1.9.8 Multimedia protocol** 

**1.9.8.1 	Uploading of the multimedia event information** 

When the terminal takes the initiative to shoot or record due to a specific event, it shall actively upload a multimedia event message after the event, which requires the platform to reply to the general response message. 

**1.9.8.2 	Uploading of multimedia data** 

The terminal shall upload the data by sending the uploading message of multimedia data. Complete multimedia data need to be preceded by the location information report message during recording, which is called location multimedia data. The platform determines the receiving timeout time according to the total number of packets. After receiving all data packets or reaching the timeout times, the platform sends  to the terminal a response message of uploading multimedia data, which confirms the receipt of all data packets or requires the terminal to retransmit the specified data packets. 

**1.9.8.3 	Immediate shooting by camera** 

The platform issues a shooting command to the terminal by sending a command message of the immediate shooting of the camera to the terminal, which requires the terminal to reply the general response message of the terminal. If real-time upload is specified, the terminal uploads the camera image / video after shooting, otherwise the image / video is stored. 

**1.9.8.4 	Start recording** 

The platform sends a recording command to the terminal by sending a command message for starting recording, which requires the terminal to reply the general response message of the terminal. If real-time upload is specified, the terminal will upload audio data after recording, otherwise the audio data will be stored. 

**1.9.8.5 	Retrieval and extraction of the multimedia data saved at the terminal** 

The platform shall obtain the multimedia data saved at the terminal by sending a retrieval message of the multimedia data storage, which requires the terminal to reply the response message of retrieval of the multimedia data storage. 

According to the retrial result, the platform can request the terminal to upload the specified multimedia data by sending upload message of the stored multimedia data, which requires the terminal to reply the general response message of the terminal. 

**1.9.9 	Transmission of general data** 

For messages that are not defined in the protocol but need to be delivered in actual use, uplink transparent transmission messages and downlink transparent transmission messages can be used for exchange of uplink data and downlink data. The terminal can compress long messages with GZIP compression algorithm and upload messages with data compression. 

**1.9.10 	Encryption protocol** 

RSA public key cryptosystem may be used for encrypted communication between the platform and the terminal. The platform can inform the terminal of its own RSA public key by sending the platform RSA public key message, and the terminal shall reply the terminal RSA public key message, and vice versa. 


# <a name="_toc161247042"></a>**2 Data format** 
## <a name="_toc161247043"></a><a name="_[0001]终端通用应答"></a>**2.1	 [0001] General response of terminal** 
[Function description]: The message data of the general response of the terminal 

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**00 01**|[Message attachment ](#_消息头)|<a name="_hlt530391502"></a>[General response schedule of terminal ](#_终端通用应答消息体数据)|<a name="_hlt491357070"></a><a name="_hlt491357071"></a>XOR|7E|

## <a name="_toc161247044"></a><a name="_[8001]平台通用应答"></a>**2.2	[8001] General response of platform** 
[Function description]: The message data of the general response of the platform 

**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**80 01**|[Message attachment ](#_消息头)|<a name="_hlt530391516"></a>[General response schedule of platform ](#_平台通用应答消息体数据)|<a name="_hlt491357236"></a><a name="_hlt491433840"></a>XOR|7E|

## <a name="_toc161247045"></a>**2.3	[0002] Terminal heartbeat** 
**[Function description]** Reporting of terminal heartbeat packet 

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**00 02**|[Message attachment ](#_消息头)|None |XOR|7E|

**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**80 01**|[Message attachment ](#_消息头)|[General response of platform ](#_平台通用应答消息体数据)|XOR|7E|


## <a name="_toc161247046"></a><a name="_[0100]终端注册"></a>**2.4   [0100] Terminal registration** 
`	`**[Function description]** Message body data of terminal registration 

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**01 00**|[Message attachment ](#_消息头)|[Message body of terminal registration ](#_终端注册消息体附表)|XOR|7E|

**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**81 00**|[Message attachment ](#_消息头)|<a name="终端注册应答消息体a"></a>[Response message body of the terminal registration ](#终端注册应答消息体b)|XOR|7E|


## <a name="_toc161247047"></a>**2.5   [0003] Terminal un-registration** 
**[Function description]** The message body of terminal un-registration is empty 

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**00 03**|[Message attachment ](#_消息头)|None |XOR|7E|

**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**80 01**|[Message attachment ](#_消息头)|[General response of platform ](#_平台通用应答消息体数据)|XOR|7E|


## <a name="_toc161247048"></a><a name="_[0102]终端鉴权"></a>**2.6   [0102] Terminal authentication** 
**[Function description]** Message body data of terminal authentication 

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**01 02**|[Message attachment ](#_消息头)|[Message body of terminal authentication ](#_附表_终端鉴权消息体)|XOR|7E|

**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**80 01**|[Message attachment ](#_消息头)|[General response of platform ](#_平台通用应答消息体数据)|XOR|7E|




## <a name="_toc161247049"></a><a name="_[0200]位置信息汇报"></a>**2.7   [0200] Location information report** 
**[Function description]** The message body of location information report consists of a list of location basic information and location additional information items 

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**02 00**|[Message attachment ](#_消息头)|[Message body of location data ](#_附表_位置数据信息体)|<a name="_hlt24447401"></a><a name="_hlt36396050"></a><a name="_hlt24447400"></a><a name="_hlt54358710"></a><a name="_hlt22569293"></a><a name="_hlt54343427"></a><a name="_hlt20736492"></a><a name="_hlt41309725"></a><a name="_hlt54343428"></a>XOR|7E|

**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**80 01**|Message attachment |General response of platform |XOR|7E|


## <a name="_toc161247050"></a><a name="_[0704]位置信息批量汇报"></a>**2.8	[0704] Batch report of location information** 
**[Function description]** Batch report of location information 

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**07 04**|Message attachment |Batch report packet of location data |XOR|7E|

<a name="_[0900]数据上行透传"></a>**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**80 01**|Message attachment |General response of platform |XOR|7E|


## <a name="_toc161247051"></a>**2.9	[020A] CAN broadcast data flow reporting**
**[Function description]** Collect CAN data flow from the bus, and customize functionality

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**02 0A**|[Message attachment ](#_消息头)|CAN broadcast data flow|<a name="_hlt22569191"></a><a name="_hlt41312797"></a><a name="_hlt531969449"></a><a name="_hlt41309746"></a>XOR|7E|

**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**80 01**|[Message attachment ](#_消息头)|General response of platform|XOR|7E|

## <a name="_toc161247052"></a><a name="_[0205]版本信息包"></a>**2.10	[0210] Reporting of BMS data flow of new energy vehicles**
[Function description] Collect BMS data flow from the bus of new energy vehicles, and customize functionality

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**02 10**|[Message attachment ](#_消息头)|New energy vehicle BMS data information body|XOR|7E|

**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**80 01**|[Message attachment ](#_消息头)|General response of platform|XOR|7E|


## <a name="_toc161247053"></a>**2.11	<a name="_toc82611064"></a> [0900] Transparent transmission of data uplink** 
<a name="_toc82611065"></a><a name="_toc161247054"></a>**Data packet of vehicle travel 						0xF1**

<a name="_toc161247055"></a><a name="_toc82611066"></a>**Data packet of vehicle fault code 					0xF2**

<a name="_toc82611067"></a><a name="_toc161247056"></a>**Data packet of vehicle sleep entry 				0xF3**

<a name="_toc161247057"></a><a name="_toc82611068"></a>**Data packet of vehicle sleep wake-up 				0xF4**

<a name="_toc82611069"></a><a name="_toc161247058"></a>**Feedback packet of MCU upgrade status 			0xF6**

<a name="_toc161247059"></a><a name="_toc82611070"></a>**Description packet of suspected collision alarm 		0xF7**

**[Function description]** Message body data of data uplink transparent transmission 

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|09 00|[Message attachment ](#_消息头)|Message body schedule of data uplink transparent transmission |XOR|7E|

**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**80 01**|[Message attachment ](#_消息头)|General response of platform |XOR|7E|
##

##
## <a name="_toc161247060"></a>**2.12	[0205] Active reporting of version information (non ministerial standard)** 
[Function description] Including software version, software release time, module model, total mileage, total fuel consumption and VIN 

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**02 05**|[Message attachment ](#_消息头)|[Packet of version information ](#_版本信息包附表)|XOR|7E|

**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**82 05**|[Message attachment ](#_消息头)|[Version packet response ](#_版本信息包应答附表)|XOR|7E|


## <a name="_toc161247061"></a>**2.13	[8103] Set terminal parameters** 
**[Function description]:** Set the message body data of terminal parameter. 

**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**81 03**|[Message attachment ](#_消息头)|[Message body schedule of terminal parameter. ](#_终端参数消息体附表)|<a name="_hlt20736631"></a><a name="_hlt24744115"></a>XOR|7E|

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**00 01**|[Message attachment ](#_消息头)|[General response of uplink ](#_[0001]终端通用应答)|XOR|7E|


## <a name="_toc161247062"></a>**2.14	[8104] Query terminal parameters** 
[Function description]: The message body of query terminal parameters is empty. 

**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**81 04**|[Message attachment ](#_消息头)|Null |XOR|7E|

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**01 04**|[Message attachment ](#_消息头)|[Message body schedule of query terminal parameter response ](#_查询终端参数应答消息体附表)|XOR|7E|


## <a name="_toc161247063"></a>**2.15	[8201] Location information query** 
[Function description]: Message body of location information query is empty. 

**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**82 01**|[Message attachment ](#_消息头)|Null |XOR|7E|

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**02 01**|[Message attachment ](#_消息头)|[Response data schedule of location information query ](#_附表_位置作息查询应答消息体数据格式)|XOR|7E|

## <a name="_toc161247064"></a><a name="_[8300]文本信息下发"></a>**2.16	[8300] Issuing of text information** 
[Function description]: Message body data of issuing of text information (SMS settings / TTS broadcast) 

**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|83 00|[Message attachment ](#_消息头)|[Message body schedule of issuing of text information ](#_附表_文本信息下发消息体)|<a name="_hlt535848832"></a><a name="_hlt535848477"></a><a name="_hlt535848993"></a><a name="_hlt535848112"></a><a name="_hlt535848778"></a><a name="_hlt535847202"></a><a name="_hlt535847518"></a><a name="_hlt535847322"></a>XOR|7E|

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**00 01**|[Message attachment ](#_消息头)|[General response of uplink ](#_[0001]终端通用应答)|XOR|7E|


## <a name="_toc161247065"></a><a name="_[6006]文本信息回复"></a>**2.17	[6006] Text information reply** 
[Function description]: Text message data on terminal equipment 

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|60 06|[Message attachment ](#_消息头)|[Message body schedule on text information ](#_附表_文本信息上发消息体)|<a name="_hlt535848989"></a><a name="_hlt535849111"></a><a name="_hlt535848829"></a><a name="_hlt535849050"></a><a name="_hlt535849086"></a><a name="_hlt535848826"></a><a name="_hlt535848965"></a><a name="_hlt535848985"></a><a name="_hlt535848857"></a>XOR|7E|

**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**80 01**|[Message attachment ](#_消息头)|[General response of platform ](#_平台通用应答消息体数据)|XOR|7E|


## <a name="_toc161247066"></a><a name="_[8105]终端控制"></a><a name="_[8105]终端控制(远程升级指令)"></a>**2.18 	 [8105] Terminal control** 
[Function description]: Data format of terminal control message body 

**[Downlink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**81 05**|[Message attachment ](#_消息头)|[Message body schedule of terminal control ](#_终端控制消息体附表)|<a name="_hlt37347188"></a><a name="_hlt54706067"></a><a name="_hlt54706132"></a><a name="_hlt54706038"></a>XOR|7E|

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**00 01**|[Message attachment ](#_消息头)|[General response of terminal ](#_附表_终端通用应答)|<a name="_hlt37347124"></a><a name="_hlt37347091"></a>XOR|7E|



**For individual types of terminal control, it is necessary to supplement the control result of 0x0105. After receiving the control command, it sends 0x0001 back to the platform to indicate that the control command is received. After control, it reports the control result to the platform through 0x0105.** 

**[Uplink]** 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**01 05**|[Message attachment ](#_消息头)|Response message body of the terminal control |XOR|7E|

## <a name="_toc161247067"></a>**2.19	 [0108] Notification of terminal upgrade result** 
[Message ID]: 0x0108.

[Function description]: After the terminal upgrade, it notifies the monitoring center through the command. 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**01 08**|[Message attachment ](#_消息头)|[Data packet of upgrade result ](#_附表_终端升级结果数据包)|<a name="_hlt17830902"></a><a name="_hlt22569118"></a><a name="_hlt20331286"></a><a name="_hlt17830901"></a>XOR|7E|
## <a name="_toc161247068"></a>**2.20	[0702] Collection of driver identity information** 
[Message ID]: 0x0702.

[Function description]: After receiving the 0x8702 command, the terminal will automatically reply 0702 collection packet of driver information, or automatically report 0702 collection packet of driver information when signing in and signing out. 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**07 02**|[Message attachment ](#_消息头)|[Schedule of driver identity information ](#_附表_驾驶员信息采集附表)|XOR|7E|

## <a name="_toc161247069"></a><a name="_toc534810536"></a>**2.21	 [8202] Temporary location tracking control** 
[Message ID]: 0x8202.

[Function description]: Message body data of temporary location tracking control 

|Identification |Function ID |Message header |Message body |Verification |Identification |
| :-: | :-: | :-: | :-: | :-: | :-: |
|7E|**82 02**|[Message attachment ](#_消息头)|[Message body schedule of temporary location tracking control ](#_附表_临时位置跟踪控制消息体)|XOR|7E|




# <a name="_toc161247070"></a>**3. Appendix I:** 
## <a name="_终端通用应答消息体数据"></a><a name="_附表_终端通用应答"></a><a name="_toc161247071"></a>**3.1 	Schedule- General response of terminal [](#_[0001]终端通用应答)**

|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Serial number of response |WORD|The corresponding serial number of the platform message |
|2|Response ID |WORD|The corresponding ID of the platform message |
|4|Result |BYTE|0: Success/confirmation; 1: Failure; 2: Information error; 3: Not supported; 4: Previous operation in progress|

## <a name="_平台通用应答消息体数据"></a><a name="_toc161247072"></a>**3.2 	Schedule- General response of platform [](#_[8001]平台通用应答)**

|<a name="_hlt37347093"></a>Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Serial number of response |WORD|The corresponding serial number of the terminal message |
|2|Response ID |WORD|The corresponding ID of the terminal message |
|4|Result |BYTE|0: Success/confirmation; 1: Failure; 2: Information error; 3: Not supported |

## <a name="_终端注册消息体数据"></a><a name="_终端注册消息体附表"></a><a name="_toc161247073"></a>**3.3 	Schedule-Message body of terminal registration [](#_[0100]终端注册)**

|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Province ID |WORD|It indicates the province ID of the vehicle installed with the terminal, with 0 reserved, and the platform shall take the default value. The province ID shall take the first two of six bits of the administrative area code specified in GB/T 2260. |
|2|City and county ID |WORD|It indicates the city and county of the vehicle installed with the terminal; with 0 reserved, and the platform shall take the default value. The city and county ID shall take the last four of six bits of the administrative area code as specified in GB/T 2260. |
|4|Manufacturer ID |BYTE[5]|5 bytes, the terminal manufacturer code |
|9|Terminal model |BYTE[20]|New Beidou 20 bytes. |
|29|Terminal ID |BYTE[7]|Seven bytes, consisting of uppercase letters and numbers. The terminal ID is defined by the manufacturer |
|36|Color of license plate |BYTE|The color of license plate shall be in accordance with 5.4.12 of JT/T 415-2006. |
|37|License plate |STRING|License plate of motor vehicle issued by the Traffic Management Department of Public Security |


## <a name="_终端注册应答消息体附表"></a><a name="_终端注册应答消息体数据"></a><a name="_toc161247074"></a>**3.4 	Schedule-Response message body of terminal registration [](#终端注册应答消息体a)**

|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Serial number of response |WORD|The serial number of the corresponding terminal registration message |
|2|Result |BYTE|0: Success; 1: The vehicle has been registered; 2: The vehicle is not in the database; 3: The terminal has been registered; 4: The terminal is not in the database. |
|3|Authentication code |STRING|The field exists only for successful authentication |


## <a name="_终端鉴权消息体数据"></a><a name="_终端鉴权消息体附表"></a><a name="_toc161247075"></a><a name="_附表_终端鉴权消息体"></a>**3.5	Schedule-Message body of terminal registration [](#_[0102]终端鉴权)**
|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Authentication code |STRING|The terminal reports the authentication code after reconnection |


## <a name="_终端参数消息体附表"></a><a name="_终端参数消息体数据"></a><a name="_toc161247076"></a><a name="_附表_终端参数消息体"></a>**3.6	Schedule-Message body of terminal parameters [](#_[8103]设置终端参数)**
|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Total number of parameters |BYTE|**N parameter items** |
|1|List of parameter items ||<p>[Parameter item format 1, ](#_参数项格式)</p><p><a name="_hlt491358173"></a>[Parameter item format 2 ](#_参数项格式), </p><p>Parameter item format 3, </p><p>. . . . . . </p><p>Parameter item format N </p>|

## <a name="_参数项格式"></a><a name="_参数项格式附表"></a><a name="_toc161247077"></a><a name="_附表_参数项格式"></a>**3.7 Schedule-Parameter item format [](#_附表_终端参数消息体)**
|Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: |
|Parameter ID |DWORD |Table of definition and description of parameter ID, see definition of terminal parameter setting for details [](#_终端参数设置各参数项定义及说明附表)|
|<a name="_hlt17118831"></a><a name="_hlt534733999"></a>Parameter length |BYTE ||
|Parameter value ||DWORD or STRING. If it is a multi value parameter, multiple parameter items with the same ID are used in the message, such as the telephone number of the dispatching center |
## <a name="_终端参数设置各参数项定义及说明"></a><a name="_终端参数设置各参数项定义及说明附表"></a><a name="_toc161247078"></a><a name="_附表_终端参数设置各参数项定义及说明"></a>**3.8 Schedule-Definition of terminal parameter setting [](#_附表_参数项格式)**
|Parameter ID |Data type |Descriptions and requirements |
| :-: | :-: | :-: |
|0x0001|DWORD|The interval of sending the terminal heartbeat, the unit is in seconds(s) |
|0x0002|DWORD|Response timeout value of TCP messages, the unit is in seconds(s) |
|0x0003|DWORD|Retransmission number of TCP messages |
|0x0004|DWORD|Response timeout value of UDP messages, the unit is in seconds(s) |
|0x0005|DWORD|Retransmission number of UDP messages |
|0x0006|DWORD|Response timeout value of SMS messages, the unit is in seconds(s) |
|0x0007|DWORD|Retransmission number of SMS messages |
|0x0008-0x000F||Reserved |
|0x0010|STRING|Primary server APN, dialing access point of wireless communications. It shall be PPP dialing numbers if the network type is CDMA. |
|0x0011|STRING|User name of wireless communication dialing of the primary server |
|0x0012|STRING|Password of wireless communication dialing of the primary server |
|0x0013|STRING|Primary server address, IP or domain name |
|0x0014|STRING|Backup server APN, dialing access point of wireless communications |
|0x0015|STRING|User name of wireless communication dialing of the backup server |
|0x0016|STRING|Password of wireless communication dialing of the backup server |
|0x0017|STRING|Backup server address, IP or domain name |
|0x0018|DWORD|Server TCP port |
|0x0019|DWORD|Server UDP port |
|0x001A-0x001F||Reserved |
|0x0020|DWORD|location report strategy, 0: Regular reporting; 1: Fixed-distance reporting; 2: Regular and fixed- distance reporting. |
|0x0021|DWORD|location report program, 0: In accordance with ACC status; 1: Firstly, determine the login status in accordance with the login status and ACC status, and then login in accordance with ACC status. |
|0x0022|DWORD|Time interval of report for driver not logged in, the unit is in seconds (s), >0 |
|0x0023-0x0026|DWORD|Reserved |
|0x0027|DWORD|Time interval of report in hibernation, the unit is in seconds (s), >0 |
|0x0028|DWORD|Time interval of report in emergency alarm, the unit is in seconds (s), >0 |
|0x0029|DWORD|Time interval of report for the default, the unit is in seconds (s), >0 |
|0x002A-0x002B|DWORD|Reserved |
|0x002C|DWORD|Distance interval of report for the default, the unit is in meters (m), >0 |
|0x002D|DWORD|Distance interval of report for the driver not logged in, the unit is in meters (m), >0 |
|0x002E|DWORD|Distance interval of report in hibernation, the unit is in meters (m), >0 |
|0x002F|DWORD|Distance interval of report in emergency alarm, the unit is in meters (m), >0 |
|0x0030|DWORD|The retransmission angle of turning points, <180° |
|0x0031-0x003F||Reserved |
|0x0040|STRING|Phone number of the monitoring platform |
|0x0041|STRING|Call the terminal phone to reset the phone number by the terminal. |
|0x0042|STRING|Call the terminal phone to restore factory settings by the terminal. |
|0x0043|STRING|SMS phone number of the monitoring platform |
|0x0044|STRING|SMS text alarm number of the receiving terminal |
|0x0045|DWORD|Answering strategy of the terminal, 0: Automatic answering; 1: Automatic answering when ACC is ON, and manual answering when ACC is OFF |
|0x0046|DWORD|The maximum duration of call, in seconds (s), 0 indicates that no call is allowed, and 0xFFFFFFFF indicates that there is no limit |
|0x0047|DWORD|The maximum duration of call in the current month, in seconds (s), 0 indicates that no call is allowed, and 0xFFFFFFFF indicates that there is no limit. |
|0x0048|STRING|Monitoring phone number |
|0x0049|STRING|Privileged SMS numbers of the monitoring platform |
|0x004A-0x004F||Reserved |
|0x0050|DWORD|Alarm mask word. It corresponds to the alarm flag in the location information report message. If the corresponding bit is 1, the corresponding alarm is masked |
|0x0051|DWORD|The alarm sending text SMS switch corresponds to the alarm flag in the position information reporting message. If the corresponding bit is 1, the text SMS will be sent when the corresponding alarm occurs |
|0x0052|DWORD|The alarm shooting switch corresponds to the alarm flag in the location information report message. If the corresponding bit is 1, the camera will shoot when the corresponding alarm occurs |
|0x0053|DWORD|The alarm shooting storage flag corresponds to the alarm flag in the location information report message. If the corresponding bit is 1, the photos shot during the corresponding alarm will be stored, otherwise it will be transmitted in real time |
|0x0054|DWORD|The key flag corresponds to the alarm flag in the location information report message. If the corresponding bit is 1, the corresponding alarm is a key alarm |
|0x0055|DWORD|Maximum speed, the unit is in Kilometers per hour(km/h)|
|0x0056|DWORD|Duration for over-speed, the unit is in seconds(s) |
|0x0057|DWORD|Time threshold for continuous driving, the unit is in seconds(s) |
|0x0058|DWORD|Accumulative driving time threshold on that day, the unit is in seconds(s) |
|0x0059|DWORD|Minimum rest time, the unit is in seconds(s) |
|0x005A|DWORD|Longest parking time, the unit is in seconds(s) |
|0x005B-0x006F||Reserved |
|0x0070|DWORD|Image/video quality, 1-10, 1 for the best |
|0x0071|DWORD|Brightness, 0-255 |
|0x0072|DWORD|Contrast, 0-127 |
|0x0073|DWORD|Saturation, 0-127 |
|0x0074|DWORD|Chromaticity, 0-255 |
|0x0075-0x007F|DWORD||
|0x0080|DWORD|The odometer readings of vehicle, 1/10km |
|0x0081|WORD|Province ID of the vehicle |
|0x0082|WORD|City ID of the vehicle |
|0x0083|STRING|License plate of motor vehicle issued by the Traffic Management Department of Public Security |
|0x0084|BYTE|The color of license plate shall be in accordance with 5.4.12 of JT/T415-2006. |
|0x0090|BYTE|Positioning mode:   0x01: GPS,0x02: BD,0x03 bi-module |
|The following ID is for the manufacturer |||
|0x2001|BYTE|Reset fault code       0x01: Clear 0x00: No clear |
|0x2002|BYTE|Clear vehicle data      0x01: Clear 0x00: No clear |
|0x2003|BYTE|Clear driving travel data  0x01: Clear 0x00: No clear |
|0x2004|DWORD|The total fuel consumption ml |
|0x2006|DWORD|Water temperature alarm parameter, unit ℃ |
|0x2007|BYTE|Schedule of rapid acceleration parameters |
|0x2008|BYTE|Schedule of rapid deceleration parameters |
|0x2009|BYTE|Schedule of sharp turning parameters |
|0x200A|WORD|Vehicle type, see the manufacturer's model table for details |
|0x200B|DWORD|Low voltage alarm parameter, unit: 0.1V |
|0x200C|DWORD|Alarm for too long idle time, unit s |
|0x200D|DWORD|Alarm for too long positioning time, unit s |
|0x200E|STRING|Schedule of trailer alarm parameters |
|<a name="_hlt16601362"></a>0x200F|BYTE|Schedule of collision alarm parameters |
|<a name="_hlt534706070"></a><a name="_hlt530407757"></a>0x2010|STRING|Schedule of privilege number |
|0x2011|DWORD|Ignition threshold voltage, unit: 0.1V |
|0x2012|WORD|<p>Mileage type (high byte), type of fuel consumption (low byte) </p><p>Mileage type: </p><p>0x00:Cancel mandatory settings </p><p>0x01: GPS</p><p>0x02: J19391</p><p>0x03: J19392</p><p>0x04: J19393</p><p>0x05: J19394</p><p>0x06: J19395</p><p>0x07: OBD instrument </p><p>0x08: OBD/ private protocol </p><p>0x09: J1939A</p><p>0x0A: J1939B</p><p>0x0B: J1939C</p><p>0x0C: J1939D</p><p>0x0D:Impulse speed </p><p>…</p><p>0xff:No change in mandatory type </p><p>Type of fuel consumption: </p><p>0x00:Cancel mandatory settings </p><p>0x01: J19391</p><p>0x02: J19392</p><p>0x03: J19393</p><p>0x04: J19394</p><p>0x05: J19395</p><p>0x06: OBD1</p><p>0x07: OBD2</p><p>…</p><p>0xff: No change in mandatory type </p>|
|0x2013|WORD|Mileage coefficient: Setting value/1000. Example:   1020 ->  1.02 |
|0x2014|WORD|Fuel consumption factor: Setting value/1000. Example:   1020 ->  1.02 |
|0x2015|WORD|<p>Oil density: </p><p>Diesel oil 0#  835</p><p>Diesel oil 10#  840</p><p>Diesel oil 20#  830</p><p>Diesel oil 35#  820</p><p>Diesel oil 50#  816</p><p>Gasoline 90#  722</p><p>Gasoline 92#  725</p><p>Gasoline 95#  737</p><p>Gasoline 98#  753</p>|
|0x2016|WORD|Fuel consumption factor at idling: Setting value/1000. Example:   1020 ->  1.02 |
|0x2017|BYTE|<p>0x01: Turn on OBD </p><p>0x00: Turn off OBD </p>|
|0x2018|BYTE|<p>The location data is sent by the equipment in the mode of first- in first- out by default </p><p>0x00: First- in first- out (by default) </p><p>0x01: Real- time priority </p>|
|0x2019|BYTE|<p>Data packets before and after a few seconds added for three- emergency alarm: The added data is mainly for 0200. </p><p>0x00-0x0A, the maximum is 10 seconds, the default is 0 seconds, that is, the function is turned off. </p>|
|0x201A|BYTE|<p>Reading instructions on fault code: </p><p>0x01: Read the OBD fault code and report it through 0900 F2. </p><p>0x00: No reading fault code: </p>|
|0x201B|STRING|Schedule of WIFI parameter |
|<a name="_hlt16601537"></a>0x201C|DWORD|The unit is seconds. The minimum sleep wake-up time is 5 minutes, that is, 300 seconds |
|0x201D|WORD|Rapid acceleration threshold, unit: mg |
|0x201E|WORD|Rapid deceleration threshold, unit: mg |
|0x201F|WORD|Sharp turn speed threshold, unit: mg |
|0x2020|BYTE[2]|<p>Emergency braking parameters: The specific principle is described in accordance with the alarm reporting </p><p>BYTE[0]: Velocity difference threshold default 9 km/h</p><p>BYTE[1]: The vehicle speed is greater than a certain speed, default 0 km/h</p>|
|0x2021|BYTE|<p>Emergency braking parameters: The specific principle is described in accordance with the alarm reporting </p><p>Velocity difference threshold default 18 km/h</p>|
|0x2022|WORD|<p>Over-speed parameters: The specific principle is described in accordance with the alarm reporting </p><p>Engine speed threshold 2400 rpm </p>|
|0x2023|WORD|<p>PTO idle parameters: The specific principle is described in accordance with the alarm reporting </p><p>Engine speed threshold 1000 rpm </p>|
|0x2024|BYTE|Log upload, 1 is on (it will be automatically closed after upload for 20 minutes) 0 is off |
|0x2025|WORD|Delay in flameout duration, unit s |
|0x2026|BYTE|Whether ACC line set is effective; Ineffective: 0x00; Effective: 0x01 |
|0x2027|WORD|The time interval of sending OBD data flow indicates that the interval between frames is not less than 70ms, and the default value is 70ms |
|0x2028|WORD|Interval of OBD data return, in 0200 data, the default is 60s |
|0x2029|BYTE[n]|Bluetooth authorization code: A string of hexadecimal characters, up to 50 bytes. W/R|
|0x202A|STRING|Bluetooth name: Greater than 8 bytes, up to 35 bytes. W/R|
|0x202B|STRING|Bluetooth MAC: Such as 44A6E5148CFE. R|


## <a name="_附表wifi参数附表"></a><a name="_toc161247079"></a>**3.9 Schedule- Schedule of WIFI parameter** 

|S/N |Contents |Number of bytes |Data type |Description |
| :-: | :-: | :-: | :-: | :- |
|0|Enable|1|ASCII|Enabling:   0: OFF 1 ON |
|1|,|1|ASCII|0x2C|
|2|SSID|Variable length |ASCII|Wifi SSID|
|3|,|1|ASCII|0x2C|
|4|Password|Variable length |ASCII |WiFi password |

## <a name="_急加速参数附表"></a><a name="_急加速参数"></a><a name="_toc161247080"></a><a name="_附表_急加速参数"></a>**3.10 Schedule- Rapid acceleration parameters [](#_附表_终端参数设置各参数项定义及说明)**
|Byte position |Contents |Number of bytes |Data type |Description |
| :-: | :-: | :-: | :-: | :-: |
|0|Level of rapid acceleration |1|BYTE|0X03: Highly sensitive; 0X02: Moderately sensitive; 0X01: Lowly sensitive, 0X00: OFF |

## <a name="_急减速参数附表"></a><a name="_急减速参数"></a><a name="_toc161247081"></a>**3.11 Schedule- Rapid deceleration parameters [](#_附表_终端参数设置各参数项定义及说明)**

|Byte position |Contents |Number of bytes |Data type |Description |
| :-: | :-: | :-: | :-: | :-: |
|0|Level of rapid deceleration |1|BYTE|0X03: Highly sensitive; 0X02: Moderately sensitive; 0X01: Lowly sensitive, 0X00: OFF |

## <a name="_急拐弯参数附表"></a><a name="_急拐弯参数"></a><a name="_toc161247082"></a>**3.12	Schedule-Sharp turning parameters [](#_附表_终端参数设置各参数项定义及说明)**

|Byte position |Contents |Number of bytes |Data type |Description |
| :-: | :-: | :-: | :-: | :-: |
|0|Level of sharp turn |1|BYTE|0X03: Highly sensitive; 0X02: Moderately sensitive; 0X01: Lowly sensitive, 0X00: OFF |


## <a name="_toc161247083"></a><a name="_附表_终端升级数据包"></a>**3.13	Schedule- data packet of terminal upgrade** 
|Byte position |Contents |Number of bytes |Data type |Description |
| :-: | :-: | :-: | :-: | :- |
|0|Length of firmware version number |1|BYTE|The length of character string of version number of the upgrade file |
|1|Firmware version number |N|BYTE[N]|Character string |
|1+n|Offset address |4|BYTE|Identifies the data offset address of the currently requested upgrade file, starting from 0 |
|1+n+4|Length of requested data |4|BYTE|Identify the length of data currently requested |

## <a name="_toc161247084"></a><a name="_附表_平台升级数据包应答"></a>**3.14	Schedule-response of platform upgrade packet** 
|Byte position |Contents |Number of bytes |Data type |Description |
| :-: | :-: | :-: | :-: | :-: |
|0|Length of firmware version number |1|BYTE|The length of character string of version number of the upgrade file |
|1|Firmware version number |N|BYTE[N]|Character string |
|1+n|Total size of upgrade file |4|Dword||
|1+n+4|Total verification of upgrade file |4|Dword||
|1+n+8|Offset address |4|Dword|Identifies the data offset address of the currently requested upgrade file, starting from 0 |
|1+n+12|Contents of upgrade packet |N|BYTE[N]|Contents of data currently distributed |

## <a name="_拖车报警参数"></a><a name="_拖车报警参数附表"></a><a name="_toc161247085"></a>**3.15	Schedule-Trailer alarm parameters [](#_附表_终端参数设置各参数项定义及说明)**

|S/N |Contents |Number of bytes |Data type |Description |
| :-: | :-: | :-: | :-: | :-: |
|0|Enable|1|ASCII|Alarm shielding: 0 no alarm 1 alarm |
|1|,|1|ASCII|0x2C|
|2|Tow Spd|Variable length |ASCII|Trailer threshold speed (unit KM/H, greater than 15KM/H) |
|3|,|1|ASCII|0x2C|
|4|Tow Interval|Variable length |ASCII |Duration of trailer condition (unit in seconds, greater than 20 seconds) |

## <a name="_碰撞报警参数包附表"></a><a name="_碰撞报警参数包"></a><a name="_toc161247086"></a>**3.16	Schedule- Collision alarm parameter packet [](#_附表_终端参数设置各参数项定义及说明)**

|Byte position |Contents |Number of bytes |Data type |Description |
| :-: | :-: | :-: | :-: | :-: |
|0|Collision level |1|BYTE|<p>0X03: Highly sensitive; </p><p>0X02: Moderately sensitive; </p><p>0X01: Lowly sensitive, </p><p>0X00: OFF </p>|

## <a name="_toc161247087"></a><a name="_附表_特权号列表"></a>**3.17	Schedule-List of privilege number [](#_附表_终端参数设置各参数项定义及说明)**
|Byte position |Contents |Number of bytes |Data type |Description |
| :-: | :-: | :-: | :-: | :- |
|0|Privilege number |11|ASCII|13866668888, indicating that the number allows the configuration of query parameters. |

## <a name="_查询终端参数应答消息体数据"></a><a name="_查询终端参数应答消息体附表"></a><a name="_toc161247088"></a>**3.18	Schedule-Message body of query terminal parameter response [](#_[0104]查询终端参数应答)**

|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Serial number of response |WORD|Serial numbers of corresponding terminal parameters query |
|2|Total number of parameters |BYTE||
|3|List of parameter items ||[Schedule of parameter item format ](#_参数项格式附表)|

## <a name="_hlt491358806"></a><a name="_终端控制消息体数据"></a><a name="_终端控制消息体附表"></a><a name="_toc161247089"></a><a name="_附表_终端控制消息体"></a>**3.19	Schedule-Message body of terminal control [](#_[8105]终端控制)**
|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Command word |BYTE|[Descriptions of the terminal control command words ](#_终端控制命令字说明附表)|
|1|Command parameter |STRING|For the format of command parameters, see the following description, and half width is used between fields"; "To separate each string field, it is to process it according to GBK code, and then form a message |
## <a name="_终端控制命令字说明"></a><a name="_终端控制命令字说明附表"></a><a name="_toc161247090"></a><a name="_附表_终端控制命令字说明"></a>**3.20	Schedule -Descriptions of the terminal control command words [](#_附表_终端控制消息体)**
|Command word |Command parameter |Descriptions and requirements |
| :-: | :-: | :-: |
|0x01|[Command parameter format ](#_命令参数格式附表)|<a name="_hlt37344838"></a>Wireless upgrade. Parameters are separated by half -width semicolons. The command is as follows: "URL address; name of dialing point; name of dialing user; dialing password; IP address; TCP port; UDP port; manufacturer ID; hardware version; firmware version; time limits of connecting to the specified server", and empty the parameter if it has no value. |
|0x04|None |Reset the terminal |
|0x05|None |Restore factory settings of the terminal |
|0x90|1byte+10byte|Fuel cut and power outage  1byte 0x01: Fuel cut            10byte reserved|
|0x91|1byte+10byte|Fuel and power on        1byte 0x01: Fuel on             10byte reserved|
|0x92|1byte+10byte|Ignition                1byte 0x01: Igniting             10byte reserved|
|0x93|1byte+10byte|Flameout               1byte 0x01: Flameout           10byte reserved|
|0xA0|1byte+10byte|Placing an order/ returning the vehicle  1byte 0x01: Placing an order; 0x00: Returning the vehicle     10byte reserved|
|0xA1|1byte+10byte|Searching a vehicle       1byte 0x01: Horn; 0x02: Light; 0x03: Horn+light;    10byte reserved|
|0xA2|1byte+10byte|Central locking          1byte 0x01: Unlocking; 0x00: Locking      10byte reserved|
|0xA3|1byte+10byte|Window               1byte 0x01: Opening the window; 0x00: Closing the window; 10byte reserved|
|0xA4|1byte+10byte|Trunk lock             1byte 0x01: Opening the trunk; 0x00: Closing the trunk    10byte reserved |
|0xA5|1byte+10byte|Air conditioner          1byte 0x01: Turning on air conditioner; 0x00: Turning off air conditioner      10byte reserved|
|0xA6|1byte+10byte|Wiper                 1byte 0x01: Turn on the wiper; 0x00: Turn off the wiper   10byte reserved|
|0xA7|1byte+10byte|Sunroof      1byte 0x01: Turn on the sunroof; 0x00: Turn off the sunroof    10byte reserved|
|0xF1|None |Start of OTA upgrade of GSM module |
## <a name="_toc161247091"></a>**3.21 Schedule- Message body of terminal control**

|Starting byte |Field|Data type|Descriptions and requirements|
| :-: | :-: | :-: | :-: |
|0|Serial number of response|WORD|The corresponding serial number of the platform message|
|1|Command parameter|BYTE[N]|Terminal control response |
## <a name="_toc161247092"></a>**3.22	Schedule- Terminal control response** 
0x0105 command needs to be answered only when the control command word is in the following table 

|Command word |Command parameter |Descriptions and requirements |
| :-: | :-: | :-: |
|0x90|1byte+10byte|Fuel cut and power outage    1byte terminal control response result        10byte reserved |
|0x91|1byte+10byte|Fuel and power on         1byte terminal control response result         10byte reserved|
|0x92|1byte+10byte|Ignition                  1byte terminal control response result        10byte reserved|
|0x93|1byte+10byte|Flameout                1byte terminal control response result         10byte reserved|
|0xA0|1byte+10byte|Placing an order/ returning the vehicle  1byte terminal control response result  10byte reserved|
|0xA1|1byte+10byte|Searching a vehicle       1byte terminal control response result          10byte reserved|
|0xA2|1byte+10byte|Central locking           1byte terminal control response result         10byte reserved|
|0xA3|1byte+10byte|Window                 1byte terminal control response result          10byte reserved|
|0xA4|1byte+10byte|Trunk lock               1byte terminal control response result          10byte reserved|
|0xA5|1byte+10byte|Air conditioner            1byte terminal control response result         10byte reserved|
|0xA6|1byte+10byte|Wiper                  1byte terminal control response result         10byte reserved|
|0xA7|1byte+10byte|Sunroof                 1byte terminal control response result        10byte reserved|
## <a name="_命令参数格式"></a><a name="_命令参数格式附表"></a><a name="_toc161247093"></a>**3.23 	Schedule- Terminal control response result**

|Control response result|Data type|Descriptions and requirements|
| :-: | :-: | :-: |
|0x00|BYTE|Control successful|
|0x01|BYTE|Control failed (command not supported/this function is not supported)|
|0x02|BYTE|Due to the vehicle is not turned off, control failed.|
|0x03|BYTE|Due to the handbrake is not engaged, control failed.|
|0x04|BYTE|Due to the vehicle speed is not zero, control failed.|
|0x05|BYTE|Due to the left front door is unlocked, control failed.|
|0x06|BYTE|Due to the right front door is unlocked, control failed.|
|0x07|BYTE|Due to the left rear door is unlocked, control failed.|
|0x08|BYTE|Due to the right rear door is unlocked, control failed.|
|0x09|BYTE|Due to the left front window is open, control failed.|
|0x0A|BYTE|Due to the right front window is open, control failed.|
|0x0B|BYTE|Due to the left rear window is open, control failed.|
|0x0C|BYTE|Due to the right rear window is open, control failed.|
|0x0D|BYTE|Due to the sunroof is not closed, control failed.|
|0x0E|BYTE|Due to the left front door is open, control failed.|
|0x0F|BYTE|Due to the right front door is open, control failed.|
|0x10|BYTE|Due to the left rear door is open, control failed.|
|0x11|BYTE|Due to the right rear door is open, control failed.|
|0x12|BYTE|Due to the front compartment hood is open, control failed.|
|0x13|BYTE|Due to the rear trunk is open, control failed.|
|0x14|BYTE|Due to the reading light is on, control failed.|
|0x15|BYTE|Due to the low beam headlights are on, control failed.|
|0x16|BYTE|Due to the high beam headlights are on, control failed.|
|0x17|BYTE|Due to the front fog lights are on, control failed.|
|0x18|BYTE|Due to the rear fog lights are on, control failed.|
|0x19|BYTE|Due to the hazard lights are on, control failed.|
|0x1A|BYTE|Due to the width lights are on, control failed.|
|0x1B|BYTE|Due to the turn signals are on, control failed.|
|0x1C|BYTE|Due to the wipers are on, control failed.|
|0x1D|BYTE|Due to the air conditioning is on, control failed.|
|0x1E|BYTE|Due to the vehicle is not in P gear, control failed.|
|0x1F|BYTE|Due to the vehicle is not in N gear, control failed.|
|0x20|BYTE|Due to the terminal commands control action to the vehicle, and the vehicle fails to respond within the timeout period, control failed.|
|0x21|BYTE|Due to the doors are not closed, control failed.|
|0x22|BYTE|Due to the door locks are not closed, control failed.|
|0x23|BYTE|Due to the windows are not closed, control failed.|
|0x24|BYTE|Due to the vehicle door is opened and then closed without the vehicle being turned off, and with the driver's door closed, causing control failure. (This command is for BMW vehicles, ensuring locking when the instrument cluster and central control are turned off.)|

## <a name="_toc161247094"></a>**3.24	Schedule- Format of command parameters [](#_附表_终端控制命令字说明)**

|Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: |
|Connection control |BYTE|<p>0x00: It enters the emergency state once it is switched to the specified supervision platform server. In this state, only the supervision platform that issues the control command can send the control command including SMS: </p><p>0x01: It is switched back to the original default monitoring platform server, and recovered to the normal state. </p>|
|Name of dialing point |STRING|<p>It is server APN, dialing access point of wireless communications. </p><p>It shall be PPP dialing numbers if the network type is CDMA. </p>|
|Name of dialing user |STRING|User name of wireless communication dialing of the server |
|Dialing password |STRING|Password of wireless communication dialing of the server |
|Address |STRING|Server address, IP or domain name |
|TCP port |WORD|TCP port of server |
|UDP port |WORD|<p>UDP port of server </p><p>Hidden function 0xAA: Switch to FTP upgrade server </p><p>Hidden function 0xBB: Switch to TCP upgrade server </p>|
|Manufacturer ID |BYTE[5]|Terminal manufacturer code |
|Authentication code of monitoring platform |STRING|The authentication code issued by the supervision platform is only used for terminal connection For authentication after the supervision platform, the original authentication code is still used when the terminal connects the original monitoring platform |
|Hardware version |STRING|The hardware version number of the terminal shall be customized by the manufacturer |
|Firmware version |STRING|The firmware version number of the terminal shall be customized by the manufacturer |
|URL address |STRING|Complete URL address |
|Time limit of connecting to the specified server |WORD|Unit: Minute (min), if the value is not 0, it means that the terminal should connect back to the original address before the expiration of the validity period after the terminal receives the instruction to upgrade or connect to the specified server. If the value is 0, it means that it is always connected to the specified server |
