#!/usr/bin/env python

# ctags file must be generated with:
#   ctags --fields=afmikKlnsStz --c-kinds=+pc -R

import ctags, re, simplejson, sys, os

def main():
    if len(sys.argv) != 3:
        print 'usage: python gen_syscalls.py tbl_file tags_file'
        sys.exit(2)
    tbl_path = sys.argv[1]
    tags_path = sys.argv[2]

    tags = ctags.CTags(tags_path)
    entry = ctags.TagEntry()

    tbl_file = open(tbl_path, 'r')
    syscalls = []
    for line in tbl_file:
        syscall = {'num': -1, 'entry': '', 'name': '', 'args': []}

        match = re.search(r'^(\w+)\t+\w+\t+(\w+)\t+(\w+)', line)
        if not match:
            continue

        num = match.group(1)
        name = match.group(2)
        entrypoint = match.group(3)
        if not tags.find(entry, entrypoint, ctags.TAG_FULLMATCH | ctags.TAG_OBSERVECASE):
            continue

        syscall['num'] = num
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
                arg = {'name': '', 'type': '', 'refcount': ''}
                arg['name'] = re.split(r'[ *]+', strarg)[-1]
                arg['type'] = strarg
                arg['refcount'] = strarg.count('*')
                args.append(arg)

        syscall['args'] = args

        syscalls.append(syscall)

    tbl_file.close()

    print simplejson.dumps({'syscalls': syscalls}, indent='\t')

if __name__ == "__main__":
    main()
