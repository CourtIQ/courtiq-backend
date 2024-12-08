class AuthService {
    constructor(tokenVerifier) {
      this.tokenVerifier = tokenVerifier;
    }
  
    async authenticate(headers) {
      const authHeader = headers.authorization || '';
      const token = authHeader.startsWith('Bearer ') ? authHeader.slice(7) : null;
  
      if (!token) {
        return null;  // Allows requests without tokens to proceed as unauthenticated
        // throw new Error('No auth header provided'); // Strict error if no token is present
      }
  
      try {
        const userClaims = await this.tokenVerifier.verify(token);
        return userClaims;
      } catch (error) {
        return null; // Silent fail, treats invalid tokens as unauthenticated
        // throw new Error('Invalid token'); // Strict error if token is invalid
      }
    }
  }
  
  module.exports = AuthService;
  