#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import print_function

import argparse
import collections
import json
import select
import sys

def args_parser():
    parser = argparse.ArgumentParser(description='Reads JSON from STDIN and flatten everything from first level children.')
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

if __name__ == '__main__':
    parser = args_parser()

    if not select.select([sys.stdin,],[],[],0.0)[0]:
        parser.print_help()
        sys.exit(0)

    args = parser.parse_args()
    output = dict((path, flatten(data, sep=args.separator)) for path, data in read().items())

    print(json.dumps(output))
