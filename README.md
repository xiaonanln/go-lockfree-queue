# go-lockfree-queue
Lock free queue in golang (multiple producers, multiple consumers)

Since I develop this package after reading through yireyun's code, The project look just like [yireyun/go-queue](https://github.com/yireyun/go-queue) and I also used some code from yireyun's project.

Difference with `yireyun/go-queue`:
* Queue interfaces are different.
* Put/Get does not run `runtime.Gosched()` if Queue is full/empty, but let the user decide what to do.
* Put/Get does not return `false` if `CAS` fail, but will retry after `runtime.Gosched()` until success.
* Make it faster for a little bit by eliminating [false sharing](https://en.wikipedia.org/wiki/False_sharing).
