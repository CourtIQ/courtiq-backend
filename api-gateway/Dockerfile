# api-gateway/Dockerfile
FROM node:18-alpine

# Set working directory
WORKDIR /app

# Copy package files first for better caching
COPY package*.json ./
RUN npm install

# Copy source files
COPY . .

# Create directory for logs
RUN mkdir -p logs && chmod -R 777 logs

# Set environment variables
ENV NODE_ENV=production \
    SERVICE_NAME=api-gateway

# Expose the API Gateway port
EXPOSE 3000

# Set health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
  CMD wget -q -O - http://localhost:3000/health || exit 1

# Start the service with proper signal handling
CMD ["node", "src/index.js"]