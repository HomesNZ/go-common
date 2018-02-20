package env

func MustGetServiceHost(serviceName string) string {
	return MustGetServiceHostname(serviceName) + ":" + MustGetServicePort(serviceName)
}

func MustGetServiceHostname(serviceName string) string {
	h := GetString(serviceName+"_SERVICE_HOST", "")
	if h != "" {
		return h
	}

	return MustGetString(serviceName + "_HOST")
}

func MustGetServicePort(serviceName string) string {
	h := GetString(serviceName+"_SERVICE_PORT", "")
	if h != "" {
		return h
	}

	return MustGetString(serviceName + "_HOST")
}
