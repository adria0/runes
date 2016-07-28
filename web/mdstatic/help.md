# Syntax help

Gopad uses 3 main formatters:

- Markdown from https://github.com/russross/blackfriday
- Graphviz from http://www.graphviz.org/Documentation.php
- UMLlet from http://www.umlet.com/

Edit this entry to see markdown source code of it

---
---

## Markdown syntax

**bold** *italic*

A table    | Age
-----------|------
Bob        | 27
Alice      | 23

```
code blocks
```

`code elements`

> cites

### Header level 3
#### Header level 4
##### Header level 5
 
Paragraphs are separated
by a blank line.

Two spaces at the end of a line leave a  
line break.

Text attributes _italic_, *italic*, __bold__, **bold**, `monospace`.

Horizontal rule:


Bullet list:

  * apples

Numbered list:

  1. apples

A [link](http://example.com).

---
---

## Graphviz syntax

```dot
a->b [color=red]
c->t
```

---
---

## UMLet sequence diagrams

```umlet:sequence
obj=Usr~usr
obj=App~app

usr->app
app->app +:ACT=CreateKeyPair
app->usr:show PIN=\nhash pbkACT
```

---
---

## About images & attaching files

* Drop an image or a file into the markdown editor to add it
  to the gopad filesystem

* When using images you can scale it easely

![100px](https://blog.golang.org/gopher/header.jpg)

![200px](https://blog.golang.org/gopher/header.jpg)

![300px](https://blog.golang.org/gopher/header.jpg)

![400px](https://blog.golang.org/gopher/header.jpg)

