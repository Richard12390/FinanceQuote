FROM node:20-alpine AS build
WORKDIR /app
COPY ../vue_web/package*.json ./
RUN npm ci
COPY ../vue_web ./
RUN npm run build

FROM nginx:1.27-alpine
COPY deploy/nginx/nginx.conf /etc/nginx/nginx.conf        
COPY --from=build /app/dist /usr/share/nginx/html
