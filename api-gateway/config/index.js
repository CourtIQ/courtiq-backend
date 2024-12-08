require('dotenv').config();

const ENV = process.env.NODE_ENV || 'development';
const PORT = parseInt(process.env.PORT, 10) || 3000;
const GRAPHQL_PLAYGROUND = process.env.GRAPHQL_PLAYGROUND === 'true';

const MONGODB_URL = process.env.MONGODB_URL;
const FIREBASE_CONFIG = JSON.parse(process.env.FIREBASE_CONFIG || '{}');
const FIREBASE_SERVICE_ACCOUNT = JSON.parse(process.env.FIREBASE_SERVICE_ACCOUNT || '{}');

const SERVICES = {
  userService: {
    name: 'user-service',
    url: process.env.USER_SERVICE_URL || 'http://localhost:8080',
  },
  relationshipService: {
    name: 'relationship-service',
    url: process.env.RELATIONSHIP_SERVICE_URL || 'http://localhost:8080',
  },
  matchupService: {
    name: 'matchup-service',
    url: process.env.MATCHUP_SERVICE_URL || 'http://localhost:8080',
  },
  equipmentService: {
    name: 'equipment-service',
    url: process.env.EQUIPMENT_SERVICE_URL || 'http://localhost:8080',
  },
};

module.exports = {
  ENV,
  PORT,
  GRAPHQL_PLAYGROUND,
  MONGODB_URL,
  FIREBASE_CONFIG,
  FIREBASE_SERVICE_ACCOUNT,
  SERVICES,
};
