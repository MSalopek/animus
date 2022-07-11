import { useState } from 'react';
import { CopyToClipboard } from 'react-copy-to-clipboard';
import { FolderAddIcon } from '@heroicons/react/solid';
import {
  TrashIcon,
  BanIcon,
  // ArrowsExpandIcon,
  DuplicateIcon,
  LightningBoltIcon,
  KeyIcon,
} from '@heroicons/react/outline';

import { ModalBtn, BtnContainer, RoundActionBtn } from '../buttons/Buttons';

export default KeysView;

function KeysView({ keys }) {
  const [isKeyModalOpen, setIsKeyModalOpen] = useState(false);

  return (
    <section className="bg-white dark:bg-gray-900">
      <div className="container px-6 py-12 mx-auto">
        <div className="flex flex-col">
          <h1 className="text-3xl font-semibold text-gray-800 dark:text-white mb-2">
            API Keys
          </h1>
          <BtnContainer>
            <ModalBtn
              Icon={KeyIcon}
              title={'Add Key'}
              action={() => {
                setIsKeyModalOpen(true);
              }}
            />
          </BtnContainer>
        </div>

        <div className="mt-4 space-y-4 lg:mt-8 border p-4 rounded-lg">
          <KeyRow disabled={true} />
          <KeyRow disabled={false} />
        </div>
      </div>
    </section>
  );
}

function KeyRow({ key, created_at, disabled, rights }) {
  return (
    <div className="py-2 px-8 bg-gray-100 rounded-lg dark:bg-gray-800">
      <div className="flex items-center justify-between w-full">
        <div>
          <span className="text-sm text-gray-400">Client Key:</span>
          <h1
            className={`font-semibold ${
              disabled ? 'text-gray-400' : 'text-gray-700'
            } dark:text-white`}
          >
            xWQWPOEK123saasdawdpow
          </h1>
        </div>

        <div>
          <span className="text-sm text-gray-400">Access Rights</span>
          <h1
            className={`font-semibold text-sm ${
              disabled ? 'text-gray-400' : 'text-gray-700'
            } dark:text-white`}
          >
            Read-write-delete
          </h1>
        </div>
        <div>
          <span className="text-sm text-gray-400">Status</span>
          <h1
            className={`text-sm font-bold ${
              disabled ? 'text-red-500' : 'text-green-700'
            }`}
          >
            {disabled ? 'Disabled' : 'Enabled'}
          </h1>
        </div>

        <BtnContainer>
          <CopyToClipboard text={'xWQWPOEK123saasdawdpow'}>
            <RoundActionBtn Icon={DuplicateIcon} />
          </CopyToClipboard>
          <RoundActionBtn Icon={disabled ? LightningBoltIcon : BanIcon} />
          <RoundActionBtn Icon={TrashIcon} />
          {/* <RoundActionBtn Icon={ArrowsExpandIcon} /> */}
        </BtnContainer>
      </div>
    </div>
  );
}
