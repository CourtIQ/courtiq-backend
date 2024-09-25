FROM node:18-alpine

WORKDIR /app

COPY package*.json ./

# Install latest dependencies and health check tools
RUN npm install && npm update && \
    apk add --no-cache curl

COPY . .

EXPOSE 4000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:4000/health || exit 1

CMD ["npm", "start"]