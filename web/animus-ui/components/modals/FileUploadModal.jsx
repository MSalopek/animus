import { useState, useCallback } from 'react';
import Router from 'next/router';

import { useDropzone } from 'react-dropzone';
import { Dialog } from '@headlessui/react';
import {
  CloudUploadIcon,
  CheckIcon,
} from '@heroicons/react/outline';
import { MAX_FILE_SIZE, MEGABYTE } from '../../util/constants';

import { Alert } from '../alert/Alert';

export default function FileUploadModal({ isOpen, setIsOpen, uploadFunc }) {
  return (
    <Dialog
      open={isOpen}
      onClose={() => setIsOpen(false)}
      className="fixed z-10 inset-0 overflow-y-auto"
    >
      <div className="flex flex-col items-center mt-60 h-screen">
        <Dialog.Overlay className="fixed inset-0 bg-black opacity-50" />
        <UploadFileBox setIsOpen={setIsOpen} uploadFunc={uploadFunc} />
      </div>
    </Dialog>
  );
}

function UploadFileBox({ setIsOpen, uploadFunc }) {
  const [errMessage, setErrMessage] = useState('');
  const [disabled, setDisabled] = useState(true);

  const onDropRejected = useCallback((fileRejections) => {
    if (fileRejections.length) {
      setDisabled(true);
      setErrMessage(
        `${
          fileRejections[0].file.name
        } rejected - ${fileRejections[0].errors[0].code.replaceAll('-', ' ')}`
      );
      return;
    }
  }, []);

  const onDropAccepted = useCallback((acceptedFiles) => {
    if (acceptedFiles.length) {
      setDisabled(false);
      return;
    }
  }, []);

  const { getRootProps, getInputProps, acceptedFiles, fileRejections } =
    useDropzone({
      multiple: true,
      minSize: 0,
      maxSize: MAX_FILE_SIZE,
      onDropRejected,
      onDropAccepted,
    });

  const upload = async () => {
    setDisabled(true);
    if (acceptedFiles.length < 1) {
      setErrMessage('No files provided.');
      return;
    }
    try {
      const res = await uploadFunc(acceptedFiles[0]);
      setDisabled(true);
      if (res.status !== 201) {
        setErrMessage('Something went wrong.');
        return;
      }
      Router.reload(window.location.pathname);
    } catch (error) {
      // The request was made and server responded with error code
      setErrMessage(
        `Error ${error.response.status}: ${error.response.data.error}`
      );
      // console.log(error);
    }
  };

  const clearErrMessage = () => {
    setErrMessage(null);
  };

  const acceptedList = (
    <ul className="flex flex-col list-group my-1 mx-4">
      {acceptedFiles.map((af) => (
        <li className="list-group-item list-group-item-success inline-flex">
          <CheckIcon className="w-6 h-6 text-green-300" />
          <span className="text-gray-600 text-sm">{af.name}</span>
          <span className="text-gray-400 text-sm pl-2">{`(${(
            af.size / MEGABYTE
          ).toPrecision(3)} MB)`}</span>
        </li>
      ))}
    </ul>
  );

  return (
    <div className="shadow-lg relative flex flex-col w-1/3 bg-white bg-clip-padding rounded-md outline-none text-current p-2">
      {errMessage && <Alert message={errMessage} onClick={clearErrMessage} />}
      <div className="flex items-center justify-between p-4 border-b border-gray-200">
        <h5 className="text-xl font-medium leading-normal text-gray-800">
          Upload a File to IPFS
        </h5>
        <p className="text-sm text-gray-400">{'(max. 25 MB)'}</p>
      </div>

      <div className="py-2 px-3 bg-white" {...getRootProps()}>
        <div className="rounded-md overflow-hidden">
          <div className="md:flex">
            <div className="w-full">
              <div className="relative h-52 rounded-md border-dashed border-2 border-blue-200 hover:bg-gray-100 flex justify-center items-center">
                <div className="absolute">
                  <div className="flex flex-col items-center">
                    <CloudUploadIcon className="w-12 h-12 text-gray-400" />
                    <p className="pt-1 tracking-wider text-center text-gray-400 group-hover:text-gray-600">
                      Drag and drop a file <br></br>
                      or click to choose
                    </p>
                  </div>
                </div>
                <input
                  type="file"
                  className="h-full w-full opacity-0 hidden"
                  {...getInputProps()}
                />
              </div>
            </div>
          </div>
        </div>
      </div>
      {acceptedFiles.length > 0 && acceptedList}
      <div className="flex flex-shrink-0 gap-2 flex-wrap items-center justify-end p-4 border-t border-gray-200 rounded-b-md">
        <button
          type="button"
          disabled={disabled}
          className={`inline-block px-6 py-2.5 ${
            disabled
              ? 'bg-gray-200'
              : 'bg-blue-600 hover:bg-blue-700 hover:shadow-lg focus:bg-blue-700'
          } text-white font-medium text-xs leading-tight uppercase rounded shadow-md focus:shadow-lg focus:outline-none focus:ring-0 active:bg-blue-800 active:shadow-lg transition duration-150 ease-in-out ml-1`}
          onClick={() => upload()}
        >
          Upload
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
