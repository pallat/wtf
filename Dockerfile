FROM cgr.dev/chainguard/go:latest as build

WORKDIR /app
ARG GIT_COMMIT

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
	-ldflags "-extldflags -static -X main.commit=$GIT_COMMIT" \
	-o api .

FROM cgr.dev/chainguard/static:latest

# USER nonroot

# ENV TINI_VERSION v0.19.0
# ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-static ./tini
# # RUN chmod +x ./tini

COPY --from=build /app/api .

EXPOSE 8080

# ENTRYPOINT ["./tini", "--"]

CMD ["/api"]
