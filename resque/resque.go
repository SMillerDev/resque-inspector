package resque

import (
	"context"
	"github.com/valkey-io/valkey-go"
	"log"
)

const Prefix = "resque:"

var Client valkey.Client
var Dsn string

func PrepareClient() {
	if Client != nil {
		return
	}
	preClient, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{Dsn}})
	if err != nil {
		log.Default().Fatal(err)
	}

	Client = preClient
}

func GetList(key string, shouldPrefix bool) []string {
	PrepareClient()
	ctx := context.Background()
	if shouldPrefix {
		key = Prefix + key
	}
	members, err := Client.Do(ctx, Client.B().Smembers().Key(key).Build()).AsStrSlice()
	if err != nil {
		log.Default().Printf("Failed to get list for \"%s\": %s", key, err)
		return make([]string, 0)
	}

	return members
}

func GetEntryCount(queue string, shouldPrefix bool) int64 {
	PrepareClient()
	ctx := context.Background()
	if shouldPrefix {
		queue = Prefix + queue
	}
	count, err := Client.Do(ctx, Client.B().Llen().Key(queue).Build()).AsInt64()
	if err != nil {
		log.Default().Printf("Failed to get entry count for \"%s\": %s", queue, err)
		return -1
	}

	return count
}

func GetEntries(queue string, shouldPrefix bool) []string {
	PrepareClient()
	ctx := context.Background()
	if shouldPrefix {
		queue = Prefix + queue
	}

	count := GetEntryCount(queue, false)
	if count <= 0 {
		return []string{}
	}

	entries, err := Client.Do(ctx, Client.B().Lrange().Key(queue).Start(0).Stop(count).Build()).AsStrSlice()
	if err != nil {
		log.Default().Printf("Failed to get entries for \"%s\": %s", queue, err)
		return []string{}
	}

	return entries
}

func GetEntry(entry string, shouldPrefix bool) string {
	PrepareClient()
	ctx := context.Background()
	if shouldPrefix {
		entry = Prefix + entry
	}
	item, err := Client.Do(ctx, Client.B().Get().Key(entry).Build()).ToString()
	if err != nil {
		log.Default().Printf("Failed to get entry for \"%s\": %s", entry, err)
		return ""
	}

	return item
}

/**

resque_failed_all_jobs() {
  local failed_count

  failed_count=$( resque_redis LLEN resque:failed )

  resque_redis LRANGE resque:failed 0 "${failed_count}"
}

resque_clear_queue() {
  local queue="$1"

  if [ -z "${queue}" ]; then
    error "Missing queue argument to clear"
  fi

  local key result

  if [ "${queue}" = 'failed' ]; then
    key='resque:failed'
  else
    key="resque:queue:${queue}"
  fi

  result=$( resque_redis DEL "${key}" )

  if [ "${result}" -eq 0 ]; then
    return 1
  fi
}
*/
