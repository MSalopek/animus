import { ArrowLeftIcon, ArrowRightIcon } from '@heroicons/react/outline';
import Link from 'next/link';

export default Pagination;

function Pagination() {
  const total = 100;
  const pages = 10;
  const currentPage = 1;
  const url = '/test/pagination';

  return (
    <div className="w-full bg-white dark:bg-gray-800">
      <div className="container flex flex-col items-center px-6 py-5 mx-auto space-y-6 sm:flex-row sm:justify-between sm:space-y-0 ">
        <div className="flex flex-col items-center">
          <span className="text-sm text-gray-700 dark:text-gray-400">
            Showing page{' '}
            <span className="font-semibold text-gray-900 dark:text-white">
              1
            </span>{' '}
            out of{' '}
            <span className="font-semibold text-gray-900 dark:text-white">
              100
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
    </div>
  );
}
