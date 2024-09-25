const services = [
    { name: 'user', url: process.env.USER_SERVICE_URL || 'http://user-service:8081/query' },
    { name: 'string', url: process.env.STRING_SERVICE_URL || 'http://string-service:8082/query' },
    { name: 'matchup', url: process.env.MATCHUP_SERVICE_URL || 'http://matchup-service:8083/query' },
    { name: 'relationship', url: process.env.RELATIONSHIP_SERVICE_URL || 'http://relationship-service:8084/query' },
  ];
  
  module.exports = { services };