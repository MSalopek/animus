import { authOptions } from '../api/auth/[...nextauth]';
import { unstable_getServerSession } from 'next-auth/next';


import Layout from '../../components/layout/Layout';
import Main from '../../components/layout/Main';
import KeysView from '../../components/view/KeysView';

import { GetKeys } from '../../service/http';

export default KeysManager;

function KeysManager({ rows }) {
  return (
    <Layout>
      <Main>
        <KeysView rows={rows}/>
      </Main>
    </Layout>
  );
}

export async function getServerSideProps(context) {
  // TODO: maybe change usage of unstable_ function if next-auth changes it
  const session = await unstable_getServerSession(
    context.req,
    context.res,
    authOptions
  );

  if (!session) {
    return {
      redirect: {
        destination: '/account/login',
        permanent: false,
      },
    };
  }

  try {
    const res = await GetKeys(session.user.accessToken);
    // put disabled keys after all enabled keys
    const data = res.data.sort((a, b) =>  a.disabled - b.disabled );

    return {
      props: {
        rows: data,
      },
    };
  } catch (error) {
    console.log("error getting Keys props", error);
    return {
      props: {
        rows: [],
      },
    };
  }
}
