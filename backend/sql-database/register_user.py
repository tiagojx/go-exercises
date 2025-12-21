#!/usr/bin/env python

import subprocess
import sys
import logging

try:
    user = {"name": sys.argv[1], "email": sys.argv[2]}
    cmd = [
        "curl", 
        "-X", 
        "POST", 
        "http://127.0.0.1:8080/register", 
        "-d", 
        '{"name": "%(name)s", "email": "%(email)s"}' % user
    ]
    proc = subprocess.run(cmd, check=True, stderr=sys.stderr, stdout=sys.stdout)
    if proc.returncode != 0:
        raise subprocess.CalledProcessError(proc.returncode, proc.stderr)
except Exception as e:
    logging.error(e)
else:
    print("Created new user in database!")
