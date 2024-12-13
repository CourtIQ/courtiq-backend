const { verifyToken } = require('./firebaseUtils');

class TokenVerifier {
  async verify(idToken) {
    console.log('Verifying token:', idToken);
    if (!idToken) throw new Error("No token provided");
    return verifyToken(idToken);
  }
}

module.exports = TokenVerifier;
