# Kubebrainz

Kubebrainz is a self loading / self updating grpc bridge to the
MusicBrainz database.

## Running

On first run, kubebrainz will download and load the musicbrainz database
into kubernetes with a running postgres instance, It will then serve api requests
that are backed by the database. KB checks every 24 hours for a new instance of
the database and then refreshes with that version.

## Usage

You can use the GetStatus method to see the state of KB - i.e. which database
version it's running and how much it's downloaded of that database,
and what state of the database is (i.e. how much of the version is downloaded
and whether we're in the loading stage).

Database loads are flipped (i.e. we load version A whilst version B is being
served and then flip to version B once it's fully loaded), meaning that we
should always be able to field an API request.

## Metrics

KB exports prometheus metrics (on port 8081 by default). Make use of the supplied
dashboard for showing these metrics in prometheus, or create your own.
