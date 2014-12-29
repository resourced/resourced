#!/usr/bin/env ruby

require 'json'

load_avgs_string = `uptime`.split('load averages:')[1].strip
load_avgs = load_avgs_string.split(' ').map(&:to_f)

puts Hash['LoadAvg1m' => load_avgs[0], 'LoadAvg5m' => load_avgs[1], 'LoadAvg15m' => load_avgs[2]].to_json
