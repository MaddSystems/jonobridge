import sys
import re
from datetime import datetime

def parse_meitrack_message(message):
    """Parse a Meitrack protocol message into its component fields."""
    try:
        # Remove any whitespace, carriage returns, and line feeds
        message = message.strip()
        
        # Check if the message starts with $$ (Meitrack protocol header)
        if not message.startswith("$$"):
            print("Error: Message does not start with $$ (Meitrack protocol header)")
            return
        
        # Extract the checksum part first
        if "*" in message:
            main_part, checksum_part = message.split("*", 1)
            # Clean up the checksum part (remove \r\n if present)
            checksum = re.sub(r'[\r\n]', '', checksum_part)
        else:
            main_part = message
            checksum = "Not found"
        
        # Now split the main part by commas
        parts = main_part.split(",")
        
        # Make sure we have enough parts for basic parsing
        if len(parts) < 16:
            print(f"Error: Message has only {len(parts)} parts, expected at least 16")
            return
        
        # Parse the header ($$X142)
        header = parts[0]
        identifier = header[2:3] if len(header) > 2 else "?"
        length = header[3:] if len(header) > 3 else "?"
        
        # Parse the rest of the fields
        imei = parts[1]
        command = parts[2]
        event_code = parts[3]
        latitude = parts[4]
        longitude = parts[5]
        
        # Parse datetime (format: YYMMDDhhmmss)
        raw_datetime = parts[6]
        try:
            parsed_datetime = datetime.strptime(raw_datetime, '%y%m%d%H%M%S')
            formatted_datetime = parsed_datetime.strftime('%Y-%m-%d %H:%M:%S')
        except:
            formatted_datetime = f"Could not parse: {raw_datetime}"
            
        status = parts[7]
        satellites = parts[8]
        gsm_signal = parts[9]
        speed = parts[10]
        direction = parts[11]
        accuracy = parts[12]
        altitude = parts[13]
        mileage = parts[14]
        runtime = parts[15]
        
        # Parse MCC|MNC|LAC|CellID
        cell_info_part = parts[16] if len(parts) > 16 else ""
        cell_info = cell_info_part.split('|')
        if len(cell_info) >= 4:
            mcc = cell_info[0]
            mnc = cell_info[1]
            lac = cell_info[2]
            cell_id = cell_info[3]
        else:
            mcc, mnc, lac, cell_id = "Unknown", "Unknown", "Unknown", "Unknown"
        
        # Make sure we have enough parts for the rest of the parsing
        port_status = parts[17] if len(parts) > 17 else "Unknown"
        
        # Parse AD values (AD1|AD2|AD3|Battery|AD5)
        ad_info_part = parts[18] if len(parts) > 18 else ""
        ad_info = ad_info_part.split('|')
        if len(ad_info) >= 5:
            ad1 = ad_info[0]
            ad2 = ad_info[1]
            ad3 = ad_info[2]
            battery = ad_info[3]
            ad5 = ad_info[4]
        else:
            ad1, ad2, ad3, battery, ad5 = "Unknown", "Unknown", "Unknown", "Unknown", "Unknown"
        
        # Event Info - check if this part exists in the message
        event_info = parts[19] if len(parts) > 19 else "Not present"
        
        # Check if there's a trailing comma (will be an empty string in parts)
        has_trailing_comma = len(parts) > 20 and parts[20] == ""
        
        # Display the parsed information in a nice format
        print("\n=== MEITRACK MESSAGE PARSING RESULTS ===")
        print(f"Raw Message: {message}")
        print(f"Parts Count: {len(parts)}")
        print("\n--- Header ---")
        print(f"Identifier: {identifier}")
        print(f"Length: {length}")
        
        print("\n--- Device Info ---")
        print(f"IMEI: {imei}")
        print(f"Command: {command}")
        print(f"Event Code: {event_code}")
        
        print("\n--- Location ---")
        print(f"Latitude: {latitude}")
        print(f"Longitude: {longitude}")
        print(f"Datetime: {raw_datetime} ({formatted_datetime})")
        print(f"Positioning Status: {status}")
        print(f"Number of Satellites: {satellites}")
        
        print("\n--- Network & Movement ---")
        print(f"GSM Signal Strength: {gsm_signal}")
        print(f"Speed: {speed}")
        print(f"Direction: {direction}")
        print(f"Accuracy/HDOP: {accuracy}")
        print(f"Altitude: {altitude}")
        print(f"Mileage: {mileage}")
        print(f"Runtime: {runtime}")
        
        print("\n--- Cell Tower Info ---")
        print(f"MCC: {mcc}")
        print(f"MNC: {mnc}")
        print(f"LAC: {lac}")
        print(f"Cell ID: {cell_id}")
        
        print("\n--- I/O & Analog ---")
        print(f"IO Port Status: {port_status}")
        print(f"AD1: {ad1}")
        print(f"AD2: {ad2}")
        print(f"AD3: {ad3}")
        print(f"Battery: {battery}")
        print(f"AD5: {ad5}")
        
        print("\n--- Additional Info ---")
        print(f"Event Info: {event_info}")
        print(f"Has Trailing Comma: {has_trailing_comma}")
        print(f"Checksum: {checksum}")
        
        # Verify checksum
        calculated_checksum = calculate_checksum(main_part)
        print(f"Calculated Checksum: {calculated_checksum}")
        print(f"Checksum Valid: {calculated_checksum == checksum}")
        
    except Exception as e:
        print(f"Error parsing message: {str(e)}")
        import traceback
        traceback.print_exc()

def calculate_checksum(source):
    """Calculate the checksum for a Meitrack message."""
    sum = 0
    for char in source:
        sum += ord(char)
    module = sum % 256
    return format(module, 'X').upper()

def unparse_meitrack_message(imei=None, event_code=None, latitude=0.0, longitude=0.0, 
                           datetime_str=None, status="A", satellites=0, gsm_signal=31,
                           speed=0, direction=0, hdop=0, altitude=0, mileage=0, runtime=0,
                           mcc="0", mnc="0", lac="0000", cell_id="0000", port_status="0000",
                           ad1="0", ad2="0", ad3="0", battery="0", ad5="0", event_info="00000000"):
    """
    Create a Meitrack protocol message with the given parameters.
    
    Returns:
        A properly formatted Meitrack message string.
    """
    # Get current datetime if not provided
    if datetime_str is None:
        datetime_str = datetime.now().strftime("%y%m%d%H%M%S")
        
    # Format the base message
    base_msg = f",{imei},AAA,{event_code},{latitude:.6f},{longitude:.6f},{datetime_str},{status},{satellites},{gsm_signal},"
    base_msg += f"{speed},{direction},{hdop},{altitude},{mileage},{runtime},"
    base_msg += f"{mcc}|{mnc}|{lac}|{cell_id},{port_status},{ad1}|{ad2}|{ad3}|{battery}|{ad5},{event_info},"
    
    # Get next identifier (simulating what Go code does)
    identifier = chr(65 + (ord('A') % 26))  # Just using 'A' for simplicity
    
    # Calculate length
    data_length = len(base_msg) + 4 + 1  # +4 for "$$id", +1 for "*"
    
    # Create header
    header = f"$${identifier}{data_length}"
    
    # Final message before checksum
    pre_output = header + base_msg + "*"
    
    # Calculate checksum
    checksum = calculate_checksum(pre_output)
    
    # Final output
    output = pre_output + checksum
    
    return output

def main():
    """Main function to handle user input and process Meitrack messages."""
    print("=== Meitrack Protocol Message Parser ===")
    print("Enter a Meitrack message to parse (or 'q' to quit):")
    
    while True:
        user_input = input("> ")
        if user_input.lower() == 'q':
            break
        
        # Check if input might be a valid Meitrack message
        if user_input.startswith("$$"):
            parse_meitrack_message(user_input)
        else:
            print("Input doesn't appear to be a valid Meitrack message. It should start with $$.")
        
        print("\nEnter another message or 'q' to quit:")

if __name__ == "__main__":
    # If a command-line argument is provided, parse it
    if len(sys.argv) > 1:
        parse_meitrack_message(sys.argv[1])
    else:
        main()
