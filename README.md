# go-lockfree-queue
Lock free queue in golang (implementing ...)

The project look just like [yireyun/go-queue](https://github.com/yireyun/go-queue), 
but with different Queue interface. I develop this package after reading through yireyun's code.

Difference with `yireyun/go-queue`:
* Put do not run `runtime.Gosched()` if Queue is full, but let the user decide what to do.
* Put do not return `false` if `CAS` fail, but will retry after `runtime.Gosched()` until success.
