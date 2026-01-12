#!/usr/bin/env python3
"""Analyze SkyWave XML files and list unique Payload Name values.

Usage:
  python3 skywave_analyze.py <file-or-directory> [--json out.json]

If a directory is provided, the script will scan all .xml files inside it.
"""
import argparse
import json
import os
import sys
import xml.etree.ElementTree as ET


def find_payload_names_in_file(path):
    names = set()
    try:
        tree = ET.parse(path)
        root = tree.getroot()
        # Find all Payload elements and extract Name attribute
        for payload in root.findall('.//Payload'):
            name = payload.get('Name')
            if name:
                names.add(name)
    except ET.ParseError as e:
        print(f"Failed to parse {path}: {e}", file=sys.stderr)
    except Exception as e:
        print(f"Error reading {path}: {e}", file=sys.stderr)
    return names


def main():
    parser = argparse.ArgumentParser(description='Extract unique Payload Name attributes from SkyWave XML files')
    parser.add_argument('path', help='File or directory to analyze')
    parser.add_argument('--json', '-j', help='Optional output JSON file to save mapping')
    args = parser.parse_args()

    paths = []
    if os.path.isdir(args.path):
        for fn in os.listdir(args.path):
            if fn.lower().endswith('.xml'):
                paths.append(os.path.join(args.path, fn))
    else:
        paths.append(args.path)

    all_names = set()
    for p in paths:
        print(f"Scanning {p}...")
        all_names.update(find_payload_names_in_file(p))

    sorted_names = sorted(all_names)

    print("\nFound Payload Names:")
    for n in sorted_names:
        print(n)

    if args.json:
        with open(args.json, 'w', encoding='utf-8') as f:
            json.dump(sorted_names, f, indent=2)
        print(f"\nSaved {len(sorted_names)} names to {args.json}")


if __name__ == '__main__':
    main()
