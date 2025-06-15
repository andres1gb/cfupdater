# Cloudflare updater

This is a very simple updated for Cloudflare DNS, useful for home servers and other dynamic IP setups. You can choose between a few of sources to detect your IP and configure
multiple DNS zones.

## Compilation

Just `make build` and check the executable in the bin directory.

## Configuration

Copy the example config file, edit it to match your configuration and you're ready to go. A cron job to update the IP whenever is needed would do the work.

## To-do

This is a very early version, and lots of work is to be done. A few improvements for the future are:

- Allow for other type of records, currently only A records are allowed
- Skip single records. Currently, if a record for a zone doesn't exist, the whole zone will be skipped without update
- Logs!!!
- Add tests
- Add other DNS providers.