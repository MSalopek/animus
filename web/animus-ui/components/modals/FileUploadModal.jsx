import { useDropzone } from 'react-dropzone';
import { Dialog } from '@headlessui/react';
import { CloudUploadIcon, CheckIcon, XIcon } from '@heroicons/react/outline';

const MB = 1048576; // in bytes
const maxFilesize = 1 * MB;

export default function FileUploadModal({ isOpen, setIsOpen }) {
  return (
    <Dialog
      open={isOpen}
      onClose={() => setIsOpen(false)}
      className="fixed z-10 inset-0 overflow-y-auto"
    >
      <div className="flex flex-col items-center justify-center h-screen">
        <Dialog.Overlay className="fixed inset-0 bg-black opacity-50" />
        <Upload setIsOpen={setIsOpen} />
      </div>
    </Dialog>
  );
}

function Upload({ setIsOpen }) {
  const { getRootProps, getInputProps, acceptedFiles, fileRejections } =
    useDropzone({ multiple: true, minSize: 0, maxSize: maxFilesize });

  const acceptedList = (
    <ul className="flex flex-col list-group my-1 mx-1">
      {acceptedFiles.map((af) => (
          <li className="list-group-item list-group-item-success inline-flex">
            <CheckIcon className="w-6 h-6 text-green-300" />
            <span className="text-gray-600 text-sm">{af.name}</span>
          </li>
        ))}
    </ul>
  );
  const rejectedList = (
    <ul className="flex flex-col list-group my-1 mx-1">
      {fileRejections.map((r) => (
          <li className="list-group-item list-group-item-fail inline-flex">
            <XIcon className="w-5 h-5 text-red-400" />
            <span className="text-gray-400 text-sm px-2">{r.file.name}</span>
            <span className="text-gray-600 text-sm">{`(${r.errors[0].code.replaceAll("-", " ")})`}</span>
          </li>
        ))}
      {/* {fileRejections.map((r) => (console.log("ALL", r)))} */}
        
    </ul>
  );
  return (
    <div className="shadow-lg relative flex flex-col w-1/3 bg-white bg-clip-padding rounded-md outline-none text-current p-2">
      <div className="flex items-center justify-between p-4 border-b border-gray-200">
        <h5
          className="text-xl font-medium leading-normal text-gray-800"
          id="exampleModalScrollableLabel"
        >
          Upload a Files to IPFS
        </h5>
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
                      Drag and drop files <br></br>
                      or click here to choose files
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
      {acceptedFiles && acceptedList}
      {fileRejections && rejectedList}
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

{
  /* <div class="flex items-center justify-center w-full"> */
}
{
  /* <label for="dropzone-file" class="flex flex-col items-center justify-center w-full h-64 border-2 border-gray-300 border-dashed rounded-lg cursor-pointer bg-gray-50 dark:hover:bg-bray-800 dark:bg-gray-700 hover:bg-gray-100 dark:border-gray-600 dark:hover:border-gray-500 dark:hover:bg-gray-600"> */
}
{
  /* <div class="flex flex-col items-center justify-center pt-5 pb-6"> */
}
{
  /* <svg class="w-10 h-10 mb-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path></svg> */
}
{
  /* <p class="mb-2 text-sm text-gray-500 dark:text-gray-400"><span class="font-semibold">Click to upload</span> or drag and drop</p> */
}
{
  /* <p class="text-xs text-gray-500 dark:text-gray-400">SVG, PNG, JPG or GIF (MAX. 800x400px)</p> */
}
{
  /* </div> */
}
{
  /* <input id="dropzone-file" type="file" class="hidden"> */
}
{
  /* </label> */
}
{
  /* </div> */
}
