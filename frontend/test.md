# MEITRACK MS06 MySQL Server Installation Guide

This guide outlines the steps to install the MEITRACK MS06 MySQL Server on a Windows Server system. Follow the instructions below to ensure a successful setup.

## 1. Preparation
Before beginning the installation, ensure the following requirements are met:
- **Operating System**: Install Windows Server 2022 Standard 64-bit English Edition or Windows Server 2025 Standard 64-bit English Edition on the server.
- **Network**: Ensure a dedicated fiber optic connection with a static IP address and a minimum bandwidth of 1000 Mbps.
- **Hardware**: Refer to the server specification recommendations in the appendix based on the number of devices (see Appendix below).

## 2. Installation
1. Locate the installation package named `m6e6mySQL.exe`.
2. Double-click the `m6e6mySQL.exe` file to start the installation process.
3. Follow the on-screen prompts to complete the installation.

## 3. Open the MT Server Manager Program
1. After installation, restart the server.
2. Upon restart, a shortcut named **MT Server Manager** will be automatically created on the desktop and launched.
3. The MT Server Manager will automatically start the associated service programs.

## 4. Log in to MS06
1. Open the MT Server Manager program.
2. Enter the provided credentials:
   - **Connection Method**: Remote connection.
   - **Account**: User 21040 (guest account).
   - **Password**: Provided by the user or Meitrack support.
3. Ensure the domain name used is complete and accurate, with no missing characters.

## 5. Open Ports
1. Configure the server firewall to open the necessary ports for MS06 MySQL Server operation.
2. Refer to the port descriptions provided in the original documentation or contact Meitrack support (info@meitrack.com) for specific port requirements.

## Appendix: Server Specification Recommendations
The following configurations are recommended based on the number of devices supported by the MS06 MySQL Server. Physical dedicated servers are required due to high I/O demands, and cloud or virtual servers are not suitable for large deployments.

### Fewer than 500 Devices
- **CPU**: Intel Xeon E-2244, 4 cores
- **Memory**: 32 GB
- **Hard Disk**: 600 GB 2K SAS 2 RAID5 (SSD strongly recommended)
- **Network**: Dedicated 1000 Mbps or higher fiber optic connection with static IP address

### 500 to 1,000 Devices
- **CPU**: Intel Xeon E-2244, 4 cores
- **Memory**: 32 GB
- **Hard Disk**: 1.2 TB 3K SAS 2 RAID5 (SSD strongly recommended)
- **Network**: Dedicated 1000 Mbps or higher fiber optic connection with static IP address

### 1,001 to 10,000 Devices
- **CPU**: Intel Xeon E-2388G, 8 cores
- **Memory**: 64 GB
- **Hard Disk**: 1.2 TB 3K SAS 1/2 RAID5 PERC H710 (SSD strongly recommended)
- **Network**: Dedicated 1000 Mbps or higher fiber optic connection with static IP address

### Storage Estimation
- **Historical Data**: Each device generates approximately 0.3 KB per record. Assuming one record every 10 seconds:
  - Daily data size per device: 2.5 MB
  - For 1,000 devices over 180 days: ~444.76 GB
- **Video Storage**: Each video channel consumes 768 MB per 60 minutes.
  - Example: For 100 devices with 4 channels, recording 3 hours daily for 15 days: ~1,350 GB
- **Recommendation**: Calculate exact storage needs based on the number of devices, channels, and retention period.

### Additional Requirements
- **Operating System**: Windows Server 2022 or 2025 Standard 64-bit English Edition.
- **Network**: Fiber optic with an independent static IP address.
- **Email Functionality** (if required):
  - Provide email address, account, password, SMTP address, SMTP port, and SSL encryption status.
- **Customer Connection**: Provide remote connection method, account (e.g., user 21040 guest), and password.

## Support
For further assistance, contact Meitrack Group at **info@meitrack.com**.

*Copyright © 2025 Meitrack Group. All rights reserved.*