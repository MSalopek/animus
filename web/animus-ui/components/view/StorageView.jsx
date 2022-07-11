import { useState } from 'react';
import { CopyToClipboard } from 'react-copy-to-clipboard';
import { DocumentAddIcon, FolderAddIcon } from '@heroicons/react/solid';
import {
  TrashIcon,
  InformationCircleIcon,
  // ArrowsExpandIcon,
  DuplicateIcon,
} from '@heroicons/react/outline';

import { ModalBtn, BtnContainer, RoundActionBtn } from '../buttons/Buttons';
import Pagination from '../pagination/Pagination';
import FileUploadModal from '../modals/FileUploadModal';
import DirectoryUploadModal from '../modals/DirectoryUploadModal';

export default StorageView;

function StorageView({ rows, total, pages }) {
  const [isFileModalOpen, setIsFileModalOpen] = useState(false);
  const [isDirModalOpen, setIsDirModalOpen] = useState(false);

  console.log('## rows', rows, total, pages);
  return (
    <section className="bg-white dark:bg-gray-900">
      <FileUploadModal
        isOpen={isFileModalOpen}
        setIsOpen={setIsFileModalOpen}
      />
      <DirectoryUploadModal
        isOpen={isDirModalOpen}
        setIsOpen={setIsDirModalOpen}
      />
      <div className="container px-6 py-12 mx-auto">
        <div className="flex flex-col">
          <h1 className="text-3xl font-semibold text-gray-800 dark:text-white mb-2">
            Storage Manager
          </h1>
          <BtnContainer>
            <ModalBtn
              Icon={DocumentAddIcon}
              title={'Add File'}
              action={() => {
                setIsFileModalOpen(true);
              }}
            />
            <ModalBtn
              Icon={FolderAddIcon}
              title={'Add Directory'}
              action={() => {
                setIsDirModalOpen(true);
              }}
            />
          </BtnContainer>
        </div>

        {rows && (
          <div className="mt-4 space-y-4 lg:mt-8 border p-4 rounded-lg">
            {rows.map((r) => (
              <StorageRow
                key={r.name}
                id={r.id}
                cid={r.cid}
                dir={r.dir}
                name={r.name}
                isPublic={r.public}
                meta={r.meta}
                stage={r.stage}
                pinned={r.pinned}
                created_at={r.created_at}
              />
            ))}
          </div>
        )}
      </div>
      <Pagination total={total} />
    </section>
  );
}

function StorageRow({
  id,
  cid,
  dir,
  name,
  isPublic,
  meta,
  stage,
  pinned,
  created_at,
}) {
  const [expanded, setExpanded] = useState(false);

  return (
    <div className="py-2 px-8 bg-gray-100 rounded-lg dark:bg-gray-800">
      <div className="flex items-center justify-between">
        <div className="flex flex-row">
          <div clasName="text-gray-700 dark:text-white">
            <span className="text-sm text-gray-400">Name:</span>
            <h1 className="font-semibold w-56 lg:w-80">{name}</h1>
          </div>

          <div clasName="text-gray-700 dark:text-white">
            <span className="text-sm text-gray-400">CID:</span>
            <div className="flex items-center text-gray-700 dark:text-white">
              <h1 className="pr-2">{cid || 'N/A'}</h1>
              {cid && (
                <CopyToClipboard text={cid || ''}>
                  <DuplicateIcon className="w-5 h-5" />
                </CopyToClipboard>
              )}
            </div>
          </div>
        </div>

        <BtnContainer>
          <RoundActionBtn
            Icon={InformationCircleIcon}
            onClick={() => setExpanded(!expanded)}
          />
          <RoundActionBtn Icon={TrashIcon} />
          {/* <RoundActionBtn Icon={ArrowsExpandIcon} /> */}
        </BtnContainer>
      </div>
      {expanded && (
        <StorageMeta
          isPinned={pinned}
          isPublic={isPublic}
          metaData={meta}
          location={stage}
          created_at={created_at}
        />
      )}
    </div>
  );
}

function StorageMeta({
  metaData,
  dir,
  isPinned,
  isPublic,
  location,
  created_at,
}) {
  return (
    <div className="grid grid-cols-2 w-1/2 py-2">
      <p className="text-gray-500 py-0.5">
        <span className="font-semibold pr-2">Created:</span>
        {new Date(created_at).toISOString().split('T')[0]}
      </p>
      <p className="text-gray-500 py-0.5">
        <span className="font-semibold pr-2">Directory:</span>
        {dir ? 'true' : 'false'}
      </p>

      <p className="text-gray-500 py-0.5">
        <span className="font-semibold pr-2">Pinned:</span>
        {isPinned ? 'true' : 'false'}
      </p>
      <p className="text-gray-500 py-0.5">
        <span className="font-semibold pr-2">Public: </span>
        {isPublic ? 'true' : 'false'}
      </p>
      <p className="text-gray-500 py-0.5">
        <span className="font-semibold pr-2">Location:</span>
        {location ? location : 'N/A'}
      </p>
      {/* <p className="text-gray-500 py-0.5 col-span-2">
        <span className="font-semibold pr-2">Metadata:</span>
        <br></br>
        {metaData ? JSON.stringify(metaData) : 'N/A'}
      </p> */}
    </div>
  );
}
