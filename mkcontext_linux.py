#!/usr/bin/env python

# Copyright 2015 The syscallinfo Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# Dependencies:
#   pip install python-ctags
# The ctags file must be generated with:
#   ctags --fields=afmikKlnsStz --c-kinds=+pc -R

import sys
import os
import re
import ctags
import simplejson

def main():
    if len(sys.argv) != 3:
        print 'usage: %s generate_context.py tbl_file tags_file' % sys.argv[0]
        sys.exit(2)
    tbl_path = sys.argv[1]
    tags_path = sys.argv[2]

    tags = ctags.CTags(tags_path)
    entry = ctags.TagEntry()

    tbl_file = open(tbl_path, 'r')
    syscalls = []
    for line in tbl_file:
        syscall = {'num': -1, 'entry': '', 'name': '', 'context':'', 'args': []}

        match = re.search(r'^(\w+)\t+\w+\t+(\w+)\t+(\w+)', line)
        if not match:
            continue

        num = match.group(1)
        name = match.group(2)
        entrypoint = match.group(3)

        # Rename stub_* entrypoints to sys_*
        entrypoint = re.sub(r'^stub_', r'sys_', entrypoint)

        if not tags.find(entry, entrypoint, ctags.TAG_FULLMATCH | ctags.TAG_OBSERVECASE):
            continue

        syscall['num'] = int(num)
        syscall['name'] = name
        syscall['entry'] = entrypoint

        found_prototype = False
        while not found_prototype:
            if(entry['kind'] == 'prototype'):
                found_prototype = True
            elif not tags.findNext(entry):
                break

        if not found_prototype:
            continue

        args = [];
        if(entry['signature'] != '(void)'):
            strargs = entry['signature'].strip('()').split(',')
            for strarg in strargs:
                strarg = strarg.strip()
                arg = {'sig': '', 'refcount': 0, 'context': ''}
                arg['sig'] = strarg
                arg['refcount'] = strarg.count('*')
                args.append(arg)

        syscall['args'] = args

        syscalls.append(syscall)

    tbl_file.close()

    print simplejson.dumps(syscalls, indent='\t')

if __name__ == "__main__":
    main()
