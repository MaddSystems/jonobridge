#!/usr/bin/env python3
"""
fetch_skywave.py
Fetch Skywave get_return_messages XML and print a concise summary.
Defaults to ACCESS_ID=70001184, PASSWORD=JEUTPKKH and uses FROMIDSKYWAVE env or --from-id.
"""
import os
import sys
import argparse
from collections import Counter
import requests
import xml.etree.ElementTree as ET

ACCESS_ID = os.getenv('SKYWAVE_ACCESS_ID', '70001184')
PASSWORD = os.getenv('SKYWAVE_PASSWORD', 'JEUTPKKH')

URL_TEMPLATE = 'https://isatdatapro.skywave.com/GLGW/GWServices_v1/RestMessages.svc/get_return_messages.xml/?access_id={access}&password={passw}&from_id={fromid}'


def fetch_and_parse(access_id, password, from_id, timeout=30):
    url = URL_TEMPLATE.format(access=access_id, passw=password, fromid=from_id)
    resp = requests.get(url, timeout=timeout)
    resp.raise_for_status()
    xml = resp.text
    root = ET.fromstring(xml)
    return root, xml


def summarize(root):
    # ErrorID
    error_id_elem = root.find('ErrorID')
    error_id = error_id_elem.text if error_id_elem is not None else None
    more_elem = root.find('More')
    more = more_elem.text.lower() == 'true' if more_elem is not None else False
    next_start_elem = root.find('NextStartID')
    next_start = next_start_elem.text if next_start_elem is not None and next_start_elem.text else None

    messages = root.findall('.//ReturnMessage')
    total = len(messages)

    per_mobile = Counter()
    per_payload = Counter()
    ids = []
    times = []
    samples = []

    for i, m in enumerate(messages):
        id_elem = m.find('ID')
        mid = m.find('MobileID')
        payload = m.find('Payload')
        msgutc = m.find('MessageUTC')
        recvutc = m.find('ReceiveUTC')
        ids.append(id_elem.text if id_elem is not None else '')
        if mid is not None and mid.text:
            per_mobile[mid.text] += 1
        plname = payload.get('Name') if payload is not None and 'Name' in payload.attrib else ''
        if plname:
            per_payload[plname] += 1
        if msgutc is not None and msgutc.text:
            times.append(msgutc.text)
        # collect a small sample for display
        if i < 10:
            # try find latitude/longitude fields
            lat = ''
            lon = ''
            if payload is not None:
                fields = payload.find('Fields')
                if fields is not None:
                    for f in fields.findall('Field'):
                        name = f.get('Name')
                        if name == 'Latitude':
                            lat = f.get('Value')
                        if name == 'Longitude':
                            lon = f.get('Value')
            samples.append({
                'ID': id_elem.text if id_elem is not None else '',
                'MobileID': mid.text if mid is not None else '',
                'Payload': plname,
                'MessageUTC': msgutc.text if msgutc is not None else '',
                'Latitude': lat,
                'Longitude': lon,
            })

    # first/last ids
    first_id = ids[0] if ids else None
    last_id = ids[-1] if ids else None
    first_time = times[0] if times else None
    last_time = times[-1] if times else None

    return {
        'error_id': error_id,
        'more': more,
        'next_start': next_start,
        'total': total,
        'per_mobile': per_mobile,
        'per_payload': per_payload,
        'first_id': first_id,
        'last_id': last_id,
        'first_time': first_time,
        'last_time': last_time,
        'samples': samples,
    }


def print_summary(s, access_id, from_id):
    print('\nSummary for access_id={} from_id={}'.format(access_id, from_id))
    print('----------------------------------------------')
    print('ErrorID:', s['error_id'])
    print('Total messages:', s['total'])
    print('More:', s['more'])
    print('NextStartID:', s['next_start'])
    print('First ID:', s['first_id'], 'Last ID:', s['last_id'])
    print('First MessageUTC:', s['first_time'], 'Last MessageUTC:', s['last_time'])
    print('\nMessages per MobileID:')
    for mid, cnt in s['per_mobile'].most_common():
        print('  {}: {}'.format(mid, cnt))
    print('\nMessages per Payload type:')
    for p, cnt in s['per_payload'].most_common():
        print('  {}: {}'.format(p, cnt))

    print('\nSample messages (up to 10):')
    for samp in s['samples']:
        print('  ID={ID} Mobile={MobileID} Payload={Payload} Time={MessageUTC} Lat={Latitude} Lon={Longitude}'.format(**samp))


def main(argv=None):
    p = argparse.ArgumentParser(description='Fetch and summarize Skywave get_return_messages')
    p.add_argument('--from-id', dest='from_id', help='from_id value (default FROMIDSKYWAVE env or 13969586728)', default=os.getenv('FROMIDSKYWAVE', '13969586728'))
    p.add_argument('--access', dest='access', help='access id', default=ACCESS_ID)
    p.add_argument('--password', dest='password', help='password', default=PASSWORD)
    p.add_argument('--raw', dest='raw', help='print raw XML', action='store_true')
    args = p.parse_args(argv)

    try:
        root, xml = fetch_and_parse(args.access, args.password, args.from_id)
    except Exception as e:
        print('Error fetching data:', e)
        sys.exit(1)

    if args.raw:
        print(xml)
        return

    s = summarize(root)
    print_summary(s, args.access, args.from_id)


if __name__ == '__main__':
    main()
