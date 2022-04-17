import Layout from "../../components/layout/Layout";
import Main from "../../components/layout/Main";
import StorageView from "../../components/view/StorageView";

export default StorageManager;

function StorageManager() {
  return (
    <Layout>
      <Main>
        <StorageView />
      </Main>
    </Layout>
  );
}
