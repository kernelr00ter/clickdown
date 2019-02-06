# ClickDown

![Git Tag](https://img.shields.io/github/tag/fdhadzh/clickdown.svg?style=flat)
![License](https://img.shields.io/github/license/fdhadzh/clickdown.svg?style=flat)
[![Go Report Card](https://goreportcard.com/badge/github.com/fdhadzh/clickdown)](https://goreportcard.com/report/github.com/fdhadzh/clickdown)
[![Docker Image Metadata](https://images.microbadger.com/badges/version/fdhadzh/clickdown.svg)](https://microbadger.com/images/fdhadzh/clickdown)

Explore vulnerable ClickHouse servers using Shodan.io

## Table of Contents

- [Disclaimer](#disclaimer)
- [Installing](#installing)
- [Usage](#usage)
  - [Configuration](#configuration)
  - [Running](#running)
  - [Using Docker](#using-docker)
- [License](#license)

### Disclaimer

When used properly, ClickDown helps protect your ClickHouse servers.
But when used improperly, ClickDown can get you sued, fired, expelled or jailed.
Reduce your risk by reading this legal guide before launching ClickDown.

### Installing

Install ClickDown by running:

```bash
go get github.com/fdhadzh/clickdown
```

and ensuring that $GOPATH/bin is added to your $PATH.

### Usage

#### Configuration

Configuration using environment variables:

| Variable                  | Description                        | Default |
|---------------------------|------------------------------------|---------|
| `DEBUG`                   | Debug flag                         | `false` |
| `MAX_WORKERS`             | Max workers count                  | `32`    |
| `CLICKHOUSE_READ_TIMEOUT` | ClickHouse read timeout in seconds | `10`    |
| `SHODAN_API_KEY`          | Shodan.io API key                  |         |

#### Running

ClickDown usage format:

```bash
clickdown <search-query>
```

> Where `search-query` is a Shodan.io search query ([Shodan.io Search Query Fundamentals](https://help.shodan.io/the-basics/search-query-fundamentals))

For example, explore ClickHouse servers from Russian country:

```bash
$ clickdown country:ru
INFO[0001] table "tsb_shard_2.app_method" (17.90 MiB) on "42.159.11.151:9000" (ClickHouse 18.16.1) 
INFO[0001] table "tsb_shard_2.app_method_aggr_min" (3.65 KiB) on "42.159.11.151:9000" (ClickHouse 18.16.1)
INFO[0001] table "axs.events" (2.10 GiB) on "46.166.165.41:9000" (ClickHouse 1.1.54394) 
INFO[0001] table "axs.transactions" (12.00 GiB) on "46.166.165.41:9000" (ClickHouse 1.1.54394) 
INFO[0001] table "tipico.bets" (1.76 GiB) on "94.130.175.22:9000" (ClickHouse 1.1.54385) 
INFO[0001] table "tipico.players" (3.77 MiB) on "94.130.175.22:9000" (ClickHouse 1.1.54385) 
INFO[0001] table "tipico.all_bets" (2.88 GiB) on "94.130.175.22:9000" (ClickHouse 1.1.54385) 
INFO[0001] table "tipico.events" (18.69 MiB) on "94.130.175.22:9000" (ClickHouse 1.1.54385) 
INFO[0001] table "tipico.tickets" (432.57 MiB) on "94.130.175.22:9000" (ClickHouse 1.1.54385)
...
```

Ensure that found server is open for everyone:

```bash
$ docker run -it --rm yandex/clickhouse-client --host 42.159.11.151 --port 9000
ClickHouse client version 19.1.6.
Connecting to 42.159.11.151:9000.
Connected to ClickHouse server version 18.16.1 revision 54412.

localhost :) SELECT count(*) FROM tsb_shard_2.app_method

SELECT count(*)
FROM tsb_shard_2.app_method 

┌─count()─┐
│  314034 │
└─────────┘

1 rows in set. Elapsed: 0.493 sec. Processed 314.03 thousand rows, 1.26 MB (637.59 thousand rows/s., 2.55 MB/s.)
```

#### Using Docker

```bash
docker run -e SHODAN_API_KEY="<shodan-api-key>" fdhadzh/clickdown
```

### License

MIT; see [LICENSE](/LICENSE) for details.