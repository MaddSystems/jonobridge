#!/usr/bin/env python3
"""
Send MVT366 message to GPS server

This program sends an MVT366 GPS tracking message to server1.gpscontrol.com.mx:8500
with updated timestamp and proper checksum calculation.
"""

import socket
import time
from datetime import datetime

def calculate_mvt366_checksum(mvt366_message):
    """
    Calculate checksum using the same algorithm as main.go
    checksum = (len(header) + len(mvt366) + 2) % 256
    """
    # Calculate length and create header
    length = len(mvt366_message) + 5
    header = f"$$H{length}"

    # Calculate checksum
    checksum = (len(header) + len(mvt366_message) + 2) % 256

    return header, checksum

def create_mvt366_message():
    """
    Create MVT366 message with current timestamp
    Based on: $$f125,02092248SKYEE75,AAA,9,19.521100,-99.211500,251002193123,A,0,31,0,361,0,21,0,0,334|50|0030|0030,0000,0|0|0|0|0,00000000,*23
    """

    # Device information
    mobile_id = "02092248SKYEE75"
    event_code = 1  # Input 9 Inactive (DigInp2Lo)

    # Coordinates (Mexico City area)
    latitude = 19.521100
    longitude = -99.211500

    # Current timestamp in YYMMDDHHMMSS format
    now = datetime.now()
    datetime_str = now.strftime("%y%m%d%H%M%S")

    # Other GPS data
    speed = 0.0
    heading = 361  # undefined
    altitude = 21.0

    # Build the MVT366 payload (without header and checksum)
    mvt366_payload = (
        f"{mobile_id},AAA,{event_code},{latitude:.6f},{longitude:.6f},"
        f"{datetime_str},A,0,31,{speed:.6f},{heading},"
        f"0.000000,{altitude:.6f},0.000000,0,"
        f"0030|0030|0030|0030|0030,,,3,,,0,0"
    )

    # Calculate header and checksum
    header, checksum = calculate_mvt366_checksum(mvt366_payload)

    # Create final message
    final_message = f"{header},{mvt366_payload}*{checksum:02X}\r\n"

    return final_message

def send_mvt366_message(host="server1.gpscontrol.com.mx", port=8500):
    """
    Send MVT366 message to GPS server via TCP
    """
    try:
        # Create message with current timestamp
        message = create_mvt366_message()

        print("=" * 60)
        print("MVT366 GPS Message Sender")
        print("=" * 60)
        print(f"Server: {host}:{port}")
        print(f"Message: {message.strip()}")
        print(f"Message length: {len(message)} bytes")
        print()

        # Create TCP socket
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(10)  # 10 second timeout

        print(f"Connecting to {host}:{port}...")

        # Connect to server
        sock.connect((host, port))
        print("✅ Connected successfully!")

        # Send message
        print("Sending message...")
        sock.send(message.encode('ascii'))

        print("✅ Message sent successfully!")

        # Try to receive response (some servers might acknowledge)
        try:
            sock.settimeout(2)  # Wait up to 2 seconds for response
            response = sock.recv(1024)
            if response:
                print(f"📨 Server response: {response.decode('ascii', errors='ignore').strip()}")
            else:
                print("📭 No response from server (this is normal)")
        except socket.timeout:
            print("📭 No response from server within timeout (this is normal)")

    except socket.error as e:
        print(f"❌ Socket error: {e}")
        return False
    except Exception as e:
        print(f"❌ Error: {e}")
        return False
    finally:
        try:
            sock.close()
            print("🔌 Connection closed.")
        except:
            pass

    print("=" * 60)
    return True

def main():
    """Main function"""
    print("Sending MVT366 GPS tracking message...")
    success = send_mvt366_message()

    if success:
        print("✅ Operation completed successfully!")
    else:
        print("❌ Operation failed!")

if __name__ == "__main__":
    main()