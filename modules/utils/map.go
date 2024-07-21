package utils

// On error returns index of value which caused error
// If no error int return value == -1
func Map[R any, I any](values []I, fn func(I) (R, error)) ([]R, int, error) {

	results := make([]R, len(values))

	for i, value := range values {
		result, err := fn(value)

		if err != nil {
			return nil, i, err
		}

		results[i] = result
	}

	return results, -1, nil
}
