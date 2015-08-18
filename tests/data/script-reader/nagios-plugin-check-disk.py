#!/usr/bin/python

import os
import sys

used_space = os.popen("df -h / | grep -v Filesystem | awk '{print $5}'").readline().strip()

if used_space < "85%":
        print "OK - %s of disk space used." % used_space
        sys.exit(0)
elif used_space == "85%":
        print "WARNING - %s of disk space used." % used_space
        sys.exit(1)
elif used_space > "85%":
        print "CRITICAL - %s of disk space used." % used_space
        sys.exit(2)
else:
        print "UKNOWN - %s of disk space used." % used_space
        sys.exit(3)
