#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import print_function

import json
import sys

EVENT_TYPE = "ServerDiskUsage"

def read():
    return json.loads(sys.stdin.read())

def reformat_output(in_dict):
    items = []

    if not in_dict:
        return items

    # the top most key is useless.
    in_dict   = in_dict.get("/du", {})
    host_json = in_dict.get("Host", None)
    hostname  = host_json.get("Name", "") if host_json else ""

    for path, du_data in in_dict.get("Data", {}).items():
        du_data["Hostname"] = hostname
        du_data["eventType"] = EVENT_TYPE

        items.append(du_data)

    return items

def run_tests():
    sys.stdout.write('Running tests')

    in_string = '''{"/du": {"UnixNano": 1.4306920221411776e+18, "Interval": "3s", "GoStruct": "Du", "Host": {"Name": "MacBook-Pro.local", "Tags": []}, "Path": "/du", "Data": {"/home": {"DeviceName": "map auto_home", "Used": 0, "InodesFree": 0, "InodesTotal": 0, "Free": 0, "InodesUsed": 0, "Path": "/home", "Total": 0}, "/dev": {"DeviceName": "devfs", "Used": 202752, "UsedPercent": 100, "InodesFree": 0, "InodesTotal": 686, "Free": 0, "InodesUsed": 686, "InodesUsedPercent": 100, "Path": "/dev", "Total": 202752}, "/Volumes/Main": {"DeviceName": "/dev/disk1s10", "Used": 198138892288.0, "UsedPercent": 93.1625119255904, "InodesFree": 3550301.0, "InodesTotal": 51924052.0, "Free": 14542032896.0, "InodesUsed": 48373751.0, "InodesUsedPercent": 93.16251166222543, "Path": "/Volumes/Main", "Total": 212680925184.0}, "/Volumes/Vagrant": {"DeviceName": "/dev/disk2s1", "Used": 82825216.0, "UsedPercent": 79.01914810472842, "InodesFree": 5369, "InodesTotal": 25588, "Free": 21991424.0, "InodesUsed": 20219, "InodesUsedPercent": 79.01750820697202, "Path": "/Volumes/Vagrant", "Total": 104816640.0}, "/Volumes/Files": {"DeviceName": "/dev/disk1s12", "Used": 50932432896.0, "UsedPercent": 47.55382632989547, "InodesFree": 13713958.0, "InodesTotal": 26148632.0, "Free": 56172371968.0, "InodesUsed": 12434674.0, "InodesUsedPercent": 47.55382231850599, "Path": "/Volumes/Files", "Total": 107104804864.0}, "/net": {"DeviceName": "map -hosts", "Used": 0, "InodesFree": 0, "InodesTotal": 0, "Free": 0, "InodesUsed": 0, "Path": "/net", "Total": 0}, "/": {"DeviceName": "/dev/disk0s2", "Used": 131473858560.0, "UsedPercent": 82.5933545127858, "InodesFree": 6700714.0, "InodesTotal": 38862822.0, "Free": 27446124544.0, "InodesUsed": 32162108.0, "InodesUsedPercent": 82.75803543036582, "Path": "/", "Total": 159182127104.0}}}}'''
    in_dict   = json.loads(in_string)
    output    = reformat_output(in_dict)

    if in_dict["/du"]["Host"]["Name"] != output[0]["Hostname"]:
        print("Error: Input hostname does not match output hostname.")
    else:
        sys.stdout.write('.')

    if EVENT_TYPE != output[0]["eventType"]:
        print("Error: Incorrect eventType.")
    else:
        sys.stdout.write('.')

    if len(in_dict["/du"]["Data"]) != len(output):
        print("Error: Data length does not match.")
    else:
        sys.stdout.write('.')


if __name__ == '__main__':
    if len(sys.argv) == 2:
        if sys.argv[1].startswith('test'):
            run_tests()
    else:
        in_dict = read()
        print(json.dumps(reformat_output(in_dict)))
