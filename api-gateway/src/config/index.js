require('dotenv').config();

const ENV = process.env.NODE_ENV || 'development';
const PORT = parseInt(process.env.PORT, 10) || 3000;
const GRAPHQL_PLAYGROUND = process.env.GRAPHQL_PLAYGROUND === 'true';

const SERVICES = {
  userService: {
    name: 'user-service',
    envUrl: process.env.USER_SERVICE_URL,
  },
  relationshipService: {
    name: 'relationship-service',
    envUrl: process.env.RELATIONSHIP_SERVICE_URL,
  },
  matchupService: {
    name: 'matchup-service',
    envUrl: process.env.MATCHUP_SERVICE_URL,
  },
  equipmentService: {
    name: 'equipment-service',
    envUrl: process.env.EQUIPMENT_SERVICE_URL,
  },
};

module.exports = {
  ENV,
  PORT,
  GRAPHQL_PLAYGROUND,
  SERVICES,
};
