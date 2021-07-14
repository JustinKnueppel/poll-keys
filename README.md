# Synchronize Remote File

This is a simple utility that a local file synchronized with a remote file

## Usage

`sync-remote-file -remote <url> -path <path/to/file> [-perms 0644] [-interval 10s] [-temp /tmp/dummy.txt]`

## Arguments

| Flag       | Type            | Default              | Required           | Description                                                           |
| ---------- | --------------- | -------------------- | ------------------ | --------------------------------------------------------------------- |
| `remote`   | remote url      | none                 | :heavy_check_mark: | Remote file location                                                  |
| `path`     | file path       | none                 | :heavy_check_mark: | Local file location                                                   |
| `perms`    | file permission | `0644`               | :x:                | File permissions for local file                                       |
| `interval` | time interval   | 0                    | :x:                | Time interval to poll remote file location - 0 will run once and exit |
| `temp`     | file path       | `/tmp/sync_file.txt` | :x:                | Initial save location for remote file                                 |
