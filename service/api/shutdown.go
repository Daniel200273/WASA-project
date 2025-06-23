package api

// Close should close everything opened in the lifecycle of the `_router`; for example, background goroutines.
func (rt *_router) Close() error {
	// Clean up the temporary uploads directory on shutdown
	if err := cleanupUploadsDirectory(); err != nil {
		rt.baseLogger.WithError(err).Warning("error cleaning up uploads directory during shutdown")
		// Don't return error as this is cleanup, not critical
	} else {
		rt.baseLogger.Info("uploads directory cleaned up successfully")
	}

	return nil
}
