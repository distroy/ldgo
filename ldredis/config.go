/*
 * Copyright (C) distroy
 */

package ldredis

import (
	"crypto/tls"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type Config struct {
	Cluster  bool
	Addrs    []string
	Addr     string
	Password string
	DB       int // not use in cluster

	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration

	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// PoolSize applies per cluster node and not for the whole cluster.
	PoolSize     int
	MinIdleConns int
	PoolTimeout  time.Duration

	TLSConfig *tls.Config
}

func (cfg *Config) toClient() *redis.Options {
	// TODO: more configs
	return &redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,

		MaxRetries:      cfg.MaxRetries,
		MinRetryBackoff: cfg.MinRetryBackoff,
		MaxRetryBackoff: cfg.MaxRetryBackoff,

		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,

		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		PoolTimeout:  cfg.PoolTimeout,

		TLSConfig: cfg.TLSConfig,

		ContextTimeoutEnabled: true,
	}
}

func (cfg *Config) toCluster() *redis.ClusterOptions {
	// TODO: more configs
	return &redis.ClusterOptions{
		Addrs:    cfg.Addrs,
		Password: cfg.Password,

		MaxRetries:      cfg.MaxRetries,
		MinRetryBackoff: cfg.MinRetryBackoff,
		MaxRetryBackoff: cfg.MaxRetryBackoff,

		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,

		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		PoolTimeout:  cfg.PoolTimeout,

		TLSConfig: cfg.TLSConfig,

		ContextTimeoutEnabled: true,
	}
}
