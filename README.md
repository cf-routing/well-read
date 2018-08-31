## Well-read

To test that the sample app returns no 502s the engineering team has run this ab
command with a post payload of both 10KB and 4KB.

### Generate a large json payload

```bash
dd if=/dev/urandom of=10k.json bs=1024 count=10
```

### Running the bench

```bash
ab -T "application/json" -n 1000 -c 100 -l -p 10k.json http://<your-app-name>.<your-domain>/api/boomerangnsq
```

### Sample AB output
```bash
Server Software:
Server Hostname:        well-read.dev-full-1.routing.cf-app.com
Server Port:            80

Document Path:          /api/boomerangnsq
Document Length:        Variable

Concurrency Level:      100
Time taken for tests:   3.683 seconds
Complete requests:      1000
Failed requests:        0 <------------
Total transferred:      174000 bytes
Total body sent:        10413000
HTML transferred:       0 bytes
Requests per second:    271.51 [#/sec] (mean)
Time per request:       368.307 [ms] (mean)
Time per request:       3.683 [ms] (mean, across all concurrent requests)
Transfer rate:          46.14 [Kbytes/sec] received
                        2761.00 kb/s sent
                        2807.13 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        1    1   0.7      1       5
Processing:    49  348 212.7    392     839
Waiting:       49  348 212.7    392     839
Total:         50  349 212.7    393     840

Percentage of the requests served within a certain time (ms)
  50%    393
  66%    487
  75%    538
  80%    556
  90%    607
  95%    664
  98%    742
  99%    803
 100%    840 (longest request)
 ```
