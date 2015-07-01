commute-api
===========

Gtfs API used by commute.sh webapp.

[![Build Status](https://travis-ci.org/helyx-io/commute-api.svg?branch=master)](https://travis-ci.org/helyx-io/commute-api)
[![Coverage Status](https://coveralls.io/repos/helyx-io/commute-api/badge.png)](https://coveralls.io/r/helyx-io/commute-api)


Benchmark
=========

Install Vegeta tool:

```
    brew install vegeta
```

Run benchmark:

```
    echo "GET http://localhost:4000/api/agencies/RATP/stops/2015-04-13/nearest?lat=48.875203299999995&lon=2.310958&distance=1000" | vegeta attack -duration=60s -rate=2 | tee results.bin | vegeta report
    vegeta report -inputs=results.bin -reporter=json > metrics.json
    cat results.bin | vegeta report -reporter=plot > plot.html
    cat results.bin | vegeta report -reporter="hist[0,100ms,200ms,300ms,400ms,500ms,600ms,700ms,800ms,900ms,1000ms]"
    open plot.html
```

