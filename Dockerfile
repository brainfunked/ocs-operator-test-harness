FROM registry.svc.ci.openshift.org/openshift/release:golang-1.13 AS builder

ENV PKG=/go/src/github.com/brainfunked/ocs-operator-test-harness/
WORKDIR ${PKG}

# compile test binary
COPY . .
RUN make

FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

COPY --from=builder /go/src/github.com/brainfunked/ocs-operator-test-harness/ocs-operator-test-harness.test ocs-operator-test-harness.test

ENTRYPOINT [ "/ocs-operator-test-harness.test" ]

