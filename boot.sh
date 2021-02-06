#!/bin/sh

/usr/bin/memcached \
  --user=${MEMCACHED_USER:-memcached} \
  --listen=${MEMCACHED_HOST:-127.0.0.1} \
  --port=${MEMCACHED_PORT:-11211} \
  --memory-limit=${MEMCACHED_MEMUSAGE:-64} \
  --conn-limit=${MEMCACHED_MAXCONN:-1024} \
  --threads=${MEMCACHED_THREADS:-4} \
  --max-reqs-per-event=${MEMCACHED_REQUESTS_PER_EVENT:-20} \
  --verbose
