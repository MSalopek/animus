import { ActivateUser } from '../../../../service/http';

export default async function activateUser(req, res) {
  const { token, email } = req.query;
  if (!token && !email) {
    return res.status(401).json({ message: 'invalid request parameters' });
  }

  const response = await ActivateUser(email, token);
  if (response.status >= 400) {
	// go to error page
    return res.status(401).json({ message: 'could not activate user' });
  } else {
    res.writeHead(307, { Location: '/account/login?verified=true' });
    res.end();
  }
}
