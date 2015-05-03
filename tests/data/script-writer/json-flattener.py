#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import print_function

import argparse
import collections
import json
import select
import sys

def args_parser():
    parser = argparse.ArgumentParser(description='Reads JSON from STDIN and flatten everything.')
    parser.add_argument('-sep', '--separator', dest='separator', default='.', help='JSON key separator after flattened (default: ".")')
    return parser

def read():
    return json.loads(sys.stdin.read())

def flatten(d, parent_key='', sep='.'):
    items = []

    for k, v in d.items():
        new_key = parent_key + sep + k if parent_key else k
        if isinstance(v, collections.MutableMapping):
            items.extend(flatten(v, new_key, sep=sep).items())
        else:
            items.append((new_key, v))

    return dict(items)

def run_tests():
    in_string = '''{"Data":{"ApiVersion":"1.16","Arch":"amd64","Containers":10,"Debug":true,"DockerRootDir":"/mnt/sda1/var/lib/docker","Driver":{"Dirs":86,"Name":"aufs","RootDir":"/mnt/sda1/var/lib/docker/aufs"},"ExecutionDriver":"native-0.2","GitCommit":"5bc2ff8","GoVersion":"go1.3.3","ID":"A7LJ:M42E:ZPYS:TG37:56A5:NMGP:B7NQ:6UJ5:JUTL:MUXB:Y5E6:6XTI","IPv4Forwarding":true,"Images":58,"IndexServerAddress":"https://index.docker.io/v1/","InitPath":"/usr/local/bin/docker","InitSha1":"","KernelVersion":"3.16.7-tinycore64","Labels":"null","MemTotal":2.106249216e+09,"MemoryLimit":true,"NGoroutines":38,"Name":"boot2docker","NumCPUs":2,"NumEventsListeners":0,"NumFileDescriptors":32,"OperatingSystem":"Boot2Docker 1.4.1 (TCL 5.4); master : 86f7ec8 - Tue Dec 16 23:11:29 UTC 2014","Os":"linux","SwapLimit":true,"Version":"1.4.1"},"GoStruct":"DockerInfoVersion","Host":{"Name":"MacBook-Pro.local","Tags":[]},"Interval":"3s","Path":"/docker/info-version","UnixNano":1430676379672586965}'''
    in_json = json.loads(in_string)
    output  = flatten(in_json)

    if in_json["Data"]["Driver"]["Name"] != output["Data.Driver.Name"]:
        print("Error: Input driver name does not match output driver name.")
    else:
        sys.stdout.write('.')

    print("All tests passed.")


if __name__ == '__main__':
    if len(sys.argv) == 2:
        if sys.argv[1] in ["test", "tests"]:
            run_tests()
    else:
        parser = args_parser()

        if not select.select([sys.stdin,],[],[],0.0)[0]:
            parser.print_help()
            sys.exit(0)

        args = parser.parse_args()
        output = flatten(read(), sep=args.separator)

        print(json.dumps(output))
