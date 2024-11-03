FROM golang:1.22 AS builder

ARG TEMPL_VERSION "latest"
ARG GO_ROUTER_VERSION "latest"
ARG PORT "3000"

WORKDIR /app

COPY .  .
COPY ./.env /app/

RUN go install github.com/a-h/templ/cmd/templ@${TEMPL_VERSION}
RUN go install github.com/jetnoli/go-router/grc@${GO_ROUTER_VERSION}
RUN templ generate

RUN grc generate-assets

# Copy and filter CSS files from view/
# This finds files and then copies them with parents, i.e. retaining directories
RUN mkdir -p build/view/pages build/view/components
RUN find ./view/pages -name '*.css' -exec cp --parents {} build \;
RUN find ./view/components -name '*.css' -exec cp --parents {} build \;

RUN go build -o app.exe

#TODO: Allow override to specify go version
FROM golang:1.22 AS runner

WORKDIR /app

COPY --from=builder app/asset_map.json /app
COPY --from=builder app/app.exe /app
COPY --from=builder app/.env /app
COPY --from=builder app/assets /app/assets
COPY --from=builder app/build/view/ /app/view

EXPOSE $PORT 

CMD ["./app.exe"]