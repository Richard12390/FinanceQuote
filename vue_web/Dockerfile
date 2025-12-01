# Build
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
ARG VITE_BACKEND_ORIGIN
ENV VITE_BACKEND_ORIGIN=${VITE_BACKEND_ORIGIN}
RUN npm run build

# Runtime
FROM nginx:1.27-alpine
WORKDIR /usr/share/nginx/html
COPY --from=builder /app/dist .
COPY ./nginx.conf.docker /etc/nginx/nginx.conf
EXPOSE 80