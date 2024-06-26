# novus
novus is a website change detector and news Atom feed generator. Novus focused on tracking changes in "sets" - tables, lists that represent discographies, filmographies, etc.

## Installation

The recommended way to install `novus` is with docker-compose:

create a `docker-compose.yml` file (this minimal example omits common settings like network, restart, etc):

```yaml
version: '3'
services:
  novus:
    container_name: novus
    image: ghcr.io/arelate/novus:latest
  environment:
      # - NV_PUBLISH_NOVUS-URL=http://NOVUS-ADDRESS
      # - NV_SYNC_NOVUS-URL=http://NOVUS-ADDRESS
    volumes:
      # backups
      - /docker/novus/backups:/usr/share/novus/backups
      # input
      - /docker/novus:/usr/share/novus/input
      #local-content
      - /docker/novus/local-content:/usr/share/novus/local-content
      #matched-content
      - /docker/novus/matched-content:/usr/share/novus/matched-content
      #output
      - /docker/novus:/usr/share/novus/output
      #redux
      - /docker/novus/redux:/usr/share/novus/redux
      # sharing timezone from the host
      - /etc/localtime:/etc/localtime:ro
    ports:
      # https://en.wikipedia.org/wiki/Acta_Diurna
      - "59222:59222"
```

After deployment you can use your favorite RSS-subscription service or app to add `novus` feed. To use that you need to know the hostname or IP-address of the server you've deployed it and use the following endpoint (also the only supported endpoint at the moment):

`http://HOSTNAME_OR_IPADDRESS:59222/atom`

## Setting up sources

[`source-example.txt`](https://github.com/boggydigital/novus/blob/main/sources-example.txt) gives and example of parameters you can use to define a source.  

## Updating data

To update your data you can use `novus` with a CLI interface. This approach can be used to schedule periodic updates (e.g. every day):

- most commonly you would run sync `docker-compose exec novus nv sync` that gets all data, diffs the changes and generates RSS updates.

NOTE: There are many ways to schedule periodic updates that are system dependent. This guide doesn't cover that topic. 