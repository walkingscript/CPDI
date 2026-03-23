package config

const HelpString string = `
To use the tool you need to specify at least a source and destination directory.
Progamy copies each file matched to conditions from src-dir to dst-dir.

Using:

	./cpdi \
		--src-directory data_1 \
		--dst-directory data2 \
		--min-file-size 0B \
		--max-file-size 1G \
		--exclude-dir-path folder_1/folder_2_excl \
		--exclude-file-path folder_1/file3_excl.txt \
		--exclude-common-names do_not_copy \
		--verbose

`
