package test

import "github.com/rh-messaging/shipshape/pkg/framework/operators"

type OperatorDeploymentWrapper struct {
	EnvVariables map[string]string
}

func (odw *OperatorDeploymentWrapper) AddEnvVar(name string, value string) *OperatorDeploymentWrapper {
	odw.EnvVariables[name] = value
	return odw
}

func (odw *OperatorDeploymentWrapper) PrepareOperator() operators.OperatorSetupBuilder {
	builder := operators.SupportedOperators[operators.OperatorTypeBroker]
	//Set image to parameter if one is supplied, otherwise use default from shipshape.
	if len(Config.OperatorImageName) != 0 {
		builder.WithImage(Config.OperatorImageName)
	}
	if Config.DownstreamBuild {
		builder.WithCommand("/home/amq-broker-operator/bin/entrypoint")
		builder.WithOperatorName("amq-broker-operator")
	}
	if Config.RepositoryPath != "" {
		// Try loading YAMLs from the repo.
		yamls, err := LoadYamls(Config.RepositoryPath)
		if err != nil {
			panic(err)
		} else {
			builder.WithYamls(yamls)
		}
	}
	if Config.AdminUnavailable {
		builder.SetAdminUnavailable()
	}

/*	if len(odw.EnvVariables) > 0 {
		for name, value := range odw.EnvVariables {
                builder.AddEnvVariable(name, value)
		}
	} */
	return builder
}
