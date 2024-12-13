// const AuthService = require('../services/authService');
// const TokenVerifier = require('../utils/tokenVerifier');

// const tokenVerifier = new TokenVerifier();

// const authService = new AuthService(tokenVerifier);

// async function authMiddleware({ req }) {
//   try {
//     const user = await authService.authenticate(req.headers);
//     return { user }; // Attach user to the context
//   } catch (err) {
//     console.error('Authentication failed:', err.message);
//     return { user: null }; // Proceed as unauthenticated
//   }
// }
// module.exports = authMiddleware;