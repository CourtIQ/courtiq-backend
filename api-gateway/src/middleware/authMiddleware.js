const AuthService = require('../services/authService');

const tokenVerifier = require('../services/tokenVerifier'); // or `new TokenVerifier()`
const authService = new AuthService(tokenVerifier);

async function authMiddleware({ req }) {
  const user = await authService.authenticate(req.headers);
  return { user };
}

module.exports = authMiddleware;
