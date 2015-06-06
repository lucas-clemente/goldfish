# Goldfish

[![Build Status](https://travis-ci.org/lucas-clemente/goldfish.svg?branch=master)](https://travis-ci.org/lucas-clemente/goldfish)

[Download](https://github.com/lucas-clemente/goldfish/releases)

A personal wiki / notes blend powered by markdown and git.

Inspired by Evernote and gollum.

## Features

- Just files on your disk. Actually, a git repo on your disk.
- A server and a web-frontend, all in a single executable. The server manages files, your browser displays them.
- Files are markdown with some extensions:
  - LaTeX `$\latex$` or, if you want it on its own line, `\[ \latex \]`)
  - Easier links `[[foo]]`, also works for images `[[foo.png]]`

Future features:

- Search
- Auto-push
- Auto-Update of files and folders
- Windows support

## Usage

If you're on linux, make sure to install either `libinotifytools` or `inotifytools`.

```bash
# Or whatever path suits you
./goldfish ~/goldfish
```

Then open [http://localhost:2345](http://localhost:2345) and start writing those markdown files :)

## Screenshot time!

![](screen.png)

Code for the page:

    # Demo Page

    ## Markdown

    Things you could do:

    - Make _important_ notes
    - Write in __strong__ letters

    ## Equations

    Both $e^\text{inline}$ and in display mode:

    \[
      e^{i \pi} = -1
    \]

    ## Syntax Highlighting

    ```ruby
    foo = Bar.new
    puts foo if foo.baz?
    ```

    ## Images

    [[fish.png]]
