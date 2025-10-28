# Download

Once we've established that there is a fresh version of the MBDB available
we trigger a download, wait for that download to finish and then load the
database into the unused slot. Thus the state of the database goes through
three stages:

1. Downloading
   GetStatus will return the percentage of the db file downloaded
1. Loading
   GetStatus will report how long the load has been running for
1. Serving
   GetStatus will report which slot the database is being served from

Since the MB DB is quite large (~6Gb at the time of writing), we need to
have sufficient tmp directory storage to be able to download it - this should be
provided by the containerisation used to server KubeBrainz. At startup KB will
check that it has at least 10Gb free for scratch and fail to load if this is
not available.

On the first pass the API won't be active until it's downloaded and loaded
into a slot, otherwise the downloads and slot loads are transparent.
