import { Dialog } from "@headlessui/react";
import { CloudUploadIcon } from "@heroicons/react/outline";

export default function FileUploadModal({ isOpen, setIsOpen }) {
  return (
    <Dialog
      open={isOpen}
      onClose={() => setIsOpen(false)}
      className="fixed z-10 inset-0 overflow-y-auto"
    >
      <div className="flex items-center justify-center h-screen">
        <Dialog.Overlay className="fixed inset-0 bg-black opacity-50" />

        <div className="shadow-lg relative flex flex-col w-96 bg-white bg-clip-padding rounded-md outline-none text-current p-2">
          <div className="flex items-center justify-between p-4 border-b border-gray-200">
            <h5
              className="text-xl font-medium leading-normal text-gray-800"
              id="exampleModalScrollableLabel"
            >
              Upload a File to IPFS
            </h5>
          </div>

          <div className="py-2 bg-white px-2">
            <div className="max-w-md mx-auto rounded-lg overflow-hidden md:max-w-xl">
              <div className="md:flex">
                <div className="w-full">
                  <div className="relative h-52 rounded-md border-dashed border-2 border-blue-200 hover:bg-gray-100 flex justify-center items-center">
                    <div className="absolute">
                      <div className="flex flex-col items-center">
                        <CloudUploadIcon className="w-12 h-12 text-gray-400" />
                        <p className="pt-1 tracking-wider text-gray-400 group-hover:text-gray-600">
                          Upload File
                        </p>
                      </div>
                    </div>
                    <input
                      type="file"
                      className="h-full w-full opacity-0"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
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
      </div>
    </Dialog>
  );
}
