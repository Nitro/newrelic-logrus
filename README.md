New Relic Logrus Hook
=====================

This is a simple logrus hook that lets existing logrus applications hook into
New Relic error reporting using the New Relic Go Agent.

New Relic relies on having errors associated witha transaction. But logrus
doesn't know which transaction its being called from. So currently this hook
plugin just creates a new transaction called `errorTxn` and reports logged
errors under that transaction. They show up just like any other reported
errors. If any fields are supplied to the log line via `WithFields`, they are
reported as custom attributes on the `errorTxn` and will be visible in New
Relic.
