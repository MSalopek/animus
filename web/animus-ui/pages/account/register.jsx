import Head from 'next/head';
import Link from 'next/link';

import { unstable_getServerSession } from 'next-auth/next';
import { authOptions } from '../api/auth/[...nextauth]';

import { RegisterUser } from '../../service/http';
import { useState } from 'react';
import { Alert } from '../../components/alert/Alert';

export default Register;

function Register() {
  const [email, setEmail] = useState('');
  const [emailValid, setEmailValid] = useState(false);

  const [password, setPassword] = useState('');
  const [confirm, setConfirm] = useState('');
  const [passwordValid, setPasswordValid] = useState(false);

  const [firstname, setFirstname] = useState('');
  const [lastname, setLastname] = useState('');
  const [username, setUsername] = useState('');

  const [passMatch, setPassMatch] = useState(false);

  const [success, setSuccess] = useState(false);
  const [err, setErr] = useState('');

  const isDisabled = () => {
    return !email || !passwordValid || !passMatch || !username;
  };
  const submit = async () => {
    try {
      const p = {
        email,
        password,
        username,
        firstname,
        lastname,
      };
      const res = await RegisterUser(p);

      if (!res.error) {
        setSuccess(true);
      }

      handleErr(res.error);
    } catch (error) {
      setErr(
        'Internal error happened. Please try again later or contact support.'
      );
    }
  };

  const clearErrMessage = () => {
    setErr(null);
  };

  return (
    <>
      <Head>
        <title>Animus Store | Register</title>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta
          name="description"
          content="Decentralized storage made easy. Manage files, seamlessly access IPFS and regain control over your private data."
        />
      </Head>
      <div className="grid justify-items-center items-center h-screen">
        {success ? (
          <Success email={email} />
        ) : (
          <div className="w-full max-w-lg mx-auto border overflow-hidden bg-white rounded-lg shadow-md dark:bg-gray-800">
            <div className="px-6 py-4">
              <h2 className="text-3xl font-bold text-center text-gray-700 dark:text-white">
                Animus Store
              </h2>

              <h2 className="text-lg font-semibold text-gray-700 capitalize text-center my-6 dark:text-white">
                Register your account
              </h2>

            {err && <Alert message={err} onClick={clearErrMessage} />}
              <form>
                <div className="mt-6">
                  <label
                    className="text-gray-700 dark:text-gray-200"
                    htmlFor="username"
                  >
                    Username
                  </label>
                  <input
                    id="username"
                    type="text"
                    placeholder="Enter username (required)"
                    className="block w-full px-4 py-2 mt-2 text-gray-700 bg-white border border-gray-200 rounded-md dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 focus:ring-blue-300 focus:ring-opacity-40 dark:focus:border-blue-300 focus:outline-none focus:ring"
                    onChange={(e) => {
                      setUsername(e.target.value);
                    }}
                  />
                </div>

                <div className="mt-6">
                  <label
                    className="text-gray-700 dark:text-gray-200"
                    htmlFor="emailAddress"
                  >
                    Email Address
                  </label>
                  <input
                    id="emailAddress"
                    type="email"
                    placeholder="Enter email (required)"
                    className="block w-full px-4 py-2 mt-2 text-gray-700 bg-white border border-gray-200 rounded-md dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 focus:ring-blue-300 focus:ring-opacity-40 dark:focus:border-blue-300 focus:outline-none focus:ring"
                    onChange={(e) => {
                      setEmail(e.target.value);
                      setEmailValid(isEmailValid(email));
                    }}
                  />
                  {email && !emailValid && (
                    <p className="text-sm text-red-600 py-2 px-1">
                      Invalid email address.
                    </p>
                  )}
                </div>

                <div className="mt-6">
                  <label
                    className="text-gray-700 dark:text-gray-200"
                    htmlFor="firstname"
                  >
                    Firstname
                  </label>
                  <input
                    id="firstname"
                    type="text"
                    placeholder="Enter firstname (optional)"
                    className="block w-full px-4 py-2 mt-2 text-gray-700 bg-white border border-gray-200 rounded-md dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 focus:ring-blue-300 focus:ring-opacity-40 dark:focus:border-blue-300 focus:outline-none focus:ring"
                    onChange={(e) => {
                      setFirstname(e.target.value);
                    }}
                  />
                </div>

                <div className="mt-6">
                  <label
                    className="text-gray-700 dark:text-gray-200"
                    htmlFor="lastname"
                  >
                    Lastname
                  </label>
                  <input
                    id="lastname"
                    type="text"
                    placeholder="Enter lastname (optional)"
                    className="block w-full px-4 py-2 mt-2 text-gray-700 bg-white border border-gray-200 rounded-md dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 focus:ring-blue-300 focus:ring-opacity-40 dark:focus:border-blue-300 focus:outline-none focus:ring"
                    onChange={(e) => {
                      setLastname(e.target.value);
                    }}
                  />
                </div>

                <div className="mt-6">
                  <label
                    className="text-gray-700 dark:text-gray-200"
                    htmlFor="password"
                  >
                    Password
                  </label>
                  <input
                    id="password"
                    type="password"
                    placeholder="Enter password (required)"
                    className="block w-full px-4 py-2 mt-2 text-gray-700 bg-white border border-gray-200 rounded-md dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 focus:ring-blue-300 focus:ring-opacity-40 dark:focus:border-blue-300 focus:outline-none focus:ring"
                    onChange={(e) => {
                      setPassword(e.target.value);
                      setPasswordValid(password.length >= 8);
                    }}
                  />
                </div>

                <div className="mt-6">
                  <label
                    className="text-gray-700 dark:text-gray-200"
                    htmlFor="passwordConfirmation"
                  >
                    Confirm Password
                  </label>
                  <input
                    id="passwordConfirmation"
                    type="password"
                    className="block w-full px-4 py-2 mt-2 text-gray-700 bg-white border border-gray-200 rounded-md dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 focus:ring-blue-300 focus:ring-opacity-40 dark:focus:border-blue-300 focus:outline-none focus:ring"
                    onChange={(e) => {
                      setConfirm(e.target.value);
                      setPassMatch(e.target.value === password);
                    }}
                  />
                  {password && !passMatch && (
                    <p className="text-sm text-red-600 py-2 px-1">
                      Passwords do not match.
                    </p>
                  )}
                </div>

                <div className="flex justify-center mt-8">
                  <button
                    className={`px-6 py-2 leading-5 text-white transition-colors duration-200 transform ${
                      isDisabled()
                        ? 'bg-gray-200'
                        : 'bg-gray-700 hover:bg-gray-600'
                    } rounded-md focus:outline-none focus:bg-gray-600`}
                    disabled={isDisabled()}
                    onClick={(e) => {
                      e.preventDefault();
                      submit();
                    }}
                  >
                    Register
                  </button>
                </div>
              </form>
            </div>

            <div className="flex items-center justify-center py-4 text-center bg-gray-50 dark:bg-gray-700">
              <span className="text-sm text-gray-600 dark:text-gray-200">
                Already have an account?{' '}
              </span>

              <Link href="/account/login">
                <a className="mx-2 text-sm font-bold text-blue-500 dark:text-blue-400 hover:underline">
                  Login
                </a>
              </Link>
            </div>
          </div>
        )}
      </div>
    </>
  );
}

function Success({ email }) {
  return (
    <div className="w-full max-w-lg h-96 mx-auto overflow-hidden border bg-white dark:bg-gray-800 rounded-lg shadow-lg">
      <div class="flex p-8 items-center gap-4">
        <svg
          class="w-9 h-9 fill-current text-green-500"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
        >
          <path d="M0 0h24v24H0V0z" fill="none" />
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm4.59-12.42L10 14.17l-2.59-2.58L6 13l4 4 8-8z" />
        </svg>
        <h2 class="font-semibold text-2xl text-gray-800">Account Created</h2>
      </div>
      <p class="px-8 text-lg text-gray-600 leading-relaxed">
        You need to activate your account before you can proceed.
        <br></br>
        <br></br>
        Please check your inbox at{' '}
        <span className="font-semibold text-gray-800"> {email} </span>(including
        spam folder) and click on the provided link.
        <br></br>
        <br></br>
        You will receive your activation code shortly.
      </p>
    </div>
  );
}

function isEmailValid(email) {
  return /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/.test(
    email
  );
}

export async function getServerSideProps(context) {
  const session = await unstable_getServerSession(
    context.req,
    context.res,
    authOptions
  );

  if (session) {
    return {
      redirect: {
        destination: '/',
        permanent: false,
      },
    };
  }

  // const resSpiritus = await GetSpiritusBySlug(slug);
  // const spiritus = resSpiritus.data;

  // const stories = resStories.data?.content;

  return {
    props: {
      session,
    },
  };
}
