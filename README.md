# Statusbar

This is a statusbar for a **linux** window managers written in **golang**.

Install **Hack** font. powerline works pretty good.

The **master** branch of this repo is same as the fork as a token of respect.

Currently it provides these details:

## Screenshot

![Screenshot](powerline.png)

- shows active keyboard layout, using **setxkbmap**.
- **network** connection details, **wifi** or **ethernet**, upload and
  download speeds.
- **cpu** temperature.
- **power** details, **AC** if on power cable, or remaining **battery**
  percentage.
- **cpu** load.
- **memory** utilization percent.
- **date** local date and time, plus one in different timezone in my case.

## Requirements

- **go** in order to compile statusbar.
- **dzen2** is the package used to render the status bar on your X11 screen

## Installation

You must have **go** installed on your system.

This repository is meant to be editable to your own needs, so fork or
clone and edit. Create your statusbar configuration:

    cp statusbar.dist.json statusbar.json

**NOTE:** the arguments for **dzen2** output formatting should be changed
on your needs.

If you run `make` it will build and move binary to
**/usr/local/bin/statusbar** and statusbar.json if available, to
**/usr/local/etc/statusbar.json**.

    make

If dependencies were not met, install them. Now you can run statusbar
which takes configuration option json as an argument:

    statusbar statusbar.json > /tmp/statusbar.log 2>&1

**NOTE:** you may change configuration properties based on your screen
layout and fonts. Statusbar logs errors to **stdout** and in case of panic
to **stderr**.

