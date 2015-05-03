#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import print_function

import json
import sys

EVENT_TYPE = "ServerDiskUsage"

def read():
    return json.loads(sys.stdin.read())

def reformat_output(in_json):
    items = []

    if not in_json:
        return items

    host_json = in_json.get("Host", None)
    hostname  = host_json.get("Name", "") if host_json else ""

    for path, du_data in in_json.get("Data", {}).items():
        du_data["Hostname"] = hostname
        du_data["eventType"] = EVENT_TYPE

        items.append(du_data)

    return items

def run_tests():
    in_string = '''{"Data":{"/":{"DeviceName":"/dev/disk0s2","Free":2.9958176768e+10,"InodesFree":7.314008e+06,"InodesTotal":3.8862822e+07,"InodesUsed":3.1548814e+07,"InodesUsedPercent":81.17993592951125,"Path":"/","Total":1.59182127104e+11,"Used":1.28961806336e+11,"UsedPercent":81.01525509314506},"/Volumes/Files":{"DeviceName":"/dev/disk1s12","Free":5.6172371968e+10,"InodesFree":1.3713958e+07,"InodesTotal":2.6148632e+07,"InodesUsed":1.2434674e+07,"InodesUsedPercent":47.55382231850599,"Path":"/Volumes/Files","Total":1.07104804864e+11,"Used":5.0932432896e+10,"UsedPercent":47.55382632989547},"/Volumes/Main":{"DeviceName":"/dev/disk1s10","Free":1.4542032896e+10,"InodesFree":3.550301e+06,"InodesTotal":5.1924052e+07,"InodesUsed":4.8373751e+07,"InodesUsedPercent":93.16251166222543,"Path":"/Volumes/Main","Total":2.12680925184e+11,"Used":1.98138892288e+11,"UsedPercent":93.1625119255904},"/Volumes/Vagrant":{"DeviceName":"/dev/disk2s1","Free":2.1991424e+07,"InodesFree":5369,"InodesTotal":25588,"InodesUsed":20219,"InodesUsedPercent":79.01750820697202,"Path":"/Volumes/Vagrant","Total":1.0481664e+08,"Used":8.2825216e+07,"UsedPercent":79.01914810472842},"/dev":{"DeviceName":"devfs","Free":0,"InodesFree":0,"InodesTotal":686,"InodesUsed":686,"InodesUsedPercent":100,"Path":"/dev","Total":202752,"Used":202752,"UsedPercent":100},"/home":{"DeviceName":"map auto_home","Free":0,"InodesFree":0,"InodesTotal":0,"InodesUsed":0,"Path":"/home","Total":0,"Used":0},"/net":{"DeviceName":"map -hosts","Free":0,"InodesFree":0,"InodesTotal":0,"InodesUsed":0,"Path":"/net","Total":0,"Used":0}},"GoStruct":"Du","Host":{"Name":"MacBook-Pro.local","Tags":[]},"Interval":"3s","Path":"/du","UnixNano":1430674834475700306}'''
    in_json = json.loads(in_string)
    output  = reformat_output(in_json)

    if in_json["Host"]["Name"] != output[0]["Hostname"]:
        print("Error: Input hostname does not match output hostname.")
    else:
        sys.stdout.write('.')

    if EVENT_TYPE != output[0]["eventType"]:
        print("Error: Incorrect eventType.")
    else:
        sys.stdout.write('.')

    if len(in_json["Data"]) != len(output):
        print("Error: Data length does not match.")
    else:
        sys.stdout.write('.')

    print("All tests passed.")


if __name__ == '__main__':
    if len(sys.argv) == 2:
        if sys.argv[1].startswith('test'):
            run_tests()
    else:
        in_json = read()
        print(json.dumps(reformat_output(in_json)))
