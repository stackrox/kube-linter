#############################
# STEP 1 Download Kube-Linter
#############################
FROM golang:alpine AS builder

# Bash and Git are required for go-build.sh
RUN apk add --update --no-cache bash>5.0 git>2.29.0 && \
    rm -rf /tmp/* /var/cache/apk/*

WORKDIR /app
COPY . .

# Execute packr
RUN go get -u github.com/gobuffalo/packr/packr && \
    /go/bin/packr

# Build kube-linter
RUN CGO_ENABLED=0 GOOS=linux scripts/go-build.sh ./cmd/kube-linter

###################################
# STEP 2 Copy the static executable
###################################
FROM scratch

COPY --from=builder /app/bin/linux/kube-linter /usr/local/bin/kube-linter

ENTRYPOINT ["/usr/local/bin/kube-linter"]
