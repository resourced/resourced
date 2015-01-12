#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import print_function

'''
usage:
cat about.txt | python stdin-stdout.py
'''

import sys

def read_in():
    lines = sys.stdin.readlines()
    for i in range(len(lines)):
        lines[i] = lines[i].replace('\n','')
    return lines

def main():
    lines = read_in()
    print('\n'.join(lines))

if __name__ == '__main__':
    main()
