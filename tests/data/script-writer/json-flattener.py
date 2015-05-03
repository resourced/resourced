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
    parser.add_argument('-d', '--data', dest='user_data', help='custom user defined data to merged into JSON output. (Format: key=value,key2=value2)')
    parser.add_argument('-t', '--test', dest='test', action="store_true", default=False, help='Run tests')
    return parser


def read():
    return json.loads(sys.stdin.read())


def read_user_data(user_data):
    if not user_data:
        return {}

    user_data_list = [x.strip() for x in args.user_data.split(',')]
    return dict([arg.split('=') for arg in user_data_list])


def flatten(d, parent_key='', sep='.'):
    items = []

    for k, v in d.items():
        new_key = parent_key + sep + k if parent_key else k
        if isinstance(v, collections.MutableMapping):
            items.extend(flatten(v, new_key, sep=sep).items())
        else:
            items.append((new_key, v))

    return dict(items)


def run_tests(user_dict):
    in_string = '''{"Data":{"ApiVersion":"1.16","Arch":"amd64","Containers":10,"Debug":true,"DockerRootDir":"/mnt/sda1/var/lib/docker","Driver":{"Dirs":86,"Name":"aufs","RootDir":"/mnt/sda1/var/lib/docker/aufs"},"ExecutionDriver":"native-0.2","GitCommit":"5bc2ff8","GoVersion":"go1.3.3","ID":"A7LJ:M42E:ZPYS:TG37:56A5:NMGP:B7NQ:6UJ5:JUTL:MUXB:Y5E6:6XTI","IPv4Forwarding":true,"Images":58,"IndexServerAddress":"https://index.docker.io/v1/","InitPath":"/usr/local/bin/docker","InitSha1":"","KernelVersion":"3.16.7-tinycore64","Labels":"null","MemTotal":2.106249216e+09,"MemoryLimit":true,"NGoroutines":38,"Name":"boot2docker","NumCPUs":2,"NumEventsListeners":0,"NumFileDescriptors":32,"OperatingSystem":"Boot2Docker 1.4.1 (TCL 5.4); master : 86f7ec8 - Tue Dec 16 23:11:29 UTC 2014","Os":"linux","SwapLimit":true,"Version":"1.4.1"},"GoStruct":"DockerInfoVersion","Host":{"Name":"MacBook-Pro.local","Tags":[]},"Interval":"3s","Path":"/docker/info-version","UnixNano":1430676379672586965}'''
    in_json = json.loads(in_string)
    output  = flatten(in_json)

    if user_dict:
        output.update(user_dict)

    if in_json["Data"]["Driver"]["Name"] != output["Data.Driver.Name"]:
        print("Error: Input driver name does not match output driver name.")
    else:
        sys.stdout.write('.')

    for key, value in user_dict.items():
        if key not in output:
            print("Error: key: {0} not in output.".format(key))
        else:
            sys.stdout.write('.')

    print("All tests passed.")


if __name__ == '__main__':
    parser    = args_parser()
    args      = parser.parse_args()
    user_dict = read_user_data(args.user_data)

    if args.test:
        run_tests(user_dict)
        sys.exit(0)

    if not select.select([sys.stdin,],[],[],0.0)[0]:
        parser.print_help()
        sys.exit(0)

    output = flatten(read(), sep=args.separator)
    if user_dict:
        output.update(user_dict)

    print(json.dumps(output))
