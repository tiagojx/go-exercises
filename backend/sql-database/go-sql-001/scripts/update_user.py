#!/usr/bin/env python

import subprocess
import sys
import logging

try:
    user = {"id": sys.argv[1], "name": sys.argv[2]}
    cmd = [
        "curl", 
        "-v", 
        "-X", 
        "PATCH", 
        "http://127.0.0.1:8080/users/%(id)s/" % user, 
        "-d", 
        #'{"name": "%(name)s", "email": "%(email)s"}' % user
        '{"name": "%(name)s"}' % user
    ]
    proc = subprocess.run(cmd, check=True, stderr=sys.stderr, stdout=sys.stdout)
    if proc.returncode != 0:
        raise subprocess.CalledProcessError(proc.returncode, proc.stderr)
except Exception as e:
    logging.error(e)
