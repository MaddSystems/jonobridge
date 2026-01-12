import re
import sys

def decode_percent_encoded(text):
    try:
        # Manual decoding of percent-encoded sequences
        def replace_percent(match):
            hex_value = match.group(1)
            try:
                # Convert hex to byte and decode as Latin-1 (since %e9=é, %e1=á in Latin-1)
                byte = bytes.fromhex(hex_value)
                char = byte.decode('latin-1')
                return char
            except Exception as e:
                print(f"Failed to decode % {hex_value}: {str(e)}")
                return '\ufffd'  # Return replacement character on failure

        # Replace all %XX sequences
        result = re.sub(r'%([0-9A-Fa-f]{2})', replace_percent, text)
        return result
    except Exception as e:
        return f"Error decoding {text}: {str(e)}"

# Ensure stdout uses UTF-8 encoding
sys.stdout.reconfigure(encoding='utf-8')

# Example usage
if __name__ == "__main__":
    test_strings = ["M%e9xico", "Tehuac%e1n"]
    for s in test_strings:
        decoded = decode_percent_encoded(s)
        unicode_points = [hex(ord(c)) for c in decoded] if isinstance(decoded, str) else []
        print(f"Original: {s} --> Decoded: {decoded} (Unicode: {unicode_points})")