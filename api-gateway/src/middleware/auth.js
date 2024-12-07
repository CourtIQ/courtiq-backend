const AuthService = require('../services/authService');
const TokenVerifier = require('../services/tokenVerifier');

const tokenVerifier = new TokenVerifier();
const authService = new AuthService(tokenVerifier);

async function authMiddleware({ req }) {
  const user = await authService.authenticate(req.headers);
  return { user };
}

module.exports = authMiddleware;