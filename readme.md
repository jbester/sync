Sync - Additional synchronization primtives
===========================================

Go code (golang) set of packages that provide synchonization primitives beyond those provided by the golang sync package.

Primitives include:

-	[Events](#events)
-	[StartGroups](#start-group)
-	[Semaphores](#semaphores)

Check out the API Documentation http://godoc.org/bitbucket.org/jbester/sync

[`events`](http://godoc.org/bitbucket.org/jbester/sync/events "API documentation") package
---------------------------------------------------------------------------------------------

The `events` package provides a single synchronization primitive the Event. An event is used to notify the occurrence of a condition to routines.

Multiple routines can wait on a condition. *All* routines unblock once the condition occurs. A routine that waits on a condition that has already occurred will not block.

The event primitive is similar to the event in the pSOS or ARINC 653 API sets.

[`startgroup`](http://godoc.org/bitbucket.org/jbester/sync/startgroup "API documentation") package
--------------------------------------------------------------------------------------------------

The `startgroup` package provides a mechanism for a collection of goroutines to wait for a release event. When released, all blocked routines simultaneously.

A typical use is when multiple routines need to know when a resource is available but do not need exclusive access to the resource.

[`semaphores`](http://godoc.org/bitbucket.org/jbester/sync/semaphores "API documentation") package
--------------------------------------------------------------------------------------------------

The `semaphores` package provides a go implementation of binary and counting semaphores.  It is designed to use atomic operations to maintain the semaphore count and a channel to signal waiting threads.


Installation
============

To install, use `go get`:

```
go get bitbucket.org/jbester/sync/...
```

This will then make the following packages available to you:

```
bitbucket.org/jbester/sync/semaphores
bitbucket.org/jbester/sync/events
bitbucket.org/jbester/sync/startgroup
```

---

Staying up to date
==================

To update to the latest version, use `go get -u bitbucket.org/jbester/sync`.

---

Contributing
============

Please feel free to submit issues, fork the repository and send pull requests!

---

Licence
=======

MIT License - see LICENSE file 
