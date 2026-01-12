import socket
from socket import SOL_SOCKET, SO_REUSEADDR
import sys
from datetime import datetime, timedelta
import time



def crc(source):
    b = 0
    for i in source:
        b = b + ord(i)
    ret = hex(b % 256)
    ret = ret.upper()
    ret = ret.replace("0X", "")
    return ret


def charcounter(source):
    c = 0
    for i in source:
        c = c + 1
    return c


class identifier:
    idCounter = 64

    def __init__(self):
        if identifier.idCounter < 123:
            identifier.idCounter += 1
        else:
            identifier.idCounter = 65


def payload(imei, eventCode, latitude, longitude, utc, status, sats, gsmStrenght, speed, direction, accuracy, altitude,
            mileage, runtime, mcc, mnc, lac, cellId, portStatus, AD1, AD2, AD3, battery, AD5, eventInfo):
    # MVT380
    imei = imei.strip()
    newIdentifier = identifier()
    mydataidentifier = str(chr(newIdentifier.idCounter))

    first_output = "," + imei + ",AAA," + eventCode + "," + latitude + "," + longitude + "," + utc + "," + status + "," + str(
        sats) + "," + str(gsmStrenght) + "," + str(speed) + "," + str(direction) + "," + str(accuracy) + "," + str(
        altitude) + "," + str(mileage) + "," + str(
        runtime) + "," + mcc + "|" + mnc + "|" + lac + "|" + cellId + "," + portStatus + "," + AD1 + "|" + AD2 + "|" + AD3 + "|" + str(
        battery) + "|" + AD5 + "," + eventInfo + ",*"
    totalchar = charcounter(first_output) + 4
    header = "$$" + mydataidentifier + str(totalchar)
    preoutput = header + first_output
    output = preoutput + crc(preoutput) + chr(13) + chr(10)
    return output


def test_parameters(server, port):
    # Original values from your code
    original_values = {
        "imei": "867630074536695",
        "eventCode": "1",
        "latitude": "19.611106",
        "longitude": "-99.028335",
        # UTC will be calculated just before sending
        "status": "A",
        "sats": "9",
        "gsmStrenght": "12",
        "speed": "98",
        "direction": "76",
        "accuracy": "1",
        "altitude": "2239",
        "mileage": "0",
        "runtime": "1348",
        "mcc": "0",
        "mnc": "0",
        "lac": "0000",
        "cellId": "0000",
        "portStatus": "0000",
        "AD1": "0000",
        "AD2": "0000",
        "AD3": "0000",
        "battery": "80",
        "AD5": "0000",
        "eventInfo": "00000000"
    }

    # Values from the example message
    test_values = {
        "imei": "867630074536695",
        "eventCode": "35",
        "latitude": "19.521142",
        "longitude": "-99.211361",
        # UTC will be calculated just before sending
        "status": "A",
        "sats": "12",
        "gsmStrenght": "31",
        "speed": "0",
        "direction": "204",
        "accuracy": "0",
        "altitude": "0",
        "mileage": "0",
        "runtime": "0",
        "mcc": "334",
        "mnc": "50",
        "lac": "2550",
        "cellId": "0000",
        "portStatus": "00000000",
        "AD1": "0000",
        "AD2": "0000",
        "AD3": "0000",
        "battery": "C",  # Note: this is not a number but a letter in the example
        "AD5": "0000",
        "eventInfo": "00000000"
    }

    # First, send a message with all original values
    print("\n=== SENDING WITH ALL ORIGINAL VALUES ===")
    current_values = original_values.copy()
    # Set current UTC time just before sending
    current_values["utc"] = datetime.utcnow().strftime('%y%m%d%H%M%S')
    output = payload(**current_values)
    print(output)
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.connect((server, port))
        s.sendall(output.encode('utf-8'))
        s.close()
        print("Message sent successfully!")
    except Exception as e:
        print(f"Error sending message: {e}")

    input("Press Enter to continue to parameter testing...")

    # Now test each parameter one by one
    for param, test_value in test_values.items():
        # Skip UTC parameter as we'll always use the current time
        if param == "utc":
            continue
            
        # Reset to original values
        current_values = original_values.copy()

        # Change only the current parameter
        current_values[param] = test_value
        
        # Always update UTC to current time
        current_values["utc"] = datetime.utcnow().strftime('%y%m%d%H%M%S')

        print(f"\n=== TESTING PARAMETER: {param} = {test_value} ===")
        output = payload(**current_values)
        print(output)

        try:
            s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            s.connect((server, port))
            s.sendall(output.encode('utf-8'))
            s.close()
            print(f"Message with {param}={test_value} sent successfully!")
        except Exception as e:
            print(f"Error sending message: {e}")

        input("Press Enter to continue to the next parameter...")

    # Finally, send a message with all test values
    print("\n=== SENDING WITH ALL TEST VALUES ===")
    current_values = test_values.copy()
    # Set current UTC time just before sending
    current_values["utc"] = datetime.utcnow().strftime('%y%m%d%H%M%S')
    output = payload(**current_values)
    print(output)
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.connect((server, port))
        s.sendall(output.encode('utf-8'))
        s.close()
        print("Message sent successfully!")
    except Exception as e:
        print(f"Error sending message: {e}")


# Main execution
url = "server1.gpscontrol.com.mx"
port = 8500

# Run the test sequence
test_parameters(url, port)