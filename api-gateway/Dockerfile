FROM node:18-alpine

WORKDIR /app

COPY package*.json ./

# Install latest dependencies and health check tools
RUN npm install && npm update && \
    apk add --no-cache curl

# Copy everything except node_modules
COPY . .

# Ensure the service account file has correct permissions if it exists
RUN chmod 600 firebase-service-account.json || true

EXPOSE 4000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:4000/health || exit 1

CMD ["npm", "start"]