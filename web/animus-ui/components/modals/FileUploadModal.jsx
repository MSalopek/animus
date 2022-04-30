import { useDropzone } from 'react-dropzone';
import { Dialog } from '@headlessui/react';
import { CloudUploadIcon, FolderIcon, CheckIcon, XIcon } from '@heroicons/react/outline';
import { MAXSIZE, MEGABYTE } from '../../util/constants';

export default function FileUploadModal({ isOpen, setIsOpen }) {
  return (
    <Dialog
      open={isOpen}
      onClose={() => setIsOpen(false)}
      className="fixed z-10 inset-0 overflow-y-auto"
    >
      <div className="flex flex-col items-center mt-60 h-screen">
        <Dialog.Overlay className="fixed inset-0 bg-black opacity-50" />
        <UploadFile setIsOpen={setIsOpen} />
      </div>
    </Dialog>
  );
}

function UploadFile({ setIsOpen }) {
  const { getRootProps, getInputProps, acceptedFiles, fileRejections } =
    useDropzone({ multiple: false, minSize: 0, maxSize: MAXSIZE });

  const acceptedList = (
    <ul className="flex flex-col list-group my-1 mx-4">
      {acceptedFiles.map((af) => (
          <li className="list-group-item list-group-item-success inline-flex">
            <CheckIcon className="w-6 h-6 text-green-300" />
            <span className="text-gray-600 text-sm">{af.name}</span>
            <span className="text-gray-400 text-sm pl-2">{`(${(af.size/MEGABYTE).toPrecision(3)} MB)`}</span>
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
            <span className="text-gray-600 text-sm">{`(${r.errors[0].code.replaceAll("-", " ")})`}</span>
          </li>
        ))}
        
    </ul>
  );
  return (
    <div className="shadow-lg relative flex flex-col w-1/3 bg-white bg-clip-padding rounded-md outline-none text-current p-2">
      <div className="flex items-center justify-between p-4 border-b border-gray-200">
        <h5
          className="text-xl font-medium leading-normal text-gray-800"
        >
          Upload a File to IPFS
        </h5>
        <p className="text-sm text-gray-400">{"(max. 25 MB)"}</p>
      </div>

      <div className="p-2 bg-white" {...getRootProps()}>
        <div className="rounded-lg overflow-hidden">
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
      {fileRejections.length > 0 && rejectedList}
      <div className="flex flex-shrink-0 gap-2 flex-wrap items-center justify-end p-4 border-t border-gray-200 rounded-b-md">
        <button
          type="button"
          className="inline-block px-6 py-2.5 bg-blue-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-blue-700 hover:shadow-lg focus:bg-blue-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-blue-800 active:shadow-lg transition duration-150 ease-in-out ml-1"
          onClick={() => setIsOpen(false)}
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

function UploadDirectory({ setIsOpen }) {
  const { getRootProps, getInputProps, acceptedFiles, fileRejections } =
    useDropzone({ multiple: true, minSize: 0, maxSize: maxFilesize });

  const acceptedList = (
    <ul className="flex flex-col list-group my-1 mx-4">
      {acceptedFiles.map((af) => (
          <li className="list-group-item list-group-item-success inline-flex">
            <CheckIcon className="w-6 h-6 text-green-300" />
            <span className="text-gray-600 text-sm">{af.name}</span>
            <span className="text-gray-400 text-sm pl-2">{`(${(af.size/MB).toPrecision(3)} MB)`}</span>
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
            <span className="text-gray-600 text-sm">{`(${r.errors[0].code.replaceAll("-", " ")})`}</span>
          </li>
        ))}
        
    </ul>
  );
  return (
    <div className="shadow-lg relative flex flex-col w-1/3 bg-white bg-clip-padding rounded-md outline-none text-current p-2">
      <div className="flex items-center justify-between p-4 border-b border-gray-200">
        <h5
          className="text-xl font-medium leading-normal text-gray-800"
          id="exampleModalScrollableLabel"
        >
          Upload a Directory to IPFS
        </h5>
      </div>

      <div className="p-2 bg-white" {...getRootProps()}>
        <div className="rounded-lg overflow-hidden">
          <div className="md:flex">
            <div className="w-full">
              <div className="relative h-52 rounded-md border-dashed border-2 border-blue-200 hover:bg-gray-100 flex justify-center items-center">
                <div className="absolute">
                  <div className="flex flex-col items-center">
                    <FolderIcon className="w-12 h-12 text-gray-400" />
                    <p className="pt-1 tracking-wider text-center text-gray-400 group-hover:text-gray-600">
                      Drag and drop folder <br></br>
                      or click here to choose
                    </p>
                  </div>
                </div>
                <input
                  type="file"
                  className="h-full w-full opacity-0 hidden"
                  directory="" webkitdirectory="" mozdirectory=""
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
          className="inline-block px-6 py-2.5 bg-blue-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-blue-700 hover:shadow-lg focus:bg-blue-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-blue-800 active:shadow-lg transition duration-150 ease-in-out ml-1"
          onClick={() => setIsOpen(false)}
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