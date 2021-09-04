# MinerHive
Standardised set of tools to help monitor and maintain CGMiner based ASICs at a glance.

# The Project

### Quick glance:
<img src="/docs/screenshots/high-level-overview.png" alt="High level overview" width="700" height="516"/>

<img src="/docs/screenshots/grafana-example.png" alt="Grafana example" width="700" height="688"/>

# Configuration
All configuration files are to be structured as YAML files.

The configuration path can be stipulated to a custom path with the `--config <path-to-file>` flag. 

Otherwise, `config.yml` file will be searched for in these directories:

```
./
$HOME/<program>/
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

# Miner compatibiltiy:
Developed against Antminer S19 (Firmware version: Mon Apr 19 16:36:50 CST 2021 / 49.0.1.3). This should work for any miner that uses the CGMiner
software implementation. If it doesn't work, create an issue - I'll help get it supported.

# Why?
I need stats and programmable functionality for a machine I own. 

# Todo:
- Testing
- Control functionality (changing pools, restarting, setting fan speeds etc...)
- DB change? (PSQL + GORM)
- REST change? (GraphQL via GQLGen)