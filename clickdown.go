package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/panjf2000/ants"
	log "github.com/sirupsen/logrus"

	"gopkg.in/ns3777k/go-shodan.v3/shodan"
)

// ErrUnreachable is the server unreachable error
var ErrUnreachable = errors.New("service is unreachable")

var (
	// GetClickHousePartsQuery is the query to obtain server version
	GetClickHouseVersionQuery = `SELECT version() AS version`

	// GetClickHousePartsQuery is the query to obtain parts (tables) size
	GetClickHousePartsQuery = `SELECT concat(database, '.', table) AS name,
	    formatReadableSize(sum(bytes)) as size
	FROM system.parts
	WHERE active
	GROUP BY name`
)

type _ClickDown struct {
	cfg    *_Config
	client *shodan.Client
	pool   *ants.Pool
}

type _ServerMeta struct {
	Server       _ClickHouseServer
	Parts        []_ClickHousePart
}

type _ClickHouseServer struct {
	Version string `db:"version"`
}

type _ClickHousePart struct {
	Name string `db:"name"`
	Size string `db:"size"`
}

func _NewClickDown(cfg *_Config) (*_ClickDown, error) {
	client := shodan.NewClient(&http.Client{}, cfg.ShodanAPIKey)
	client.SetDebug(cfg.Debug)
	pool, err := ants.NewPool(cfg.MaxWorkers)
	if err != nil {
		return nil, err
	}
	return &_ClickDown{cfg, client, pool}, nil
}

func (cd *_ClickDown) Run(query string) error {
	var wg sync.WaitGroup

	page := 1

	for {
		opts := &shodan.HostQueryOptions{
			Query:  query,
			Page:   page,
			Minify: true,
		}
		resp, err := cd.client.GetHostsForQuery(context.Background(), opts)
		if err != nil {
			if strings.HasPrefix(err.Error(), "Request rate limit") {
				log.Debug("rate limiter triggered")
				time.Sleep(time.Second * 3)
				continue
			}
			return err
		}
		for _, match := range resp.Matches {
			wg.Add(1)

			endpoint := fmt.Sprintf("%s:%d", match.IP.String(), match.Port)
			task := cd.WrapCheckServer(endpoint)
			err := cd.pool.Submit(func() {
				task()
				wg.Done()
			})
			if err != nil {
				return err
			}
		}

		if len(resp.Matches) != resp.Total && len(resp.Matches) != 0 {
			page++
		} else {
			break
		}
	}

	wg.Wait()

	return nil
}

func (cd *_ClickDown) WrapCheckServer(endpoint string) func() {
	wrapper := func() {
		meta, err := cd.CheckServer(endpoint)
		if err != nil {
			log.Debugf("failed to check server: %v", err)
			return
		}
		version := fmt.Sprintf("ClickHouse %s", meta.Server.Version)
		for _, part := range meta.Parts {
			log.Infof("table \"%s\" (%s) on \"%s\" (%s)", part.Name,
				part.Size, endpoint, version)
		}
	}
	return wrapper
}

func (cd *_ClickDown) CheckServer(endpoint string) (*_ServerMeta, error) {
	opts := url.Values{}
	opts.Set("read_timeout", strconv.Itoa(cd.cfg.ClickHouseReadTimeout))
	if cd.cfg.Debug {
		opts.Set("debug", "true")
	}
	dsn := fmt.Sprintf("tcp://%s?%s", endpoint, opts.Encode())
	log.Debugf("try open using dsn \"%s\"", dsn)
	db, err := sqlx.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Errorf("failed to close db: %v", err)
		}
	}()
	if err := db.Ping(); err != nil {
		return nil, ErrUnreachable
	}

	parts := []_ClickHousePart{}
	err = db.Select(&parts, GetClickHousePartsQuery)
	if err != nil {
		return nil, err
	}

	server := _ClickHouseServer{}
	row := db.QueryRowx(GetClickHouseVersionQuery)
	if err != nil {
		return nil, err
	}
	if err := row.StructScan(&server); err != nil {
		return nil, err
	}

	meta := &_ServerMeta{
		Server: server,
		Parts:  parts,
	}
	return meta, nil
}

func (cd *_ClickDown) Shutdown() error {
	return cd.pool.Release()
}
