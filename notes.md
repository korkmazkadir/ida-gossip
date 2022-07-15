# Notes to myself

* The following script lists md5 digest of each file in a folder. This is a good method to see duplicate files.

```ls |  cat | while read in ; do md5 $in ; done```
