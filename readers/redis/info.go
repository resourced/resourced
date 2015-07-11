package redis

import (
	"encoding/json"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/resourced/resourced/readers"
)

func init() {
	readers.Register("RedisInfo", NewRedisInfo)
}

func NewRedisInfo() readers.IReader {
	r := &RedisInfo{}
	r.Data = make(map[string]string)

	return r
}

type RedisInfo struct {
	Data map[string]string
	Base
}

func (r *RedisInfo) Run() error {
	err := r.initConnection()
	if err != nil {
		return err
	}

	data, err := redis.String(connections[r.HostAndPort].Do("INFO"))
	if err != nil {
		return err
	}

	/*
	   {"aof_current_rewrite_time_sec":"-1","aof_enabled":"0","aof_last_bgrewrite_status":"ok","aof_last_rewrite_time_sec":"-1","aof_last_write_status":"ok","aof_rewrite_in_progress":"0","aof_rewrite_scheduled":"0","arch_bits":"64","blocked_clients":"0","client_biggest_input_buf":"0","client_longest_output_list":"0","config_file":"/usr/local/etc/redis.conf","connected_clients":"1","connected_slaves":"0","db0":"keys=325,expires=0,avg_ttl=0","evicted_keys":"0","expired_keys":"0","gcc_version":"4.2.1","hz":"10","instantaneous_input_kbps":"0.00","instantaneous_ops_per_sec":"0","instantaneous_output_kbps":"0.00","keyspace_hits":"0","keyspace_misses":"0","latest_fork_usec":"0","loading":"0","lru_clock":"15307932","master_repl_offset":"0","mem_allocator":"libc","mem_fragmentation_ratio":"1.56","multiplexing_api":"kqueue","os":"Darwin 14.1.0 x86_64","process_id":"61455","pubsub_channels":"0","pubsub_patterns":"0","rdb_bgsave_in_progress":"0","rdb_changes_since_last_save":"0","rdb_current_bgsave_time_sec":"-1","rdb_last_bgsave_status":"ok","rdb_last_bgsave_time_sec":"-1","rdb_last_save_time":"1424593316","redis_build_id":"70633d1af7244f5e","redis_git_dirty":"0","redis_git_sha1":"00000000","redis_mode":"standalone","redis_version":"2.8.19","rejected_connections":"0","repl_backlog_active":"0","repl_backlog_first_byte_offset":"0","repl_backlog_histlen":"0","repl_backlog_size":"1048576","role":"master","run_id":"7e5c1f3da4b7ac0048259bf97c66120ed1555822","sync_full":"0","sync_partial_err":"0","sync_partial_ok":"0","tcp_port":"6379","total_commands_processed":"9","total_connections_received":"12","total_net_input_bytes":"140","total_net_output_bytes":"16897","uptime_in_days":"0","uptime_in_seconds":"760","used_cpu_sys":"0.46","used_cpu_sys_children":"0.00","used_cpu_user":"0.17","used_cpu_user_children":"0.00","used_memory":"1298192","used_memory_human":"1.24M","used_memory_lua":"35840","used_memory_peak":"1298192","used_memory_peak_human":"1.24M","used_memory_rss":"2031616"}
	*/

	for _, line := range strings.Split(data, "\r\n") {
		if strings.Contains(line, ":") {
			keyAndValue := strings.Split(line, ":")
			r.Data[keyAndValue[0]] = keyAndValue[1]
		}
	}

	return nil
}

// ToJson serialize Data field to JSON.
func (r *RedisInfo) ToJson() ([]byte, error) {
	return json.Marshal(r.Data)
}
