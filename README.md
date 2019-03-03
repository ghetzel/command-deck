![CommandDeck Demo](contrib/logo-tagline.png?raw=true)

# Overview

Command Deck is a simple utility inspired by [Powerline](https://github.com/powerline/powerline) for creating simple, customizable terminal prompts that take advantage of Unicode and high-color terminal environments to create functional, beautiful, and data-rich command prompts.

# Examples

![Sample prompts](contrib/demo.png?raw=true)


# Installation

Retrieve the program using Golang's `go get` utility:

`go get github.com/ghetzel/command-deck`

And add the following to the end of your `~/.bashrc` file:

```bash
if which cdeck 2>&1 > /dev/null; then
    PROMPT_COMMAND='export LASTEXIT=$?; PS1=$(LASTEXIT=$LASTEXIT cdeck)'
fi
```

Note that `$LASTEXIT` is not special in any way.  We capture the last exit status via `$?` and store it in this variable, which is explicitly referenced in the configuration file.  If you don't need to get the last exit status, you can simply use `PROMPT_COMMAND='PS1=$(cdeck)'` instead.

### "How do I get that cool looking chevron separator?"

<figure>
    <img src="contrib/chevron.png?raw=true" alt="uE0B0" />
    <figcaption>"The Chevron"</figcaption>
</figure>

The "chevron" is actually a custom full-height character stored in the Private Use Area of pre-modified fonts.  You can find the fonts Powerline uses here: https://github.com/powerline/fonts.  The examples shown here are using _DejaVu Sans Mono for Powerline (Book)_.

# Usage

CommandDeck is configured in two ways: via positional arguments in _eval_ mode (using the `--eval`/`-e` flag), or via a YAML configuration file located at `~/.config/cdeck.yml`.

## Expressions

Expressions allow for external commands to be evaluated for the purpose of using their output while building prompt segments.  Expressions are strings that contain literal text and _command substitutions_, which are surrounded by `{!` and `!}`.  When a command substitution is encountered, the command specified will be run and its standard output will be used in place of the substitution.  For example, the expression `Hello, {! whoami !}!` would expand to "Hello, ghetzel!".

## Config File

A configuration file is used to specify the individual segments that make up the CommandDeck output.  Expressions are accepted in the `fg`, `bg`, `if`, `except`, and `expr` fields.  This sample configuration file shows the options that produced the example prompt above:

```yaml
segments:
# Will output the value of the `$LASTEXIT` environment variable with a background color of red and
# a foreground color of *bold* white (ANSI color #15 + "b" for bold), but only if the command `test
# ${LASTEXIT}" -ne 0` exits with status 0.
-   name:    lastexit
    bg:      red
    fg:      '15+b'
    if:      'test "${LASTEXIT}" -ne 0'
    expr:    '{! echo "${LASTEXIT:-?}" !}'

# Output the current HH:MM:SS time by executing `date +%H:%M:%S`.  Foreground is grey (ANSI #252,
# "Grey82"), background is darker grey (ANSI #240, "Grey35").
-   name:    time
    fg:      252
    bg:      240
    expr:    "{! date +%H:%M:%S !}"

# Output the current username and hostname as "user@host" by running `whoami` and `hostname -s`
# (respectively).  Foreground is white (ANSI #15), background is selected on-the-fly by reading `~/
# .ps1hostcolor`, and falling back to navy blue (ANSI #4).
-   name:    user
    fg:      15
    bg:      '{! cat ~/.ps1hostcolor || echo 4 !}'
    expr:    "{! whoami !}@{! hostname -s !}"

# Output the current directory as returned from `pwd`, substiuting the first occurrence of the
# string contained in `${HOME}` with a tilde (~).  Foreground is white, background is chartreuse
# (ANSI #64, "Chartreuse4").
-   name:    path
    fg:      15
    bg:      64
    expr:    '{! pwd | sed "s|${HOME}|~|" !}'

# Conditionally output a segment showing the current Git branch.  Foreground is white, background
# is determined by the exit status of `test -z "$(git status -s)"`.  If it exits successfully,
# there are no modified files in the current repository and the background color will be ANSI #24
# (DeepSkyBlue4).  If the command exits non-zero, there _are_ uncommitted changes and the
# background will be ANSI #9 (Red).  The commands must exit within 250ms or the segment will not be
# shown.  This is a protection measure against terminal slowdown caused by entering large,
# unindexed Git repositories.
-   name:    git
    if:      'git rev-parse --abbrev-ref HEAD'
    fg:      15
    bg:      '{! test -z "$(git status -s)" && echo 24 || echo 9 !}'
    expr:    "{! git rev-parse --abbrev-ref HEAD !}"
    timeout: 250ms
```

## Command Line (Eval Mode)

You can also invoke CommandDeck to print one-off evaluations for various purposes.  This approach is useful in a wide variety of circumstances that use text formatted inputs for various purposed (for example, the [i3 window manager's](https://i3wm.org/docs/userguide.html#_configuring_i3bar) `i3-bar` project could be configured this way.)

#### Invocation

```bash
cdeck -e -T ' ' -s '|' \
    '::MEM {! free -w -g -h | grep "^Mem" | tr -s " " | cut -d" " -f8 !}' \
    "::CPU {! echo $(( $(cat /sys/class/thermal/thermal_zone0/temp) / 1000 )) !}$(printf '\u00B0')C" \
    "::LOAD {! cat /proc/loadavg | cut -d' ' -f1-3 !}" \
    "::{! date '+%Y-%m-%d %H:%M:%S' !}"
```

#### Output

```
 MEM 8.6G | CPU 39°C | LOAD 0.61 0.40 0.32 | 2006-01-02 15:04:05
```

Each positional parameter is specified as a "foreground:background:expression" triple.  If either the foreground or background is omitted, the default color for the terminal will be used.  The expression is interpreted the same way as expressions in the config file.