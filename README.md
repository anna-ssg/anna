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
|  |--page.html (This file is necessary and cannot be omitted)
|  |--posts.html (This file is necessary to create a 'Posts' section) 
|  |--partials/
|     |--header.html
|  ....
|  |--config.yml (This file is necessary and cannot be omitted)
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
   |--posts.html
   |--posts/
      |--post1.html
   |--static/
      |--image1.jpg
      |--fonts/
   ....
```

## Description of the directory structure

- The markdown content for the site is stored in content/. It can contain subdirectories as the folder is recursively rendered
- Static assets such as images and fonts are stored in static/
- The layout of the site is configured using html files in layout/

   - The 'config.yml' file stores the configuration of the site and includes details such as the baseURL
   - The 'page.html' file defines the layout of a basic page of the site
   - The 'posts.html' file defines the layout of a page displaying all the posts of the site
   - The layout files can be composed of smaller html files which are stored in the partials/ folder

#### Layout

The layout files can access the following rendered data from the markdown files:

- {{.Body}} : Returns the markdown body rendered to HTML
- {{.Frontmatter.[Tagname]}} : Returns the value of the frontmatter tag
   - Example: {{.Frontmatter.Title}} : Returns the value of the title tag
- {{.Layout.[Tagname]}}: Returns the particular configuration detail of the page
   - Example: {{.Layout.Navbar}} : Returns a string slice with the names of all the navbar elements

## Notes

1. Images: To add images, add it to the 'static/' folder or a subdirectory under it. Use "static/[imagename.format]" as the image link format in the markdown files.

2. CSS: CSS can be added in the following ways:

- In an external file in the 'static/' directory and linked to the layout files
    - To link the stylesheet, use the baseURL along with the relative path

       Example: `<link rel="stylesheet" href="{{.Layout.BaseURL}}static/style.css">`
- Placed inside `<style></style>` tags in the `<head></head>` of the layout files
- Inline with the html elements

3. Frontmatter: Metadata such as the title of the page can be added as frontmatter to the markdown files in the YAML format. Currently, the following tags are supported:

- title : The title of the current page
- date: The date of the current page

4. config.yml: This file stores additional information regarding the layout

- navbar: Stores the links to be added to the navbar (same name as the markdown files)
- baseURL: Stores the base URL of the site

Sample config.yml:

```yml
navbar:
  - about
  - posts

baseURL: http://localhost:8000/
# Replace this with the actual canonical-url of your site.

# baseURL tells search-engines (SEO), web-crawlers (robots.txt) so people can discover your site on the internet. It's also embeded in your sitemap / atom feed and can be used to change metadata about your site. 
```

## Flags

```
Usage:
  ssg [flags]

Flags:
  -a, --addr string   ip address to serve rendered content to (default "8000")
  -d, --draft         renders draft posts
  -h, --help          help for ssg
  -s, --serve         serve the rendered content
```
