#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import print_function

import sys
import json
import collections

def read():
    return json.loads(sys.stdin.read())

def flatten(d, parent_key='', sep='_'):
    items = []

    for k, v in d.items():
        new_key = parent_key + sep + k if parent_key else k
        if isinstance(v, collections.MutableMapping):
            items.extend(flatten(v, new_key, sep=sep).items())
        else:
            items.append((new_key, v))

    return dict(items)

if __name__ == '__main__':
    print(json.dumps(flatten(read())))
