# api-gateway/Dockerfile
FROM node:18-alpine

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install

# Copy source files
COPY . .

# Expose the necessary port
EXPOSE 3000

# Default command
CMD ["npm", "start"]
