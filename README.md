# GRPC boilerplate

## How the configuration handled

config, DB client, logger will self initialize if it imported based on the configuration provided. If configuration block missing in configuration source(config.yml), If will silently fail to initialize.
