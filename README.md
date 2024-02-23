# Static Site Generator

### Directory structure

The ssg currently requires the following directory structure
```text
ssg/
|--content/  
|  |--index.md (This file is necessary and cannot be omitted)  
|  |--about.md  
|    ....  
|
|--layout/  
|  |--layout.html (This file is necessary and cannot be omitted)
|
|--static/
|  |--image1.jpg
|    ....
|
|--rendered/ (This directory is created by the ssg)
   |--index.html  
   |--about.html  
   |--static/
      |--image1.jpg
       ....
```

### Notes:
1. Images: To add images, add it to the 'static/' folder or a subdirectory under it. Use "./static/[imagename.format]" as the image link format in the markdown files.

### Flags

1. -serve=true (default: false) : Serves the rendered content on the browser

2. -addr=":6060" (default: ":8000") : Specifies the address over which the rendered files are served
