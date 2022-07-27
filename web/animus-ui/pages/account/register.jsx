import Head from 'next/head';
import Link from 'next/link';

import { EyeIcon, EyeOffIcon } from '@heroicons/react/outline';

import { unstable_getServerSession } from 'next-auth/next';
import { authOptions } from '../api/auth/[...nextauth]';

import { useForm } from 'react-hook-form';

import { RegisterUser } from '../../service/http';
import { useState, useRef } from 'react';
import { Alert } from '../../components/alert/Alert';

export default Register;

function Register () {
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
    getValues,
  } = useForm({ mode: 'onChange' });

  // for password confirmation
  const pwRef = useRef({});
  pwRef.current = watch('password', '');

  const [success, setSuccess] = useState(false);
  const [requestErr, setRequestErr] = useState('');

  const [showPass, setShowPass] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);

  const isEnabled = () => {
    // if any are undefined return false
    return getValues(['username', 'email', 'password', 'confirm', 'tos']).every(
      v => !!v
    );
  };

  const onSubmit = async data => {
    try {
      const p = {
        email: data.email,
        password: data.password,
        username: data.username,
        firstname: data.firstname,
        lastname: data.lastname,
      };
      const res = await RegisterUser(p);

      if (!res.error) {
        setSuccess(true);
      }
    } catch (error) {
      if (error?.response?.status === 409) {
        setRequestErr('User already exists. Did you forget your password?')
      } else {
        setRequestErr(
          'Internal error happened. Please try again later or contact support.'
        );
      }
    }
  };

  const clearErrMessage = () => {
    setRequestErr(null);
  };

  return (
    <>
      <Head>
        <title>Animus Store | Register</title>
        <meta name='viewport' content='width=device-width, initial-scale=1' />
        <meta
          name='description'
          content='Decentralized storage made easy. Manage files, seamlessly access IPFS and regain control over your private data.'
        />
      </Head>
      <div className='grid justify-items-center items-center h-screen p-0 sm:p-8'>
        {/* <pre>{JSON.stringify(watch(), null, 2)}</pre> */}
        {success ? (
          <Success email={getValues('email')} />
        ) : (
          <div className='w-full max-w-lg mx-auto border bg-white rounded-lg shadow-md dark:bg-gray-800'>
            <div className='px-6 py-4'>
              <h2 className='text-3xl font-bold text-center text-gray-700 dark:text-whitem mt-4'>
                Animus Store
              </h2>

              <h2 className='text-lg font-semibold text-gray-700 capitalize text-center my-6 dark:text-white'>
                Register your account
              </h2>

              {requestErr && (
                <Alert message={requestErr} onClick={clearErrMessage} />
              )}
              <form onSubmit={handleSubmit(onSubmit)}>
                <div className='mt-6'>
                  <label
                    className='text-gray-700 dark:text-gray-200'
                    htmlFor='username'
                  >
                    Username*
                  </label>
                  <input
                    id='username'
                    type='text'
                    placeholder='Enter username'
                    className='block w-full px-4 py-2 mt-2 text-gray-700 bg-white border border-gray-200 rounded-md dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 focus:ring-blue-300 focus:ring-opacity-40 dark:focus:border-blue-300 focus:outline-none focus:ring'
                    {...register('username', {
                      required: 'Username is required',
                      minLength: 3,
                      maxLength: 64,
                    })}
                  />
                  {errors.username && errors.username.type === 'required' && (
                    <p className='text-sm text-red-600 py-2 px-1'>
                      {errors.username.message}
                    </p>
                  )}
                  {errors.username &&
                    (errors.username.type === 'maxLength' ||
                      errors.username.type === 'minLength') && (
                      <p className='text-sm text-red-600 py-2 px-1'>
                        Username must be between 3 and 64 characters long
                      </p>
                    )}
                </div>

                <div className='mt-3'>
                  <label
                    className='text-gray-700 dark:text-gray-200'
                    htmlFor='emailAddress'
                  >
                    Email Address*
                  </label>
                  <input
                    id='email'
                    type='email'
                    placeholder='Enter email'
                    className='block w-full px-4 py-2 mt-2 text-gray-700 bg-white border border-gray-200 rounded-md dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 focus:ring-blue-300 focus:ring-opacity-40 dark:focus:border-blue-300 focus:outline-none focus:ring'
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

                <div className='mt-3'>
                  <label
                    className='text-gray-700 dark:text-gray-200'
                    htmlFor='password'
                  >
                    Password*
                  </label>
                  <div className='relative'>
                    <input
                      id='password'
                      type={showPass ? 'text' : 'password'}
                      placeholder='Enter password'
                      className='block w-full px-4 py-2 mt-2 text-gray-700 bg-white border border-gray-200 rounded-md dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 focus:ring-blue-300 focus:ring-opacity-40 dark:focus:border-blue-300 focus:outline-none focus:ring'
                      {...register('password', {
                        required: true,
                        validate: v =>
                          v.length >= 8 ||
                          'Password must at least 8 characters long',
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
                  {errors.password && (
                    <p className='text-sm text-green-600 py-2 px-1'>
                      {errors.password.message}
                    </p>
                  )}
                </div>

                <div className='mt-6'>
                  <label
                    className='text-gray-700 dark:text-gray-200'
                    htmlFor='confirm'
                  >
                    Confirm Password*
                  </label>
                  <div className='relative'>
                    <input
                      id='confirm'
                      type={showConfirm ? 'text' : 'password'}
                      placeholder='Re-enter password'
                      className='block w-full px-4 py-2 mt-2 text-gray-700 bg-white border border-gray-200 rounded-md dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 focus:ring-blue-300 focus:ring-opacity-40 dark:focus:border-blue-300 focus:outline-none focus:ring'
                      {...register('confirm', {
                        required: true,
                        validate: v =>
                          v === pwRef.current || 'Passwords do not match',
                      })}
                    />
                    <button
                      type='button'
                      className='absolute inset-y-0 right-0 pr-3 items-center text-sm leading-5'
                      onClick={e => {
                        e.preventDefault();
                        setShowConfirm(prev => !prev);
                      }}
                    >
                      {showConfirm ? (
                        <EyeIcon className='w-5 h-5 text-gray-600' />
                      ) : (
                        <EyeOffIcon className='w-5 h-5 text-gray-600' />
                      )}
                    </button>
                  </div>
                  {errors.confirm && (
                    <p className='text-sm text-red-600 py-2 px-1'>
                      {errors.confirm.message}
                    </p>
                  )}
                </div>

                <div className='relative flex w-full items-center pt-7 pb-5'>
                  <div className='flex-grow border-t border-gray-200'></div>
                  <p className='p-1.5 mx-2 text-gray-400'>Optional</p>
                  <div className='flex-grow border-t border-gray-200'></div>
                </div>

                <div>
                  <label
                    className='text-gray-700 dark:text-gray-200'
                    htmlFor='firstname'
                  >
                    Firstname
                  </label>
                  <input
                    id='firstname'
                    type='text'
                    placeholder='Enter firstname'
                    className='block w-full px-4 py-2 mt-2 text-gray-700 bg-white border border-gray-200 rounded-md dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 focus:ring-blue-300 focus:ring-opacity-40 dark:focus:border-blue-300 focus:outline-none focus:ring'
                    {...register('firstname')}
                  />
                </div>

                <div className='mt-3'>
                  <label
                    className='text-gray-700 dark:text-gray-200'
                    htmlFor='lastname'
                  >
                    Lastname
                  </label>
                  <input
                    id='lastname'
                    type='text'
                    placeholder='Enter lastname'
                    className='block w-full px-4 py-2 mt-2 text-gray-700 bg-white border border-gray-200 rounded-md dark:bg-gray-800 dark:text-gray-300 dark:border-gray-600 focus:border-blue-400 focus:ring-blue-300 focus:ring-opacity-40 dark:focus:border-blue-300 focus:outline-none focus:ring'
                    {...register('lastname')}
                  />
                </div>
                <div className='mt-8 flex items-center'>
                  <input
                    id='tos'
                    {...register('tos')}
                    type='checkbox'
                    value={true}
                    placeholder=''
                  />
                  <label
                    className='text-gray-700 dark:text-gray-200 text-xs px-2'
                    htmlFor='tos'
                  >
                    I have read and accept the{' '}
                    <Link href='/terms-and-conditions.html'>
                      <a className="text-gray-600 font-semibold underline">Terms of Service</a>
                    </Link>
                  </label>
                </div>

                <div className='flex justify-center mt-8'>
                  <button
                    type='submit'
                    className={`px-6 py-2 leading-5 text-white transition-colors duration-200 transform ${
                      !isEnabled()
                        ? 'bg-gray-200'
                        : 'bg-gray-700 hover:bg-gray-600'
                    } rounded-md focus:outline-none focus:bg-gray-600`}
                    disabled={!isEnabled()}
                  >
                    Register
                  </button>
                </div>
              </form>
            </div>

            <div className='flex items-center justify-center py-4 text-center bg-gray-50 dark:bg-gray-700'>
              <span className='text-sm text-gray-600 dark:text-gray-200'>
                Already have an account?{' '}
              </span>

              <Link href='/account/login'>
                <a className='mx-2 text-sm font-bold text-blue-500 dark:text-blue-400 hover:underline'>
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

function Success ({ email }) {
  return (
    <div className='w-full max-w-lg h-96 mx-auto overflow-hidden border bg-white dark:bg-gray-800 rounded-lg shadow-lg'>
      <div className='flex p-8 items-center gap-4'>
        <svg
          className='w-9 h-9 fill-current text-green-500'
          xmlns='http://www.w3.org/2000/svg'
          viewBox='0 0 24 24'
        >
          <path d='M0 0h24v24H0V0z' fill='none' />
          <path d='M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm4.59-12.42L10 14.17l-2.59-2.58L6 13l4 4 8-8z' />
        </svg>
        <h2 className='font-semibold text-2xl text-gray-800'>Account Created</h2>
      </div>
      <p className='px-8 text-lg text-gray-600 leading-relaxed'>
        You need to activate your account before you can proceed.
        <br></br>
        <br></br>
        Please check your inbox at{' '}
        <span className='font-semibold text-gray-800'> {email} </span>(including
        spam folder) and click on the provided link.
        <br></br>
        <br></br>
        You will receive your activation code shortly.
      </p>
    </div>
  );
}

function isEmailValid (email) {
  return /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/.test(
    email
  );
}

export async function getServerSideProps (context) {
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

  return {
    props: {
      session,
    },
  };
}
