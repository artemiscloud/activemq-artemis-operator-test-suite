package bdw

func (bdw *BrokerDeploymentWrapper) WithCPULimit(cpu string) *BrokerDeploymentWrapper {
	bdw.ResourcesLimits.cpu = cpu
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithMemLimit(mem string) *BrokerDeploymentWrapper {
	bdw.ResourcesLimits.mem = mem
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithCPURequest(cpu string) *BrokerDeploymentWrapper {
	bdw.ResourcesRequests.cpu = cpu
	return bdw
}

func (bdw *BrokerDeploymentWrapper) WithMemRequest(mem string) *BrokerDeploymentWrapper {
	bdw.ResourcesRequests.mem = mem
	return bdw
}
