import { unstable_getServerSession } from 'next-auth/next';
import { authOptions } from './api/auth/[...nextauth]';

import Layout from '../components/layout/Layout';
import Main from '../components/layout/Main';
import StorageView from '../components/view/StorageView';
import { GetUserStorage } from '../service/http';

const PAGE_SIZE = 10;

export default StorageManager;

function StorageManager({ rows, total, pages, currentPage }) {
  return (
    <Layout>
      <Main>
        <StorageView
          rows={rows}
          total={total}
          pages={pages}
          currentPage={currentPage}
        />
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
    let res;
    let { page } = context.query;
    if (page) {
      page = parseInt(page);
      // use limit and offset to get page
      res = await GetUserStorage(
        session.user.accessToken,
        PAGE_SIZE,
        (page - 1) * PAGE_SIZE
      );
    } else {
      res = await GetUserStorage(session.user.accessToken, 10, 0);
    }
    const data = res.data;

    return {
      props: {
        total: data.total,
        pages: Math.ceil(data.total / PAGE_SIZE),
        currentPage: page ? page : 1,
        rows: data.rows,
      },
    };
  } catch (error) {
    console.log('error getting Storage props', error);
    return {
      props: {
        total: 0,
        pages: 1,
        currentPage: 1,
        rows: [],
      },
    };
  }
}
