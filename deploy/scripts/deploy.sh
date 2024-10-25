#!/bin/bash

# Exit on any error
set -e

# Default environment
ENV=${1:-development}
PROJECT_ID="court-iq-api"
REGION="us-central1"
VERSION=$(git rev-parse --short HEAD)

echo "🚀 Deploying to $ENV environment..."

# Build and push images
echo "📦 Building and pushing Docker images..."

# API Gateway
echo "Building api-gateway..."
docker build -t $REGION-docker.pkg.dev/$PROJECT_ID/court-iq-repo/api-gateway:$VERSION \
    -f api-gateway/Dockerfile ./api-gateway
docker push $REGION-docker.pkg.dev/$PROJECT_ID/court-iq-repo/api-gateway:$VERSION

# Microservices
for service in user-service string-service matchup-service relationship-service; do
    echo "Building $service..."
    docker build -t $REGION-docker.pkg.dev/$PROJECT_ID/court-iq-repo/$service:$VERSION \
        -f deploy/Dockerfile ./$service \
        --build-arg SERVICE_NAME=$service \
        --build-arg SERVICE_PORT=8080
    docker push $REGION-docker.pkg.dev/$PROJECT_ID/court-iq-repo/$service:$VERSION
done

# Ensure the kubernetes configs directory exists
mkdir -p deploy/k8s/overlays/$ENV

# Apply Kubernetes configurations
echo "🎮 Applying Kubernetes configurations for $ENV environment..."
kubectl apply -k deploy/k8s/overlays/$ENV

# Wait for deployments to be ready
echo "⏳ Waiting for deployments to be ready..."
kubectl wait --for=condition=available deployment --all -n court-iq --timeout=300s || true

echo "✅ Deployment complete to $ENV environment!"

# Show status
echo "📊 Deployment Status:"
kubectl get pods -n court-iq
kubectl get services -n court-iq
kubectl get ingress -n court-iq