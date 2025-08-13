# order-frontend/Dockerfile
FROM nginx:1.29-alpine

# Copy built files
COPY public /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]