
# $$B153,867630074536695,AAA,35,19.521142,-99.211361,250311222802,V,12,4,0,0,0.0,0,0,0,334|50|2550|000000,00000000,0000|0000|0000|12.0|0000,00000000,,1,0000*E5

import re

def parse_message(message):
    # Remove checksum (everything after '*')
    message = message.split('*')[0]

    # Split by commas
    parts = message.split(',')

    if len(parts) < 22:  # Ensure we have the expected number of fields
        print("Invalid message format!")
        return None

    # Extract values based on known positions
    data = {
        "Message Type": parts[0],  # e.g., $$B153
        "IMEI": parts[1],  # e.g., 867630074536695
        "Command": parts[2],  # e.g., AAA
        "Event Code": parts[3],  # e.g., AAA
        "Latitude": parts[4],  # e.g., 19.521142
        "Longitude": parts[5],  # e.g., -99.211361
        "DateTime (DDMMYYHHMMSS)": parts[6],  # e.g., 250311222802
        "Positioning Status": parts[7],  # e.g., V (Valid/Invalid GPS)
        "Number of Satellites": parts[8],  # e.g., 12
        "GSM Signal Strength": parts[9],  # e.g., 4
        "Speed": parts[10],  # e.g., 35
        "Direction": parts[11],  # e.g., 0
        "HDOP": parts[12],  # e.g., 0
        "Altitude": parts[13],  # e.g., 0.0
        "Mileage": parts[14],  # e.g., 0
        "Runtime": parts[15],  # e.g., 0
        "MCC|MNC|LAC|Cell ID": parts[16],  # e.g., 334|50|2550|000000
        "IO Port Status": parts[17],  # e.g., 00000000
        "AD Values": parts[18],  # e.g., 0000|0000|0000|12.0|0000

    }

    return data

# Get input from the user
message = input("Enter the message string: ")

# Parse and print values
parsed_data = parse_message(message)

if parsed_data:
    print("\nParsed Data:")
    for key, value in parsed_data.items():
        print(f"{key}: {value}")