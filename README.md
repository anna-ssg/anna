# ssg prototype

### Directory structure

The ssg currently requires the following directory structure
```text
ssg/
|--content/  
|   |-- index.md (This file is necessary and cannot be omitted)  
|   |-- about.md  
|     .....  
|--layout/  
|   |-- layout.html  
|--rendered/  
|   |-- index.html  
|   |-- about.html  
```

### Flags

1. -serve=true (default: false) : Serves the rendered content on the browser

2. -addr=":6060" (default: ":8000") : Specifies the address over which the rendered files are served
