import re

# Hardcoded credentials
ACCESS_ID = 70001184
PASSWORD = "JEUTPKKH"

def parse_device_info(text):
    """
    Parse the device information text into a dictionary.
    """
    lines = text.strip().split('\n')
    data = {}
    current_key = None

    for line in lines:
        line = line.strip()
        if not line:
            continue
        # Check if line contains a colon, indicating key-value
        if ':' in line:
            parts = line.split(':', 1)
            key = parts[0].strip()
            value = parts[1].strip()
            # Clean up value if it has extra spaces or dashes
            value = re.sub(r'\s+', ' ', value)
            data[key] = value
            current_key = key
        else:
            # If no colon, it might be continuation or header
            if current_key:
                data[current_key] += ' ' + line
            else:
                # Perhaps treat as a header
                data[line] = ''

    return data

# Example usage
device_text = """
ORBCOMM®
Model No: ST9101
SAT FCC ID: XGS - ST9100
SAT IC ID: 11881A - ST9100 CONTAINS:
CELL FCC ID: XPYUBX21BE01
CELL IC ID:8595A - UBX21BE01
BLE FCC ID: XGS - UNNB30
BLE IC ID: 11881A - UNNB30
9 - 32VDC 1A MAX
ST 9 P/N: ST9101 - F01 B
ST 9 S/N: 02092247SKY6A70
IMEI #: 353500723153153
SAT ID:02092247SKY6A70
SIM 1 ID: 8944500609232378913
SIM 2 ID: 89011703274176088427
"""

parsed_data = parse_device_info(device_text)
print("Parsed Device Information:")
for key, value in parsed_data.items():
    print(f"{key}: {value}")

print(f"\nHardcoded Credentials:\nAccess ID: {ACCESS_ID}\nPassword: {PASSWORD}")
