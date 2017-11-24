Sync - Additional synchronization primtives
================================

Go code (golang) set of packages that provide synchonization primitives beyond those provided by the golang sync package.

Primitives include:

  * [Events](#events)
  * [StartGroups](#start-group)
  * [Semaphores](#semaphores)

Check out the API Documentation http://godoc.org/github.com/jbester/sync

[`events`](http://godoc.org/github.com/jbester/sync/events "API documentation") package
---------------------------------------------------------------------------------------------

The `events` package provides a single synchronization primitive the Event. An event is used to notify the occurrence of a condition to routines.
                     
 Multiple routines can wait on a condition. _All_ routines unblock once the condition occurs. A routine that waits on a condition that has already occurred will not block.
 
 The event primitive is similar to the event in the pSOS or ARINC 653 API sets.

[`startgroup`](http://godoc.org/github.com/jbester/sync/startgroup "API documentation") package
---------------------------------------------------------------------------------------

The `startgroup` package provides a mechanism for a collection of goroutines to wait for a release event.
When released, all blocked routines simultaneously.

A typical use is when multiple routines need to know when a resource is available but do
not need exclusive access to the resource.

[`semaphores`](http://godoc.org/github.com/jbester/sync/semaphores "API documentation") package
-------------------------------------------------------------------------------------------

The `semaphores` package provides a go implementation of binary and counting semaphores.
The underlying implementation is built on the channel primitive.  As such, it doesn't
offer any advantages over using a channel except for readability.


Installation
============

To install, use `go get`:

    * Latest version: go get github.com/jbester/sync

This will then make the following packages available to you:

    github.com/jbester/sync/semaphores
    github.com/jbester/sync/events
    github.com/jbester/sync/startgroup

------

Staying up to date
==================

To update to the latest version, use `go get -u github.com/jbester/sync`.

------


Contributing
============

Please feel free to submit issues, fork the repository and send pull requests!

------

Licence
=======
MIT License - see LICENSE file 
