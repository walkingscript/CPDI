# My Utils Repository

## rcpd

### Build:
```
go build -o ../../bin/rcpd.exe utils/cmd/rcpd
```

### Using:
```
rcpd -src dir_1_src_ -dst dir_2_dst_ -v --min-file-size 1B --max-file-size 1G --exclude-dirs dir_1_src_\\4 .5555 --exclude-files dir_1_src_\\5\\fgh54.txt dir_1_src_\\1\\sfdsdv.txt
```
### Flags:
```
-src - where dir should be copied from
-dst - where dir should be copied to
-v - verbose output
--min-file-size - if file will be smaller than this size it will not be copied
--max-file-size - if file will be larger than this size it will not be copied
--exclude-dirs - after this param following many dir pathes, they can be absolute, relative, or just as name(for cases if you don't want to copy dirs with such name at all)
--exclude-files - after this param following many file pathes, they can be absolute, relative, or just as name(for cases if you don't want to copy files with such name at all)
```

