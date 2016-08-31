# Syntax help

### Basic formatting


use `**bold**` for **bold**

use `*italic*` for *italic*

### Tables

A table    | Age
-----------|------
Bob        | 27
Alice      | 23


```
A table    | Age
-----------|------
Bob        | 27
Alice      | 23
```

### Code blocks

```
code blocks
```


````
```
code blocks
```
````

### Code elements

`code elements`

```
`code elements`
```

### Cites

> cites

```
> cites
```

### Paragraphs

Paragraphs are separated
by a blank line.

Two spaces at the end of a line leave a  
line break.

``` 
Paragraphs are separated
by a blank line.

Two spaces at the end of a line leave a  
line break.
```

### Horizontal rule

---

```
---
```

### Lists

Bullet list:

  * apples

Numbered list:

  1. apples

```
Bullet list:

  * apples

Numbered list:

  1. apples
```

### Links

A [link](http://example.com).

```
A [link](http://example.com).
```

### Graphviz syntax

```dot
a->b [color=red]
c->t
```

````
```dot
a->b [color=red]
c->t
```
````

### Js sequence diagrams

```jsseq||scale=0.7
Title: Here is a title
A->B: Normal line
B-->C: Dashed line
C->>D: Open arrow
D-->>A: Dashed open arrow

Note left of A: Note to the\n left of A
Note right of A: Note to the\n right of A
Note over A: Note over A
Note over A,B: Note over both A and B
```

````
```jsseq||scale=0.7
Title: Here is a title
A->B: Normal line
B-->C: Dashed line
C->>D: Open arrow
D-->>A: Dashed open arrow

Note left of A: Note to the\n left of A
Note right of A: Note to the\n right of A
Note over A: Note over A
Note over A,B: Note over both A and B
```
````

Additional attributes:

- theme=<hand> or <simple>

### Goat diagrams

```goat
.---.       .-.        .-.       .-.                                       .-.
| A +----->| 1 +<---->| 2 |<----+ 4 +------------------.                  | 8 |
'---'       '-'        '+'       '-'                    |                  '-'
                       |         ^                     |                   ^
                       v         |                     v                   |
                      .-.      .-+-.        .-.      .-+-.      .-.       .+.       .---.
                     | 3 +---->| B |<----->| 5 +---->| C +---->| 6 +---->| 7 |<---->| D |
                      '-'      '---'        '-'      '---'      '-'       '-'       '---'
```

````
```goat
.---.       .-.        .-.       .-.                                       .-.
| A +----->| 1 +<---->| 2 |<----+ 4 +------------------.                  | 8 |
'---'       '-'        '+'       '-'                    |                  '-'
                       |         ^                     |                   ^
                       v         |                     v                   |
                      .-.      .-+-.        .-.      .-+-.      .-.       .+.       .---.
                     | 3 +---->| B |<----->| 5 +---->| C +---->| 6 +---->| 7 |<---->| D |
                      '-'      '---'        '-'      '---'      '-'       '-'       '---'
```
````

### Js flowchart

```jsflow||scale=0.8
st=>start: Start:>http://www.google.com[blank]
e=>end:>http://www.google.com
op1=>operation: My Operation
sub1=>subroutine: My Subroutine
cond=>condition: Yes
or No?:>http://www.google.com
io=>inputoutput: catch something...

st->op1->cond
cond(yes)->io->e
cond(no)->sub1(right)->op1
```

````
```jsflow||scale=0.8
st=>start: Start:>http://www.google.com[blank]
e=>end:>http://www.google.com
op1=>operation: My Operation
sub1=>subroutine: My Subroutine
cond=>condition: Yes
or No?:>http://www.google.com
io=>inputoutput: catch something...

st->op1->cond
cond(yes)->io->e
cond(no)->sub1(right)->op1
```
````


### About images & attaching files

* Drop an image or a file into the markdown editor to add it
  to the gopad filesystem

* When using images you can scale it easely

![100px](https://blog.golang.org/gopher/header.jpg)

![200px](https://blog.golang.org/gopher/header.jpg)

![300px](https://blog.golang.org/gopher/header.jpg)

![400px](https://blog.golang.org/gopher/header.jpg)

```
![100px](https://blog.golang.org/gopher/header.jpg)

![200px](https://blog.golang.org/gopher/header.jpg)

![300px](https://blog.golang.org/gopher/header.jpg)

![400px](https://blog.golang.org/gopher/header.jpg)
```