# Events Aggregator

A tool we use at CookiesHQ to start our [monthly events blog posts](http://cookieshq.co.uk/posts/category/events/), written in Go. It gathers all the tech events in and around Bristol, UK for the next calendar month from the [Meetup](http://www.meetup.com) API.

It was written as part of [CookiesLab](http://cookieshq.co.uk/posts/introducing-the-cookieslab-or-why-do-we-book-time-off-for-our-team-members/) as a way to experiment with Go.

## Installation

If you haven't already, [install Go](https://golang.org/doc/install) - you can also install from Homebrew (`brew install golang`) and do the necessary setup (generally: set `$GOPATH` in your `.(bash|zsh)rc`, then create `bin` and `src` directories in said `$GOPATH`).

1. Fetch the code:

  ```sh
  $ go get github.com/cookieshq/events_aggregator
  ```

2. Set `MEETUP_API_KEY` ENV var with your [Meetup.com API key](https://secure.meetup.com/meetup_api/key/): `export MEETUP_API_KEY=abc1234`
3. To run, either:

  ```sh
  $ cd $GOPATH/src/github.com/cookieshq/events_aggregator
  $ go run event_aggregator.go
  ```

  or:

  ```sh
  $ go install github.com/cookieshq/events_aggregator
  $ $GOPATH/bin/event_aggregator
  ```

## Modifying

At the moment, the API parameters are hardcoded in a `map[string]string` in `main()` - alter/add to this if you want to change the params submitted to the API.

To change the presentation of each event, you can edit the `text/template` template, also in `main()`.

## Further information

Check out the Meetup [API docs](http://www.meetup.com/meetup_api/docs/2/open_events/).
