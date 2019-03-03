![CommandDeck Demo](contrib/logo-tagline.png?raw=true)

# Overview

Command Deck is a simple utility inspired by [Powerline](https://github.com/powerline/powerline) for creating simple, customizable terminal prompts that take advantage of Unicode and high-color terminal environments to create functional, beautiful, and data-rich command prompts.

# Installation

`go get github.com/ghetzel/command-deck`

# Usage

CommandDeck is configured in two ways: via positional arguments in _eval_ mode (using the `--eval`/`-e` flag), or via a YAML configuration file located at `~/.config/cdeck.yml`.

## Config File

A sample configuration file shows an example of all of the options for a sample prompt:

```yaml
segments:
-   name:    lastexit
    bg:      red
    fg:      '15+b'
    if:      'test "${LASTEXIT}" -ne 0'
    expr:    '${! echo "${LASTEXIT:-?}" !}'

-   name:    time
    fg:      252
    bg:      240
    expr:    "${! date +%H:%M:%S !}"

-   name:    user
    fg:      15
    bg:      '${! cat ~/.ps1hostcolor || echo 4 !}'
    expr:    "${! whoami !}@${! hostname -s !}"

-   name:    path
    fg:      15
    bg:      64
    expr:    '${! pwd | sed "s|${HOME}|~|" !}'

-   name:    git
    if:      "git rev-parse --abbrev-ref HEAD"
    fg:      15
    bg:      "${! git diff-files --quiet && echo 24 || echo 9 !}"
    expr:    "${! git rev-parse --abbrev-ref HEAD !}"
    timeout: 250ms
```

Here is a description of what each entry in the above `segments` array does:

### `lastexit`

Will output the value of the `$LASTEXIT` environment variable with a background color of red and a foreground color of *bold* white (ANSI color #15 + "b" for bold), but only if the command `test "${LASTEXIT}" -ne 0` exits with status 0.

### `time`

Output the current HH:MM:SS time by executing `date +%H:%M:%S`.  Foreground is grey (ANSI #252, "Grey82"), background is darker grey (ANSI #240, "Grey35").

### `user`

Output the current username and hostname as "user@host" by running `whoami` and `hostname -s` (respectively).  Foreground is white (ANSI #15), background is selected on-the-fly by reading `~/.ps1hostcolor`, and falling back to navy blue (ANSI #4).

### `path`

Output the current directory as returned from `pwd`, substiuting the first occurrence of the string contained in `${HOME}` with a tilde (~).  Foreground is white, background is chartreuse (ANSI #64, "Chartreuse4").

### `git`

Conditionally output a segment showing the current Git branch.  Foreground is white, background is determined by the exit status of `git diff-files --quiet`.  If it exits successfully, there are no modified files in the current repository and the background color will be ANSI #24 (DeepSkyBlue4).  If the command exits non-zero, there _are_ uncommitted changes and the background will be ANSI #9 (Red).  The commands must exit within 250ms or the segment will not be shown.  This is a protection measure against terminal slowdown caused by entering large, unindexed Git repositories.