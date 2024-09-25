# Use the official Node.js image as the base image
FROM node:18-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the package.json and package-lock.json (if available)
COPY package*.json ./

# Install the dependencies
RUN npm install --production

# Copy the rest of your application code
COPY . .

# Expose the port that will be defined by the environment variable
EXPOSE 4000

# Start the API Gateway
CMD ["npm", "start"]
