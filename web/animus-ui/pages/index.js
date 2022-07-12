import { unstable_getServerSession } from 'next-auth/next';
import { authOptions } from './api/auth/[...nextauth]';

import Layout from '../components/layout/Layout';
import Main from '../components/layout/Main';
import StorageView from '../components/view/StorageView';
import { GetUserStorage } from '../service/http';

export default StorageManager;

function StorageManager({ rows, total, pages }) {
  return (
    <Layout>
      <Main>
        <StorageView rows={rows} total={total} pages={pages}/>
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
    const res = await GetUserStorage(session.user.accessToken);
    const data = res.data;

    return {
      props: {
        total: data.total,
        pages: Math.ceil(data.total / data.rows.length),
        rows: data.rows,
      },
    };
  } catch (error) {
    console.log("error getting Storage props", error);
    return {
      props: {
        total: 1,
        pages: 1,
        rows: [],
      },
    };
  }
}
