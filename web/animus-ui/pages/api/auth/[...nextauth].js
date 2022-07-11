import NextAuth from 'next-auth';
import CredentialsProvider from 'next-auth/providers/credentials';
import jwt_decode from 'jwt-decode';
import { LoginCredentials } from '../../../service/http';
import { format } from 'prettier';

export const authOptions = {
  providers: [
    CredentialsProvider({
      // The name to display on the sign in form (e.g. 'Sign in with...')
      name: 'Credentials',
      // The credentials is used to generate a suitable form on the sign in page.
      // You can specify whatever fields you are expecting to be submitted.
      // e.g. domain, username, password, 2FA token, etc.
      // You can pass any HTML attribute to the <input> tag through the object.
      credentials: {
        email: { label: 'Email', type: 'text', placeholder: 'Email' },
        password: { label: 'Password', type: 'password' },
      },
      async authorize(credentials, req) {
        // You need to provide your own logic here that takes the credentials
        // submitted and returns either a object representing a user or value
        // that is false/null if the credentials are invalid.
        // e.g. return { id: 1, name: 'J Smith', email: 'jsmith@example.com' }
        // You can also use the `req` object to obtain additional parameters
        // (i.e., the request IP address)
        try {
          const res = await LoginCredentials(
            credentials.email,
            credentials.password
          );
          // If no error and we have user data, return it
          if (res?.data) {
            res.data.user = credentials.email;
            return res.data;
          }
          // Return null if user data could not be retrieved
          return null;
        } catch (err) {
          console.log('could not authorize LoginCredentials:', err.message);
          throw err;
        }
      },
    }),
  ],
  secret: process.env.JWT_SECRET,
  callbacks: {
    async jwt({ token, user, account }) {
      if (account && user) {
        return {
          ...token,
          accessToken: user.token,
          // refreshToken: user.refresh_token,
        };
      }

      return token;
    },

    async session({ session, token }) {
      const data = jwt_decode(token.accessToken);
      session.user.accessToken = token.accessToken;
      session.user.email = data.email;
      session.user.name = null;
      session.user.image = null;
      return session;
    },
  },
  // Enable debug messages in the console if you are having problems
  debug: process.env.NODE_ENV === 'development',
};

export default NextAuth(authOptions);
