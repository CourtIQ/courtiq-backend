const { verifyToken } = require('../utils/firebaseUtils');

class TokenVerifier {
  async verify(idToken) {
    if (!idToken) throw new Error("No token provided");
    return verifyToken(idToken);
  }
}

module.exports = TokenVerifier;
