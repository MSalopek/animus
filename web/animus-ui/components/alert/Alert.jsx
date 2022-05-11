import {
  InformationCircleIcon
} from '@heroicons/react/outline';

export function Alert({ message, onClick }) {
  return (
    <div
      className="bg-red-100 border-l-4 border-red-500 rounded-md text-teal-900 px-2 py-3 shadow-md"
      role="alert"
	  onClick={() => onClick()}
    >
      <div className="flex">
        <div className="py-2 px-1 mr-2">
          <InformationCircleIcon className="w-7 h-7 text-red-500" />
        </div>
        <div>
          <p className="font-bold">Oops!</p>
          <p className="text-sm">
            {message}
          </p>
        </div>
      </div>
    </div>
  );
}
