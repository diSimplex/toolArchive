# The Discovery RESTful interface

The Discovery RESTful interface will be responsible for

  - remote Nursery discovery

  - periodic heartbeats

  - Nursery load averages

All three responsiblities will be managed using the same messages.

The discovery messages will be periodic HTTPS POST JSON messages sent at 
short *random* intervals.

The discovery messages will be sent to the configured, well-known, address 
and port of the Primary Nursery.

The contents of the discovery messages will be:

  1. the Nursery's host name

  2. the Nursery's RESTful HTTP port

  3. the current up/paused state of the Nursery

  4. the current number of ConTeXt processes running (including those who
     are simply waiting)

  5. the relative speed of the Nursery's host machine

  6. a sequence of load averages of the Nursery's host machine, over a
     number of intervals (1, 5, and 15 minutes for linux)

The contents of the discovery message's reply will be a JSON message 
containing a list of all currently known Nurseries' discovery messages.

## Principles

- The Primary Nursery is not configured with any of the secondary 
  Nurseries.

- The Secondary Nurseries are configured with the Primary Nursery's address 
  and discovery port.

## Questions

