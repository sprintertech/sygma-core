package utils

// TrySendError attempts to send error over the given channel.
// It prevents blocking when a channel is busy or nil.
func TrySendError(errChn chan error, err error) bool {
	select {
	case errChn <- err:
		return true
	default:
		return false
	}
}
