# CPDI

## Build:
```
go build -o cpdi cmd/cpdi/main.go
```

## Using:

```
./cpdi \
	--src-directory data_1 \
	--dst-directory data2 \
	--min-file-size 0B \
	--max-file-size 1G \
	--exclude-dir-path folder_1/folder_2_excl \
	--exclude-file-path folder_1/file3_excl.txt \
	--exclude-common-names do_not_copy \
	--verbose
```

## Example:
```
./cpdi \
    --src-directory "/media/user/disk/folder_1" \
    --dst-directory /home/user/Desktop/dst_dir \
    --min-file-size 100M \
    --max-file-size 1G \
    --exclude-dir-path inner_folder_1:folder_1/inner_folder_2 \
    --exclude-file-path some_folder/file1.txt \
    --exclude-common-names any_name:file_name_1:dir_name_1 \
    --verbose
```
