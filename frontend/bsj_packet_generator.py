import binascii
import struct

def calculate_checksum(data):
    """
    Calculate the BSJ protocol checksum (XOR of all bytes)
    """
    checksum = 0
    for byte in data:
        checksum ^= byte
    return bytes([checksum])

def generate_bsj_command(command_str, terminal_phone, serial_number):
    """
    Generates a BSJ-EG01 text message command packet (Message ID: 0x8300).
    
    Args:
        command_str: The command content string (e.g., "<SPBSJ*P:BSJGPS*C:30>")
        terminal_phone: The terminal's phone number (BCD encoded)
        serial_number: Message serial number (0-65535)
    
    Returns:
        The complete command packet as a hexadecimal string
    """
    try:
        # Start flag bit
        flag_bit = b'\x7e'
        
        # Message ID (0x8300 for text message delivery)
        message_id = struct.pack('>H', 0x8300)
        
        # Message body (command string)
        message_body = b'\x00' + command_str.encode('gbk')  # 0x00 = normal flag (not emergency)
        body_length = len(message_body)
        
        # Message body properties
        # bits 0-9: message body length
        # bits 10-12: encryption (0 = not encrypted)
        # bit 13: subpacket flag (0 = no subpackets)
        body_props = struct.pack('>H', body_length)
        
        # Terminal phone number (BCD[6])
        if len(terminal_phone) != 12:
            # Pad to 12 digits
            terminal_phone = terminal_phone.zfill(12)
        phone_bytes = bytes.fromhex(''.join([terminal_phone[i:i+2] for i in range(0, len(terminal_phone), 2)]))
        
        # Message serial number
        serial_bytes = struct.pack('>H', serial_number)
        
        # Header (without check code)
        header = message_id + body_props + phone_bytes + serial_bytes
        
        # Full data for checksum calculation
        data_for_checksum = header + message_body
        
        # Calculate checksum
        check_code = calculate_checksum(data_for_checksum)
        
        # Construct full packet before escaping
        raw_packet = header + message_body + check_code
        
        # Apply escaping rules:
        # 0x7e -> 0x7d 0x02
        # 0x7d -> 0x7d 0x01
        escaped_packet = bytearray()
        for byte in raw_packet:
            if byte == 0x7e:
                escaped_packet.extend(b'\x7d\x02')
            elif byte == 0x7d:
                escaped_packet.extend(b'\x7d\x01')
            else:
                escaped_packet.append(byte)
        
        # Full packet with flags
        full_packet = flag_bit + bytes(escaped_packet) + flag_bit
        
        return {
            'hex': binascii.hexlify(full_packet).decode('ascii'),
            'packet': full_packet
        }
        
    except Exception as e:
        print(f"An error occurred: {e}")
        return None
