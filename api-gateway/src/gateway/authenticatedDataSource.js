const { RemoteGraphQLDataSource } = require('@apollo/gateway');

class AuthenticatedDataSource extends RemoteGraphQLDataSource {
  willSendRequest({ request, context }) {
    if (context.user) {
      // Convert the user object to a JSON string
      const userJson = JSON.stringify(context.user);
      // Optionally base64-encode to avoid special characters in headers
      const encoded = Buffer.from(userJson).toString('base64');
      request.http.headers.set('x-user', encoded);
    }
  }
}



module.exports = AuthenticatedDataSource;