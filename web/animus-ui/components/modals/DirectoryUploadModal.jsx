import { useState, useCallback } from 'react';
import Router from 'next/router';

import { useDropzone } from 'react-dropzone';
import { Dialog } from '@headlessui/react';
import { FolderAddIcon, CheckIcon, XIcon } from '@heroicons/react/outline';
import { MAXFILESIZE, MAXTOTALSIZE, MEGABYTE } from '../../util/constants';

import { Alert } from '../alert/Alert';

export default function DirectoryUploadModal({
  isOpen,
  setIsOpen,
  uploadFunc,
}) {
  return (
    <Dialog
      open={isOpen}
      onClose={() => setIsOpen(false)}
      className="fixed z-10 inset-0 overflow-y-auto"
    >
      <div className="flex flex-col items-center mt-60 h-screen">
        <Dialog.Overlay className="fixed inset-0 bg-black opacity-50" />
        <UploadDirectory setIsOpen={setIsOpen} uploadFunc={uploadFunc} />
      </div>
    </Dialog>
  );
}

// TODO: define on DROP -> restrict max uploaded dir size to 100mb
function UploadDirectory({ setIsOpen, uploadFunc }) {
  const [errMessage, setErrMessage] = useState('');
  const [disabled, setDisabled] = useState(true);
  const [dirName, setDirname] = useState('');

  const onDropAccepted = useCallback((acceptedFiles) => {
    if (!acceptedFiles.length) {
      return;
    }

    const sum = 0;
    acceptedFiles.forEach((f) => {
      sum += f.size;
    });

    if (sum > MAXTOTALSIZE) {
      setDisabled(true);
      setErrMessage('directory size exceeds 100Mb');
    }

    setDisabled(false);
  }, []);

  const onDropRejected = useCallback((fileRejections) => {
    if (fileRejections.length) {
      setDisabled(true);
      setErrMessage('one or multiple files were rejected');
      return;
    }
  }, []);

  const { getRootProps, getInputProps, acceptedFiles, fileRejections } =
    useDropzone({
      // NOTE:
      // must use multiple=true and set directory="" webkitdirectory="" mozdirectory=""
      // on the <input type="file"> for directory/multifile upload to work
      multiple: true,
      minSize: 0,
      maxSize: MAXFILESIZE,
      onDropRejected,
      onDropAccepted,
    });

  const clearErrMessage = () => {
    setErrMessage(null);
  };

  const upload = async () => {
    setDisabled(true);
    if (acceptedFiles.length < 1) {
      setErrMessage('No files provided.');
      return;
    }
    try {
      const res = await uploadFunc(acceptedFiles, dirName);
      setDisabled(true);
      if (res.status !== 201) {
        setErrMessage('Something went wrong.');
        return;
      }
      Router.reload(window.location.pathname);
    } catch (error) {
      // The request was made and server responded with error code
      if (error?.response) {
        setDisabled(true);
        setErrMessage(
          `Error ${error.response.status}: ${error.response.data.error}`
        );
        return;
      }
      setErrMessage('Something went wrong.');
    }
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
  const rejectedList = (
    <ul className="flex flex-col list-group my-1 mx-4">
      {fileRejections.map((r) => (
        <li className="list-group-item list-group-item-fail inline-flex">
          <XIcon className="w-5 h-5 text-red-400" />
          <span className="text-gray-400 text-sm px-2">{r.file.name}</span>
          <span className="text-gray-600 text-sm">{`(${r.errors[0].code.replaceAll(
            '-',
            ' '
          )})`}</span>
        </li>
      ))}
    </ul>
  );
  return (
    <div className="shadow-lg relative flex flex-col w-1/3 bg-white bg-clip-padding rounded-md outline-none text-current p-2">
      {errMessage && <Alert message={errMessage} onClick={clearErrMessage} />}
      <div className="flex items-center justify-between p-4 border-b border-gray-200">
        <h5 className="text-xl font-medium leading-normal text-gray-800">
          Upload a Folder to IPFS
        </h5>
      </div>
      <div className="p-2">
        <label
          className="font-semibold text-gray-700 pr-2"
          htmlFor="folder-name"
        >
          Folder Name:
        </label>
        <input
          id="folder-name"
          type="text"
          placeholder="Enter name (min 3 chars)"
          className="px-4 w-96 py-2 mt-2 text-gray-700 placeholder-gray-400 bg-white border rounded-md focus:border-blue-400 focus:outline-none"
          value={dirName}
          onChange={(e) => {
            setDirname(e.target.value);
          }}
        />
      </div>
      <div className="p-2 bg-white" {...getRootProps()}>
        <div className="rounded-lg overflow-hidden">
          <div className="md:flex">
            <div className="w-full">
              <div className="relative h-52 rounded-md border-dashed border-2 border-blue-200 hover:bg-gray-100 flex justify-center items-center">
                <div className="absolute">
                  <div className="flex flex-col items-center">
                    <FolderAddIcon className="w-12 h-12 text-gray-400" />
                    <p className="pt-1 tracking-wider text-center text-gray-400 group-hover:text-gray-600">
                      Drag and drop a folder <br></br>
                      or click to choose
                    </p>
                  </div>
                </div>
                <input
                  type="file"
                  className="h-full w-full opacity-0 hidden"
                  directory=""
                  webkitdirectory=""
                  mozdirectory=""
                  {...getInputProps()}
                />
              </div>
            </div>
          </div>
        </div>
      </div>
      {acceptedFiles.length > 0 && acceptedList}
      {fileRejections.length > 0 && rejectedList}
      <div className="flex flex-shrink-0 gap-2 flex-wrap items-center justify-end p-4 border-t border-gray-200 rounded-b-md">
        <button
          type="button"
          disabled={disabled || dirName.length < 3}
          className={`inline-block px-6 py-2.5 ${
            disabled || dirName.length < 3
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
