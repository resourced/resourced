#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import print_function

import sys

def read_json_lines():
    lines = sys.stdin.readlines()
    for i in range(len(lines)):
        lines[i] = lines[i].replace('\n','')
    return lines
