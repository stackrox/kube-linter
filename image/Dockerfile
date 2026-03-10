FROM scratch

ARG TARGETPLATFORM

COPY $TARGETPLATFORM/kube-linter /

ENTRYPOINT ["/kube-linter"]
