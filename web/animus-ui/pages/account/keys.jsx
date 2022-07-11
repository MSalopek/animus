import Layout from "../../components/layout/Layout";
import Main from "../../components/layout/Main";
import KeysView from "../../components/view/KeysView";

export default KeysManager;

function KeysManager() {
  return (
    <Layout>
      <Main>
        <KeysView />
      </Main>
    </Layout>
  );
}

