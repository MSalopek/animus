import Layout from "../../components/layout/Layout";
import Main from "../../components/layout/Main";
import StorageView from "../../components/view/StorageView";

export default ChildLayout;

function ChildLayout() {
  return (
    <Layout>
      <Main>
        <StorageView />
      </Main>
    </Layout>
  );
}
