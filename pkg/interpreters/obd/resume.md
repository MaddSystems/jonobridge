




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

## <a name="_附表_驾驶员信息采集附表"></a><a name="_toc161247095"></a>**3.25	Schedule-Collection of driver information** 

<table><tr><th colspan="1" valign="bottom">Starting byte </th><th colspan="1" valign="bottom">Field </th><th colspan="1" valign="bottom">Data type </th><th colspan="1" valign="bottom">Descriptions and requirements </th></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td></tr>
<tr><td colspan="1" rowspan="2" valign="bottom">0</td><td colspan="1" rowspan="2" valign="bottom">Status </td><td colspan="1" rowspan="2" valign="bottom">BYTE</td><td colspan="1" valign="bottom">0x01: Qualification certificate IC card is inserted (When the driver is on duty);</td></tr>
<tr><td colspan="1" valign="bottom"></td></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom">0x02: Qualification certificate IC card is pulled out (When the driver is off duty); </td></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td></tr>
<tr><td colspan="1" rowspan="2" valign="bottom">1</td><td colspan="1" rowspan="2" valign="bottom">Time </td><td colspan="1" rowspan="2" valign="bottom">BCD[6]</td><td colspan="1" valign="bottom">Time for inserting / pulling out the card, YY-MM-DD-hh-mm-ss; </td></tr>
<tr><td colspan="1" valign="bottom"></td></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom">The following fields are only valid in the state of 0x01 and filled. </td></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom">0x00: IC card is successfully read; </td></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom">0x01: Failing in reading because the card key fails to pass authentication; </td></tr>
<tr><td colspan="1" rowspan="2" valign="bottom">7</td><td colspan="1" rowspan="2" valign="bottom">Reading results of IC card </td><td colspan="1" rowspan="2" valign="bottom">BYTE</td><td colspan="1" valign="bottom">0x02: Reading failed because the card has been locked; </td></tr>
<tr><td colspan="1" rowspan="2" valign="bottom">0x03: Failing in reading because the card has been pulled out; </td></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom">0x04: Failing in reading because of error in data verification. </td></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom">The following fields are valid when the reading result of IC card is equal to 0x00. </td></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td></tr>
<tr><td colspan="1" valign="bottom">8</td><td colspan="1" valign="bottom">The length of the driver's name </td><td colspan="1" valign="bottom">BYTE</td><td colspan="1" valign="bottom">n</td></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td></tr>
<tr><td colspan="1" valign="bottom">9</td><td colspan="1" valign="bottom">Name of the driver </td><td colspan="1" valign="bottom">STRING</td><td colspan="1" valign="bottom">Name of the driver </td></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td></tr>
<tr><td colspan="1" valign="bottom">9+n</td><td colspan="1" valign="bottom">Qualification certificate code </td><td colspan="1" valign="bottom">STRING</td><td colspan="1" valign="bottom">The length is 20 bits, and it shall be supplemented by 0x00 if it is insufficient. </td></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td></tr>
<tr><td colspan="1" valign="bottom">29+n</td><td colspan="1" valign="bottom">The length of the name of certifying authority </td><td colspan="1" valign="bottom">BYTE</td><td colspan="1" valign="bottom">m</td></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td></tr>
<tr><td colspan="1" valign="bottom">30+n</td><td colspan="1" valign="bottom">Name of the certifying authority </td><td colspan="1" valign="bottom">STRING</td><td colspan="1" valign="bottom">Name of the certifying authority issuing qualification certificate </td></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td></tr>
<tr><td colspan="1" valign="bottom">30+n+m</td><td colspan="1" valign="bottom">Validity period of certificate </td><td colspan="1" valign="bottom">BCD[4]</td><td colspan="1" valign="bottom">YYYYMMDD</td></tr>
<tr><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td><td colspan="1" valign="bottom"></td></tr>
</table>

## <a name="_toc161247096"></a><a name="_toc534810610"></a><a name="_附表_临时位置跟踪控制消息体"></a>**3.26	Schedule- Message body of temporary location tracking control** 
|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Time interval |WORD|The unit is (s), and 0 represents stopping tracking. There is no need for subsequent field when tracking is stopped. |
|2|Valid time of location tracking |DWORD|The unit is (s), after receiving the location tracking control message, the terminal sends the location report according to the time interval in the message before the expiration of the validity period |

## <a name="_toc161247097"></a><a name="_附表_终端升级结果数据包"></a>**3.27	Schedule- data packet of terminal upgrade result** 
|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Upgrade type |BYTE|<p>0x00: Terminal </p><p>0x12: IC card reader for road transportation license </p><p>0x52: Beidou </p><p>0x2A:RCM module </p>|
|1|Upgrade results |BYTE|<p>0x00: Success </p><p>0x01: Failure (time out) </p><p>0x02: Cancel </p><p>0x03:  NEODOWNLOAD:NULL</p><p>0x04:  NEODOWNLOAD:FAIL</p><p>0x05:  NEOUPDATE:FAIL</p><p>0x06:  NEOUPDATE:NULL</p><p>0xF0: No upgrade for the same version </p><p>0xF1: Error in attribute of upgrade version </p><p>0xF2: Error in verification of upgrade version </p><p>0xF3: Upgrade file does not exist </p>|


## <a name="_toc161247098"></a><a name="_附表_位置作息查询应答消息体数据格式"></a>**3.28	Schedule- Message body of data format of location information query response** 
|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Serial number of response |WORD|Serial numbers of corresponding information inquiry message |
|2|Location information report ||See message body of location data |

## <a name="_toc161247099"></a><a name="_附表_位置数据批量汇报"></a>**3.29	Schedule-Batch report packet of location data** 
|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Number of data items (packets) |WORD|Number of location data items (packets) included **N**, > 0 |
|2|Type of Location data item |BYTE|<p>0: Normal batch data </p><p>1: Supplementary report of blind spot </p>|
|3|Data items of location report ||<p>[Data item format of location batch report (packet 1) ](#_附表_位置汇报数据项数据格式)</p><p>[Data item format of location batch report (packet 2) ](#_附表_位置汇报数据项数据格式)</p><p>[Data item format of location batch report (packet 3) ](#_附表_位置汇报数据项数据格式)</p><p>…</p><p>[Data item format of location batch report (packet **N**) ](#_附表_位置汇报数据项数据格式)</p>|

**Note: Upload multiple packets at one time** 

## <a name="_toc161247100"></a><a name="_附表_位置汇报数据项数据格式"></a>**3.30	Schedule-Data item format of location batch report** 
|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Length of location report data body |Word|Length of location report data body, N |
|<a name="_位置数据信息体附表"></a><a name="_位置数据信息"></a>2|Location report message body |BYTE[n]|Message body of location data |

## <a name="_toc161247101"></a><a name="_附表_位置数据信息体"></a>**3.31	Schedule-Location report message body** 
|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Alarm mark |DWORD|See the schedule of definition of alarm mark bit for details [](#_报警标志位定义附表)|
|4|Status |DWORD|See the schedule of definition of status mark bit for details [](#_状态标志位定义附表)|
|8|Latitude |DWORD|Latitude value in the unit of degrees is multiplied by the sixth power of 10, accurate to one millionth of a degree. |
|12|Longitude |DWORD|Latitude value in the unit of degrees is multiplied by the sixth power of 10, accurate to one millionth of a degree. |
|16|Elevation |WORD|Altitude, in unit of meters (m) |
|18|Speed |WORD|1/10km/h|
|20|Direction |WORD|0-359, the due north is 0, clockwise |
|22|Time |BCD[6]|YY-MM-DD-hh-mm-ss (GMT + 8 equipment reporting adopts Beijing time benchmark) |
|28|List of additional information of position |nByte|See list of additional information of position for details [](#_附表_位置附加信息表)|

## <a name="_hlt20736495"></a><a name="_hlt41309726"></a><a name="_hlt36396052"></a><a name="_hlt24447438"></a><a name="_hlt22569295"></a><a name="_状态标志位定义附表"></a><a name="_状态位定义"></a><a name="_toc161247102"></a>**3.32	Schedule- Definition of status mark bit** 

|Bit |Status |
| :-: | :-: |
|0  |0: ACC off; 1: ACC on|
|1|0: Not positioned 1: Positioned|
|2  |0: Northern latitude: 1: Southern latitude |
|3   |0: East longitude; 1: West longitude |
|4|0: Operating status: 1: Outage status |
|5|0: Longitude and latitude are not encrypted by confidentiality plug-in; l: Longitude and latitude are encrypted by confidentiality plug-in |
|6-9|Reserved |
|10    |0: Normal oil-way of the vehicle: 1: Disconnection of oil-way of the vehicle |
|11  |0: Normal circuits of the vehicle; 1: Disconnection of circuits of the vehicle |
|12|0: Doors unlocked; 1: Doors locked |
|13|0: Normal mode;   1: Maintenance mode |
|14|0: WIFI off; 1: WIFI on |
|15|0: Module 433 of tire pressure is normal; 1: Module 433 of tire pressure is abnormal |
|16|0: Bluetooth is normal; 1: Bluetooth is abnormal |
|17|0: The bucket lifting status of the hopper car is not lifted, 1: The bucket lifting status of the hopper car is lifted |
|18-31||

## <a name="_报警标志位定义"></a><a name="_报警标志位定义附表"></a><a name="_toc161247103"></a>**3.33	Schedule-Definition of alarm mark bits [](#_附表_位置数据信息体)**

|Bit |Definition |Processing instructions |
| :-: | :-: | :-: |
|0|1: The emergency warning is triggered after the alarm switch is touched |Clear it after receiving the response. |
|1|1: Over-speed alarm |The mark shall be maintained until the alarm condition is removed. |
|2|1: Fatigue driving |The mark shall be maintained until the alarm condition is removed. |
|3|1: Hazard warning |Clear it after receiving the response. |
|4|1: GNSS module failure |The mark shall be maintained until the alarm condition is removed. |
|5|1: GNSS antenna is unconnected or cut |The mark shall be maintained until the alarm condition is removed. |
|6|1: GNSS antenna is short circuited |The mark shall be maintained until the alarm condition is removed. |
|7|1: Under-voltage of main power of terminal |The mark shall be maintained until the alarm condition is removed. |
|8|1: Power down of main power of terminal |The mark shall be maintained until the alarm condition is removed. |
|9|1: LCD or display failure of terminal |The mark shall be maintained until the alarm condition is removed. |
|10|1: TTS module failure |The mark shall be maintained until the alarm condition is removed. |
|11|1: Camera failure |The mark shall be maintained until the alarm condition is removed. |
|12|Reserved ||
|13|Overspeed warning|The mark shall be maintained until the alarm condition is removed. |
|14-17|||
|18|1: Cumulative driving time-out of the day |The mark shall be maintained until the alarm condition is removed. |
|19|1: Over-time parking |The mark shall be maintained until the alarm condition is removed. |
|20|1: Access regions |Clear it after receiving the response. |
|21|1: Incoming and outgoing routes |Clear it after receiving the response. |
|22|1: Insufficient / too long driving time on the section |Clear it after receiving the response. |
|23|1: Route deviation alarm |The mark shall be maintained until the alarm condition is removed. |
|24|1: Vehicle VSS failure |The mark shall be maintained until the alarm condition is removed. |
|25|1: Abnormal fuel volume of the vehicle |The mark shall be maintained until the alarm condition is removed. |
|26|1: The vehicle is stolen(by vehicle anti-theft device) |The mark shall be maintained until the alarm condition is removed. |
|27|1: Illegal ignition of the vehicle |Clear it after receiving the response |
|28|1: Illegal displacement of the vehicle |Clear it after receiving the response |
|29-31|Reserved ||

<a name="_位置附加信息附表"></a><a name="_位置附加信息项格式"></a><a name="_附表_位置附加信息"></a>
## <a name="_toc161247104"></a><a name="_附表_位置附加信息表"></a>**3.34	Schedule-List of additional information of position [](#_附表_位置数据信息体)**
|Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: |
|Additional information ID |BYTE|1-255|
|Length of additional information |BYTE|1-255|
|Additional information ||[Schedule of definition of additional information ](#_附加信息定义附表)|

## <a name="_附加信息定义"></a><a name="_toc161247105"></a><a name="_附表_附加信息定义"></a><a name="_附加信息定义附表"></a>**3.35	Schedule- Definition of additional information** 
|Additional information ID (1byte) |<p>Additional information </p><p>Length (1byte) </p>|Descriptions and requirements |
| :-: | :-: | :-: |
|0xE1|2byte|Rotation speed, unit: RPM; Offset: 0; Scope 0-8000. (Special purpose) |
|0xEA|Nbyte|Data packet includes sub ID (2BYTE), length (1BYTE) + data (NBYTE) schedule of basic data flow [](#_附表_基础数据流)|
|0xEB|Nbyte|Data packet includes sub ID (2BYTE), length (1BYTE) + data (NBYTE) extended data flow of a car [](#_轿车扩展数据流\<一\>附表)|
|<a name="_hlt54343437"></a><a name="_hlt534369304"></a><a name="_hlt54358713"></a>0xEC|Nbyte|Data packet includes sub ID (2BYTE), length (1BYTE) + data (NBYTE) extended data flow of a truck [](#_货车扩展数据流\<一\>附表)|
|<a name="_hlt27212223"></a>0xED|Nbyte|Data packet includes sub ID (2BYTE), length (1BYTE) + data (NBYTE) data items of new energy vehicles [](#_附表_新能源汽车数据项\<一\>)|
|0xEE|Nbyte|Data packet includes sub ID (2BYTE), length (1BYTE) + data (NBYTE) schedule of extended peripheral data items [](#_附表_扩展外设数据流)|
|<a name="_hlt22569302"></a><a name="_hlt20736497"></a><a name="_hlt24447446"></a><a name="_hlt20734247"></a><a name="_hlt22569385"></a>0xFA|Nbyte|Data packet includes sub ID (2BYTE), length (1BYTE) + data (NBYTE) schedule of alarm command ID and description [](#_附表_报警命令id及描述数据包)|
|<a name="_hlt54601415"></a><a name="_hlt41309734"></a><a name="_hlt534379278"></a><a name="_hlt534399030"></a><a name="_hlt534399046"></a><a name="_hlt534379392"></a>0xFB|Nbyte|Data packet includes sub ID (2BYTE), length (1BYTE) + data (NBYTE) data flow of base station, reporting when GPS is not positioned, customization [](#_附表_基站数据流)|
|<a name="_hlt54706322"></a><a name="_hlt54601641"></a>0xFC|Nbyte|Data packet includes sub ID (2BYTE), length (1BYTE) + data (NBYTE) data packet includes sub ID (2BYTE), length (1BYTE) + data (NBYTE) Wifi data flow, reporting when GPS is not positioned, customization|
|...|...|Others reserved |

Additional ID: 

0XEA: The following corresponding data item represents the basic data item, with the maximum length of 255; 

0XEB: The following corresponding data item represents the data item of a car, with the maximum length of 255; 

0XEC: The following corresponding data item represents the data item of a truck, with the maximum length of 255; 

0XED: The following corresponding data item represents the data item of new energy car, with the maximum length of 255; 

0XEE: The following corresponding data item represents the peripheral data item, with the maximum length of 255; 

0XFA: The following corresponding data item represents the alarm event ID, with the maximum length of 255; 

0XFB: The following corresponding data item represents the data flow of base station, with the maximum length of 255; 

0XFC: The following corresponding data item represents the Wifi data flow, with the maximum length of 255;





<a name="_超速报警附加信息消息体数据"></a><a name="_超速报警附加信息附表"></a>
## <a name="_toc161247106"></a><a name="_附表_基础数据流"></a>**3.36	Schedule-Basic data flow [](#_附表_附加信息定义)**
<table><tr><th colspan="1"><a name="_hlt22569299"></a><b>Functional ID domain</b> </th><th colspan="1"><b>Function ID [2]</b> </th><th colspan="1"><b>Length [1]</b> </th><th colspan="1"><b>Function</b> </th><th colspan="1"><b>Unit</b> </th><th colspan="1"><b>Description</b> </th></tr>
<tr><td colspan="1" rowspan="25">0x0001-0x0FFF</td><td colspan="1">0x0001</td><td colspan="1">4</td><td colspan="1" valign="top">Reserved </td><td colspan="1"></td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0002</td><td colspan="1">4</td><td colspan="1">Reserved </td><td colspan="1"></td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0003</td><td colspan="1">5</td><td colspan="1">Data of total mileage </td><td colspan="1">Meter </td><td colspan="1"><a name="总里程数据格式表a"></a>[Format table of total mileage ](#_附表_基础数据项：总里程格式表)</td></tr>
<tr><td colspan="1"><a name="_hlt36396361"></a><a name="_hlt36396430"></a>0x0004</td><td colspan="1">5</td><td colspan="1">Data of total fuel consumption </td><td colspan="1">ml </td><td colspan="1">[Format table of total fuel consumption ](#_附表_基础数据项：总耗油量格式表)</td></tr>
<tr><td colspan="1">0x0005</td><td colspan="1">4</td><td colspan="1">Total run time </td><td colspan="1">Second </td><td colspan="1">Cumulative total duration of vehicle operation </td></tr>
<tr><td colspan="1">0x0006</td><td colspan="1">4</td><td colspan="1">Total flameout time </td><td colspan="1">Second </td><td colspan="1">Cumulative total duration of vehicle flameout </td></tr>
<tr><td colspan="1">0x0007</td><td colspan="1">4</td><td colspan="1">Total idle time </td><td colspan="1">Second </td><td colspan="1">Cumulative total duration of vehicle idling </td></tr>
<tr><td colspan="1"><s>0x0008</s></td><td colspan="1"><s>N</s></td><td colspan="1"><s>Sheet of mileage data</s> </td><td colspan="1"></td><td colspan="1"><s>Mileage reference packet 60 bytes</s> </td></tr>
<tr><td colspan="1"><s>0x0009</s></td><td colspan="1"><s>N</s></td><td colspan="1"><s>Sheet of fuel consumption data</s> </td><td colspan="1"></td><td colspan="1"><s>Fuel consumption reference packet 35 bytes</s> </td></tr>
<tr><td colspan="1">0x0010</td><td colspan="1">N</td><td colspan="1">Accelerometer </td><td colspan="1"></td><td colspan="1">[Accelerometer ](#_附表_基础数据项：加速度表)</td></tr>
<tr><td colspan="1">0x0011</td><td colspan="1">20</td><td colspan="1">Sheet of vehicle status </td><td colspan="1"></td><td colspan="1">[Sheet of vehicle status ](#_附表_基础数据项：车辆状态表)</td></tr>
<tr><td colspan="1">0x0012</td><td colspan="1">2</td><td colspan="1">Vehicle voltage </td><td colspan="1">0\.1V</td><td colspan="1">0-36V</td></tr>
<tr><td colspan="1">0x0013</td><td colspan="1">1</td><td colspan="1">Built-in battery voltage of terminal </td><td colspan="1">0\.1V</td><td colspan="1">0-5V</td></tr>
<tr><td colspan="1">0x0014</td><td colspan="1">1</td><td colspan="1">CSQ value </td><td colspan="1"></td><td colspan="1">Strength of network signal </td></tr>
<tr><td colspan="1">0x0015</td><td colspan="1">2</td><td colspan="1">Model ID </td><td colspan="1"></td><td colspan="1">Sheet of Model ID </td></tr>
<tr><td colspan="1">0x0016</td><td colspan="1">1</td><td colspan="1">OBD protocol type </td><td colspan="1"></td><td colspan="1">[Sheet of protocol type ](#_附表_基础数据项：协议类型表)</td></tr>
<tr><td colspan="1">0x0017</td><td colspan="1">2</td><td colspan="1">Driving cycle label </td><td colspan="1"></td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0018</td><td colspan="1">1</td><td colspan="1">The number of satellites by GPS </td><td colspan="1"></td><td colspan="1">The number of satellites by GPS </td></tr>
<tr><td colspan="1">0x0019</td><td colspan="1">2</td><td colspan="1">Positional accuracy by GPS </td><td colspan="1">0\.01</td><td colspan="1">Positional accuracy by GPS </td></tr>
<tr><td colspan="1">0x001A</td><td colspan="1">1</td><td colspan="1">Average signal-to-noise ratio of GPS </td><td colspan="1">db</td><td colspan="1">Average signal-to-noise ratio of GPS </td></tr>
<tr><td colspan="1">0x001B</td><td colspan="1">1</td><td colspan="1">Antenna status of GPS </td><td colspan="1"></td><td colspan="1"><p>0: Antenna is normal 1: Open circuit of antenna </p><p>2: Short circuit of antenna (module support required) </p><p>Note: Supported by TBOX products only. </p></td></tr>
<tr><td colspan="1" rowspan="2">0x001D</td><td colspan="1" rowspan="2">1</td><td colspan="1" rowspan="2">Device pull-out status (customized) </td><td colspan="1" rowspan="2"></td><td colspan="1">0x02: Before the first positioning after the equipment is unplugged or powered on </td></tr>
<tr><td colspan="1"><p>Not 0x02: Others </p><p>Note: Avoid the straight line of GPS track from the factory location to the installation location caused by the delivery of the factory test point to the installation location by the customer for the first time </p></td></tr>
<tr><td colspan="1">0x001E</td><td colspan="1">4</td><td colspan="1">Accumulated mileage </td><td colspan="1">Meter </td><td colspan="1">When the mileage type in 0003 total mileage data is instrument mileage, it can only be accurate to 1KM or 10KM, which is not conducive to mileage statistics. The cumulative mileage is added for mileage statistics on the platform </td></tr>
<tr><td colspan="1"><s>0x001F</s></td><td colspan="1"><s>4</s></td><td colspan="1"><s>Instant fuel consumption</s> </td><td colspan="1"><p><s>0.01</s></p><p><s>L/100km</s></p></td><td colspan="1"><s>Dedicated version</s> </td></tr>
<tr><td colspan="1" rowspan="2"></td><td colspan="1">0x0020</td><td colspan="1">2</td><td colspan="1">Ignition type </td><td colspan="1"></td><td colspan="1"><p>BIT0:   1: ACC line ignition </p><p>BIT1:   1: Security monitoring ignition </p><p>BIT2:   1: GPS speed </p><p>BIT3:   1: Voltage </p><p>BIT4:   1: Engine speed</p><p>BIT5:   1: ACC interruption ignition</p><p></p></td></tr>
<tr><td colspan="1">0x0021</td><td colspan="1">4</td><td colspan="1">Carbon emission (g) </td><td colspan="1">g </td><td colspan="1">The cumulative carbon emission g is counted from the installed equipment </td></tr>
<tr><td colspan="1"></td><td colspan="1">0x0022</td><td colspan="1">2</td><td colspan="1">Roll angular velocity (special purpose) </td><td colspan="1">0\.1dps</td><td colspan="1"><p>Bit15 indicates positive and negative, 0: Positive direction; 1: Negative direction </p><p>Bit0-14, indicating value, accuracy 0.1 </p><p>Eg; The upload value 0x80FF indicates that the direction is negative, and the angle size is equal to 25.5dps, 0x8000 indicates that the direction is negative, 0x00FF=255, 255/10=25.5dps </p></td></tr>
<tr><td colspan="1"></td><td colspan="1">0x0023</td><td colspan="1">2</td><td colspan="1">Pitch angular velocity (special purpose) </td><td colspan="1">0\.1dps</td><td colspan="1"><p>Bit15 indicates positive and negative, 0: Positive direction; 1: Negative direction </p><p>Bit0-14, indicating value, accuracy 0.1 </p></td></tr>
<tr><td colspan="1"></td><td colspan="1">0x0024</td><td colspan="1">2</td><td colspan="1">Yaw angular velocity (special purpose) </td><td colspan="1">0\.1dps</td><td colspan="1"><p>Bit15 indicates positive and negative, 0: Positive direction; 1: Negative direction </p><p>Bit0-14, indicating value, accuracy 0.1 </p></td></tr>
<tr><td colspan="1"></td><td colspan="1">0x0025</td><td colspan="1">5</td><td colspan="1">Cumulative mileage 2 (only for SEEWORLD)</td><td colspan="1">Meter</td><td colspan="1">Cumulative mileage 2 format table</td></tr>
</table>


## <a name="_轿车扩展数据流<一>附表"></a><a name="_toc161247107"></a>**3.37	Schedule- Extended data flow of car [](#_附表_附加信息定义)**

<table><tr><th colspan="1"><p><b>Functional IC domain</b> </p><p></p></th><th colspan="1"><p><b>Function</b> </p><p><b>ID[2]</b> </p></th><th colspan="1"><b>Length [1]</b> </th><th colspan="1"><b>Function</b> </th><th colspan="1"><b>Unit</b> </th><th colspan="1"><b>Description</b> </th></tr>
<tr><td colspan="1" rowspan="33"><p>Data item of car (<b>common</b>) </p><p>[0x6001-0x6FFF]</p></td><td colspan="1">0x60C0</td><td colspan="1">2</td><td colspan="1">Speed </td><td colspan="1">rpm</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0 - 8000 </td></tr>
<tr><td colspan="1">0x60D0</td><td colspan="1">1</td><td colspan="1">Vehicle speed </td><td colspan="1">Km/h</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0 - 240 </td></tr>
<tr><td colspan="1">0x62F0</td><td colspan="1">2</td><td colspan="1">Remaining fuel </td><td colspan="1"><p>%</p><p>L</p></td><td colspan="1"><p>Remaining fuel, unit L或% </p><p>Bit15 = = 0% OBD is percentage </p><p>`      `==1 unit L </p><p>Displayed value is uploaded value / 10 </p></td></tr>
<tr><td colspan="1">0x6050</td><td colspan="1">1</td><td colspan="1">Coolant temperature </td><td colspan="1">℃</td><td colspan="1">Accuracy: 1℃ deviation: -40.0℃ scope: -40.0℃ - +210℃ </td></tr>
<tr><td colspan="1">0x60F0</td><td colspan="1">1</td><td colspan="1">Intake temperature </td><td colspan="1">℃</td><td colspan="1">Accuracy: 1℃ deviation: -40.0℃ scope: -40.0℃ - +210℃ </td></tr>
<tr><td colspan="1">0x60B0</td><td colspan="1">1</td><td colspan="1">Intake (manifold absolute) pressure </td><td colspan="1">kPa</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0 - 250kpa </td></tr>
<tr><td colspan="1">0x6330</td><td colspan="1">1</td><td colspan="1">Atmospheric pressure </td><td colspan="1">kPa</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0 - 250kpa</td></tr>
<tr><td colspan="1">0x6460</td><td colspan="1">1</td><td colspan="1">Ambient temperature </td><td colspan="1">℃</td><td colspan="1">Accuracy: 1℃ deviation: -40.0℃ scope: -40.0℃ - +210℃ </td></tr>
<tr><td colspan="1">0x6490</td><td colspan="1">1</td><td colspan="1">Position of accelerator pedal </td><td colspan="1">%</td><td colspan="1">Accuracy: 1 deviation: 0 scope:0% - 100%</td></tr>
<tr><td colspan="1">0x60A0</td><td colspan="1">2</td><td colspan="1">Fuel pressure </td><td colspan="1">kPa</td><td colspan="1">Accuracy: 1 deviation: 0 scope:0 - 500kpa</td></tr>
<tr><td colspan="1">0x6014</td><td colspan="1">1</td><td colspan="1">State of fault code </td><td colspan="1"></td><td colspan="1">The effective scope is 0 - 1, "0" indicates unlit and "1" indicates lit. </td></tr>
<tr><td colspan="1">0X6010</td><td colspan="1">1</td><td colspan="1">Number of fault codes </td><td colspan="1">Piece </td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0-255 </td></tr>
<tr><td colspan="1">0x6100</td><td colspan="1">2</td><td colspan="1">Air flow </td><td colspan="1">g/s</td><td colspan="1">Accuracy: 0.1 deviation: 0 scope: 0-6553.5 </td></tr>
<tr><td colspan="1">0x6110</td><td colspan="1">2</td><td colspan="1">Absolute throttle position </td><td colspan="1">%</td><td colspan="1">Accuracy: 0.1 deviation: 0 scope: 0-6553.5 </td></tr>
<tr><td colspan="1">0x61F0</td><td colspan="1">2</td><td colspan="1">The time since engine start </td><td colspan="1">sec</td><td colspan="1">Accuracy: 1 deviation: 0 </td></tr>
<tr><td colspan="1">0x6210</td><td colspan="1">4</td><td colspan="1">Fault mileage </td><td colspan="1">Km</td><td colspan="1">Accuracy: 1 deviation: 0 </td></tr>
<tr><td colspan="1">0x6040</td><td colspan="1">1</td><td colspan="1">Calculated load value </td><td colspan="1">%</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0% - 100% </td></tr>
<tr><td colspan="1">0x6070</td><td colspan="1">2</td><td colspan="1">Long-term fuel trim (cylinder banks 1 and 3) </td><td colspan="1">%</td><td colspan="1">Accuracy: 0.1 deviation: 0 scope: 0 -6553.5 </td></tr>
<tr><td colspan="1">0x60E0</td><td colspan="1">2</td><td colspan="1">Ignition timing advance angle of the first cylinder </td><td colspan="1">%</td><td colspan="1">Accuracy: 0.1 deviation: -64 </td></tr>
<tr><td colspan="1">0x6901</td><td colspan="1">1</td><td colspan="1">Wear of front brake pad (special purpose) </td><td colspan="1"></td><td colspan="1">0 normal / otherwise, the corresponding data is displayed, unit: Level </td></tr>
<tr><td colspan="1">0x6902</td><td colspan="1">1</td><td colspan="1">Wear of rear brake pad (special purpose) </td><td colspan="1"></td><td colspan="1">0 normal / otherwise, the corresponding data is displayed, unit: Level </td></tr>
<tr><td colspan="1">0x6903</td><td colspan="1">1</td><td colspan="1">Level of brake fluid (special purpose) </td><td colspan="1"></td><td colspan="1">0: Abnormal 1: Normal others: Not available </td></tr>
<tr><td colspan="1">0x6904</td><td colspan="1">2</td><td colspan="1">Oil level (special purpose) </td><td colspan="1"><p>MM</p><p>%</p></td><td colspan="1"><p>Unit MM or % </p><p>Bit15 = = 0% </p><p>`      `==1 unit MM </p><p>After the highest BIT is removed, the accuracy is 0.1 </p></td></tr>
<tr><td colspan="1">0x6905</td><td colspan="1">4</td><td colspan="1"><p>Left front tire pressure BYTE1 (special purpose) </p><p>Right front tire pressure BYTE2 </p><p>Left rear tire pressure BYTE3 </p><p>Right rear tire pressure BYTE4 </p></td><td colspan="1">bar</td><td colspan="1"><p>0xFF: Abnormal; Other values: Unit: bar, precision: 0.1 </p><p>0xFF: Abnormal; Other values: Unit: bar, precision: 0.1 </p><p>0xFF: Abnormal; Other values: Unit: bar, precision: 0.1 </p><p>0xFF: Abnormal; Other values: Unit: bar, precision: 0.1 </p></td></tr>
<tr><td colspan="1">0x6906</td><td colspan="1">2</td><td colspan="1">Coolant level (special purpose) </td><td colspan="1"></td><td colspan="1">Accuracy: 1 deviation: -48 </td></tr>
<tr><td colspan="1">0x6907</td><td colspan="1">4</td><td colspan="1">Mileage (special purpose) </td><td colspan="1">km</td><td colspan="1">Accuracy: 0.1 deviation: 0 </td></tr>
<tr><td colspan="1">0x6060</td><td colspan="1">2</td><td colspan="1">Short- term fuel trim (cylinder banks 1 and 3) (special purpose) </td><td colspan="1"></td><td colspan="1"></td></tr>
<tr><td colspan="1">0x6340</td><td colspan="1">4</td><td colspan="1">Equivalent ratio (lambda) and current of B1-S1 linear or broadband oxygen sensor (special purpose) </td><td colspan="1"><p></p><p>N/A</p><p>mA</p></td><td colspan="1"><p>4 bytes are ABCD respectively </p><p>Equivalent ratio= (A*256+B)*2/65535</p><p>Current= (C*256+D)*8/65535</p></td></tr>
<tr><td colspan="1">0x6430</td><td colspan="1">1</td><td colspan="1">Absolute load value (special purpose) </td><td colspan="1">%</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0% - 100% </td></tr>
<tr><td colspan="1">0x6680</td><td colspan="1">1</td><td colspan="1">Intake air temperature sensor (special purpose) </td><td colspan="1"></td><td colspan="1"></td></tr>
<tr><td colspan="1">0x66f0</td><td colspan="1">1</td><td colspan="1">Turbocharger compressor inlet pressure (special purpose) </td><td colspan="1"></td><td colspan="1"></td></tr>
<tr><td colspan="1">0x6C11</td><td colspan="1">4</td><td colspan="1">Mileage between services (special purpose) </td><td colspan="1">km</td><td colspan="1">Accuracy: 1 deviation: 0 </td></tr>
<tr><td colspan="1">0x6C12</td><td colspan="1">1</td><td colspan="1">Cumulative collision times, sum of front, rear, left and right collision times (special purpose) </td><td colspan="1">Times </td><td colspan="1"></td></tr>
<tr><td colspan="1" rowspan="4"></td><td colspan="1">0x6F01</td><td colspan="1">0x6F01</td><td colspan="1">AEB1 data flow (Special purpose)</td><td colspan="1"></td><td colspan="1">4-byte CAN ID plus an 8-byte data flow, forwarding terminal without parsing</td></tr>
<tr><td colspan="1">0x6F02</td><td colspan="1">0x6F02</td><td colspan="1">AEB2 data flow (Special purpose)</td><td colspan="1"></td><td colspan="1">4-byte CAN ID plus an 8-byte data flow, forwarding terminal without parsing</td></tr>
<tr><td colspan="1">0x6F03</td><td colspan="1">0x6F03</td><td colspan="1">AEB3 data flow (Special purpose)</td><td colspan="1"></td><td colspan="1">4-byte CAN ID plus an 8-byte data flow, forwarding terminal without parsing</td></tr>
<tr><td colspan="1">0x6F04</td><td colspan="1">0x6F04</td><td colspan="1">AEB4 data flow (Special purpose)</td><td colspan="1"></td><td colspan="1">4-byte CAN ID plus an 8-byte data flow, forwarding terminal without parsing</td></tr>
</table>

1


## <a name="_轿车扩展数据流<二>附表"></a><a name="_货车扩展数据流<一>附表"></a><a name="_toc161247108"></a>**3.38	Schedule- Extended data flow of truck [](#_附表_附加信息定义)**

<table><tr><th colspan="1"><a name="_hlt27215632"></a><a name="_hlt27212216"></a><b>Functional ID domain</b> </th><th colspan="1"><b>Function ID[2]</b> </th><th colspan="1"><b>Length [1]</b> </th><th colspan="1"><b>Function</b> </th><th colspan="1"><b>Unit</b> </th><th colspan="1"><b>Description</b> </th></tr>
<tr><td colspan="1" rowspan="60"><p>Truck data item </p><p>0x5001-0x6FFF</p></td><td colspan="1">0x60C0</td><td colspan="1">2</td><td colspan="1">OBD speed </td><td colspan="1">rpm</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0 - 8000 </td></tr>
<tr><td colspan="1">0x60D0</td><td colspan="1">1</td><td colspan="1">OBD speed </td><td colspan="1">Km/h</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0 - 240 </td></tr>
<tr><td colspan="1">0x62f0</td><td colspan="1">2</td><td colspan="1">OBD remaining fuel </td><td colspan="1"><p>%</p><p>L</p></td><td colspan="1"><p>Remaining fuel, unit L or%</p><p>Bit15 = = 0% OBD is percentage </p><p>= = 1 unit L</p><p>Displayed value is the uploaded value / 10 </p></td></tr>
<tr><td colspan="1">0x6050</td><td colspan="1">1</td><td colspan="1">OBD coolant temperature </td><td colspan="1"><a name="ole_link7"></a><a name="ole_link8"></a>℃</td><td colspan="1">Accuracy: 1℃ deviation: -40.0℃ scope: -40.0℃ - +210℃ </td></tr>
<tr><td colspan="1">0x60F0</td><td colspan="1">1</td><td colspan="1">OBD intake temperature </td><td colspan="1">℃</td><td colspan="1">Accuracy: 1℃ deviation: -40.0℃ scope: -40.0℃ - +210℃ </td></tr>
<tr><td colspan="1">0x60B0</td><td colspan="1">1</td><td colspan="1">OBD intake (manifold absolute) pressure </td><td colspan="1">kPa</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0 - 250kpa. 0x60B0 or 0x50B0 can be selected in the original protocol </td></tr>
<tr><td colspan="1">0x50B0</td><td colspan="1">2</td><td colspan="1">OBD intake (manifold absolute) pressure </td><td colspan="1">kPa</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0 - 500kPa, for truck, 0x60B0 or 0x50B0 can be selected </td></tr>
<tr><td colspan="1">0x6330</td><td colspan="1">1</td><td colspan="1">OBD atmospheric pressure </td><td colspan="1">kPa</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0 - 125kpa </td></tr>
<tr><td colspan="1">0x6460</td><td colspan="1">1</td><td colspan="1">OBD ambient temperature </td><td colspan="1">℃</td><td colspan="1">Accuracy: 1℃ deviation: -40.0℃ scope: -40.0℃ - +210℃ </td></tr>
<tr><td colspan="1">0x6490</td><td colspan="1">1</td><td colspan="1"><p>OBD position of accelerator pedal </p><p>(Throttle pedal) </p></td><td colspan="1">%</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0% - 100% </td></tr>
<tr><td colspan="1">0x60A0</td><td colspan="1">2</td><td colspan="1">OBD fuel pressure </td><td colspan="1">kPa</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0 - 500kpa </td></tr>
<tr><td colspan="1">0x6010</td><td colspan="1">1</td><td colspan="1">OBD number of fault codes </td><td colspan="1">Pcs</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0 - 255 </td></tr>
<tr><td colspan="1">0x5001</td><td colspan="1">1</td><td colspan="1">OBD clutch switch </td><td colspan="1">　</td><td colspan="1">0x00/0x01 OFF/ON </td></tr>
<tr><td colspan="1">0x5002</td><td colspan="1">1</td><td colspan="1">OBD brake switch </td><td colspan="1">　</td><td colspan="1">0x00/0x01 OFF/ON </td></tr>
<tr><td colspan="1">0x5003</td><td colspan="1">1</td><td colspan="1">OBD parking brake switch </td><td colspan="1">　</td><td colspan="1">0x00/0x01 OFF/ON </td></tr>
<tr><td colspan="1">0x5004</td><td colspan="1">1</td><td colspan="1">OBD throttle position: </td><td colspan="1">%</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0% - 100% </td></tr>
<tr><td colspan="1">0x5005</td><td colspan="1">2</td><td colspan="1"><p>OBD utilization rate of oil </p><p>(Fuel flow of engine) </p></td><td colspan="1">L/h</td><td colspan="1">Accuracy: 0.05L/h offset: 0 scope: 0 - 3212.75L／h </td></tr>
<tr><td colspan="1">0x5006</td><td colspan="1">2</td><td colspan="1">OBD fuel temperature </td><td colspan="1">℃</td><td colspan="1">Accuracy: 0.03125℃ offset: -273.0℃ scope: -273.0℃ - +1734.96875℃ </td></tr>
<tr><td colspan="1">0x5007</td><td colspan="1">2</td><td colspan="1">OBD engine oil temperature </td><td colspan="1">℃</td><td colspan="1">Accuracy: 0.03125℃ offset: -273.0℃ scope: -273.0℃ - +1734.96875℃ </td></tr>
<tr><td colspan="1">0x5008</td><td colspan="1">1</td><td colspan="1">OBD pressure of engine lubricating oil </td><td colspan="1">kPa</td><td colspan="1">Accuracy: 4 deviation: 0 scope: 0 - 1000kpa </td></tr>
<tr><td colspan="1">0x5009</td><td colspan="1">1</td><td colspan="1">OBD position of brake pedal </td><td colspan="1">%</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0% - 100% </td></tr>
<tr><td colspan="1">0x500A　</td><td colspan="1">2</td><td colspan="1">OBD air flow </td><td colspan="1">g/s</td><td colspan="1">Accuracy: 0.1 deviation: 0 scope: 0 - 6553.5 </td></tr>
<tr><td colspan="1">0x5101</td><td colspan="1">1</td><td colspan="1">Net output torque of engine </td><td colspan="1">%</td><td colspan="1">Accuracy: 1 deviation: -125 scope: -125% - +125% </td></tr>
<tr><td colspan="1">0x5102</td><td colspan="1">1</td><td colspan="1">Friction torque </td><td colspan="1">%</td><td colspan="1">Accuracy: 1 deviation: -125 scope: -125% - +125%</td></tr>
<tr><td colspan="1">0x5103</td><td colspan="1">2</td><td colspan="1">Output value of SCR upstream NOx sensor </td><td colspan="1">ppm</td><td colspan="1">Accuracy: 0.05 deviation: -200 scope: -200 - +3012.75ppm</td></tr>
<tr><td colspan="1">0x5104</td><td colspan="1">2</td><td colspan="1">Output value of SCR downstream NOx sensor </td><td colspan="1">ppm</td><td colspan="1">Accuracy: 0.05 deviation: -200 scope: -200 - +3012.75ppm</td></tr>
<tr><td colspan="1">0x5105</td><td colspan="1">1</td><td colspan="1">Residual reagent </td><td colspan="1">%</td><td colspan="1">Accuracy: 0.4 deviation: 0 scope: 0% - 100% </td></tr>
<tr><td colspan="1">0x5106</td><td colspan="1">2</td><td colspan="1">Air inflow </td><td colspan="1">Kg/h</td><td colspan="1">Accuracy: 0.05 deviation: 0 scope: 0 - 3212.75 Kg/h </td></tr>
<tr><td colspan="1">0x5107</td><td colspan="1">2</td><td colspan="1">Inlet temperature of SCR </td><td colspan="1">℃</td><td colspan="1">Accuracy: 0.03125℃ offset: -273.0℃ scope: -273.0℃ - +1734.96875℃ </td></tr>
<tr><td colspan="1">0x5108</td><td colspan="1">2</td><td colspan="1">Outlet temperature of SCR </td><td colspan="1">℃</td><td colspan="1">Accuracy: 0.03125℃ offset: -273.0℃ scope: -273.0℃ - +1734.96875℃ </td></tr>
<tr><td colspan="1">0x5109</td><td colspan="1">2</td><td colspan="1">DPF pressure difference </td><td colspan="1">kPa</td><td colspan="1">Accuracy: 0.1 deviation: 0 scope: 0 - 6425.5 kPa </td></tr>
<tr><td colspan="1">0x510A</td><td colspan="1">1</td><td colspan="1">Mode of engine torque </td><td colspan="1"></td><td colspan="1"><p>0: Over-speed failure </p><p>1: Speed control </p><p>2: Torque control </p><p>3: Speed / torque control </p><p>9: Normal </p></td></tr>
<tr><td colspan="1"><s>0x510B</s></td><td colspan="1"><s>1</s></td><td colspan="1"><s>Throttle pedal</s> </td><td colspan="1"><s>1%</s></td><td colspan="1"><s>The displayed value is equal to the uploaded value multiplied by 0.4</s> </td></tr>
<tr><td colspan="1">0x510C</td><td colspan="1">1</td><td colspan="1">Temperature of urea tank </td><td colspan="1">℃</td><td colspan="1">Accuracy: 1℃ deviation: -40.0℃ scope: -40.0℃ - +210℃ </td></tr>
<tr><td colspan="1">0x510D</td><td colspan="1">4</td><td colspan="1">Actual urea injection </td><td colspan="1">ml/h</td><td colspan="1">Accuracy: 0.01 deviation: 0 scope: 0 - 42949672.95 ml/h </td></tr>
<tr><td colspan="1">0x510E</td><td colspan="1">4</td><td colspan="1">Cumulative urea consumption </td><td colspan="1">g</td><td colspan="1">Accuracy: 1 deviation: 0 scope: 0 - 4294967295g </td></tr>
<tr><td colspan="1">0x510F</td><td colspan="1">2</td><td colspan="1">DPF exhaust temperature </td><td colspan="1">℃</td><td colspan="1">Accuracy: 0.03125℃ offset: -273.0℃ scope: -273.0℃ - +1734.96875℃ </td></tr>
<tr><td colspan="1"><s>0x5110</s></td><td colspan="1"><s>2</s></td><td colspan="1"><s>Fuel flow of engine</s> </td><td colspan="1"><s>L/H</s></td><td colspan="1"><s>The displayed value is equal to the uploaded value multiplied by 0.05</s> </td></tr>
<tr><td colspan="1">0x5111</td><td colspan="1">1</td><td colspan="1">OBD diagnostic protocol </td><td colspan="1"></td><td colspan="1">Valid scope is 0 - 2, "0" indicates IOS15765, "1" indicates IOS27145, "2" indicates SAEJ1939 and "0xFE" means invalid. </td></tr>
<tr><td colspan="1">0x5112</td><td colspan="1">1</td><td colspan="1">MIL status </td><td colspan="1"></td><td colspan="1">The valid scope is 0 - 1, "0" indicates unlit and "1" indicates lit. "0xFF" indicates invalid </td></tr>
<tr><td colspan="1">0x5113</td><td colspan="1">2</td><td colspan="1">Diagnostic support status </td><td colspan="1"></td><td colspan="1"><p>Each bit is defined as follows: </p><p>1 Catalyst monitoring Status </p><p>2 Heated catalyst monitoring Status </p><p>3 Evaporative system monitoring Status </p><p>4 Secondary air system monitoring Status </p><p>5 A/C system refrigerant monitoring Status </p><p>6 Exhaust Gas Sensor monitoring Status </p><p>7 Exhaust Gas Sensor heater monitoring Status </p><p>8 EGR/VVT system monitoring </p><p>9 Cold start aid system monitoring Status </p><p>10 Boost pressure control system monitoring Status </p><p>11 Diesel Particulate Filter (DPF) monitoring Status </p><p>12 NOx converting catalyst and/or NOx adsorber monitoring Status </p><p>13 NMHC converting catalyst monitoring Status </p><p>14 Misfire monitoring support </p><p>15 Fuel system monitoring support </p><p>16 Comprehensive component monitoring support </p><p>Meaning of each bit: 0= no support; 1=support; </p></td></tr>
<tr><td colspan="1">0x5114</td><td colspan="1">2</td><td colspan="1">Diagnostic ready status </td><td colspan="1"></td><td colspan="1"><p>Each bit is defined as follows: </p><p>1 Catalyst monitoring Status </p><p>2 Heated catalyst monitoring Status </p><p>3 Evaporative system monitoring Status </p><p>4 Secondary air system monitoring Status </p><p>5 A/C system refrigerant monitoring Status </p><p>6 Exhaust Gas Sensor monitoring Status </p><p>7 Exhaust Gas Sensor heater monitoring Status </p><p>8 EGR/VVT system monitoring </p><p>9 Cold start aid system monitoring Status </p><p>10 Boost pressure control system monitoring Status </p><p>11 Diesel Particulate Filter (DPF) monitoring Status </p><p>12 NOx converting catalyst and/or NOx adsorber monitoring Status </p><p>13 NMHC converting catalyst monitoring Status </p><p>14 Misfire monitoring support </p><p>15 Fuel system monitoring support </p><p>16 Comprehensive component monitoring support </p><p>Meaning of each bit: 0 = test completed or no support; 1=test not completed </p></td></tr>
<tr><td colspan="1">0x5115</td><td colspan="1">17</td><td colspan="1">Vehicle identification number (VIN) </td><td colspan="1">ASCII</td><td colspan="1">As the unique identifier for identification, vehicle identification number consists of 17 digit codes and shall conform to the provisions of 4.5 of GB16735. </td></tr>
<tr><td colspan="1">0x5116</td><td colspan="1">18</td><td colspan="1">Software calibration identification number </td><td colspan="1"></td><td colspan="1">Defined by the manufacturer, the software calibration identification number is composed of letters or numbers and the character "0" is added after it if it is insufficient. </td></tr>
<tr><td colspan="1">0x5117</td><td colspan="1">18</td><td colspan="1">Calibration verification number (CVN) </td><td colspan="1"></td><td colspan="1">Defined by the manufacturer, the calibration verification number is composed of letters or numbers and the character "0" is added after it if it is insufficient. </td></tr>
<tr><td colspan="1">0x5118</td><td colspan="1">36</td><td colspan="1">IUPR value </td><td colspan="1"></td><td colspan="1">Refer to G11 of SAE J1979-DA table for definitions </td></tr>
<tr><td colspan="1">0x511A</td><td colspan="1">2</td><td colspan="1">Coefficient of light adsorption </td><td colspan="1">0\.01m<sup>-1</sup></td><td colspan="1"></td></tr>
<tr><td colspan="1">0x511B</td><td colspan="1">2</td><td colspan="1">Opacity </td><td colspan="1">0\.1%</td><td colspan="1"></td></tr>
<tr><td colspan="1">0x511C</td><td colspan="1">2</td><td colspan="1">Particle concentration (mass flow) </td><td colspan="1">Mg/m<sup>3</sup></td><td colspan="1"></td></tr>
<tr><td colspan="1">0x511F</td><td colspan="1">1</td><td colspan="1">Real-time load of engine </td><td colspan="1">%</td><td colspan="1">0-100%</td></tr>
<tr><td colspan="1">0x5201</td><td colspan="1">2</td><td colspan="1">Current powder pressure (special purpose) </td><td colspan="1">0\.01Mpa</td><td colspan="1"></td></tr>
<tr><td colspan="1">0x5202</td><td colspan="1">2</td><td colspan="1">Current left travel pressure (special purpose) </td><td colspan="1">0\.01Mpa</td><td colspan="1"></td></tr>
<tr><td colspan="1">0x5203</td><td colspan="1">2</td><td colspan="1">Current right travel pressure (special purpose) </td><td colspan="1">0\.01Mpa</td><td colspan="1"></td></tr>
<tr><td colspan="1">0x5204</td><td colspan="1">2</td><td colspan="1">Current powder speed (special purpose) </td><td colspan="1">1rpm</td><td colspan="1"></td></tr>
<tr><td colspan="1">0x5205</td><td colspan="1">1</td><td colspan="1">Current alarm of fuel level (special purpose) </td><td colspan="1"></td><td colspan="1">0: Normal 1: Alarm </td></tr>
<tr><td colspan="1">0x5206</td><td colspan="1">1</td><td colspan="1">Left and right steering of powder handle (special purpose) </td><td colspan="1"></td><td colspan="1">0: Left turn 1: Right turn </td></tr>
<tr><td colspan="1">0x5207</td><td colspan="1">1</td><td colspan="1">Gear status (special purpose) </td><td colspan="1"></td><td colspan="1">0: Neutral gear 1: Forward gear 2: Reverse gear </td></tr>
<tr><td colspan="1">0x5208</td><td colspan="1">1</td><td colspan="1">Locked state (special purpose) </td><td colspan="1"></td><td colspan="1">0: Travel 1: Travel lock </td></tr>
<tr><td colspan="1">0x5209</td><td colspan="1">1</td><td colspan="1">Agricultural machinery status (special purpose) </td><td colspan="1"></td><td colspan="1">0: Standby: Work </td></tr>
<tr><td colspan="1">0x520A</td><td colspan="1">2</td><td colspan="1">Total operation time of powder engine (special purpose) </td><td colspan="1">0\.1H</td><td colspan="1">Unit: Hour </td></tr>
<tr><td colspan="1" rowspan="4"></td><td colspan="1">0x520B</td><td colspan="1">1</td><td colspan="1">Coolant low level alarm</td><td colspan="1"></td><td colspan="1">0: Normal 1: Alarm </td></tr>
<tr><td colspan="1">0x520C</td><td colspan="1">1</td><td colspan="1">Engine oil low level alarm</td><td colspan="1"></td><td colspan="1">0: Normal 1: Alarm </td></tr>
<tr><td colspan="1">0x520D</td><td colspan="1">1</td><td colspan="1">Air pressure warning indicator</td><td colspan="1"></td><td colspan="1">0x00/0x01 OFF/ON</td></tr>
<tr><td colspan="1">0x520E</td><td colspan="1">1</td><td colspan="1">Engine oil low level alarm</td><td colspan="1"></td><td colspan="1">0x00/0x01 OFF/ON</td></tr>
</table>

<a name="_货车扩展数据流<二>附表"></a> 


## <a name="_附表_新能源汽车数据项<一>"></a><a name="_toc161247109"></a>**3.39	Schedule-Data flow of new energy vehicle [](#_附表_附加信息定义)**

<table><tr><th colspan="1"><b>Functional ID domain</b> </th><th colspan="1"><b>Function ID[2]</b> </th><th colspan="1"><b>Length [1]</b> </th><th colspan="1"><b>Function</b> </th><th colspan="1"><b>Unit</b> </th><th colspan="1"><b>Description</b> </th></tr>
<tr><td colspan="1" rowspan="45" valign="top"><p></p><p></p><p></p><p></p><p></p><p></p><p></p><p></p><p></p><p></p><p></p><p></p><p></p><p></p><p></p><p></p><p>Data items of new energy vehicles </p><p>[0x7001-0x7FFF]</p><p></p><p></p></td><td colspan="1">0x7001</td><td colspan="1">4</td><td colspan="1">Mileage </td><td colspan="1">0\.1 km</td><td colspan="1">Displayed value is uploaded value / 10 </td></tr>
<tr><td colspan="1">0x7002</td><td colspan="1">1</td><td colspan="1">Remaining Battery (SOC)</td><td colspan="1">%</td><td colspan="1">0% - 100%</td></tr>
<tr><td colspan="1">0x7003</td><td colspan="1">1</td><td colspan="1">Vehicle speed </td><td colspan="1">Km/h</td><td colspan="1">0 - 240</td></tr>
<tr><td colspan="1">0x7004</td><td colspan="1">1</td><td colspan="1">Charging state </td><td colspan="1"></td><td colspan="1"><p>0x0: Initial value </p><p>0x1: Not charged </p><p>0x2: AC charging </p><p>0x3: DC charging </p><p>0x4: Charging completed </p><p>0x5: Driving charging </p><p>0x6: Parking charging </p><p>0x7: Invalid value </p></td></tr>
<tr><td colspan="1">0x7005</td><td colspan="1">1</td><td colspan="1">State of charging pile </td><td colspan="1"></td><td colspan="1"><p>0x01: Inserted </p><p>0x00: Not inserted </p></td></tr>
<tr><td colspan="1">0x7006</td><td colspan="1">2</td><td colspan="1">Charging and discharging current of power battery </td><td colspan="1">0\.01A</td><td colspan="1"><p>0x0-0xFFFF</p><p>Offset-32767</p><p>Charging in the positive direction</p><p>Discharging in the negative direction</p></td></tr>
<tr><td colspan="1">0x7007</td><td colspan="1">2</td><td colspan="1">Maximum voltage of single cell </td><td colspan="1">0\.01V</td><td colspan="1">0x0-0xFFFF/100</td></tr>
<tr><td colspan="1">0x7008</td><td colspan="1">2</td><td colspan="1">Minimum voltage of single cell </td><td colspan="1">0\.01V</td><td colspan="1">0x0-0xFFFF/100</td></tr>
<tr><td colspan="1">0x7009</td><td colspan="1">2</td><td colspan="1">Current speed of drive motor </td><td colspan="1">Rpm</td><td colspan="1"><p>Offset-32767</p><p>For positive motor forward rotation</p><p>For negative motor reverse rotation</p></td></tr>
<tr><td colspan="1">0x700a</td><td colspan="1">2</td><td colspan="1">Rated torque of drive motor </td><td colspan="1">Nm</td><td colspan="1"></td></tr>
<tr><td colspan="1">0x700b</td><td colspan="1">1</td><td colspan="1">Current temperature of drive motor </td><td colspan="1">C</td><td colspan="1">Uploaded value minus 40 </td></tr>
<tr><td colspan="1">0x700c</td><td colspan="1">2</td><td colspan="1">DC bus voltage, total voltage </td><td colspan="1">0\.1V</td><td colspan="1">0x0-0xFFFF/10</td></tr>
<tr><td colspan="1">0x700d</td><td colspan="1">2</td><td colspan="1">DC bus current, total current </td><td colspan="1">0\.01A</td><td colspan="1"><p>Offset-500A </p><p>Discharging is positive, charging is negative </p><p>0x0-0xFFFF/100-500</p></td></tr>
<tr><td colspan="1">0x700e</td><td colspan="1">2</td><td colspan="1">Available energy of power battery </td><td colspan="1">0\.01Kwh</td><td colspan="1">0x0-0xFFFF</td></tr>
<tr><td colspan="1">0x700f</td><td colspan="1">2</td><td colspan="1">Total power battery voltage</td><td colspan="1">0\.01V</td><td colspan="1">0x0-0xFFFF</td></tr>
<tr><td colspan="1"><s>0x7021</s></td><td colspan="1"><s>2</s></td><td colspan="1"><s>Voltage of No. 1 single cell</s> </td><td colspan="1"><s>0.01V</s></td><td colspan="1"></td></tr>
<tr><td colspan="1"><s>0x7022</s></td><td colspan="1"><s>2</s></td><td colspan="1"><s>Voltage of No. 2 single cell</s> </td><td colspan="1"><s>0.01V</s></td><td colspan="1"></td></tr>
<tr><td colspan="1"><s>0x7023</s></td><td colspan="1"><s>2</s></td><td colspan="1"><s>Voltage of No. 3 single cell</s> </td><td colspan="1"><s>0.01V</s></td><td colspan="1"></td></tr>
<tr><td colspan="1"><s>0x7024</s></td><td colspan="1"><s>2</s></td><td colspan="1"><s>Voltage of No. 4 single cell</s> </td><td colspan="1"><s>0.01V</s></td><td colspan="1"></td></tr>
<tr><td colspan="1"><s>0x7025</s></td><td colspan="1"><s>2</s></td><td colspan="1"><s>Voltage of No. 5 single cell</s> </td><td colspan="1"><s>0.01V</s></td><td colspan="1"></td></tr>
<tr><td colspan="1"><s>0x7026</s></td><td colspan="1"><s>2</s></td><td colspan="1"><s>Voltage of No. 6 single cell</s> </td><td colspan="1"><s>0.01V</s></td><td colspan="1"></td></tr>
<tr><td colspan="1"><s>0x7027</s></td><td colspan="1"><s>2</s></td><td colspan="1"><s>Voltage of No. 7 single cell</s> </td><td colspan="1"><s>0.01V</s></td><td colspan="1"></td></tr>
<tr><td colspan="1"><s>0x7028</s></td><td colspan="1"><s>2</s></td><td colspan="1"><s>Voltage of No. 8 single cell</s> </td><td colspan="1"><s>0.01V</s></td><td colspan="1"></td></tr>
<tr><td colspan="1"><s>0x7029</s></td><td colspan="1"><s>2</s></td><td colspan="1"><s>Voltage of No. 9 single cell</s> </td><td colspan="1"><s>0.01V</s></td><td colspan="1"></td></tr>
<tr><td colspan="1"><s>0x702A</s></td><td colspan="1"><s>1</s></td><td colspan="1"><s>Voltage of No. 10 single cell</s> </td><td colspan="1"><s>0.01V</s></td><td colspan="1"></td></tr>
<tr><td colspan="1">0x702B</td><td colspan="1">1</td><td colspan="1">BMS heartbeat information </td><td colspan="1"></td><td colspan="1">0-255 cycle count </td></tr>
<tr><td colspan="1">0x702C</td><td colspan="1">1</td><td colspan="1">Code of single cell with the highest voltage </td><td colspan="1"></td><td colspan="1"></td></tr>
<tr><td colspan="1">0x702D</td><td colspan="1">1</td><td colspan="1">Code of single cell with the lowest voltage </td><td colspan="1"></td><td colspan="1"></td></tr>
<tr><td colspan="1">0x702E</td><td colspan="1">1</td><td colspan="1">Total number of single cells </td><td colspan="1"></td><td colspan="1"></td></tr>
<tr><td colspan="1">0x702F</td><td colspan="1">1</td><td colspan="1">Total number of temperature probes: </td><td colspan="1"></td><td colspan="1"></td></tr>
<tr><td colspan="1">0x7030</td><td colspan="1">1</td><td colspan="1">Maximum temperature value </td><td colspan="1">C</td><td colspan="1">Uploaded value minus 40</td></tr>
<tr><td colspan="1">0x7031</td><td colspan="1">1</td><td colspan="1">Code of single probe with the highest temperature </td><td colspan="1"></td><td colspan="1"></td></tr>
<tr><td colspan="1">0x7032</td><td colspan="1">1</td><td colspan="1">Minimum temperature value </td><td colspan="1">C</td><td colspan="1">Uploaded value minus 40</td></tr>
<tr><td colspan="1">0x7033</td><td colspan="1">1</td><td colspan="1">Code of single probe with the lowest temperature </td><td colspan="1"></td><td colspan="1"></td></tr>
<tr><td colspan="1">0x7034</td><td colspan="1">4</td><td colspan="1">Alarm information </td><td colspan="1"></td><td colspan="1"></td></tr>
<tr><td colspan="1">0x7035</td><td colspan="1">1</td><td colspan="1">Temperature of the first probe </td><td colspan="1">C</td><td colspan="1">Uploaded value minus 40</td></tr>
<tr><td colspan="1">0x7036</td><td colspan="1">1</td><td colspan="1">Temperature of the second probe </td><td colspan="1">C</td><td colspan="1">Uploaded value minus 40</td></tr>
<tr><td colspan="1">0x7037</td><td colspan="1">1</td><td colspan="1">Temperature of the third probe </td><td colspan="1">C</td><td colspan="1">Uploaded value minus 40</td></tr>
<tr><td colspan="1">0x7038</td><td colspan="1">1</td><td colspan="1">Temperature of the fourth probe </td><td colspan="1">C</td><td colspan="1">Uploaded value minus 40</td></tr>
<tr><td colspan="1">0x7039</td><td colspan="1">1</td><td colspan="1">Temperature of the fifth probe </td><td colspan="1">C</td><td colspan="1">Uploaded value minus 40</td></tr>
<tr><td colspan="1">0x703A</td><td colspan="1">1</td><td colspan="1">Temperature of the sixth probe </td><td colspan="1">C</td><td colspan="1">Uploaded value minus 40</td></tr>
<tr><td colspan="1">0x703B</td><td colspan="1">1</td><td colspan="1">Temperature of the seventh probe </td><td colspan="1">C</td><td colspan="1">Uploaded value minus 40</td></tr>
<tr><td colspan="1">0x703C</td><td colspan="1">1</td><td colspan="1">Temperature of the eighth probe </td><td colspan="1">C</td><td colspan="1">Uploaded value minus 40</td></tr>
<tr><td colspan="1">0x703D</td><td colspan="1">1</td><td colspan="1">Temperature of the ninth probe </td><td colspan="1">C</td><td colspan="1">Uploaded value minus 40</td></tr>
<tr><td colspan="1">0x703E</td><td colspan="1">1</td><td colspan="1">Temperature of the tenth probe </td><td colspan="1">C</td><td colspan="1">Uploaded value minus 40</td></tr>
<tr><td colspan="1" rowspan="8"></td><td colspan="1">0x703F</td><td colspan="1">1</td><td colspan="1">Current battery temperature</td><td colspan="1">℃</td><td colspan="1">Uploaded value minus 40</td></tr>
<tr><td colspan="1">0x7040</td><td colspan="1">1</td><td colspan="1">Vehicle status</td><td colspan="1"></td><td colspan="1">0-Flameout  Start</td></tr>
<tr><td colspan="1">0x7041</td><td colspan="1">2</td><td colspan="1">Insulation resistance</td><td colspan="1"></td><td colspan="1">0x0-0xFFFF</td></tr>
<tr><td colspan="1">0x7042</td><td colspan="1">1</td><td colspan="1">Battery health state</td><td colspan="1"></td><td colspan="1">0-100</td></tr>
<tr><td colspan="1">0x7043</td><td colspan="1">2</td><td colspan="1">Maximum single voltage</td><td colspan="1">0\.01V</td><td colspan="1">0x0-0xFFFF/100</td></tr>
<tr><td colspan="1">0x7044</td><td colspan="1">2</td><td colspan="1">Maximum single voltage</td><td colspan="1">0\.01V</td><td colspan="1">0x0-0xFFFF/100</td></tr>
<tr><td colspan="1">0x7045</td><td colspan="1">2</td><td colspan="1">Unit pressure difference</td><td colspan="1">0\.01V</td><td colspan="1">0x0-0xFFFF/100</td></tr>
<tr><td colspan="1">0x7046</td><td colspan="1">1</td><td colspan="1">Power gear</td><td colspan="1"></td><td colspan="1"><p>0-1-2 Required to determine if the vehicle is engine off/ignited/engine started</p><p>How does it differ from the vehicle status above?</p><p>Or is it the vehicle gear position?</p></td></tr>
</table>

## <a name="_toc161247110"></a><a name="_附表_扩展外设数据流"></a>**3.40	Schedule-Extended peripheral data flow [](#_附表_附加信息定义)**
<table><tr><th colspan="1"><b>Functional ID domain</b> </th><th colspan="1"><p><b>Function</b> </p><p><b>ID[2]</b> </p></th><th colspan="1"><b>Length [1]</b> </th><th colspan="1"><b>Function</b> </th><th colspan="1"><b>Unit</b> </th><th colspan="1"><b>Description</b> </th></tr>
<tr><td colspan="1" rowspan="11">Peripheral data items <br>0x3001-0x4FFF</td><td colspan="1">0x3001</td><td colspan="1">1</td><td colspan="1">Forward and reverse state </td><td colspan="1">　</td><td colspan="1"><p>0x00 (stop) </p><p>0x01 (positive rotation) </p><p>0x02 (reversal rotation) </p></td></tr>
<tr><td colspan="1">0x3002</td><td colspan="1">2</td><td colspan="1">Probe temperature circuit (1) </td><td colspan="1">0\.1℃</td><td colspan="1">Starting temperature - 40.0 ℃, uploaded value minus 40 </td></tr>
<tr><td colspan="1">0x3003</td><td colspan="1">2</td><td colspan="1">Probe temperature circuit (2) </td><td colspan="1">0\.1℃</td><td colspan="1">Starting temperature - 40.0 ℃, uploaded value minus 40 </td></tr>
<tr><td colspan="1">0x3004</td><td colspan="1">2</td><td colspan="1">Probe temperature circuit (3) </td><td colspan="1">0\.1℃</td><td colspan="1">Starting temperature - 40.0 ℃, uploaded value minus 40 </td></tr>
<tr><td colspan="1">0x3005</td><td colspan="1">2</td><td colspan="1">Probe temperature circuit (4) </td><td colspan="1">0\.1℃</td><td colspan="1">Starting temperature - 40.0 ℃, uploaded value minus 40 </td></tr>
<tr><td colspan="1">0x3006</td><td colspan="1">N</td><td colspan="1">Tire pressure data </td><td colspan="1">　</td><td colspan="1">[See Sheet of tire pressure data ](#_附表_外设数据项：动态组包数据__胎压数据表)</td></tr>
<tr><td colspan="1"><a name="_hlt493015822"></a><a name="_hlt493015818"></a><a name="_hlt531958183"></a><a name="_hlt22569388"></a>0x3007</td><td colspan="1">N</td><td colspan="1">Bracelet data packet </td><td colspan="1">　</td><td colspan="1">See bracelet data packet (not available) </td></tr>
<tr><td colspan="1">0x3008</td><td colspan="1">25</td><td colspan="1">H600 video status information </td><td colspan="1"></td><td colspan="1"><a name="h600视频状态信息表"></a><a name="车辆状态表"></a>[See Sheet of H600 video status information ](#_附表_基础数据项:__h600视频状态信息表)</td></tr>
<tr><td colspan="1"><a name="_hlt1674896"></a><a name="_hlt1674957"></a>0x3009</td><td colspan="1">11</td><td colspan="1">H600 input signal </td><td colspan="1"></td><td colspan="1"><a name="见h600输入信号量表"></a>[See Sheet of H600 input signal ](#_附表_基础数据项:  h600输入信号量)</td></tr>
<tr><td colspan="1">0x300A</td><td colspan="1">N</td><td colspan="1">Data packet of load sensor </td><td colspan="1"></td><td colspan="1">[See Data sheet of load sensor ](#_附表_载重传感器数据表)</td></tr>
<tr><td colspan="1"><a name="_hlt20736500"></a>0x300B</td><td colspan="1">N</td><td colspan="1">External oil sensing data </td><td colspan="1"></td><td colspan="1">[See Sheet of external oil sensing data ](#_附表_外接油感数据表)</td></tr>
<tr><td colspan="1" rowspan="2"></td><td colspan="1"><a name="_hlt17983523"></a>0x300C</td><td colspan="1">N</td><td colspan="1">Fire truck 6 channels data collection (special purpose)</td><td colspan="1">%</td><td colspan="1">See sheet of fire truck 6 channels data collection</td></tr>
<tr><td colspan="1">0x300D</td><td colspan="1">8</td><td colspan="1">Temperature sensor data</td><td colspan="1"></td><td colspan="1"><p>12 bytes: Temperature, accuracy 0.1</p><p>34 bytes: Humidity, accuracy 0.1</p><p>56 bytes: Voltage, accuracy 0.01</p><p>78 bytes: Disassembly status and signal strength, </p><p>high byte indicates disassembly status, 0xFF, undetached; </p><p>0x00, detached; </p><p>Low byte indicates signal strength, signed number, reading is signal strength in dBm.</p></td></tr>
</table>


## <a name="_附表_报警命令id及描述数据包"></a><a name="_toc161247111"></a><a name="_附表_报警命令id及描述数据流"></a>**3.41	Schedule-Alarm command ID and description items [](#_附表_附加信息定义)**
<table><tr><th colspan="1"><a name="_hlt54601412"></a><a name="_hlt54601422"></a><b>Functional ID domain</b> </th><th colspan="1"><p><b>Function</b> </p><p><b>ID[2]</b> </p></th><th colspan="1"><b>Length [1]</b> </th><th colspan="1"><b>Function</b> </th><th colspan="1"><b>Instruction</b> </th></tr>
<tr><td colspan="1" rowspan="41" valign="top"><p></p><p></p><p>0x0001-0x0500</p></td><td colspan="1">0x0001</td><td colspan="1">0</td><td colspan="1">Ignition report </td><td colspan="1" rowspan="2"><p>The above data cannot be reported at the same time, </p><p>Only one alarm can be reported </p></td></tr>
<tr><td colspan="1">0x0002</td><td colspan="1">0</td><td colspan="1">Flameout report </td></tr>
<tr><td colspan="1">0x0003</td><td colspan="1">0</td><td colspan="1">Security report </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0004</td><td colspan="1">0</td><td colspan="1">Disarming report </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0005</td><td colspan="1">0</td><td colspan="1">Door opening </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0006</td><td colspan="1">0</td><td colspan="1">Door closing </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0007</td><td colspan="1">0</td><td colspan="1">System startup </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0101</td><td colspan="1">0</td><td colspan="1">Trailer alarm </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0102</td><td colspan="1">0</td><td colspan="1">Too-long positioning alarm </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0103</td><td colspan="1">0</td><td colspan="1">Terminal pull-out alarm (Main power outage)</td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0104</td><td colspan="1">0</td><td colspan="1">Terminal insertion alarm (Main power restored)</td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0105</td><td colspan="1">0</td><td colspan="1">Low-voltage alarm </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0106</td><td colspan="1">X</td><td colspan="1">[Too-long idling alarm ](#_附表_报警描述：怠速报警描述)</td><td colspan="1"></td></tr>
<tr><td colspan="1"><a name="_hlt534401299"></a>0x0107</td><td colspan="1">X</td><td colspan="1">[Over-speed alarm ](#_附表_报警描述：超速报警描述)</td><td colspan="1"></td></tr>
<tr><td colspan="1"><a name="_hlt534400014"></a>0x0108</td><td colspan="1">X</td><td colspan="1">[Fatigue driving alarm ](#_附表_报警描述：疲劳驾驶报警描述)</td><td colspan="1"></td></tr>
<tr><td colspan="1"><a name="_hlt534401156"></a><a name="_hlt534732811"></a>0x0109</td><td colspan="1">X</td><td colspan="1">[Water temperature alarm ](#_附表_报警描述：水温过高报警描述)</td><td colspan="1"></td></tr>
<tr><td colspan="1">0x010A</td><td colspan="1">0</td><td colspan="1">High-speed coasting alarm </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x010B</td><td colspan="1">0</td><td colspan="1">Fuel consumption unsupported alarm </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x010C</td><td colspan="1">0</td><td colspan="1">OBD unsupported alarm </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x010D</td><td colspan="1">0</td><td colspan="1">Low water temperature and high speed </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x010E</td><td colspan="1">0</td><td colspan="1">Bus no-sleep alarm </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x010F</td><td colspan="1">0</td><td colspan="1">Illegal door opening </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0110</td><td colspan="1">0</td><td colspan="1">Illegal ignition </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0111</td><td colspan="1">0</td><td colspan="1">Rapid acceleration alarm </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0112</td><td colspan="1">0</td><td colspan="1">Rapid deceleration alarm </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0113</td><td colspan="1">0</td><td colspan="1">Sharp turn alarm </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0114</td><td colspan="1">0</td><td colspan="1">Collision warning </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0115</td><td colspan="1">0</td><td colspan="1">Abnormal vibration alarm </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0201</td><td colspan="1">0</td><td colspan="1">TTS module fault alarm </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0202</td><td colspan="1">0</td><td colspan="1">FLASH fault alarm </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0203</td><td colspan="1">0</td><td colspan="1">TTS module fault alarm </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0204</td><td colspan="1">0</td><td colspan="1">3D sensor fault alarm </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0205</td><td colspan="1">0</td><td colspan="1">TTS module fault alarm </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0206</td><td colspan="1">0</td><td colspan="1">Alarm for temperature sensor fault </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0301</td><td colspan="1">0</td><td colspan="1">Reminder of the security glass not closed </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0302</td><td colspan="1">0</td><td colspan="1">Reminder of unsuccessful locking </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0303</td><td colspan="1">0</td><td colspan="1">Reminder of failure of timeout security </td><td colspan="1"></td></tr>
<tr><td colspan="1">0x0401</td><td colspan="1">0</td><td colspan="1">Emergency braking </td><td colspan="1"><p>If the current speed is greater than a certain speed, the speed of the vehicle of the next second is less than that of the previous second, and the difference between them exceeds the threshold, the emergency braking alarm is triggered; See 8103 for parameter settings </p><p><b>Velocity difference threshold: 9km/h</b> </p><p><b>The current speed is greater than a certain speed: 0km/h</b> </p></td></tr>
<tr><td colspan="1">0x0402</td><td colspan="1">0</td><td colspan="1">Emergency braking </td><td colspan="1"><p>When the speed of the vehicle in the next second is less than that in the previous second and the difference exceeds the threshold, the emergency braking alarm is triggered; See 8103 for parameter settings </p><p><b>Velocity difference threshold: 18km/h</b> </p></td></tr>
<tr><td colspan="1">0x0403</td><td colspan="1">0</td><td colspan="1">Over speed </td><td colspan="1"><p>The alarm is triggered when the speed is greater than the set threshold. See 8103 for parameter settings </p><p><b>Engine speed threshold: 2400 rpm</b> </p></td></tr>
<tr><td colspan="1">0x0404</td><td colspan="1">0</td><td colspan="1">PTO idling </td><td colspan="1"><p>When the speed is greater than the specified threshold at idle speed, the alarm is triggered. See 8103 for parameter settings </p><p><b>Engine speed threshold: 1000 rpm</b> </p></td></tr>
<tr><td colspan="1"></td><td colspan="1">0x0405</td><td colspan="1">0</td><td colspan="1">OBD connector (4-5) unplugged</td><td colspan="1"></td></tr>
<tr><td colspan="1"></td><td colspan="1">0x0406</td><td colspan="1">0</td><td colspan="1">OBD connector (4-5) plugged in</td><td colspan="1"></td></tr>
<tr><td colspan="1"></td><td colspan="1">0x0407</td><td colspan="1">0</td><td colspan="1">Main unit removal alarm (probe)</td><td colspan="1"></td></tr>
<tr><td colspan="1"></td><td colspan="1">0x0408</td><td colspan="1">0</td><td colspan="1">Main unit box opened alarm (light sensor)</td><td colspan="1"></td></tr>
<tr><td colspan="1"></td><td colspan="1">0x0409</td><td colspan="1">0</td><td colspan="1">New energy charging status change</td><td colspan="1"><b>Used for real-time updating of the charging status of new energy vehicles</b></td></tr>
<tr><td colspan="1"></td><td colspan="1">0x040A</td><td colspan="1">0</td><td colspan="1">Low battery alarm</td><td colspan="1"><b>Used for real-time updating of the charging status of new energy vehicles</b></td></tr>
<tr><td colspan="1"></td><td colspan="1">0x040B</td><td colspan="1">0</td><td colspan="1">OBD1 connector (6-14) unplugged</td><td colspan="1"></td></tr>
<tr><td colspan="1"></td><td colspan="1">0x040C</td><td colspan="1">0</td><td colspan="1">OBD2 connector (1-9) unplugged</td><td colspan="1"></td></tr>
<tr><td colspan="1"></td><td colspan="1">0x040D</td><td colspan="1">0</td><td colspan="1">OBD2 connector (3-11) unplugged</td><td colspan="1"></td></tr>
<tr><td colspan="1"></td><td colspan="1">0x040E</td><td colspan="1">0</td><td colspan="1">OBD2 connector (11-12) unplugged</td><td colspan="1"></td></tr>
</table>


## <a name="_附表_基站数据流"></a><a name="_toc161247112"></a>**3.42	Schedule -Data flow of base station** 

|Contents |Number of bytes |Data type |Description |
| :-: | :-: | :-: | :-: |
|time  bcd[6]|6|byte|Trigger time YY-MM-DD-hh-mm-ss(GMT+8 time), BCD code |
|mcc|n|string|Code of mobile user's country, such as: 460 |
|,|1|byte|Half-width comma separator in English |
|mnc|n|string|Mobile network number, China Mobile: 0; China Unicom: 1 |
|,|1|byte|Half-width comma separator in English |
|base num|1|byte|Number of messages in community, 0-9 |
|,|1|byte|Half-width comma separator in English |
|lac [1]|n|string|Location area code, value scope: 0-65535 |
|,|1|byte|Half-width comma separator in English |
|cellid [1]|n|string|Number of base station community, value scope: 0-65535, 0-268435455, wherein 0,65535,268435455 is not used. When the community number is greater than 65535, it is a 3G base station. |
|,|1|byte|Half-width comma separator in English |
|signal [1]|n|string|Signal intensity, 1-31 |
|. . . . . . |||<p>This field is only available when base num is greater than 1, LAC [2], Cellid [2], Signal [2], separated by commas </p><p>The format is as shown in the gray part above, and the information of multiple communities is accumulated in turn. </p>|

## <a name="_toc8556"></a><a name="_toc161247113"></a><a name="_附表_基础数据项：加速度表"></a>**3.43	Schedule-Basic data flow: [Accelerometer** ](#_附表_基础数据流)**
<table><tr><th colspan="1">Total length </th><th colspan="1">Byte sequence </th><th colspan="1">Contents </th><th colspan="1">Byte </th><th colspan="1">Type </th><th colspan="1">Unit </th><th colspan="1">Instruction </th></tr>
<tr><td colspan="1" rowspan="8">N bytes</td><td colspan="1">0</td><td colspan="1">Number of acquisition points in the last 1 second </td><td colspan="1">2</td><td colspan="1">u16</td><td colspan="1"></td><td colspan="1">4 points in 1 second by default </td></tr>
<tr><td colspan="1">2</td><td colspan="1">Interval of acquisition points in the last 1 second </td><td colspan="1">2</td><td colspan="1">u16</td><td colspan="1">Millisecond </td><td colspan="1"></td></tr>
<tr><td colspan="1">4</td><td colspan="1">Acceleration Mean 1</td><td colspan="1">2</td><td colspan="1">u16</td><td colspan="1">mg</td><td colspan="1">Average acceleration of the first acquisition point </td></tr>
<tr><td colspan="1">6</td><td colspan="1">Acceleration Mean 2</td><td colspan="1">2</td><td colspan="1">u16</td><td colspan="1">mg</td><td colspan="1">Average acceleration of the second acquisition point </td></tr>
<tr><td colspan="1">8</td><td colspan="1">Acceleration Mean 3</td><td colspan="1">2</td><td colspan="1">u16</td><td colspan="1">mg</td><td colspan="1">Average acceleration of the third acquisition point </td></tr>
<tr><td colspan="1">10</td><td colspan="1">Acceleration Mean 4</td><td colspan="1">2</td><td colspan="1">u16</td><td colspan="1">mg</td><td colspan="1">Average acceleration of the fourth acquisition point </td></tr>
<tr><td colspan="1">N</td><td colspan="1">Acceleration Mean N</td><td colspan="1">2</td><td colspan="1">u16</td><td colspan="1">mg</td><td colspan="1">Average acceleration of the Nth acquisition point </td></tr>
<tr><td colspan="1">N*2+2</td><td colspan="1">Acceleration Total Max</td><td colspan="1">2</td><td colspan="1">u16</td><td colspan="1">mg</td><td colspan="1">Maximum acceleration in 1 second </td></tr>
</table>


## <a name="_附表_基础数据项：动态组包数据__总里程格式表"></a><a name="_toc161247114"></a><a name="_附表_基础数据项：总里程格式表"></a>**3.44	Schedule-Basic data items: <a name="总里程格式表"></a>[Format table of total mileage** ](#_附表_基础数据流)**
<table><tr><th colspan="1"><a name="_hlt492821122"></a><a name="_hlt493861740"></a><a name="_hlt492572475"></a><a name="_hlt493861706"></a><a name="_hlt502776655"></a><a name="_hlt493015776"></a><a name="_hlt493861755"></a><a name="_hlt493015739"></a><a name="_hlt492821098"></a><a name="_hlt493861736"></a><a name="_hlt493015738"></a>Item </th><th colspan="1">Byte sequence </th><th colspan="1">Length </th><th colspan="1">Algorithm index </th><th colspan="1">Algorithm name </th></tr>
<tr><td colspan="1" rowspan="12">Mileage type </td><td colspan="1" rowspan="12">0</td><td colspan="1" rowspan="12">1</td><td colspan="1">0x01</td><td colspan="1">Total GPS mileage (cumulatively) </td></tr>
<tr><td colspan="1">0x02</td><td colspan="1">Other 1 [J1939 mileage algorithm 1] </td></tr>
<tr><td colspan="1">0x03</td><td colspan="1">Other 2 [J1939 mileage algorithm 2] </td></tr>
<tr><td colspan="1">0x04</td><td colspan="1">Other 3 [J1939 mileage algorithm 3] </td></tr>
<tr><td colspan="1">0x05</td><td colspan="1">Other 4 [J1939 mileage algorithm 4] </td></tr>
<tr><td colspan="1">0x06</td><td colspan="1">Other 5 [J1939 mileage algorithm 5] </td></tr>
<tr><td colspan="1">0x07</td><td colspan="1">OBD instrument mileage </td></tr>
<tr><td colspan="1">0x08</td><td colspan="1">OBD speed mileage </td></tr>
<tr><td colspan="1">0x09</td><td colspan="1">Other 6 [J1939 mileage algorithm 6] </td></tr>
<tr><td colspan="1">0x0A</td><td colspan="1">Other 7 [J1939 mileage algorithm 7] </td></tr>
<tr><td colspan="1">0x0B</td><td colspan="1">Other 8 [J1939 mileage algorithm 8] </td></tr>
<tr><td colspan="1">0x0C</td><td colspan="1">Other 9 [J1939 mileage algorithm 9] </td></tr>
<tr><td colspan="1">Total mileage </td><td colspan="1">1</td><td colspan="1">4</td><td colspan="2">Unit: Meter </td></tr>
</table>


## <a name="_附表_基础数据项：动态组包数据__总耗油量格式表"></a><a name="_toc161247115"></a><a name="_附表_基础数据项：总耗油量格式表"></a>**3.45	Schedule-Basic data items: Cumulative mileage 2 format table**

<table><tr><th colspan="1">Item </th><th colspan="1">Byte sequence </th><th colspan="1">Length </th><th colspan="1">Algorithm index </th><th colspan="1">Algorithm name </th></tr>
<tr><td colspan="1" rowspan="3">Accumulation type</td><td colspan="1" rowspan="3">0</td><td colspan="1" rowspan="3">1</td><td colspan="1">0x01</td><td colspan="1">GPS Speed Accumulation</td></tr>
<tr><td colspan="1">0x02</td><td colspan="1">OBD Speed Accumulation</td></tr>
<tr><td colspan="1">0x03</td><td colspan="1">OBDSpeed Accumulation</td></tr>
<tr><td colspan="1">Cumulative mileage 2</td><td colspan="1">1</td><td colspan="1">4</td><td colspan="2">Unit: Meter</td></tr>
</table>

##
## <a name="_toc161247116"></a>**3.46	Schedule-Basic data items: [Format table of total fuel consumption** ](#_附表_基础数据流)**
<table><tr><th colspan="1">Item </th><th colspan="1">Byte sequence </th><th colspan="1">Length </th><th colspan="1">Type of fuel consumption </th><th colspan="1">Algorithm name </th></tr>
<tr><td colspan="1" rowspan="7">Type of fuel consumption </td><td colspan="1" rowspan="7">0</td><td colspan="1" rowspan="7">1</td><td colspan="1">0x01</td><td colspan="1">J1939 fuel consumption algorithm 1 </td></tr>
<tr><td colspan="1">0x02</td><td colspan="1">J1939 fuel consumption algorithm 2</td></tr>
<tr><td colspan="1">0x03</td><td colspan="1">J1939 fuel consumption algorithm 3</td></tr>
<tr><td colspan="1">0x04</td><td colspan="1">J1939 fuel consumption algorithm 4</td></tr>
<tr><td colspan="1">0x05</td><td colspan="1">J1939 fuel consumption algorithm 5</td></tr>
<tr><td colspan="1">0x0B</td><td colspan="1">OBD fuel consumption algorithm 1 </td></tr>
<tr><td colspan="1">0x0C</td><td colspan="1">OBD fuel consumption algorithm 2</td></tr>
<tr><td colspan="1">Total fuel consumption </td><td colspan="1">1</td><td colspan="1">4</td><td colspan="2">Unit: ML </td></tr>
</table>



## <a name="_附表_基础数据项：动态组包数据__加速度表"></a><a name="_toc161247117"></a>**3.47	Schedule-Basic data items: [Accelerometer** ](#_附表_基础数据流)**

<table><tr><th colspan="1">Total length </th><th colspan="1">Byte sequence </th><th colspan="1">Contents </th><th colspan="1">Byte </th><th colspan="1">Type </th><th colspan="1">Unit </th><th colspan="1">Instruction </th></tr>
<tr><td colspan="1" rowspan="8">N bytes </td><td colspan="1">0</td><td colspan="1">Number of acquisition points </td><td colspan="1">2</td><td colspan="1">u16</td><td colspan="1"></td><td colspan="1"></td></tr>
<tr><td colspan="1">2</td><td colspan="1">Interval of acquisition points </td><td colspan="1">2</td><td colspan="1">u16</td><td colspan="1">Millisecond </td><td colspan="1"></td></tr>
<tr><td colspan="1">4</td><td colspan="1">Acceleration Mean 1</td><td colspan="1">2</td><td colspan="1">u16</td><td colspan="1">mg</td><td colspan="1">Average acceleration of the first acquisition point </td></tr>
<tr><td colspan="1">6</td><td colspan="1">Acceleration Mean 2</td><td colspan="1">2</td><td colspan="1">u16</td><td colspan="1">mg</td><td colspan="1">Average acceleration of the second acquisition point </td></tr>
<tr><td colspan="1">8</td><td colspan="1">Acceleration Mean 3</td><td colspan="1">2</td><td colspan="1">u16</td><td colspan="1">mg</td><td colspan="1">Average acceleration of the third acquisition point </td></tr>
<tr><td colspan="1">10</td><td colspan="1">Acceleration Mean 4</td><td colspan="1">2</td><td colspan="1">u16</td><td colspan="1">mg</td><td colspan="1">Average acceleration of the fourth acquisition point </td></tr>
<tr><td colspan="1">N</td><td colspan="1">Acceleration Mean N</td><td colspan="1">2</td><td colspan="1">u16</td><td colspan="1">mg</td><td colspan="1">Average acceleration of the Nth acquisition point </td></tr>
<tr><td colspan="1">N+2</td><td colspan="1">Acceleration Total Max</td><td colspan="1">2</td><td colspan="1">u16</td><td colspan="1">mg</td><td colspan="1">Maximum acceleration value within acquisition time </td></tr>
</table>


## <a name="_toc161247118"></a><a name="_附表_基础数据项：协议类型表"></a>**3.48	Schedule-Basic data items: [Sheet of protocol type** ](#_附表_基础数据流)**
<table><tr><th colspan="1"></th><th colspan="1">Value </th><th colspan="1">Protocol type </th></tr>
<tr><td colspan="1" rowspan="11">OBD sheet of protocol type </td><td colspan="1">0X11</td><td colspan="1">CAN 11_500</td></tr>
<tr><td colspan="1">0X12</td><td colspan="1">CAN 11_250</td></tr>
<tr><td colspan="1">0X13</td><td colspan="1">CAN 29_500_EX</td></tr>
<tr><td colspan="1">0X14</td><td colspan="1">CAN 29_250_EX</td></tr>
<tr><td colspan="1">0X20</td><td colspan="1">KWP2000</td></tr>
<tr><td colspan="1">0X30</td><td colspan="1">KWP2000M</td></tr>
<tr><td colspan="1">0X40</td><td colspan="1">ISO9141</td></tr>
<tr><td colspan="1">0X50</td><td colspan="1">VPW</td></tr>
<tr><td colspan="1">0X60</td><td colspan="1">PWM </td></tr>
<tr><td colspan="1">0X70</td><td colspan="1">PRIVATE</td></tr>
<tr><td colspan="1">0XF0</td><td colspan="1">J1939</td></tr>
</table>


## <a name="_附表_基础数据项：动态组包数据__车辆状态表"></a><a name="_toc161247119"></a><a name="_附表_基础数据项：车辆状态表"></a>**3.49	Schedule-Basic data items: [Sheet of vehicle status** ](#_附表_基础数据流)**
<table><tr><th colspan="1"><b>Segment sequence</b> </th><th colspan="1"><b>Subsequence</b> </th><th colspan="1"><b>Contents</b> </th><th colspan="1"><b>Number of words</b> </th><th colspan="1"><b>Data type</b> </th><th colspan="1"><b>Accuracy</b> </th><th colspan="1"><b>Description</b> </th></tr>
<tr><td colspan="1" rowspan="2"><b>State mask</b> </td><td colspan="1" rowspan="2"><b>1</b></td><td colspan="1" rowspan="2"><b>State mask</b> </td><td colspan="1" rowspan="2"><b>10</b></td><td colspan="1" rowspan="2"><b>u8</b></td><td colspan="1" rowspan="2"></td><td colspan="1"><b>State mask of vehicle</b> </td></tr>
<tr><td colspan="1"><b>It indicates whether the following 10 types of vehicles are supported or not</b> </td></tr>
<tr><td colspan="1" rowspan="56">State field </td><td colspan="1" rowspan="8">1</td><td colspan="1" rowspan="8">Safety status </td><td colspan="1" rowspan="8">1</td><td colspan="1" rowspan="8">u8</td><td colspan="1" rowspan="8"></td><td colspan="1">Bit0 1/0  ON/OFF     ACC status </td></tr>
<tr><td colspan="1">Bit1 1 / 0 arming / disarming  arming / disarming status </td></tr>
<tr><td colspan="1">Bit2 1 / 0 press / release the foot brake </td></tr>
<tr><td colspan="1">Bit3 1 / 0 press / release the throttle </td></tr>
<tr><td colspan="1">Bit4 1 / 0 pull up / down the handbrake </td></tr>
<tr><td colspan="1">Bit5 1 / 0 insert / release the main safety belt </td></tr>
<tr><td colspan="1">Bit6 1 / 0 insert / release the auxiliary safety belt </td></tr>
<tr><td colspan="1">Bit7 1/0  ON/OFF   Engine state </td></tr>
<tr><td colspan="1" rowspan="8">2</td><td colspan="1" rowspan="8">Door state </td><td colspan="1" rowspan="8">1</td><td colspan="1" rowspan="8">u8</td><td colspan="1" rowspan="8"></td><td colspan="1">Bit0 1/0  on/off      LF  </td></tr>
<tr><td colspan="1">Bit1 1/0  on/off      RF  </td></tr>
<tr><td colspan="1">Bit2 1/0  on/off      LB  </td></tr>
<tr><td colspan="1">Bit3 1/0  on/off      RB  </td></tr>
<tr><td colspan="1">Bit4 1/0  on/off      TRUNK   </td></tr>
<tr><td colspan="1">Bit5 1/0  on/off      engine hood </td></tr>
<tr><td colspan="1">Bit6 1/0 reserved </td></tr>
<tr><td colspan="1">Bit7 1/0 reserved </td></tr>
<tr><td colspan="1" rowspan="8">3</td><td colspan="1" rowspan="8">Lock state </td><td colspan="1" rowspan="8">1</td><td colspan="1" rowspan="8">u8</td><td colspan="1" rowspan="8"></td><td colspan="1">Bit0 1 / 0 lock / unlock     LF </td></tr>
<tr><td colspan="1">Bit1 1 / 0 lock / unlock     RF </td></tr>
<tr><td colspan="1">Bit2 1 / 0 lock / unlock     LB </td></tr>
<tr><td colspan="1">Bit3 1 / 0 lock / unlock     RB </td></tr>
<tr><td colspan="1">Bit4 1/0  reserved</td></tr>
<tr><td colspan="1">Bit5 1/0  reserved</td></tr>
<tr><td colspan="1">Bit6 1/0  reserved</td></tr>
<tr><td colspan="1">Bit7 1/0  reserved</td></tr>
<tr><td colspan="1" rowspan="8">4</td><td colspan="1" rowspan="8">Window state </td><td colspan="1" rowspan="8">1</td><td colspan="1" rowspan="8">u8</td><td colspan="1" rowspan="8"></td><td colspan="1">Bit0 1 / 0 on / off        LF </td></tr>
<tr><td colspan="1">Bit1 1 / 0 on / off        RF </td></tr>
<tr><td colspan="1">Bit2 1 / 0 on / off        LB </td></tr>
<tr><td colspan="1">Bit3 1 / 0 on / off        RB </td></tr>
<tr><td colspan="1">Bit4 1 / 0 on / off        sunroof switch </td></tr>
<tr><td colspan="1">Bit5 1 / 0 on / off        signal left </td></tr>
<tr><td colspan="1">Bit6 1 / 0 on / off        signal right </td></tr>
<tr><td colspan="1">Bit7 1 / 0 on / off        reading light </td></tr>
<tr><td colspan="1" rowspan="8">5</td><td colspan="1" rowspan="8">Light state 1 </td><td colspan="1" rowspan="8">1</td><td colspan="1" rowspan="8">u8</td><td colspan="1" rowspan="8"></td><td colspan="1">Bit0 1 / 0 on / off        low beam </td></tr>
<tr><td colspan="1">Bit1 1 / 0 on / off         high beam </td></tr>
<tr><td colspan="1">Bit2 1 / 0 on / off         front fog light </td></tr>
<tr><td colspan="1">Bit3 1 / 0 on / off         rear fog light </td></tr>
<tr><td colspan="1">Bit4 1 / 0 on / off         hazard light </td></tr>
<tr><td colspan="1">Bit5 1 / 0 on / off         backup light </td></tr>
<tr><td colspan="1">Bit6 1 / 0 on / off         AUTO light </td></tr>
<tr><td colspan="1">Bit7 1 / 0 on / off          width light </td></tr>
<tr><td colspan="1" rowspan="8">6</td><td colspan="1" rowspan="8">Switch state A </td><td colspan="1" rowspan="8">1</td><td colspan="1" rowspan="8">u8</td><td colspan="1" rowspan="8"></td><td colspan="1">Bit0 1/0  ON/OFF         oil alarm </td></tr>
<tr><td colspan="1">Bit1 1/0  ON/OFF         fuel alarm </td></tr>
<tr><td colspan="1">Bit2 1/0  ON/OFF          wiper </td></tr>
<tr><td colspan="1">Bit3 1/0  ON/OFF         horn </td></tr>
<tr><td colspan="1">Bit4 1/0  ON/OFF        air conditioner </td></tr>
<tr><td colspan="1">Bit5 1/0  ON/OFF   rearview mirror state </td></tr>
<tr><td colspan="1">Bit6 1/0    reserved</td></tr>
<tr><td colspan="1">Bit7 1/0    reserved</td></tr>
<tr><td colspan="1" rowspan="5">7</td><td colspan="1" rowspan="5">Switch state B </td><td colspan="1" rowspan="5">1</td><td colspan="1" rowspan="5">u8</td><td colspan="1" rowspan="5"></td><td colspan="1">Bit0- Bit3Reserved </td></tr>
<tr><td colspan="1">Bit4-BIT7Gear </td></tr>
<tr><td colspan="1">==0 P   ==1 R   ==2 N  ==3 D   ==4  1</td></tr>
<tr><td colspan="1">==5 2   ==6 3   ==7 4  ==8 5   ==9  6</td></tr>
<tr><td colspan="1">==10 M  ==11 S  ==12 B ==15   Non- existent </td></tr>
<tr><td colspan="1">8</td><td colspan="1">Reserved </td><td colspan="1">1</td><td colspan="1">u8</td><td colspan="1"></td><td colspan="1">Reserved </td></tr>
<tr><td colspan="1">9</td><td colspan="1">Reserved </td><td colspan="1">1</td><td colspan="1">u8</td><td colspan="1"></td><td colspan="1">Reserved</td></tr>
<tr><td colspan="1">10</td><td colspan="1">Reserved </td><td colspan="1">1</td><td colspan="1">u8</td><td colspan="1"></td><td colspan="1">Reserved</td></tr>
</table>


## <a name="_附表_报警描述：怠速报警描述"></a><a name="_toc161247120"></a>**3.50 Schedule- Alarm description: Description of idle alarm [](#_附表_报警命令id及描述数据流)**

|Byte sequence |Item |Length |Unit |Description |
| :-: | :-: | :-: | :-: | :-: |
|0|Attributes of idle alarm |1||<p>0x00: Alarm removal; Content items with the following data </p><p>0x01: Alarm triggering; Content items without the following data </p>|
|1|Alarm duration |2|Second ||
|3|Fuel consumption at idle |2|ML||
|5|Maximum idle speed |2|RPM||
|7|Minimum idle speed |2|RPM||

## <a name="_附表_报警描述：超速报警描述"></a><a name="_toc161247121"></a>**3.51	Schedule- Alarm description: Description of over-speed alarm [](#_附表_报警命令id及描述数据流)**

|Byte sequence |Item |Length |Unit |Description |
| :-: | :-: | :-: | :-: | :-: |
|0|Attributes of over-speed alarm |1||<p>0x00: Alarm removal; Content items with the following data </p><p>0x01: Alarm triggering; Content items without the following data </p>|
|1|Alarm duration |2|Second ||
|3|Maximum over-speed |2|0\.1KM/H||
|5|Average speed |2|0\.1KM/H||
|7|Over-speed distance |2|Meter ||

## <a name="_附表_报警描述：疲劳驾驶报警描述"></a><a name="_toc161247122"></a>**3.52	Schedule- Alarm description: Description of fatigue driving alarm [](#_附表_报警命令id及描述数据流)**

|Byte sequence |Item |Length |Unit |Description |
| :-: | :-: | :-: | :-: | :-: |
|0|Attributes of fatigue alarm |1||<p>0x00: Alarm removal; Content items with the following data </p><p>0x01: Alarm triggering; Content items without the following data </p>|
|1|Alarm duration |2|Second ||

## <a name="_附表_报警描述：水温过高报警描述"></a><a name="_toc161247123"></a>**3.53	Schedule- Alarm description: Alarm description of high-water temperature [](#_附表_报警命令id及描述数据流)**

|Byte sequence |Item |Length |Unit |Description |
| :-: | :-: | :-: | :-: | :-: |
|0|Attributes of water temperature alarm |1||<p>0x00: Alarm removal; Content items with the following data </p><p>0x01: Alarm triggering; Content items without the following data </p>|
|1|Alarm duration |4|Second ||
|5|Maximum temperature |2|0\.1° ||
|7|Average temperature |2|0\.1°||



## <a name="_附表_基础数据项:__h600视频状态信息表"></a><a name="_toc161247124"></a>**3.54	Schedule-extended peripheral data: H600 Sheet of video status information [](#_附表 扩展外设数据流)**

|Bit |Definition |Instruction |
| :-: | :-: | :-: |
|1|Total number of channels |Number of camera channels (1-4) |
|2|Request for the talkback |0: No request for the talkback 1: The device is initiating a request for the talkback |
|3|Real-time video |0: Not connected, non-zero: The video is being transmitted online, bit0 channel 1, bit1 channel 2, bit2 channel 3 and bit3 channel 4 |
|4|Talkback state |0: Not started, 1: In talkback |
|5|Video playback |0: Not started, non-zero: The channel is being played back remotely bit0 channel 1, bit1 channel 2, bit2 channel 3 and bit3 channel 4 |
|6|SD1 status |0: Non- existent, 1: Normal,  0xff: Disk error |
|7|SD2 status |0: Non- existent, 1: Normal,  0xff: Disk error |
|8|HDD status |0: Non- existent, 1: Normal,  0xff: Disk error |
|9|USB flash disk status |0: Non- existent, 1: Normal,  0xff: Disk error |
|10|EMMC status |0: Non- existent, 1: Normal,  0xff: Disk error |
|11|Working disk |0xff: No working disk, 0: SD1 is the working disk, 1: SD2 is the working disk, 2: The hard disk is a working disk |
|12|Video status |0: All videos are normal, non-zero: Channel video loss exception: bit0 channel 1, bit1 channel 2, bit2 channel 3, bit3 channel 4 |
|13|Video occlusion |0: All videos are normal, non-zero: Channel video loss exception: bit0 channel 1, bit1 channel 2, bit2 channel 3, bit3 channel 4 |
|14|Channel video recording |0: No video recording, 1: Timed video recording, 2: Manual video recording, 3: Alarm video recording |
|15|Channe2 video recording|0: No video recording, 1: Timed video recording, 2: Manual video recording, 3: Alarm video recording |
|16|Channe3 video recording|0: No video recording, 1: Timed video recording, 2: Manual video recording, 3: Alarm video recording |
|17|Channe4 video recording|0: No video recording, 1: Timed video recording, 2: Manual video recording, 3: Alarm video recording |
|18|Channe5 video recording|0: No video recording, 1: Timed video recording, 2: Manual video recording, 3: Alarm video recording |
|19|Channe6 video recording|0: No video recording, 1: Timed video recording, 2: Manual video recording, 3: Alarm video recording |
|20|Channe7 video recording|0: No video recording, 1: Timed video recording, 2: Manual video recording, 3: Alarm video recording |
|21|Channe8 video recording|0: No video recording, 1: Timed video recording, 2: Manual video recording, 3: Alarm video recording |
|22|Disaster recovery video recording |0: No video recording, non-zero: Channel video recording: bit0 channel 1, bit1 channel 2, bit2 channel 3, bit3 channel 4 |
|23|emmc video recording |0: No video recording, non-zero: Channel video recording: bit0 channel 1, bit1 channel 2, bit2 channel 3, bit3 channel 4 |
|24|Authorization status |0: Unauthorized, 1: Authorization |
|25|AV output |<p>The upper 4 digits indicate how many pictures are displayed, and the lower 4 digits indicate the enlarged serial number </p><p>0x11-0x16 single image </p><p>0x20:2 image, 0x40:4 image, 0x60:6 image, 0x90:9 image </p><p>For single image, 0x11: Channel 1 amplification 0x12: Channel 2 amplification </p><p>0x16: Channel 6 amplification </p>|


## <a name="_toc161247125"></a><a name="_附表_基础数据项:__h600输入信号量"></a>**3.55	Schedule-extended peripheral data: H600 input signal [](#_附表 扩展外设数据流)**
|Bit |Definitions |Instruction |
| :-: | :-: | :-: |
|0  |Signal 1 |Brake signal (high trigger) 1: Trigger 0: No trigger|
|1|Signal 2 |Low beam signal (high trigger) 1: Trigger 0: No trigger|
|2|Signal 3 |High beam signal (high trigger) 1: Trigger 0: No trigger|
|3|Signal 4 |Left turn signal (high trigger) 1: Trigger 0: No trigger|
|4|Signal 5 |Right turn signal (high trigger) 1: Trigger 0: No trigger|
|5|Signal 6 |Custom high 1 signal (high trigger) 1: Trigger 0: No trigger|
|6|Signal 7 |Custom high 2 signal (high trigger) 1: Trigger 0: No trigger|
|7|Signal 8 |Robbery alarm signal (low trigger) 1: Trigger 0: No trigger|
|8|Signal 9 |Door signal (low trigger) 1: Trigger 0: No trigger|
|9|Signal 10 |Custom low 1 signal (low trigger) 1: Trigger 0: No trigger|
|10|Signal 11 |Custom low 2 signal (low trigger) 1: Trigger 0: No trigger |

## <a name="_附表_外设数据项：动态组包数据__胎压数据表"></a><a name="_toc161247126"></a>**3.53	Schedule-extended peripheral data: [Sheet of tire pressure data** ](#_附表 扩展外设数据流)**

<table><tr><th colspan="1">Total length </th><th colspan="1">Byte sequence </th><th colspan="1">Type </th><th colspan="1">Length </th><th colspan="1">Content </th><th colspan="1">Description </th></tr>
<tr><td colspan="1" rowspan="8">4+4*N</td><td colspan="1">0</td><td colspan="1">u32</td><td colspan="1">4</td><td colspan="1">Tire mask </td><td colspan="1"><p>BIT31-BIT0 high position in front and low position in back </p><p>BIT31: No. 1 tire (if it is 1, it is followed by tire pressure byte, otherwise it is empty) </p><p>BIT30: No. 2 tire (if it is 1, it is followed by tire pressure byte, otherwise it is empty) </p><p></p><p>BIT0 : No. 32 tire (if it is 1, it is followed by tire pressure byte, otherwise it is empty)</p></td></tr>
<tr><td colspan="1">4</td><td colspan="1">u16</td><td colspan="1">2</td><td colspan="1">Tire pressure of No. X tire </td><td colspan="1">Unit 1 Kpa </td></tr>
<tr><td colspan="1">6</td><td colspan="1">u8</td><td colspan="1">1</td><td colspan="1">Tire temperature of No. X tire </td><td colspan="1">Unit 1 C The displayed value is the uploaded value minus 40 C </td></tr>
<tr><td colspan="1">7</td><td colspan="1">u8</td><td colspan="1">1</td><td colspan="1">Status of No. X tire </td><td colspan="1"><p>BYTE </p><p>BIT7: Rapid air leakage </p><p>BIT6: Slow air leakage </p><p>BIT5: Low power </p><p>BIT4: High temperature </p><p>BIT3: High pressure </p><p>BIT2: Low pressure </p><p></p><p>Others reserved </p></td></tr>
<tr><td colspan="1">...</td><td colspan="1">...</td><td colspan="1">...</td><td colspan="1">...</td><td colspan="1">...</td></tr>
<tr><td colspan="1">N</td><td colspan="1">u16</td><td colspan="1">2</td><td colspan="1">Tire pressure of No. X+N tire </td><td colspan="1">Unit 1 Kpa </td></tr>
<tr><td colspan="1">N+2</td><td colspan="1">u8</td><td colspan="1">1</td><td colspan="1">Tire temperature of No. X+N tire </td><td colspan="1">Unit 1 C The displayed value is the uploaded value minus 40 C </td></tr>
<tr><td colspan="1">N+3</td><td colspan="1">u8</td><td colspan="1">1</td><td colspan="1">Status of No. X+N tire </td><td colspan="1"><p>BYTE </p><p>BIT7: Rapid air leakage </p><p>BIT6: Slow air leakage </p><p>BIT5: Low power </p><p>BIT4: High temperature </p><p>BIT3: High pressure </p><p>BIT2: Low pressure </p><p></p><p>Others reserved </p></td></tr>
</table>

Note: Only the tire with the tire mask set is followed by the tire pressure byte 

For example: When 0x80000000 is the mask, it is only followed by the tire pressure of tire No.1.  

For example:  0x88000000    0x0082 0x10 0x00   0x0096 0x11 0x80 

Represent 

No. 1 tire, tire pressure 130kpa, tire temperature 16C, status: Normal; 

No. 5 tire, tire pressure 150kpa, tire temperature 17C, status: Rapid air leakage 


## <a name="_附表_载重传感器数据表"></a><a name="_toc161247127"></a>**3.57	Schedule- Data sheet of load sensor** 

|Field |Data type |Field |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|BYTE|Type |1: Unit -ton; 2: Unit- kg |
|1|Word|Rated load |Rated load, filled with 0 if it is not available |
|2|Word|Current load |Current load |
|3|BYTE|Original data type |<p>If there is no original data, the fields and the following fields can be omitted </p><p>1-CNHOYUN </p>|
|` `4|BYTE[N]|Original data ||

## <a name="_toc161247128"></a><a name="油感数据_ap表"></a><a name="_附表_外接油感数据表"></a>**3.58	Schedule-Sheet of external oil sensing data [](#_附表 扩展外设数据流)**
<table><tr><th colspan="1">Total length </th><th colspan="1">Byte sequence </th><th colspan="1">Type </th><th colspan="1">Length </th><th colspan="1">Content </th><th colspan="1">Description </th></tr>
<tr><td colspan="1" rowspan="8">3+5*N</td><td colspan="1">0</td><td colspan="1">u8</td><td colspan="1">1</td><td colspan="1">Effective sign </td><td colspan="1">0: Invalid            1: Valid (oil sense online) </td></tr>
<tr><td colspan="1">1</td><td colspan="1">u8</td><td colspan="1">1</td><td colspan="1">Type of oil sense </td><td colspan="1">0: Xide ultrasonic oil sense;   1: Omnicomm oil sense 2: Differential pressure sensor</td></tr>
<tr><td colspan="1">2</td><td colspan="1">u8</td><td colspan="1">1</td><td colspan="1">Oil sensing mask </td><td colspan="1"><p>BIT7-BIT0 high position in front and low position in back </p><p>BIT6: No. 1 oil sense (if it is 1, it is followed by oil sense byte, otherwise it is empty) </p><p>BIT5: No. 2 oil sense (if it is 1, it is followed by oil sense byte, otherwise it is empty) </p><p></p><p>BIT0: No. 8 oil sense (if it is 1, it is followed by oil sense byte, otherwise it is empty)</p></td></tr>
<tr><td colspan="1">3</td><td colspan="1">u8 </td><td colspan="1">1</td><td colspan="1">No. 1 oil sense unit </td><td colspan="1">1: 0.1mm      2: 0.1%     3:0.1ml</td></tr>
<tr><td colspan="1">3+1</td><td colspan="1">u32</td><td colspan="1">4</td><td colspan="1">No. 1 oil sense value </td><td colspan="1">Actual value </td></tr>
<tr><td colspan="1">…</td><td colspan="1">…</td><td colspan="1">…</td><td colspan="1"> </td><td colspan="1">…</td></tr>
<tr><td colspan="1">3+(N-1)*5</td><td colspan="1">u8</td><td colspan="1">1</td><td colspan="1">No. 1 oil sense unit </td><td colspan="1">1: 0.1mm      2: 0.1%     3:0.1ml</td></tr>
<tr><td colspan="1">3+(N-1)*5+1</td><td colspan="1">u32</td><td colspan="1">4</td><td colspan="1">No. N oil sense value </td><td colspan="1">Actual value </td></tr>
</table>

## <a name="_toc161247129"></a><a name="_附表_消防车6路数据采集数据表"></a><a name="_附表_adc采集数据表"></a>3.59 **Schedule - Sheet of fire truck 6 channels data collection**

<table><tr><th colspan="1">Total length </th><th colspan="1">Byte sequence </th><th colspan="1">Type </th><th colspan="1">Length </th><th colspan="1">Content </th><th colspan="1">Description </th></tr>
<tr><td colspan="1" rowspan="6">12</td><td colspan="1">0</td><td colspan="1">WORD</td><td colspan="1">2</td><td colspan="1">Fire truck water tank level</td><td colspan="1">Actual value, unit 0.01% (divide upload value by 100 for xx.xx%)</td></tr>
<tr><td colspan="1">2</td><td colspan="1">WORD</td><td colspan="1">2</td><td colspan="1">Fire truck foam tank level</td><td colspan="1">Actual value, unit 0.01% (divide upload value by 100 for xx.xx%)</td></tr>
<tr><td colspan="1">4</td><td colspan="1">WORD</td><td colspan="1">2</td><td colspan="1">Reserve</td><td colspan="1"></td></tr>
<tr><td colspan="1">6</td><td colspan="1">WORD</td><td colspan="1">2</td><td colspan="1">Reserve</td><td colspan="1"></td></tr>
<tr><td colspan="1">8</td><td colspan="1">WORD</td><td colspan="1">2</td><td colspan="1">Reserve</td><td colspan="1"></td></tr>
<tr><td colspan="1">10</td><td colspan="1">WORD</td><td colspan="1">2</td><td colspan="1">Reserve</td><td colspan="1"></td></tr>
</table>


## <a name="_文本信息下发消息体数据"></a><a name="_文本信息下发消息体附表"></a><a name="_版本信息包附表"></a><a name="_toc161247130"></a><a name="vin码数据包附表"></a>**3.60	Schedule- Version information packet [](#_[0205]版本信息包)**

|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Version number of terminal software |STRING[14]|<p>Software version number:    HLM200\_V201001<br>HL ------Product name</p><p>M200\_------Terminal name code </p><p>V201 ------ software major version number, release version </p><p>001--------Software minor version number, submitted as internal test </p>|
|14|Version date of terminal software |STRING[10]|Software date: November 19, 2018 |
|24|CPU ID No. |BYTE[12]||
|36|GSM TYPE Name|STRING[15]|GSM model:  |
|51|GSM IMEI No. |STRING[15]|GSM IMEI No. |
|66|SIM card IMSI No. |STRING[15]|IMSI No. of terminal SIM card |
|81|SIM card ICCID |STRING[20]|ICCID No. of terminal SIM card |
|101|Car Type|WORD|Car series/ car model ID |
|103|VIN|STRING[17]|Vehicle VIN |
|120|Total mileage |DWORD|Cumulative total vehicle mileage or vehicle instrument mileage with the terminal installed (m) |
|124|Total fuel consumption |DWORD|Cumulative total fuel consumption of the vehicle with the terminal installed (ml) |
## <a name="_版本信息包应答附表"></a><a name="_toc161247131"></a>**3.61	Schedule- Version information packet response [](#_[8205]版本信息包应答)**

|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Current time of platform |BYTE[6]|<p>YY-MM-DD-hh-mm-ss (BCD code ) Beijing Time GMT+08:00 </p><p>For instance: 0x19,0x01,0x28,0x18,0x10,0x30 </p><p>Beijing Time on 19/1/28 18:10:30 </p>|
|6|Model ID |WORD|If the vehicle model does not need to be set, it is to fill 0 |
|8|Emissions |WORD|Unit in ml. if no setting is required, it is to fill 0 |
|10|Whether to upgrade or not |BYTE|Upgrade for 0x55 rather than others |
## <a name="_toc161247132"></a><a name="_附表_文本信息下发消息体"></a>**3.62	Schedule- Message body of issuing of text information [](#_[8300]文本信息下发)**
|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Marking |BYTE|[Schedule of meaning of the text information mark bits ](#_文本信息、标志位含义附表)|
|<a name="_hlt530400277"></a><a name="_hlt535849001"></a><a name="_hlt491434778"></a><a name="_hlt535847241"></a>1|Text information |STRING|The maximum length is 102 bytes, encoded by GBK. |
## <a name="_文本信息、标志位含义"></a><a name="_文本信息、标志位含义附表"></a><a name="_toc161247133"></a>**3.63	Schedule - meaning of the text information mark bits** 

|Bit |Marking |
| :-: | :-: |
|0|1: Emergency (for sending text messages) |
|1|Reserved |
|2|1: Displayed on the terminal display |
|3|1: Broadcast and reading of Terminal TTS |
|4|1: Ad-screen display |
|5|1: HUD text data transparent transmission |
|6-7|Reserved |
## <a name="_附表_文本信息上发消息体"></a><a name="_toc161247134"></a>**3.64	[Schedule- Message body of issuing of text information** ](#_[6006]文本信息回复)**

|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Marking |BYTE|'0' stands for TXT\_ BG2312, '1' is TXT\_ UNICODE |
|1|Marks |STRING|The default is "\* prompt \*", which takes 6 bytes |
|7|Text information |STRING|the maximum length is 1024 bytes, encoded by GBK. |

## <a name="_数据上、下行透传消息体附表"></a><a name="_车辆控制消息体附表"></a><a name="_事件项组成数据附表"></a><a name="_事件设置消息体数据"></a><a name="_事件项组成数据"></a><a name="_发候选答案消息附表"></a><a name="_信息点播菜单设置消息体数据"></a><a name="_车辆控制消息体数据"></a><a name="_事件设置消息体附表"></a><a name="_信息点播菜单设置消息体附表"></a><a name="_问下发候选答案消息组成"></a><a name="_toc161247135"></a><a name="_附表_数据上行透传消息体"></a>**3.65	Schedule-Message body of data uplink transparent transmission [](#_[0900]数据上行透传)**
|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0|Type of transparent transmission message |BYTE|Sheet of definition of type of transparent transmission message |
|1|The content of transparent transmission message |[N]BYTE|Corresponding the content of the message |

## <a name="_数据压缩上报消息体附表"></a><a name="_平台rsa公钥消息体附表"></a><a name="_自定义数据包附表"></a><a name="_子功能id附表"></a><a name="_toc161247136"></a><a name="_附表_透传消息类型定义"></a>**3.66	Schedule- Definition of type of transparent transmission message** 
|Type of transparent transmission message |The content of transparent transmission message |Descriptions and requirements |
| :-: | :-: | :-: |
|0xF1|Driving travel data (sent out with engine off) |[Data packet of driving travel ](#_驾驶行程数据包附表)|
|<a name="_hlt530401321"></a><a name="_hlt531969450"></a>0xF2|Fault code data (sent out with status changes) |[Data packet of fault code ](#_故障码数据包附表)|
|<a name="_hlt530401334"></a>0xF3|Sleep entry (sent in sleep mode) |[Data packet of sleep entry ](#_休眠进入数据包附表)|
|<a name="_hlt530401323"></a>0xF4|Sleep wake- up (sent out of sleep mode) |[Data packet of sleep wake-up ](#_休眠唤醒数据包附表)|
|<a name="_hlt530401331"></a>0xF5|Compact data packet of vehicle GPS (truck version) |Temporarily not joined |
|0xF6|Feedback packet of MCU upgrade status |[Feedback packet of MCU upgrade status ](#_附表mcu升级状态反馈包)|
|<a name="_hlt20737767"></a><a name="_hlt22569197"></a>0xF7|Description packet of suspected collision alarm |[Description packet of suspected collision alarm ](#_附表_碰撞汇总描述包_f7)|
||||
<a name="_hlt41312875"></a><a name="_hlt54358858"></a>
## <a name="_驾驶行程数据包附表"></a><a name="_toc161247137"></a><a name="_附表_驾驶行程数据包_f1"></a>**3.67a	Schedule-Data packet of driving travel F1 [](#_附表_透传消息类型定义)**
|Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: |
|Information ID |WORD||
|Information length |BYTE||
|Information contents ||[Dynamic information sheet of driving travel data ](#_附表_驾驶行程数据动态信息附表)|

## <a name="_hlt534382247"></a><a name="_附表_驾驶行程数据动态信息附表"></a><a name="_toc161247138"></a>**3.68	Schedule-Dynamic information sheet of driving travel data [](#_附表_驾驶行程数据包 f1)**

|Information ID |Information length |Information contents |Type |Description ||
| :-: | :-: | :-: | :-: | :- | :- |
|0x0001|6|ACC ON TimeBCD[6]|u8|YY-MM-DD-hh-mm-ss (GMT+8time), BCD code ||
|0x0002|6|ACC OFF TimeBCD[6]|u8|YY-MM-DD-hh-mm-ss (GMT+8time), BCD code ||
|0x0003|4|ACC ON latitude |u32|Unit: 0.000001 degree, Bit31 = 0 / 1 north latitude / south latitude ||
|0x0004|4|ACC ON latitude |u32|Unit: 0.000001 degree, Bit31 = 0 / 1 east longitude / west longitude ||
|0x0005|4|ACC OFF纬度|u32|Unit: 0.000001 degree, Bit31 = 0 / 1 north latitude / south latitude ||
|0x0006|4|ACC OFF经度|u32|Unit: 0.000001 degree, Bit31 = 0 / 1 east longitude / west longitude ||
|0x0007|2|Trip Mark|u16|Driving cycle label ||
|0x0008|1|Trip Distance Type|u8|<p>Total mileage type of one driving cycle:  </p><p>0x01 	Total GPS mileage (cumulatively)	 </p><p>0x02 	Other 1 [J1939 mileage algorithm 1] </p><p>0x03 	Other 2 [J1939 mileage algorithm 2] </p><p>0x04 	Other 3 [J1939 mileage algorithm 3] </p><p>0x05 	Other 4 [J1939 mileage algorithm 4] </p><p>0x06 	Other 5 [J1939 mileage algorithm 5] </p><p>0x07 	OBD instrument mileage	 </p><p>0x08 	OBD speed mileage	 </p><p>0x09 	Other 6 [J1939 mileage algorithm 6] </p><p>0x0A 	Other 7 [J1939 mileage algorithm 7] </p><p>0x0B 	Other 8 [J1939 mileage algorithm 8] </p><p>0x0C 	Other 9 [J1939 mileage algorithm 9] </p>||
|0x0009|4|Trip Distance|u32|Total mileage of one driving cycle, unit meter ||
|0x000A|4|Trip Fuel Consum|u32|Total fuel consumption of one driving cycle, unit ml ||
|0x000B|4|Trip Duration Total|u32|Total duration of one driving cycle, unit second ||
|0x000C|2|Trip Overspeed Duration |u16|Cumulative duration of over-speed of one driving cycle, unit second ||
|0x000D|2|Trip OverSpd Times|u16|Over-speed times of one driving cycle, unit times ||
|0x000E|1|Trip Speed Average|u8|Average speed of one driving cycle, unit KM/H ||
|0x000F|1|Trip Speed Maximum|u8|Maximum speed of one driving cycle, unit KM/H ||
|0x0010|4|Trip Idle Duration|u32|Idle time of one driving cycle, unit second ||
|0x0011|1|Trip Mask of Braking |u8|Whether the number of foot brakes of one driving cycle is supported or not, 1 is supported ||
|0x0012|2|Trip Number of Braking|u16|Total times of foot brake of one driving cycle, unit times ||
|0x0013|4|Trip Accelerate times|u32|Times of rapid acceleration of one driving cycle ||
|0x0014|4|Trip Decelerate times|u32|Times of rapid deceleration of one driving cycle ||
|0x0015|4|Trip Sharp turn times|u32|Times of sharp turns of one driving cycle ||
|0x0016|4|Trip Miles Spd less than 20Km/H|u32|Mileage with speed of-20Km/H, unit: m ||
|0x0017|4|Trip Miles Spd between 20-40Km/H|u32|Mileage with speed of 20-40Km/H, unit: m ||
|0x0018|4|Trip Miles Spd between 40-60Km/H|u32|Mileage with speed of 40-60Km/H, unit: m ||
|0x0019|4|Trip Miles Spd between 60-80Km/H|u32|Mileage with speed of 60-80Km/H, unit: m ||
|0x001A|4|Trip Miles Spd between 80-100Km/H|u32|Mileage with speed of 80-100Km/H, unit: m ||
|0x001B|4|Trip Miles Spd between 100-120Km/H|u32|Mileage with speed of 100-120Km/H, unit: m ||
|0x001C|4|Trip Miles Spd Over 120Km/H|u32|Mileage with speed above 120Km/H, unit: m ||
|0x001D |4 |Fuel consumption at idle speed |u32|Fuel consumption value at idle speed in one travel, unit: ML ||

## <a name="_故障码数据包附表"></a><a name="_toc161247139"></a>**3.69	Schedule-Data packet of fault codes F2 [](#_附表_数据上行透传消息体)**

|Byte position |Contents |Number of bytes |Data type |Description |
| :-: | :-: | :-: | :-: | :- |
|0|TIME BCD[6]|6|u8|YY-MM-DD-hh-mm-ss (GMT+8time) |
|6|Latitude |4|u32|Unit: 0.000001 degree, Bit31 = 0 / 1 north latitude / south latitude |
|10|Longitude |4|u32|Unit: 0.000001 degree, Bit31 = 0 / 1 east longitude / west longitude |
|14|Dtc Num|1|u8|0 indicates no fault code, and non-0 indicates the number of fault codes |
|15|Dtc1 ID|4|BYTE|ID number of the first fault code: 4 bytes |
|19|Dtc2 ID|4|BYTE|ID number of the second fault code: 4 bytes |
|23|Dtc3 ID|4|BYTE|ID number of the third fault code: 4 bytes |
|…|…|…|…|… |

Instructions: One fault code number consists of 4 bytes: 

If the protocol type is not 0xF0 (i.e. not J1939 protocol), it is system ID, fault byte 1, fault byte 2 and fault byte 3 respectively; 

<a name="_hlt359220704"></a>If the protocol type is 0XF0, the first three bytes are fault code bytes and the fourth byte is status of fault code.  

<a name="_自定义数据包应答附表"></a>
## <a name="_休眠进入数据包附表"></a><a name="_toc161247140"></a>**3.70	Schedule- Data packet of sleep entry F3 [](#_附表_透传消息类型定义)**

|Byte position |Contents |Number of bytes |Data type |Description ||
| :-: | :-: | :-: | :-: | :-: | :- |
|0|Time  BCD[6]|6|u8|Sleep entry time YY-MM-DD-hh-mm-ss (GMT+8time), BCD code ||

## <a name="_休眠唤醒数据包附表"></a><a name="_toc161247141"></a>**3.71	Schedule-Data packet of sleep wake-up F4 [](#_附表_透传消息类型定义)**

|Byte position |Contents |Number of bytes |Data type |Description |
| :-: | :-: | :-: | :-: | :-: |
|0|Time  BCD[6]|6|u8|Sleep wake-up time YY-MM-DD-hh-mm-ss (GMT+8time), BCD code |
|6|Wake Type|1|u8|<p>Heartbeat 		    	 0X01</p><p>CAN1		0X02</p><p>Low voltage 		0X04</p><p>` `G-SENSOR  	0X08</p><p>ACC interruption 	0X10</p><p>` `GSM		     0X20</p><p>Voltage threshold up to the standard 0X40</p><p>Voltage fluctuation 0X80</p>|
|7|Bat Vol|2|u16|Bus voltage |
|9|Accel Total|2|u16|Vibration acceleration value |


## <a name="_附表mcu升级状态反馈包"></a><a name="_toc161247142"></a>**3.72	Schedule-Feedback packet of MCU upgrade status F6** 

|Byte position |Contents |Number of bytes |Data type |Description ||
| :-: | :-: | :-: | :-: | :-: | :- |
|0|Status after upgrading |1|u8|<p>0x00 succeess </p><p>0x01 same software version number </p><p>0x02 error in upgrading parameter format </p><p>0x03 FTP login timeout </p><p>0x04 download timeout </p><p>0x05 error in file verification </p><p>0x06 error in file type </p><p>0x07 no file </p><p>0x08 Other errors </p>||








## <a name="_附表_碰撞汇总描述包_f7"></a><a name="_toc161247143"></a>**3.73	Schedule-Description packet of suspected collision alarm F7** 
**After triggering the collision, it is to collect the collision data at fixed time points before and after the collision at the specified frequency, and then report it to the platform through F7 command** 

**Notes: When the equipment supports F7 collision reporting, the 0x0114 collision command in 0200 alarm data does not need to be reported again, and the F7 command shall prevail** 

|Byte position |Contents |Number of bytes |Data type |Description |
| :-: | :-: | :-: | :-: | :- |
|0|TIME BCD[6] |6|unsigned char|YY-MM-DD-hh-mm-ss (GMT + 8 time), the time of collision |
|6|Latitude |4|unsigned int|Unit: 0.000001 degree, Bit31 = 0 / 1 north latitude / south latitude, latitude at the time of collision |
|10|Longitude |4|unsigned int|Unit: 0.000001 degree, Bit31 = 0 / 1 east longitude / west longitude, longitude at the time of collision |
|14|Frequency of collection |4|unsigned int|<p>The frequency of collecting raw data shall ensure that the data before and after the collision can be collected for a total of 20 seconds. The acquisition frequency by default is: 500 ms, which can be modified </p><p><1> . The frequency is once every 1000 milliseconds, and the total number of acquisition is 20 times, totaling 20 seconds </p><p><2> . The frequency is once every 500 milliseconds, and the total number of acquisition is 40 times, totaling 20 seconds </p><p><3> . The frequency is once every 250 milliseconds, and the total number of acquisition is 80 times, totaling 20 seconds </p>|
|18|Collision level |1|unsigned char|<p>0x00: Minor level </p><p>0x01: Moderate level </p><p>0x02: Severe level </p>|
|19+((N-1)\*7)+0|Acceleration of X-axis [1] |2|signed short int|N=1; Unit: mg ; Scope -32768 - 32768 |
|19+((N-1)\*7)+2|Acceleration of Y-axis [1] |2|signed short int|N=1; Unit: mg; Scope-32768 - 32768 |
|19+((N-1)\*7)+4|Acceleration of Z-axis [1] |2|signed short int|N=1; Unit: mg; Scope-32768 - 32768 |
|19+((N-1)\*7)+6|Speed [1] |1|unsigned char|N=1; Unit: km/h |
|19+((N-1)\*7)+0|Acceleration of X-axis [2] |2|signed short int|N=2; Unit: mg; Scope-32768 - 32768 |
|19+((N-1)\*7)+2|Acceleration of Y-axis [2] |2|signed short int|N=2; Unit: mg; Scope-32768 - 32768 |
|19+((N-1)\*7)+4|Acceleration of Z-axis [2] |2|signed short int|N=2; Unit: mg; Scope-32768 - 32768 |
|19+((N-1)\*7)+6|Speed [2] |1|unsigned char|N=2; Unit: km/h |
|...|...|...|...|...|
|...|...|...|...|...|
|...|...|...|...|...|
|...|...|...|...|...|
|19+((N-1)\*7)+0|Acceleration of X-axis [N] |2|signed short int|N=N; Unit: mg; Scope -32768 - 32768 |
|19+((N-1)\*7)+2|Acceleration of Y-axis [N] |2|signed short int|N=N; Unit: mg; Scope -32768 - 32768 |
|19+((N-1)\*7)+4|Acceleration of Z-axis [N] |2|signed short int|N=N; Unit: mg; Scope -32768 - 32768 |
|19+((N-1)\*7)+6|Speed [N] |1|unsigned char|N=N; Unit: km/h |

## <a name="_toc161247144"></a><a name="_附表_can广播数据流"></a>3.74 	 **<a name="_toc3396"></a>Schedule - CAN broadcast data flow** 
|Content |Number of bytes |Data type |Description |
| :-: | :-: | :-: | :-: |
|Time |6 |BCD[6] |YY-MM-DD-hh-mm-ss (GMT + 8 equipment reporting adopts Beijing time benchmark) |
|Total number of CAN ID |2 |Word |The total number of CAN ID data collected on the bus, how many lists correspond to the following, quantity of N |
|CAN List[1] |12 |BYTE[12] |Totally 12 bytes, with the first 4 bytes representing the CAN ID, and the following 8 bytes representing the corresponding data flow. |
|. . . ||||
|CAN List[N] |12 |BYTE[12] |Totally 12 bytes, with the first 4 bytes representing the CAN ID, and the following 8 bytes representing the corresponding data flow. |
## <a name="_toc161247145"></a><a name="_附表_新能源汽车bms数据信息体"></a><a name="_附表_wifi数据流"></a>3.75 	**Schedule - New energy vehicle BMS data information body** 
|Starting byte |Field |Data type |Descriptions and requirements |
| :-: | :-: | :-: | :-: |
|0 |Time |BCD[6] |YY-MM-DD-hh-mm-ss (GMT + 8 equipment reporting adopts Beijing time benchmark) |
|6 |BMS data content length |WORD ||
|8 |GPS data content |nByte |<p>See details in: [Data flow of new energy vehicle ](file:///D:\项目\2024\3月\李盼\022403071\1-原文\旧版\比较--车葫芦科技网关通信协议.docx#_附表_新能源汽车BMS数据流)</p><p>Data format: Data packet includes sub ID (2BYTE), length (NBYTE)+ data (NBYTE) </p>|
## <a name="_toc161247146"></a><a name="_附表_新能源汽车bms数据流"></a>3.76 	Schedule - Data flow of new energy vehicle 
|<p><a name="_toc14543"></a>**Function** </p><p>**ID[2]** </p>|**Length[2]** |**Function** |**Unit** |**Description** |
| :-: | :-: | :-: | :-: | :-: |
|0x0001 |N |Single battery pack voltage ||[Single battery pack voltage data sheet ](file:///D:\项目\2024\3月\李盼\022403071\1-原文\旧版\比较--车葫芦科技网关通信协议.docx#_附表_新能源汽车BMS数据流：单体电池组电压数据表)|
|0x0002 |N |Single battery pack temperature ||[Single battery pack temperature data sheet ](file:///D:\项目\2024\3月\李盼\022403071\1-原文\旧版\比较--车葫芦科技网关通信协议.docx#_附表_新能源汽车BMS数据流：单体电池组温度数据表)|
||||||

## <a name="_toc161247147"></a><a name="_附表_新能源汽车bms数据流：单体电池组电压数据表"></a>3.77 	**Schedule - Data flow of new energy vehicle Single battery pack voltage data sheet** 
<table><tr><th colspan="1">Total length </th><th colspan="1">Byte sequence </th><th colspan="1">Type </th><th colspan="1">Content </th><th colspan="1">Description </th></tr>
<tr><td colspan="1" rowspan="7">16+2*N </td><td colspan="1">0 </td><td colspan="1">DWORD </td><td colspan="1">Single battery pack voltage mask 0 </td><td colspan="1"><p>BIT31-BIT0  high position in front and low position in back </p><p>BIT31: No.1 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p><p>BIT30: No.2 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p><p>.................. </p><p>BIT0: No. 32 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p></td></tr>
<tr><td colspan="1">4 </td><td colspan="1">DWORD </td><td colspan="1">Single battery pack voltage mask 1 </td><td colspan="1"><p>BIT31-BIT0  high position in front and low position in back </p><p>BIT31: No. 33 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p><p>BIT30: No. 34 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p><p>.................. </p><p>BIT0: No. 64 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p></td></tr>
<tr><td colspan="1">8 </td><td colspan="1">DWORD </td><td colspan="1">Single battery pack voltage mask 2 </td><td colspan="1"><p>BIT31-BIT0  high position in front and low position in back </p><p>BIT31: No. 65 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p><p>BIT30: No. 66 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p><p>.................. </p><p>BIT0: No. 96 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p></td></tr>
<tr><td colspan="1">12 </td><td colspan="1">DWORD </td><td colspan="1">Single battery pack voltage mask 3 </td><td colspan="1"><p>BIT31-BIT0  high position in front and low position in back </p><p>BIT31: No. 96 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p><p>BIT30: No. 66 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p><p>.................. </p><p>BIT0: No. 128 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p></td></tr>
<tr><td colspan="1">16 </td><td colspan="1">WORD </td><td colspan="1">No. 1 Single battery pack voltage </td><td colspan="1"><p>Unit 0.01V </p><p>Range: 0 - 65534 (numerical offset 32767, indicating -327.67V to +327.67V) </p><p>Descriptions: (upload value -32767)/100V </p></td></tr>
<tr><td colspan="1">. . . </td><td colspan="1">WORD </td><td colspan="1">. . . </td><td colspan="1"><p>Unit 0.01V </p><p>Range: 0 - 65534 (numerical offset 32767, indicating -327.67V to +327.67V) </p><p>Descriptions: (upload value -32767)/100V </p></td></tr>
<tr><td colspan="1">N </td><td colspan="1">WORD </td><td colspan="1">No. N Single battery pack voltage </td><td colspan="1"><p>Unit 0.01V </p><p>Range: 0 - 65534 (numerical offset 32767, indicating -327.67V to +327.67V) </p><p>Descriptions: (upload value -32767)/100V </p></td></tr>
</table>

<a name="_附表_新能源汽车bms数据流：单体电池组温度数据表"></a>Notes: Only single battery pack with the voltage mask set will have subsequent single battery pack voltage bytes. 

For example: When the mask is 0x80000000 0x00000000 0x00000000 0x00000000, only the voltage of the No.1 single battery pack will follow.  

For example: When the mask is 0x88000000 0x00000000 0x00000000 0x00000000 0x8158 0x7EA6, the voltages of the No.1 and No.5 single battery packs will follow, as shown below: 

No.1 single battery pack voltage +3.45V 

No.5 single battery pack voltage -3.45V 
## <a name="_toc161247148"></a>3.78 	**Schedule - BMS Data flow of new energy vehicles Single battery pack temperature data sheet** 
<table><tr><th colspan="1">Total length </th><th colspan="1">Byte sequence </th><th colspan="1">Type </th><th colspan="1">Content </th><th colspan="1">Description </th></tr>
<tr><td colspan="1" rowspan="7">16+2*N </td><td colspan="1">0 </td><td colspan="1">DWORD </td><td colspan="1">Single battery pack temperature mask 0 </td><td colspan="1"><p>BIT31-BIT0  high position in front and low position in back </p><p>BIT31: No.1 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p><p>BIT30: No.2 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p><p>.................. </p><p>BIT0: No.32 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p></td></tr>
<tr><td colspan="1">4 </td><td colspan="1">DWORD </td><td colspan="1">Single battery pack temperature mask 1 </td><td colspan="1"><p>BIT31-BIT0  high position in front and low position in back </p><p>BIT31: No.33 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p><p>BIT30: No.34 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p><p>.................. </p><p>BIT0: No.64 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p></td></tr>
<tr><td colspan="1">8 </td><td colspan="1">DWORD </td><td colspan="1">Single battery pack temperature mask 2 </td><td colspan="1"><p>BIT31-BIT0  high position in front and low position in back </p><p>BIT31: No.65 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p><p>BIT30: No.66 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p><p>.................. </p><p>BIT0: No.96 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p></td></tr>
<tr><td colspan="1">12 </td><td colspan="1">DWORD </td><td colspan="1">Single battery pack temperature mask 3 </td><td colspan="1"><p>BIT31-BIT0  high position in front and low position in back </p><p>BIT31: No.96 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p><p>BIT30: No.66 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p><p>.................. </p><p>BIT0: No.128 Single battery pack (If it is 1, then there are subsequent bytes for individual battery cells; otherwise, it is empty) </p></td></tr>
<tr><td colspan="1">16 </td><td colspan="1">WORD </td><td colspan="1">No.1 single battery pack temperature </td><td colspan="1"><p>Unit 0.1℃ </p><p>Range: 0 -- 2400 (numerical offset is 40℃, denoting -40℃-+200℃); minimum measuring unit 0.1℃. </p><p>Descriptions: (upload value -400)/10℃ </p></td></tr>
<tr><td colspan="1">. . . </td><td colspan="1">WORD </td><td colspan="1">. . . </td><td colspan="1"><p>Unit 0.1℃ </p><p>Range: 0 -- 2400 (numerical offset is 40℃, denoting -40℃-+200℃); minimum measuring unit 0.1℃. </p><p>Descriptions: (upload value -400)/10℃ </p></td></tr>
<tr><td colspan="1">N </td><td colspan="1">WORD </td><td colspan="1">No.N single battery pack temperature </td><td colspan="1"><p>Unit 0.1℃ </p><p>Range: 0 -- 2400 (numerical offset is 40℃, denoting -40℃-+200℃); minimum measuring unit 0.1℃. </p><p>Descriptions: (upload value -400)/10℃ </p></td></tr>
</table>

Notes: Only single battery pack with the temperature mask set will have subsequent single battery pack temperature bytes. 

For example: When the mask is 0x80000000 0x00000000 0x00000000 0x00000000, only the temperature of the No.1 single battery pack will follow.  

For example: When the mask is 0x88000000 0x00000000 0x00000000 0x00000000 0x0000 0x0960, the temperatures of the No.1 and No.5 single battery packs will follow, as shown below: 

No.1 single battery pack temperature -40.0℃; 

No.5 single battery pack temperature +200.0℃; 
## <a name="_toc161247149"></a>3.79 	**Schedule - Wifi data flow** 
|Content |Number of bytes |Data type |Description |
| :-: | :-: | :-: | :-: |
|wifi num |1 |byte |Wifi hotspot count |
|ecn[0] |n |string |Wifi hotspot encryption method, reserved, fixed as "-" |
|, |1 |byte |Half-width comma separator in English |
|ssid[0] |n |string |Wifi hotspot name, reserved, fixed as "-" |
|, |1 |byte |Half-width comma separator in English |
|rssi[0] |1 |byte |Wifi hotspot signal strength, unit: dBm  |
|, |1 |byte |Half-width comma separator in English |
|mac [0] |n |string |Wifi hotspot MAC address, e.g., "1C:20:DB:8D:D7:80" |
|, |1 |byte |Half-width comma separator in English |
|channel [0] |1 |byte |Wifi hotspot channel used, varies depending on the module returned range |
|ecn[1] |n |string |Wifi hotspot encryption method, reserved, fixed as "-" |
|, |1 |byte |Half-width comma separator in English |
|ssid[1] |n |string |Wifi hotspot name, reserved, fixed as "-" |
|, |1 |byte |Half-width comma separator in English |
|rssi[1] |1 |byte |Wifi hotspot signal strength, unit: dBm  |
|, |1 |byte |Half-width comma separator in English |
|mac [1] |n |string |Wifi hotspot MAC address, e.g., "1C:20:DB:8D:D7:80" |
|, |1 |byte |Half-width comma separator in English |
|channel [1] |1 |byte |Wifi hotspot channel used, varies depending on the module returned range |
|. . . . . . |||<p>When "wifi num" is greater than 1, this field is present, enc[1], ssid[1], rssi[1], mac[1], channel[1] separated by commas </p><p>The format is as shown in the gray area above, and multiple wifi hotspot information is cumulatively added in sequence. </p>|
##


<a name="_附表_基站数据包_f8"></a>
# <a name="_toc161247150"></a>**4. Appendix 2: Examples** 
## <a name="_toc161247151"></a>**4.1	Examples of escape functions** 
**/\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\***

**/\* Function name:  void JT\_EscapeData(u16 InLen,u8 \*InBuf,u16 \*OutLen,u8 \*OutBuf)**

**Note: Escape** 

**\* InBuf    :Input data that need to be escaped** 

**\* InLen    :Input the length of data that need to be escaped** 

**\* OutBuf   :Output needs data to be escaped** 

**\* OutLen   :Output the length of data that need to be escaped** 

**\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*/**

**void JT\_EscapeData (u16 InLen,u8 \*InBuf, u16 \*OutLen, u8 \*OutBuf)**

**{**

`	`**u16 i=0;**

`	`**u16 Len=0;**	

`	`**// Escape**	 

`	`**for(i=0;i<InLen;i++)**

`	`**{**

`		`**if(InBuf[i]==0x7E)**

`		`**{**

`			`**OutBuf[Len++]=0x7D;**

`			`**OutBuf[Len++]=0x02;**

`		`**}**

`		`**else if(InBuf[i]==0x7D)**

`		`**{**

`			`**OutBuf[Len++]=0x7D;**

`			`**OutBuf[Len++]=0x01;**

`		`**}**

`		`**else**

`		`**{**

`			`**OutBuf[Len++]=InBuf[i];**

`		`**}**

`	`**}**

`	`**\*OutLen=Len;**

**}**
## <a name="_toc161247152"></a>**4.2	Examples of un-escape functions** 
**/\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\***

**/\* Function name:  void JT\_UnEscapeData (u16 InLen,u8 \*InBuf,u16 \*OutLen,u8 \*OutBuf)**

**Note: Input the data that needs to be un-escaped, and the original data can be output after un-escape** 

**\* InBuf    :Input the data that needs to be un-escaped** 

**\* InLen    :Input the length of data that needs to be un-escaped** 

**\* OutBuf   :Output the data that needs to be un-escaped** 

**\* OutLen   :Output the length of data that needs to be un-escaped** 

**\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*/**

**void JT\_UnEscapeData (u16 InLen,u8 \*InBuf,u16 \*OutLen,u8 \*OutBuf)**

**{**

`	`**u16 i=0;**

`	`**u16 ValidPos=0;**

`	`**if(InBuf[0]!=0x7E)**

`		`**return 0;**	

`	`**OutBuf[ValidPos++]=0x7E;**

`	`**for(i=1;i<InLen;i++)**

`	`**{**

`		`**if(InBuf[i]==0x7D)**

`		`**{**

`			`**if(InBuf[i+1]==0x01)**

`			`**{**

`				`**OutBuf[ValidPos++]=0x7D;**

`				`**i++;**

`			`**}**

`			`**else if(InBuf[i+1]==0x02)**

`			`**{**

`				`**OutBuf[ValidPos++]=0x7E;**

`				`**i++;**

`			`**}**

`			`**else  return 0;**

`		`**}**

`		`**else**

`		`**{**

`			`**OutBuf[ValidPos++]=InBuf[i];**

`		`**}**

`		`**if(InBuf[i]==0x7E)**

`		`**{**

`			`**break;**

`		`**}**

`	`**}**

`	`**if(i==InLen)**

`		`**return 0;**	

`	`**\*OutLen=ValidPos;**

`	`**return (i+1);**

**}**
## <a name="_toc161247153"></a>**4.3	[0200] Details of analysis of location data analysis** 

## <a name="_toc161247154"></a>**4.4	[0900] Details of analysis of transparent transmission of uplink data** 

## <a name="_toc161247155"></a>**4.5	[8300] Details of analysis of text information data** 





