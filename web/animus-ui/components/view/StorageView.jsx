import { useState } from "react";
import { CopyToClipboard } from "react-copy-to-clipboard";
import { DocumentAddIcon, FolderAddIcon } from "@heroicons/react/solid";
import {
  TrashIcon,
  InformationCircleIcon,
  // ArrowsExpandIcon,
  DuplicateIcon,
} from "@heroicons/react/outline";

import { ModalBtn, BtnContainer, RoundActionBtn } from "../buttons/Buttons";
import Pagination from "../pagination/Pagination";
import FileUploadModal from "../modals/FileUploadModal";

export default StorageView;

function StorageView() {
  const [isFileModalOpen, setIsFileModalOpen] = useState(false);

  return (
    <section className="bg-white dark:bg-gray-900">
        <FileUploadModal
        isOpen={isFileModalOpen}
        setIsOpen={setIsFileModalOpen}
      />
      <div className="container px-6 py-12 mx-auto">
        <div className="flex flex-col">
          <h1 className="text-3xl font-semibold text-gray-800 dark:text-white mb-2">
            Storage Manager
          </h1>
          <BtnContainer>
            <ModalBtn
              Icon={DocumentAddIcon}
              title={"Add File"}
              action={() => {
                setIsFileModalOpen(true);
              }}
            />
            <ModalBtn
              Icon={FolderAddIcon}
              title={"Add Directory"}
              action={() => {
                setIsFileModalOpen(false);
              }}
            />
          </BtnContainer>
        </div>

        <div className="mt-4 space-y-4 lg:mt-8 border p-4 rounded-lg">
          <StorageRow />
        </div>
      </div>
      <Pagination />
    </section>
  );
}

function StorageRow() {
  const [expanded, setExpanded] = useState(false);

  return (
    <div className="py-2 px-8 bg-gray-100 rounded-lg dark:bg-gray-800">
      <div className="flex items-center justify-between w-full">
        <h1 className="font-semibold text-lg text-gray-700 dark:text-white">
          Filename.txt
          <span className="text-sm px-2 text-gray-400">(11.12.2022)</span>
        </h1>

        <div className="flex items-center text-gray-700 dark:text-white">
          <span className="pr-2">
            QmVaTLF6H23B7tno8REwSFxT8J21aWEnaSTpTzNM6sRi6b
          </span>
          <CopyToClipboard
            text={"QmVaTLF6H23B7tno8REwSFxT8J21aWEnaSTpTzNM6sRi6b"}
          >
            <RoundActionBtn Icon={DuplicateIcon} />
          </CopyToClipboard>
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
          isPinned={true}
          isPublic={false}
          metaData="{ad;fkasdf;klams;dfklma;sdlkfma;slkdmfa;skdmf;alskdmf;alskmdf;alksmd;flaksmd;flkams;dlkfma;slkdmf;alksdmf;alskmdf;alksmdf}"
        />
      )}
    </div>
  );
}

function StorageMeta({ metaData, isPinned, isPublic }) {
  return (
    <div className="grid grid-cols-2 w-3/4 pb-4">
      <p className="text-gray-500 py-1">
        <span className="font-semibold pr-2">Pinned:</span>
        {isPinned ? "true" : "false"}
      </p>
      <p className="text-gray-500 py-1">
        <span className="font-semibold pr-2">Public: </span>
        {isPublic ? "true" : "false"}
      </p>
      <p className="text-gray-500 py-1 col-span-2">
        <span className="font-semibold pr-2">Metadata:</span>
        <br></br>
        {metaData ? metaData : "N/A"}
      </p>
    </div>
  );
}
