package resque

import (
	"context"
	"github.com/valkey-io/valkey-go"
	"log"
	"strings"
	"time"
)

const Prefix = "resque:"

var Client valkey.Client
var Dsn string

func GetList(key string) []string {
	prepareClient()
	ctx := context.Background()
	key = ensurePrefix(key)
	members, err := Client.DoCache(ctx, Client.B().Smembers().Key(key).Cache(), time.Minute).AsStrSlice()
	if err != nil {
		log.Default().Printf("Failed to get list for \"%s\": %s", key, err)
		return make([]string, 0)
	}

	return members
}

func GetEntryCount(key string) int64 {
	prepareClient()
	ctx := context.Background()
	key = ensurePrefix(key)
	count, err := Client.Do(ctx, Client.B().Llen().Key(key).Build()).AsInt64()
	if err != nil {
		log.Default().Printf("Failed to get entry count for \"%s\": %s", key, err)
		return -1
	}

	return count
}

func GetEntries(key string, start int64, limit int64) []string {
	prepareClient()
	ctx := context.Background()
	key = ensurePrefix(key)

	count := GetEntryCount(key)
	if count <= start {
		return []string{}
	}
	if count <= limit {
		limit = count
	}

	entries, err := Client.Do(ctx, Client.B().Lrange().Key(key).Start(start).Stop(count).Build()).AsStrSlice()
	if err != nil {
		log.Default().Printf("Failed to get entries for \"%s\": %s", key, err)
		return []string{}
	}

	return entries
}

func GetEntry(key string) string {
	entry, err := getEntry(key)
	if err != nil {
		log.Default().Printf("Failed to get entry for \"%s\": %s", entry, err)
		return ""
	}

	return entry
}

func GetEntryOrNil(key string) string {
	entry, _ := getEntry(key)

	return entry
}

func prepareClient() {
	if Client != nil {
		return
	}
	preClient, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{Dsn}})
	if err != nil {
		log.Default().Fatal(err)
	}

	Client = preClient
}

func ensurePrefix(key string) string {
	if strings.HasPrefix(key, Prefix) {
		return key
	}

	return Prefix + key
}

func getEntry(key string) (string, error) {
	prepareClient()
	ctx := context.Background()
	key = ensurePrefix(key)
	item, err := Client.Do(ctx, Client.B().Get().Key(key).Build()).ToString()
	return item, err
}

func Clear(key string) error {
	prepareClient()
	ctx := context.Background()
	key = ensurePrefix(key)
	err := Client.Do(ctx, Client.B().Del().Key(key).Build()).Error()

	return err
}

func Retry(queue string, payload string) error {
	prepareClient()
	ctx := context.Background()
	key := ensurePrefix(queue)

	return Client.Do(ctx, Client.B().Rpush().Key(key).Element(payload).Build()).Error()
}
