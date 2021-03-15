New Relic Logrus Hook
=====================

This is a simple logrus hook that lets existing logrus applications hook into
New Relic error reporting using the New Relic Go Agent.

New Relic relies on having errors associated witha transaction. But logrus
doesn't know which transaction its being called from. So, you have two
options:
 * Pass nothing and a new New Relic transaction will be created and the
   error associatied with that.
 * Use `WithFields(logrus.Fields{"txn": tx}).Error(...)` and pass in
   your existing transaction. This will be used instead.

If you use the second strategy, you should also use the included `Formatter` to
remove the `txn` field from the logs. The `Formatter` can also masque any
fields you wish to prevent from being logged in the clear.

In either case, if any fields are supplied to the log line via `WithFields`,
they are reported as custom attributes on the `errorTxn` and will be visible in
New Relic.

Usage
-----

You can install this like any other logrus hook. Assuming that `application`
is your `newrelic.Application` from the Go agent, you can "hook" it up like
this:

```
log.AddHook(
	newrelic_logrus.NewNewRelicLogrusHook(
		application,
		[]log.Level{log.ErrorLevel, log.FatalLevel},
	),
)
```

The following will enable filtered `txn` from your logs, and will also allow
filtering of any fields entitled `password`. You may specify any underlying
logrus formatter you like. In this case we use the normal `TextFormatter`.

```
formatter := NewFormatter([]string{"password"}, &log.TextFormatter{})
log.SetFormatter(formatter)
```
