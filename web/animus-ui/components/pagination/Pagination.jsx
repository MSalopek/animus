import Link from 'next/link';
import { ArrowLeftIcon, ArrowRightIcon } from '@heroicons/react/outline';

export default Pagination;

function Pagination({ url, currentPage, pages }) {
  return (
    <div className="flex flex-col items-center sm:flex-row sm:justify-between">
      <div className="flex flex-col items-center">
        <span className="text-sm text-gray-700 dark:text-gray-400">
          Showing page{' '}
          <span className="font-semibold text-gray-900 dark:text-white">
            {currentPage}
          </span>{' '}
          out of{' '}
          <span className="font-semibold text-gray-900 dark:text-white">
            {pages}
          </span>{' '}
        </span>
        <div className="inline-flex mt-2 xs:mt-0 gap-x-2">
          {currentPage !== 1 ? (
            <Link href={`${url}?page=${currentPage - 1}`}>
              <a className="inline-flex gap-x-2 items-center py-2 px-4 text-sm font-medium text-gray-6=700 hover:bg-gray-50 rounded border border-gray-400">
                <ArrowLeftIcon className="w-5 h-5" />
                Prev
              </a>
            </Link>
          ) : (
            <a
              role="link"
              aria-disabled={true}
              className="inline-flex gap-x-2 items-center py-2 px-4 text-sm font-medium text-gray-300 rounded border border-gray-200"
            >
              <ArrowLeftIcon className="w-5 h-5" />
              Prev
            </a>
          )}
          {currentPage !== pages ? (
            <Link href={`${url}?page=${currentPage + 1}`}>
              <a className="inline-flex gap-x-2 items-center py-2 px-4 text-sm font-medium text-gray-700 hover:bg-gray-50 rounded border border-gray-400">
                Next
                <ArrowRightIcon className="w-5 h-5" />
              </a>
            </Link>
          ) : (
            <a
              role="link"
              aria-disabled={true}
              className="inline-flex gap-x-2 items-center py-2 px-4 text-sm font-medium text-gray-300 rounded border border-gray-200"
            >
              Next
              <ArrowRightIcon className="w-5 h-5" />
            </a>
          )}
        </div>
      </div>
    </div>
  );
}
