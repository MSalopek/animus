import { useState } from 'react';
import { Dialog } from '@headlessui/react';
import {
  DuplicateIcon,
  EyeIcon,
  EyeOffIcon,
} from '@heroicons/react/outline';

import { CopyToClipboard } from 'react-copy-to-clipboard';

import { Alert } from '../alert/Alert';

export default function CreateKeyModal({
  isOpen,
  setIsOpen,
  currentKeys,
  setCurrentKeys,
  createFunc,
}) {
  return (
    <Dialog
      open={isOpen}
      onClose={() => setIsOpen(false)}
      className="fixed z-10 inset-0 overflow-y-auto"
    >
      <div className="flex flex-col items-center mt-60 h-screen">
        <Dialog.Overlay className="fixed inset-0 bg-black opacity-50" />
        <AddKeyBox
          setIsOpen={setIsOpen}
          currentKeys={currentKeys}
          setCurrentKeys={setCurrentKeys}
          createFunc={createFunc}
        />
      </div>
    </Dialog>
  );
}

function AddKeyBox({ setIsOpen, currentKeys, setCurrentKeys, createFunc }) {
  const [errMessage, setErrMessage] = useState('');
  const [disabled, setDisabled] = useState(false);

  const [key, setKey] = useState('');
  const [secret, setSecret] = useState('');
  const [hidden, setHidden] = useState(true);

  const create = async () => {
    setDisabled(true);
    try {
      const res = await createFunc();
      if (res.status !== 201) {
        setErrMessage('Something went wrong.');
        return;
      }

      if (res.data) {
        setKey(res.data.client_key);
        setSecret(res.data.client_secret);
        // update existing key list
        const newKeys = [...currentKeys];
        newKeys.push(res.data);
        setCurrentKeys(newKeys.sort((a, b) => a.disabled - b.disabled));
      }
    } catch (error) {
      if (error.response.status === 400) {
        // show server response error to user
        setErrMessage(
          error.response.data?.error
            ? `Cannot create key: ${error.response.data?.error}`
            : 'Something went wrong.'
        );
        return;
      }
      setErrMessage('Something went wrong.');
    }
  };

  const clearErrMessage = () => {
    setErrMessage(null);
  };

  return (
    <div className="shadow-lg relative flex flex-col w-1/3 bg-white bg-clip-padding rounded-md outline-none text-current p-2">
      {errMessage && <Alert message={errMessage} onClick={clearErrMessage} />}
      <div className="p-4">
        <h5 className="text-2xl font-medium leading-normal text-gray-800 border-b border-gray-200">
          New API Key
        </h5>
        <p className="font-bold text-red-600 py-2">Important Notes</p>
        <p className="text-gray-400 py-2">
          Your API Key Secret will be displayed only once. <br></br>The Secret
          cannot be retrieved after you close this window.<br></br> If you lose
          your Secret you should delete the API Key and create a new one.
          <br></br>
          <br></br>
          Click Create to create a new API Key.
        </p>
        <p className="text-gray-500 font-bold py-2">
          API Keys allow access to your data - please keep them stored securely
          to prevent unauthorized access.
        </p>
      </div>

      {key && secret ? (
        <div className="px-8 py-4">
          <div className="flex gap-2 py-2 items-center">
            <p className="font-semibold text-gray-700 w-1/5 text-lg">
              Client Key:
            </p>
            <p className="w-3/5">{key}</p>
            <CopyToClipboard text={key}>
              <DuplicateIcon className="w-5 h-5 cursor-pointer" />
            </CopyToClipboard>
          </div>
          <div className="flex gap-2 py-2 items-center">
            <p className="font-semibold text-gray-700 w-1/5 text-lg">
              Client Secret:
            </p>
            <p className="break-all w-3/5 text-sm">
              {hidden ? '*'.repeat(48) : secret}
            </p>
            <button onClick={() => setHidden((prev) => !prev)}>
              {hidden ? (
                <EyeIcon className="w-5 h-5" />
              ) : (
                <EyeOffIcon className="w-5 h-5" />
              )}
            </button>
            <CopyToClipboard text={secret}>
              <DuplicateIcon className="w-5 h-5 cursor-pointer" />
            </CopyToClipboard>
          </div>
        </div>
      ) : (
        <></>
      )}

      <div className="flex flex-shrink-0 gap-2 flex-wrap items-center justify-end p-4 border-t border-gray-200 rounded-b-md">
        <button
          type="button"
          disabled={disabled}
          className={`inline-block px-6 py-2.5 ${
            disabled
              ? 'bg-gray-200'
              : 'bg-blue-600 hover:bg-blue-700 hover:shadow-lg focus:bg-blue-700'
          } text-white font-medium text-xs leading-tight uppercase rounded shadow-md focus:shadow-lg focus:outline-none focus:ring-0 active:bg-blue-800 active:shadow-lg transition duration-150 ease-in-out ml-1`}
          onClick={() => create()}
        >
          Create
        </button>
        <button
          type="button"
          className="inline-block px-6 py-2.5 bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-purple-700 hover:shadow-lg focus:bg-purple-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 ease-in-out"
          onClick={() => setIsOpen(false)}
        >
          Close
        </button>
      </div>
    </div>
  );
}
