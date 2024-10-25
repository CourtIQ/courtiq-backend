const services = [
  {
    name: 'users',
    url: process.env.USER_SERVICE_URL || 'http://user-service:8081/query'
  },
  {
    name: 'strings',
    url: process.env.STRING_SERVICE_URL || 'http://string-service:8082/query'
  }
];

module.exports = { services };