import { useState } from 'react';
import Router from 'next/router';
import { useSession } from 'next-auth/react';

import { CopyToClipboard } from 'react-copy-to-clipboard';
import { DocumentAddIcon, FolderAddIcon } from '@heroicons/react/solid';
import {
  DownloadIcon,
  DuplicateIcon,
  TrashIcon,
  InformationCircleIcon,
} from '@heroicons/react/outline';

import { ModalBtn, BtnContainer, RoundActionBtn } from '../buttons/Buttons';
import Pagination from '../pagination/Pagination';
import FileUploadModal from '../modals/FileUploadModal';
import DirectoryUploadModal from '../modals/DirectoryUploadModal';
import { DeleteStorage, UploadDirectory, UploadFile } from '../../service/http';

export default StorageView;

function StorageView({ rows, total, pages, currentPage }) {
  const { data: session, status } = useSession();

  const [isFileModalOpen, setIsFileModalOpen] = useState(false);
  const [isDirModalOpen, setIsDirModalOpen] = useState(false);

  const uploadFile = async (file) => {
    return await UploadFile(session.user.accessToken, file);
  };

  const uploadDir = async (files, dirname) => {
    return await UploadDirectory(session.user.accessToken, files, dirname);
  };

  const deleteRow = async (id) => {
    await DeleteStorage(session.user.accessToken, id);
    Router.reload(window.location.pathname);
  };

  return (
    <section className="bg-white dark:bg-gray-900 min-h-screen">
      <FileUploadModal
        isOpen={isFileModalOpen}
        setIsOpen={setIsFileModalOpen}
        uploadFunc={uploadFile}
      />
      <DirectoryUploadModal
        isOpen={isDirModalOpen}
        setIsOpen={setIsDirModalOpen}
        uploadFunc={uploadDir}
      />
      <div className="container px-6 pt-12 mx-auto">
        <div className="flex flex-col">
          <h1 className="text-3xl font-semibold text-gray-800 dark:text-white mb-2">
            Storage Manager
          </h1>
          <div className="flex w-full justify-between">
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
                title={'Add Folder'}
                action={() => {
                  setIsDirModalOpen(true);
                }}
              />
            </BtnContainer>
            <Pagination url={'/'} currentPage={currentPage} pages={pages} />
          </div>
        </div>

        {rows && (
          <div className="space-y-2 lg:mt-8 border p-4 rounded-lg">
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
                deleteFunc={deleteRow}
              />
            ))}
          </div>
        )}
      </div>
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
  deleteFunc,
}) {
  const [expanded, setExpanded] = useState(false);

  return (
    <div className="py-1 px-8 bg-gray-100 rounded-lg dark:bg-gray-800">
      <div className="flex items-center justify-between">
        <div className="grid grid-cols-2 w-3/4">
          <div className="text-gray-700 dark:text-white">
            <span className="text-sm text-gray-400">Name:</span>
            <h1 className="font-semibold break-all">{name}</h1>
          </div>

          <div className="text-gray-700 dark:text-white">
            <span className="text-sm text-gray-400">CID:</span>
            <div className="flex flex-row sm:flex-col md:flex-row lg:flex-row xl:flex-row items-center text-gray-700 dark:text-white">
              <h1 className="pr-2 break-all">{cid || 'N/A'}</h1>
              {cid && (
                <span>
                  <CopyToClipboard text={cid || ''}>
                    <DuplicateIcon className="w-5 h-5" />
                  </CopyToClipboard>
                </span>
              )}
            </div>
          </div>
        </div>

        <BtnContainer>
          <RoundActionBtn Icon={DownloadIcon}></RoundActionBtn>
          <RoundActionBtn
            Icon={InformationCircleIcon}
            onClick={() => setExpanded(!expanded)}
          />
          <RoundActionBtn Icon={TrashIcon} onClick={() => deleteFunc(id)} />
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
