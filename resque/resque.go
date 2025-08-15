package resque

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/valkey-io/valkey-go"
)

const Prefix = "resque:"

var Client valkey.Client
var Dsn string

var Debug bool

func GetList(key string) []string {
	prepareClient()
	ctx := context.Background()
	key = ensurePrefix(key)
	if Debug {
		log.Printf("[Resque] GetList \"%s\"", key)
	}
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
	count, err := Client.DoCache(ctx, Client.B().Llen().Key(key).Cache(), time.Second*5).AsInt64()
	if err != nil {
		log.Default().Printf("Failed to get entry count for \"%s\": %s", key, err)
		return -1
	}

	if Debug {
		log.Printf("[Resque] GetEntryCount \"%s\": %d", key, count)
	}

	return count
}

func GetEntries(key string, start int64, limit int64) []string {
	prepareClient()
	ctx := context.Background()
	key = ensurePrefix(key)

	count := GetEntryCount(key)
	if count < start {
		return []string{}
	}
	if count <= limit {
		limit = count
	}

	if Debug {
		log.Printf("[Resque] GetEntries \"%s\": %d until %d out of %d", key, start, limit, count)
	}
	entries, err := Client.Do(ctx, Client.B().Lrange().Key(key).Start(start).Stop(limit).Build()).AsStrSlice()
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
	preClient, err := valkey.NewClient(valkey.MustParseURL(Dsn))
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

	if Debug {
		log.Println("[Resque] GetEntry", key)
	}

	item, err := Client.Do(ctx, Client.B().Get().Key(key).Build()).ToString()
	return item, err
}

func Clear(key string) error {
	prepareClient()
	ctx := context.Background()
	key = ensurePrefix(key)

	if Debug {
		log.Println("[Resque] Clear", key)
	}

	_, err := Client.Do(ctx, Client.B().Del().Key(key).Build()).AsBool()
	if Debug && err != nil {
		log.Printf("[Resque] Clearing '%s' failed: %s", key, err)
	}

	return err
}

func Delete(queue string, element string) error {
	prepareClient()
	ctx := context.Background()
	key := ensurePrefix(queue)

	if Debug {
		log.Printf("[Resque] Delete %s in queue %s", element, key)
	}

	err := Client.Do(ctx, Client.B().Lrem().Key(key).Count(1).Element(element).Build()).Error()
	if err != nil {
		log.Default().Printf("Failed to delete entry: %s", element)
		log.Default().Println(err)
	} else {
		log.Default().Printf("Deleted entry: %s", element)
	}

	return err
}

func Queue(queue string, element string) error {
	prepareClient()
	ctx := context.Background()
	key := ensurePrefix(queue)

	if Debug {
		log.Println("[Resque] Queueing on: ", key)
		log.Println("[Resque] Payload: ", element)
	}

	err := Client.Do(ctx, Client.B().Rpush().Key(key).Element(element).Build()).Error()

	if err != nil {
		log.Default().Printf("Failed to queue entry: %s", element)
		log.Default().Println(err)
	} else {
		log.Default().Printf("Queued entry: %s", element)
	}

	return err
}
