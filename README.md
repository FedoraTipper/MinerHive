# AntHive
Webcrawler API provider for ssh-locked Antminers

# The Project


# Configuration
All configuration files are to be structured as YAML files.

The configuration path can be stipulated to a custom path with the `--config <path-to-file>` flag. 

Otherwise, `config.yml` file will be searched for in these directories:

```
./
./configs/<program-name>/
./config/<program-name>/
$HOME/.config/<program-name>/

e.g. ./config/minercrawler/config.yml
```

## Logging
List of available logging levels:
- debug
- info (default)
- warning
- error
- fatal

Developed against Antminer S19 (Firmware version: Mon Apr 19 16:36:50 CST 2021 / 49.0.1.3)

# Why?
I need stats and programmable functionality for a machine I own. 
