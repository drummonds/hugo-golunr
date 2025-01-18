# `hugo-golunr`, a golang alternative to [hugo-lunr](https://www.npmjs.com/package/hugo-lunr)

As you probably don't like installing node, npm and a ton of packages into your CI, which generates
a static hugo page, I created this golang implementation of `hugo-lunr`. It generates a lunrjs
search index from the current working directory. 

## Installing

`go install github.com/riesinger/hugo-golunr@latest`

## Usage 

```sh
cd /path/to/your/site
hugo-golunr
```

Pretty easy, huh? After running `hugo-golunr`, you'll see a `search_index.json` file in your
`./static` directory. Just load that in your theme.

## Tidying
I have moved the code to a standard format for a command line format and added a mod.ls

## Change History

### V1.1.0 2023-05-29
- Fixing _index.md being resolved as ../_ instead of ../  

### V1.0.0 2023-03-31  

- Breaking change change location to be relative so can reuse json with local build as well as remote.
