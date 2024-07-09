# FROM node:alpine as frontend_builder
# WORKDIR /app
# COPY frontend/package.json frontend/yarn.lock ./frontend/
# COPY themes ./themes
# RUN yarn --cwd ./frontend install --frozen-lockfile --non-interactive --production=false
# COPY frontend ./frontend
# RUN yarn --cwd ./frontend build --env output=/app/live/public

FROM golang:alpine as backend_builder
RUN apk add --no-cache git
ARG project
ARG version
# ARG github_token
WORKDIR /app
COPY . .
# RUN git config --global url.https://$github_token@github.com/.insteadOf https://github.com/
# ENV GOPRIVATE="github.com/toyz/portfolio-stats-service"
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.VersionHash=${version}" -o bin ./cmd/${project}/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=backend_builder /app/bin /server
# COPY --from=frontend_builder /app/live /themes/live
# COPY --from=backend_builder /app/themes /themes
# COPY --from=backend_builder /app/themes/flatline /themes/flatline
CMD ["/server"]