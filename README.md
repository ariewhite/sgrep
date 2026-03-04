# sgrep — Roadmap / TODO

## Phase 1 — Basic functionality

- [x] Read from stdin
- [x] Search for a substring
- [x] Highlight matches with ANSI colors
- [x] Print line numbers

**Correct exit codes:**

- [x] `0` — match found
- [x] `1` — no match
- [x] `2` — error

## Phase 2 — Input handling

- [*] Support `-f <file>`
- [*] Automatically read from stdin when no file is provided
- [*] Support multiple files:

```sh
sgrep pattern file1.txt file2.txt
```

- [x] Print filename before the line (like `grep`)

## Phase 3 — Search features
- [x] Fix regular exp. and double output
- [x] Fix command: sgrep <file> (loop)
- [x] Fix highlighting for regex matches
- [x] Support regular expressions
- [x] Flag `-v` (invert match)
- [x] Flag `-i` (ignore case)
- [x] Flag `-c` (count matches)

- [ ] Redesign the arch, to input - matcher - result - formatter - output

- [ ] Flag `-l` (print only file names)
- [ ] Flag `-w` (match whole word)
- [ ] Flag `-o` (print only matched part)

## Phase 4 — Output improvements

- [ ] Disable colors when output is not a terminal
- [ ] Flag `--color=auto|always|never`
- [ ] Support `-n` (line number as a separate flag)
- [ ] Align output
- [ ] Support context:
  - [ ] `-A <n>` (after)
  - [ ] `-B <n>` (before)
  - [ ] `-C <n>` (context)

## Phase 5 — Performance & Robustness

- [ ] Increase `bufio.Scanner` buffer (>64KB)
- [ ] Handle very long lines
- [ ] Use `bufio.Reader` instead of `Scanner` when appropriate
- [ ] Profiling (`pprof`)
- [ ] Minimize allocations
- [ ] Benchmark with large files (100MB+)
- [ ] Parallel search across multiple files

## Phase 6 — Architecture & Clean Code

- [ ] Split into components:
  - [ ] CLI parsing
  - [ ] Input handling
  - [ ] Search engine
  - [ ] Output formatting
- [ ] Define a `Searcher` interface
- [ ] Use dependency injection for testability
- [ ] Unit tests
  - [ ] Table-driven tests
  - [ ] Golden tests

## Phase 7 — Advanced / Systems level

- [ ] `mmap` file reading
- [ ] Zero-copy optimizations
- [ ] Streaming regexp engine
- [ ] Worker pool for file processing
- [ ] Graceful cancelation via `context.Context`
- [ ] Handle `SIGINT`
- [ ] Fuzz testing

## Phase 8 — Production polish

- [ ] Makefile
- [ ] CI
- [ ] Cross-compilation
- [ ] Release binary
- [ ] README with examples
- [ ] Versioning
- [ ] Man page
