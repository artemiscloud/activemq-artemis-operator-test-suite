run: 
	ginkgo -r -keepGoing ./test/... -- \
		-operator-image registry.redhat.io/amq7/amq-broker-rhel7-operator:latest \
		-broker-image registry.redhat.io/amq7/amq-broker:latest \
		-broker-image-old registry.redhat.io/amq7/amq-broker:7.5-4 \
		-broker-version 7.6.0 -broker-version-old 7.5.0 \
		-broker-name amq-broker -debug-run

build:
	go build ./test/...
