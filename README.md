```
  __ _  ___   __ _ _   _  ___ _   _  ___ 
 / _` |/ _ \ / _` | | | |/ _ \ | | |/ _ \
| (_| | (_) | (_| | |_| |  __/ |_| |  __/
 \__, |\___/ \__, |\__,_|\___|\__,_|\___|
 |___/          |_|                      
 ```

Poc of golang server using queue

# Usefull Commands

Add item to queue
```bash
curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"type":"mytype","content":"mycontent"}' \
     http://localhost:8080/queue
```

Monitor memory usage:
```bash
ps goqueue -i VmRss --watch
```

Benchmark
```bash
autocannon -c 10 -p 1 -b '{"type": "mytype", "content": "mycontent"}' -H 'Content-Type: application/json' http://localhost:8080/queue
```
