import { useEffect, useState } from 'react';
import { useSession } from 'next-auth/react';

import { CopyToClipboard } from 'react-copy-to-clipboard';
import {
  TrashIcon,
  BanIcon,
  // ArrowsExpandIcon,
  DuplicateIcon,
  LightningBoltIcon,
  KeyIcon,
} from '@heroicons/react/outline';

import { ModalBtn, BtnContainer, RoundActionBtn } from '../buttons/Buttons';
import {
  CreateKey,
  DeleteKey,
  DisableKey,
  EnableKey,
} from '../../service/http';
import CreateKeyModal from '../modals/CreateKeyModal';

const accessRights = {
  r: 'Read',
  rw: 'Read, Write',
  rwd: 'Read, Write, Delete',
};

export default KeysView;

function KeysView({ rows }) {
  const { data: session, status } = useSession();

  const [keys, setKeys] = useState(rows);

  const [isKeyModalOpen, setIsKeyModalOpen] = useState(false);

  const disableKey = async (id) => {
    if (status !== 'authenticated') {
      return;
    }

    const res = await DisableKey(session.user.accessToken, id);

    if (res.data) {
      const newKeys = [...keys];
      const replaceIdx = keys.findIndex((elem) => elem.id === res.data.id);
      newKeys[replaceIdx] = res.data;
      setKeys(newKeys);
    }
  };

  const enableKey = async (id) => {
    if (status !== 'authenticated') {
      return;
    }

    const res = await EnableKey(session.user.accessToken, id);
    if (res.data) {
      const newKeys = [...keys];
      const replaceIdx = keys.findIndex((elem) => elem.id === res.data.id);
      newKeys[replaceIdx] = res.data;
      setKeys(newKeys);
    }
  };

  const deleteKey = async (id) => {
    if (status !== 'authenticated') {
      return;
    }

    const res = await DeleteKey(session.user.accessToken, id);

    if (res.status === 204) {
      setKeys((prev) => prev.filter((k) => k.id !== id));
    }
  };

  const createKey = async () => {
    return await CreateKey(session.user.accessToken);
  };

  return (
    <section className="bg-white dark:bg-gray-900">
      <CreateKeyModal
        isOpen={isKeyModalOpen}
        setIsOpen={setIsKeyModalOpen}
        currentKeys={keys}
        setCurrentKeys={setKeys}
        createFunc={createKey}
      />
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

        {keys && keys.length ? (
          <div className="mt-4 space-y-4 lg:mt-8 border p-4 rounded-lg">
            {keys.map((k) => (
              <KeyRow
                key={k.client_key}
                client_key={k.client_key}
                rights={k.rights}
                disabled={k.disabled}
                id={k.id}
                disableFunc={disableKey}
                enableFunc={enableKey}
                deleteFunc={deleteKey}
              />
            ))}
          </div>
        ) : (
          <></>
        )}
      </div>
    </section>
  );
}

function KeyRow({
  id,
  client_key,
  disabled,
  rights,
  disableFunc,
  enableFunc,
  deleteFunc,
}) {
  return (
    <div className="py-2 px-8 bg-gray-100 rounded-lg dark:bg-gray-800">
      <div className="flex items-center justify-between">
        <div className="grid grid-cols-4 w-4/5">
          <div className="col-span-2">
            <span className="text-sm text-gray-400">Client Key:</span>
            <h1
              className={`font-semibold ${
                disabled ? 'text-gray-400' : 'text-gray-700'
              } dark:text-white`}
            >
              {client_key}
            </h1>
          </div>

          <div className="place-items-center">
            <span className="text-sm text-gray-400">Access Rights</span>
            <h1
              className={`font-semibold text-sm ${
                disabled ? 'text-gray-400' : 'text-gray-700'
              } dark:text-white`}
            >
              {accessRights[rights]}
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
        </div>

        <BtnContainer>
          <CopyToClipboard text={client_key}>
            <RoundActionBtn Icon={DuplicateIcon} />
          </CopyToClipboard>
          {disabled ? (
            <RoundActionBtn
              Icon={LightningBoltIcon}
              onClick={() => enableFunc(id)}
            />
          ) : (
            <RoundActionBtn Icon={BanIcon} onClick={() => disableFunc(id)} />
          )}
          <RoundActionBtn Icon={TrashIcon} onClick={() => deleteFunc(id)} />
          {/* <RoundActionBtn Icon={ArrowsExpandIcon} /> */}
        </BtnContainer>
      </div>
    </div>
  );
}
