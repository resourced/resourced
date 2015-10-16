#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import print_function

'''
usage:
cat about.json | python stdin-stdout.py
'''

import sys

if __name__ == '__main__':
    print(sys.stdin.read())
