![Docker Publish](https://github.com/iantal/rk/workflows/Docker%20Publish/badge.svg) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# RK - Repository Keeper

This service handles the upload functionality. The git repository of the project will be uploaded as a zip archive and stored on the local storage (disk). 

The zip file should contain hidden directories and files as well. For example the .git directory. 

Command to generate the zip file

```
zip archiveName.zip -r .* * -x "../*"
```

To upload the zip, use the following curl command or it's equivalent:

`curl localhost:8002/api/v1/projects/<projectName> --data-binary @<projectName>.zip`