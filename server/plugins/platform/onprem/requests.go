package onprem

const MetricsCommand = `sh -c '
echo __IVORY_CPU__; head -n 1 /proc/stat;
echo __IVORY_MEM__; grep -E "MemTotal|MemAvailable" /proc/meminfo;
echo __IVORY_NET__; cat /proc/net/dev'`

const ProcessesCommand = `ps -eo pid,user,pcpu,pmem,rss,nlwp,comm,args --no-headers --sort=-pcpu | head -n 100`
