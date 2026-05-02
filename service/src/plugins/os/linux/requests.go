package linux

const MetricsCommand = `sh -lc '
echo __IVORY_CPU__; head -n 1 /proc/stat;
echo __IVORY_MEM__; grep -E "MemTotal|MemAvailable" /proc/meminfo;
echo __IVORY_NET__; cat /proc/net/dev'`
