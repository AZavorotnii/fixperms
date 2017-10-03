About
=====

This is basic solution if you messed permissions for files in some folder.
For example, I backed up my files on NTFS disk. All files permissions got reset to 777 after restoring. That was nasty.


Build
=====

  go build -o fixperms main.go


Run
===

Display planned changes, but don't change:

  ./fixperms -root ~/path/to/folder -test


Do changes with same details in output:

  ./fixperms -root ~/path/to/folder -verbose


Do changes silently and output only errors:

  ./fixperms -root ~/path/to/folder
