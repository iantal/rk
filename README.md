# rk
Repository Keeper

The zip file should contain hidden directories and files as well. For example .git is required by `ld`.

Command to generate the zip file

```
zip archiveName.zip -r .* * -x "../*"
```