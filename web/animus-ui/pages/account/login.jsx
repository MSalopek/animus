import { useState } from 'react';

import Head from 'next/head';
import Link from 'next/link';
import Router from 'next/router';

import { EyeIcon, EyeOffIcon } from '@heroicons/react/outline';

import { useForm } from 'react-hook-form';

import { unstable_getServerSession } from 'next-auth/next';
import { authOptions } from '../api/auth/[...nextauth]';
import { signIn } from 'next-auth/react';
import { AlertSuccess } from '../../components/alert/Alert';

export default Login;

function Login ({ verified }) {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
    getValues,
  } = useForm();

  const [err, setErr] = useState('');
  const [showPass, setShowPass] = useState(false);
  const [showVerifiedBanner, setShowVerifiedBanner] = useState(verified);

  const handleErr = errString => {
    if (!errString) {
      return;
    }

    if (
      errString.includes('status code 404') ||
      errString.includes('status code 400')
    ) {
      setErr('Email or password are not correct.');
    } else {
      setErr(
        'Unable to log in at this time. Please try again later or contact support.'
      );
    }
  };

  const onSubmit = async data => {
    try {
      const res = await signIn('credentials', {
        redirect: false,
        email: data.email,
        password: data.password,
      });

      if (!res.error) {
        return Router.push('/');
      }

      handleErr(res.error);
    } catch (error) {
      setErr(
        'Internal error happened. Please try again later or contact support.'
      );
    }
  };

  const clearSuccess = () => {
    setShowVerifiedBanner(false);
  };

  const isEnabled = () => {
    // if any are undefined return false
    return getValues(['email', 'password']).every(
      v => !!v
    );
  };

  return (
    <>
      <Head>
        <title>Animus Store | Login</title>
        <meta name='viewport' content='width=device-width, initial-scale=1' />
        <meta
          name='description'
          content='Decentralized storage made easy. Manage files, seamlessly access IPFS and regain control over your private data.'
        />
      </Head>
      <div className='grid justify-items-center items-center h-screen bg-gray-50'>
        <div className='w-full max-w-lg mx-auto overflow-hidden bg-white rounded-lg shadow-md dark:bg-gray-800'>
          {/* <pre>{JSON.stringify(watch(), null, 2)}</pre> */}
          <div className='px-6 py-4'>
            {showVerifiedBanner && (
              <div className='py-2'>
                <AlertSuccess
                  message={'Email verification complete.'}
                  onClick={clearSuccess}
                />
              </div>
            )}
            <h2 className='text-3xl font-bold text-center text-gray-700 dark:text-white'>
              Animus Store
            </h2>

            <p className='mt-1 text-center text-gray-500 dark:text-gray-400'>
              Login
            </p>

            <form onSubmit={handleSubmit(onSubmit)}>
              <div className='w-full mt-4'>
                <input
                  className='block w-full px-4 py-2 mt-2 text-gray-700 placeholder-gray-500 bg-white border rounded-md dark:bg-gray-800 dark:border-gray-600 dark:placeholder-gray-400 focus:border-blue-400 dark:focus:border-blue-300 focus:ring-opacity-40 focus:outline-none focus:ring focus:ring-blue-300'
                  type='email'
                  placeholder='Email Address'
                  aria-label='Email Address'
                  {...register('email', {
                    required: true,
                    validate: v =>
                      isEmailValid(v) || 'A valid email address is required',
                  })}
                />
                {errors.email && (
                  <p className='text-sm text-red-600 py-2 px-1'>
                    {errors.email.message}
                  </p>
                )}
              </div>

              <div className='w-full mt-4 relative'>
                <input
                  className='block w-full px-4 py-2 mt-2 text-gray-700 placeholder-gray-500 bg-white border rounded-md dark:bg-gray-800 dark:border-gray-600 dark:placeholder-gray-400 focus:border-blue-400 dark:focus:border-blue-300 focus:ring-opacity-40 focus:outline-none focus:ring focus:ring-blue-300'
                  type={showPass ? 'text' : 'password'}
                  placeholder='Password'
                  aria-label='Password'
                  {...register('password', {
                    required: true,
                  })}
                />
                <button
                  type='button'
                  className='absolute inset-y-0 right-0 pr-3 items-center text-sm leading-5'
                  onClick={e => {
                    e.preventDefault();
                    setShowPass(prev => !prev);
                  }}
                >
                  {showPass ? (
                    <EyeIcon className='w-5 h-5 text-gray-600' />
                  ) : (
                    <EyeOffIcon className='w-5 h-5 text-gray-600' />
                  )}
                </button>
              </div>
              {err && <p className='text-sm text-red-600 py-2 px-1'>{err}</p>}

              <div className='flex items-center justify-between mt-4'>
                <a
                  href='#'
                  className='text-sm text-gray-600 dark:text-gray-200 hover:text-gray-500'
                >
                  Forget Password?
                </a>

                <button
                  type='submit'
                  className="px-4 py-2 leading-5 text-white transition-colors duration-200 transform bg-gray-700 hover:bg-gray-600 rounded focus:outline-none"
                >
                  Login
                </button>
              </div>
            </form>
          </div>

          <div className='flex items-center justify-center py-4 text-center bg-gray-50 dark:bg-gray-700'>
            <span className='text-sm text-gray-600 dark:text-gray-200'>
              Don&apos;t have an account?{' '}
            </span>

            <Link href='/account/register'>
              <a className='mx-2 text-sm font-bold text-blue-500 dark:text-blue-400 hover:underline'>
                Register
              </a>
            </Link>
          </div>
        </div>
      </div>
    </>
  );
}

function isEmailValid (email) {
  return /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/.test(
    email
  );
}

export async function getServerSideProps (context) {
  const { verified } = context.query;
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
      verified: verified ? true : false,
    },
  };
}
