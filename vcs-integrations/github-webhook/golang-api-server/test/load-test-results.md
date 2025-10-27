Load testing the GitHub webhook endpoint written in golang with [wrk tool](https://github.com/wg/wrk).

Golang app is containerized and deployed in Amazon Elastic Container Service (ECS)
Only one ECS task is run. Launch mode is Fargate (serverless).
Load tests are run from my machine, ~300km away from eu-west-3 AWS region (Paris). Network roundtrip is ~20ms.

# Observations
- server processes the request in <5ms, this includes:
    - validating HMAC signature
    - parsing json payload
    - publishing a message to SQS
- network roundtrip is 20ms
- TCP connection initialization adds one additional roundtrip
- TLS1.3 handshake adds one additional roundtrip
- This means that if we enable TLS and assume that GitHub creates a new connection each time it delivers a new event **total latency is 3x numbers observed during our load tests**. 
- with **.5 vCPU**, that would mean **~120req/sec**
- with **1 vPCU**, **~300req/sec** 
- with **2 vPCUs**, **~500req/sec**
- latency will heavily depend on the network roundtrip between GitHub servers and our server. Looks like GitHub webhook request originate from San Francisco. This is a 150ms roundtrip to Paris. That could justify deploying our webhook server in the us-west-1 AWS region, closer to GitHub's servers.

# Pricing (in eu-west-3 AWS region = Paris)
- ECS Fargate: 0.5vCPU + 1GB = 21$/mo
- ECS Fargate: 1vCPU + 2GB = 42$/mo
- ECS Fargate: 2vCPU + 4GB = 84$/mo

Comparison with EC2 instance (on demand):
- **t4g.micro: 2vCPU + 1GB = 7$/mo**
- t4g.medium: 2vCPU + 4GB = 27$/mo

Notes:
- EC2 is way cheaper and offers more options in terms of vPCU and memory.
- ECS Fargate is 3x more expensive than EC2 instances.

(new) **ECS Managed Instance** run ECS containers on EC2 instances with EC2 pricing +10-15%.

# Tests details
## Task resources: .5 vCPU | 1 GB
### 50 connections, 2 threads, 1 minute
**CPU utilization = 100%**, memory utilization = 2%

'''
wrk "http://35.181.56.38:8000/event" -s ./test/wrk_post_binary.lua --latency -t 2 -c 50 -d 1m --timeout 10
Running 1m test @ http://35.181.56.38:8000/event
  2 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   138.62ms   56.22ms 599.26ms   78.98%
    Req/Sec   181.96     64.42   360.00     64.86%
  Latency Distribution
     50%  111.10ms
     75%  188.73ms
     90%  202.04ms
     99%  297.33ms
  21746 requests in 1.00m, 3.13MB read
Requests/sec:    362.06
Transfer/sec:     53.39KB
'''

### 25 connections, 2 threads, 1 minute
**CPU utilization = 100%**, memory utilization = 2%

'''
wrk "http://35.181.56.38:8000/event" -s ./test/wrk_post_binary.lua --latency -t 2 -c 25 -d 1m --timeout 10
Running 1m test @ http://35.181.56.38:8000/event
  2 threads and 25 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    65.34ms   33.10ms 204.14ms   58.95%
    Req/Sec   184.97     39.15   290.00     64.92%
  Latency Distribution
     50%   80.82ms
     75%   90.65ms
     90%   96.49ms
     99%  122.65ms
  22133 requests in 1.00m, 3.19MB read
Requests/sec:    368.49
Transfer/sec:     54.34KB
'''

### 10 connections, 2 threads, 1 minute
**CPU utilization = 100%**, memory utilization = 2%

'''
wrk "http://35.181.56.38:8000/event" -s ./test/wrk_post_binary.lua --latency -t 2 -c 10 -d 1m --timeout 10
Running 1m test @ http://35.181.56.38:8000/event
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    23.20ms    6.34ms 134.64ms   87.76%
    Req/Sec   217.70     24.83   260.00     72.50%
  Latency Distribution
     50%   20.86ms
     75%   22.97ms
     90%   31.28ms
     99%   48.23ms
  26030 requests in 1.00m, 3.75MB read
Requests/sec:    433.46
Transfer/sec:     63.92KB
'''

### 5 connections, 2 threads, 1 minute
**CPU utilization = 53%**, memory utilization = 2%

'''
rk "http://35.181.56.38:8000/event" -s ./test/wrk_post_binary.lua --latency -t 2 -c 5 -d 1m --timeout 10
Running 1m test @ http://35.181.56.38:8000/event
  2 threads and 5 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    21.34ms    5.32ms 196.23ms   95.08%
    Req/Sec    94.58     10.70   121.00     82.00%
  Latency Distribution
     50%   20.35ms
     75%   21.13ms
     90%   22.78ms
     99%   41.88ms
  11318 requests in 1.00m, 1.63MB read
Requests/sec:    188.47
Transfer/sec:     27.79KB
'''

## Task resources: 1 vCPU | 2 GB
### 10 connections, 2 threads, 1 minute
**CPU utilization = 48%**, memory utilization = <1%

'''
wrk "http://35.180.132.212:8000/event" -s ./test/wrk_post_binary.lua --latency -t 2 -c 10 -d 1m --timeout 10
Running 1m test @ http://35.180.132.212:8000/event
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    20.59ms    3.06ms  85.01ms   96.84%
    Req/Sec   244.42     14.52   282.00     83.83%
  Latency Distribution
     50%   20.14ms
     75%   20.84ms
     90%   21.62ms
     99%   36.20ms
  29231 requests in 1.00m, 4.21MB read
Requests/sec:    486.77
Transfer/sec:     71.78KB
'''

### 20 connections, 2 threads, 1 minute
**CPU utilization = 98%**, memory utilization = <1%

'''
wrk "http://35.180.132.212:8000/event" -s ./test/wrk_post_binary.lua --latency -t 2 -c 20 -d 1m --timeout 10
Running 1m test @ http://35.180.132.212:8000/event
  2 threads and 20 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    22.07ms    5.02ms  97.00ms   90.82%
    Req/Sec   456.65     44.30   535.00     70.75%
  Latency Distribution
     50%   20.51ms
     75%   22.36ms
     90%   26.45ms
     99%   44.12ms
  54577 requests in 1.00m, 7.86MB read
Requests/sec:    908.96
Transfer/sec:    134.04KB
'''

### 40 connections, 2 threads, 1 minute
**CPU utilization = 100%**, memory utilization = <1%

'''
wrk "http://35.180.132.212:8000/event" -s ./test/wrk_post_binary.lua --latency -t 2 -c 40 -d 1m --timeout 10
Running 1m test @ http://35.180.132.212:8000/event
  2 threads and 40 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    53.39ms   26.61ms 313.43ms   61.19%
    Req/Sec   378.72     72.37   590.00     72.42%
  Latency Distribution
     50%   58.08ms
     75%   74.02ms
     90%   85.22ms
     99%  109.90ms
  45277 requests in 1.00m, 6.52MB read
Requests/sec:    753.83
Transfer/sec:    111.16KB
'''

## Task resources: 2 vCPU | 4 GB
### Baseline: 1 connections, 1 threads, 1 minute
**CPU utilization = 3%**, memory utilization = <1%

'''
wrk "http://15.237.211.43:8000/event" -s ./test/wrk_post_binary.lua --latency -t 1 -c 1 -d 1m --timeout 10
Running 1m test @ http://15.237.211.43:8000/event
  1 threads and 1 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    20.97ms    3.37ms  64.76ms   97.06%
    Req/Sec    47.88      4.97    60.00     80.83%
  Latency Distribution
     50%   20.48ms
     75%   20.94ms
     90%   21.33ms
     99%   37.49ms
  2873 requests in 1.00m, 423.66KB read
Requests/sec:     47.84
Transfer/sec:      7.05KB
'''

### 30 connections, 2 threads, 1 minute
**CPU utilization = 75%**, memory utilization = <1%

'''
wrk "http://15.237.211.43:8000/event" -s ./test/wrk_post_binary.lua --latency -t 2 -c 30 -d 1m --timeout 10
Running 1m test @ http://15.237.211.43:8000/event
  2 threads and 30 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    23.38ms    5.03ms 240.51ms   91.71%
    Req/Sec   646.71     55.23   747.00     74.00%
  Latency Distribution
     50%   22.18ms
     75%   24.26ms
     90%   27.63ms
     99%   40.81ms
  77280 requests in 1.00m, 11.13MB read
Requests/sec:   1287.11
Transfer/sec:    189.80KB
'''

### 40 connections, 2 threads, 1 minute
**CPU utilization = 88%**, memory utilization = <1%

'''
wrk "http://15.237.211.43:8000/event" -s ./test/wrk_post_binary.lua --latency -t 2 -c 40 -d 1m --timeout 10
Running 1m test @ http://15.237.211.43:8000/event
  2 threads and 40 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    30.23ms    9.47ms 259.04ms   82.71%
    Req/Sec   668.39     90.88     0.89k    65.94%
  Latency Distribution
     50%   27.58ms
     75%   33.83ms
     90%   42.20ms
     99%   62.66ms
  79809 requests in 1.00m, 11.49MB read
Requests/sec:   1328.95
Transfer/sec:    195.97KB
'''

### 50 connections, 2 threads, 1 minute
**CPU utilization = 91%**, memory utilization = <1%

'''
wrk "http://15.237.211.43:8000/event" -s ./test/wrk_post_binary.lua --latency -t 2 -c 50 -d 1m --timeout 10
Running 1m test @ http://15.237.211.43:8000/event
  2 threads and 50 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    40.58ms   14.27ms 144.23ms   73.33%
    Req/Sec   621.12     99.59     0.91k    67.75%
  Latency Distribution
     50%   37.29ms
     75%   47.73ms
     90%   60.07ms
     99%   86.95ms
  74246 requests in 1.00m, 10.69MB read
Requests/sec:   1236.23
Transfer/sec:    182.30KB
'''