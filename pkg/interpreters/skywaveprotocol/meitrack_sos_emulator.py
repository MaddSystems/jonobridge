import socket
import time
from datetime import datetime

class MeitrackSimulator:
    def __init__(self, host, port, imei):
        self.host = host
        self.port = port
        self.imei = imei
        self.socket = None
        self.mileage = 5181284  # meters
        self.runtime = 3528723  # seconds
        
    def connect(self):
        """Connect to the GPS server"""
        self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.socket.connect((self.host, self.port))
        print(f"Connected to {self.host}:{self.port}")
        
    def disconnect(self):
        """Close connection"""
        if self.socket:
            self.socket.close()
            print("Disconnected")
    
    def calculate_checksum(self, data):
        """Calculate checksum (sum of all bytes excluding checksum and ending)"""
        checksum = sum(ord(c) for c in data)
        return f"{checksum % 256:02X}"
    
    def build_packet(self, event_code, io_status, latitude=19.521055, longitude=-99.211760):
        """Build a Meitrack protocol packet"""
        # Current timestamp
        now = datetime.now()
        timestamp = now.strftime("%y%m%d%H%M%S")
        
        # Packet components
        data_identifier = "Q"  # Example identifier
        command_type = "AAA"
        positioning_status = "A"  # Valid GPS
        satellites = "6"
        gsm_signal = "26"
        speed = "0"
        direction = "0"
        hdop = "1.3"
        altitude = "2311"
        
        # Base station info: MCC|MNC|LAC|CI
        base_station = "334|50|75F4|00BE2934"
        
        # Analog input values (5 values minimum)
        # Battery voltage ~4.0V, External power ~56V
        analog_values = "0000|0000|0000|0198|04B0"
        
        # Assisted event info (empty for most events)
        assisted_info = ""
        
        # Customized data (empty)
        custom_data = ""
        
        # Protocol version
        protocol_version = "3"
        
        # Fuel percentage (empty)
        fuel = ""
        
        # Temperature sensors (empty)
        temp = ""
        
        # Acceleration values (empty)
        max_accel = ""
        max_decel = ""
        
        # Build content (everything after the first comma)
        content = (f"{self.imei},{command_type},{event_code},"
                  f"{latitude},{longitude},{timestamp},"
                  f"{positioning_status},{satellites},{gsm_signal},"
                  f"{speed},{direction},{hdop},{altitude},"
                  f"{self.mileage},{self.runtime},{base_station},"
                  f"{io_status},{analog_values},{assisted_info},"
                  f"{custom_data},{protocol_version},{fuel},{temp},"
                  f"{max_accel},{max_decel}")
        
        # Calculate data length (from first comma to \r\n)
        data_length = len(content) + 3  # +3 for *XX
        
        # Build packet without checksum first
        packet_without_checksum = f"$${data_identifier}{data_length},{content}*"
        
        # Calculate checksum
        checksum = self.calculate_checksum(packet_without_checksum[2:])  # Exclude $$
        
        # Complete packet
        packet = f"{packet_without_checksum}{checksum}\r\n"
        
        return packet
    
    def send_packet(self, packet):
        """Send packet to server"""
        self.socket.sendall(packet.encode('ascii'))
        print(f"Sent: {packet.strip()}")
        
        # Try to receive response
        try:
            self.socket.settimeout(2.0)
            response = self.socket.recv(1024)
            if response:
                print(f"Received: {response.decode('ascii').strip()}")
        except socket.timeout:
            print("No response received (timeout)")
        
        # Update counters
        self.runtime += 3
    
    def simulate_sos_event(self):
        """Simulate complete SOS button press and release cycle"""
        print("\n=== Starting SOS Event Simulation ===\n")
        
        # 1. Send normal position (Event 35, no SOS)
        print("1. Sending normal position...")
        normal_packet = self.build_packet(
            event_code=35,
            io_status="0200"  # Input 2 ON (ignition), Input 1 OFF
        )
        self.send_packet(normal_packet)
        time.sleep(2)
        
        # 2. SOS Button Pressed (Event 1 - Input 1 Active)
        print("\n2. SOS BUTTON PRESSED!")
        sos_active_packet = self.build_packet(
            event_code=1,
            io_status="0300"  # Input 1 ON (SOS), Input 2 ON (ignition)
        )
        self.send_packet(sos_active_packet)
        time.sleep(3)
        
        # 3. Keep SOS active for a few seconds (simulate holding button)
        print("\n3. SOS still active...")
        sos_still_active = self.build_packet(
            event_code=1,
            io_status="0300"
        )
        self.send_packet(sos_still_active)
        time.sleep(2)
        
        # 4. SOS Button Released (Event 9 - Input 1 Inactive)
        print("\n4. SOS BUTTON RELEASED!")
        sos_inactive_packet = self.build_packet(
            event_code=9,
            io_status="0200"  # Input 1 OFF, Input 2 ON (ignition)
        )
        self.send_packet(sos_inactive_packet)
        time.sleep(2)
        
        # 5. Return to normal tracking
        print("\n5. Returning to normal tracking...")
        final_normal = self.build_packet(
            event_code=35,
            io_status="0200"
        )
        self.send_packet(final_normal)
        
        print("\n=== SOS Event Simulation Complete ===\n")

def main():
    # Configuration
    SERVER_HOST = "server1.gpscontrol.com.mx"
    SERVER_PORT = 8500
    IMEI = "864606045840682"  # From your example
    
    # Create simulator
    simulator = MeitrackSimulator(SERVER_HOST, SERVER_PORT, IMEI)
    
    try:
        # Connect to server
        simulator.connect()
        
        # Simulate SOS event cycle
        simulator.simulate_sos_event()
        
    except Exception as e:
        print(f"Error: {e}")
    finally:
        # Disconnect
        simulator.disconnect()

if __name__ == "__main__":
    main()