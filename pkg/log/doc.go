/*
Package log provides a unified way to log messages across LoungeUp applications.

The default [Logger] is pre-configured and ready to use.

	// Log a debug message with the default logger.
	log.Default().Debug("Initializing application", slog.String("version", "0.0.1"))

	if err := startServer(); err != nil {
		// Log an error message with the default logger.
		log. Default().Error("Could not start server",
			slog.String("error", err.Error()),
			slog.String("version", "0.0.1"),
		)
	}

You can automatically send log messages to Datadog.
Use the pre-defined methods to add the 'formattedMessage' log attribute.

	// Log a debug message with the default logger and automatically add a 'formattedMessage' attribute.
	log.Default().FormattedDebug("Initializing application", slog.String("version", "0.0.1"))
	log.Default().FormattedError("Could not initialize application", slog.String("version", "0.0.1"))

You can use the [Adapter] to pass our pre-configured logger to external libraries.

	// Use the adapter of the default logger with the 'go-res' library.
	res.NewService("example").SetLogger(log.Default().Adapter)

	// Use the adapter of the default logger with the 'badger' library.
	badger.Open(badger.DefaultOptions("/tmp/").WithLogger(log.Default().Adapter))
*/
package log
