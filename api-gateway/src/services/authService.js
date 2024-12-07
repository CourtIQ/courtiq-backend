class AuthService {
    constructor(tokenVerifier) {
      this.tokenVerifier = tokenVerifier;
    }
  
    async authenticate(headers) {
      const authHeader = headers.authorization || '';
      const token = authHeader.startsWith('Bearer ') ? authHeader.slice(7) : null;
  
      if (!token) {
        // Returning null allows requests without tokens to proceed (unauthenticated)
        return null;
      
        // Throw an error if no token is provided (strict enforcement)
        // throw new Error('No auth header provided');
      }

      try {
        const userClaims = await this.tokenVerifier.verify(token);
        return userClaims; 
      } catch (error) {
        // Verification failed
        return null; // or throw new Error('Invalid token');
      }
    }
  }
  
  module.exports = AuthService;