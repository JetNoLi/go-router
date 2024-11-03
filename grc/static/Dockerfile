FROM golang:1.22 as builder

ARG TEMPL_VERSION = "latest"
ARG GO_ROUTER_VERSION = "latest"
ARG PORT = 3000

WORKDIR /app

COPY .  .

COPY ./.env /app/

#TODO: Allow Templ Version Override
RUN go install github.com/a-h/templ/cmd/templ@${TEMPL_VERSION}
RUN templ generate

# Copy and filter CSS files from view/
# This finds files and then copies them with parents, i.e. retaining directories
RUN mkdir -p build/view/pages build/view/components
RUN find ./view/pages -name '*.css' -exec cp --parents {} build \;
RUN find ./view/components -name '*.css' -exec cp --parents {} build \;

#TODO: Allow way to pass in executable name
RUN go build -o app.exe

#TODO: Allow override to specify go version
FROM golang:1.22 as runner

WORKDIR /app

COPY --from=builder app/$app.exe /app
COPY --from=builder app/.env /app
COPY --from=builder app/assets /app/assets
COPY --from=builder app/build/view/ /app/view

EXPOSE $PORT 

#TODO: Use Executable Name Here
CMD ["./app.exe"]