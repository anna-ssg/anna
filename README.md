# Static Site Generator

## Directory structure

The ssg currently requires the following directory structure

```text
ssg/
|--content/  
|  |--index.md (This file is necessary and cannot be omitted)  
|  |--about.md  
|  |--posts/
|     |--post1.md
|  ....  
|
|--layout/  
|  |--layout.html (This file is necessary and cannot be omitted)
|  |--posts.html (This file is necessary to create a 'Posts' section) 
|  |--config.yml (This file is necessary for the navbar and 'Posts' section) 
|
|--static/
|  |--image1.jpg
|  |--fonts
|     |--font1.ttf
|  |--images
|     |--image2.png
|  ....
|
|--rendered/ (This directory is created by the ssg)
   |--index.html
   |--about.html
   |--posts.html (Generated by the ssg)
   |--posts/
      |--post1.html
   |--static/
      |--image1.jpg
      |--fonts/
   ....
```

### Layout

The layout files can access the following rendered data from the markdown files:

- {{.Body}} : Returns the markdown body rendered to HTML
- {{.Frontmatter.[Tagname]}} : Returns the value of the frontmatter tag
   - Example: {{.Frontmatter.Title}} : Returns the value of the title tag
- {{.Layout.[Tagname]}}: Returns the particular configuration detail of the page
   - Example: {{.Layout.Navbar}} : Returns a string slice with the names of all the navbar elements

### Notes

1. Images: To add images, add it to the 'static/' folder or a subdirectory under it. Use "static/[imagename.format]" as the image link format in the markdown files.

2. CSS: CSS can be added in the following ways:

- In an external file in the 'static/' directory and linked to layout.html
- Placed inside `<style></style>` tags in the `<head></head>` of layout.html
- Inline with the html elements

3. Frontmatter: Metadata such as the title of the page can be added as frontmatter to the markdown files in the YAML format. Currently, the following tags are supported:

- title : The title of the current page
- date: The date of the current page

4. config.yml: This file stores additional information regarding the layout

- navbar: Stores the links to be added to the navbar (same name as the markdown files)
- posts: Stores the posts to be listed in "posts.html" (same name as the markdown files in the posts/ folder)
- baseURL: Stores the base URL of the site

Sample config.yml:

```yml
navbar:
  - about
  - posts

posts:
  - post1

baseURL: http://localhost:8000/
```

### Flags

1. -serve=true (default: false) : Serves the rendered content on the browser

2. -addr=":6060" (default: ":8000") : Specifies the address over which the rendered files are served
