package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	IDKey      = "next.url.id"
	ShortenKey = "short:%s:url"
	UrlKey     = "url:%s:short"
	DetailKey  = "short:%s:detail"
)

type RedisCli struct {
	Cli *redis.Client
}

type ShortLinkDetail struct {
	URL      string `json:"url"`
	CreateAt string `json:"create_at"`
	Expired  int64  `json:"exp"`
}

func NewRedisCli(addr string, passwd string, db int) *RedisCli {

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}
	return &RedisCli{
		Cli: client,
	}

}

func toSha(origin string) string {
	sum := sha256.Sum256([]byte(origin))
	return fmt.Sprintf("%x", sum)
}

func (c *RedisCli) Shorten(url string, exp int64) (string, error) {
	// check url exists
	h := toSha(url)

	ctx := context.Background()
	if short, err := c.Cli.Get(ctx, fmt.Sprintf(UrlKey, h)).Result(); err != redis.Nil {
		if err != nil {
			panic(err)
			return "", err
		}
		return short, nil
	}

	// increase the id
	if id, err := c.Cli.Incr(ctx, IDKey).Result(); err != nil {
		panic(err)
		return "", err
	} else {
		strOfId := Base62Encode(id)
		// save the relationship
		_, err = c.Cli.Set(ctx, fmt.Sprintf(ShortenKey, strOfId), url, time.Duration(exp)*time.Minute).Result()
		if err != nil {
			panic(err)
			return "", err
		}
		_, err = c.Cli.Set(ctx, fmt.Sprintf(UrlKey, h), strOfId, time.Duration(exp)*time.Minute).Result()
		if err != nil {
			panic(err)
			return "", err
		}

		detail := &ShortLinkDetail{
			URL:      url,
			CreateAt: time.Now().Format(time.UnixDate),
			Expired:  exp,
		}
		var b []byte
		b, err = json.Marshal(detail)
		if err != nil {
			panic(err)
			return "", err
		}
		_, err = c.Cli.Set(ctx, fmt.Sprintf(DetailKey, strOfId), string(b), time.Duration(exp)*time.Minute).Result()
		if err != nil {
			panic(err)
			return "", err
		}

		return "/" + strOfId, nil
	}

}

func (c *RedisCli) UnShorten(shortlink string) (string, error) {
	ctx := context.Background()
	if url, err := c.Cli.Get(ctx, fmt.Sprintf(ShortenKey, shortlink)).Result(); err != redis.Nil {
		if err != nil {
			panic(err)
			return "", err
		}
		return url, nil
	}
	return "", nil
}

func (c *RedisCli) Detail(shortlink string) (interface{}, error) {
	ctx := context.Background()
	if str, err := c.Cli.Get(ctx, fmt.Sprintf(DetailKey, shortlink)).Result(); err != redis.Nil {
		if err != nil {
			panic(err)
			return "", err
		}
		var detail ShortLinkDetail
		if err = json.Unmarshal([]byte(str), &detail); err != nil {
			panic(err)
			return "", err
		}
		return detail, nil
	}
	return "", nil
}
