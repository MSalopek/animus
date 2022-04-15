import Footer from "../../components/layout/Footer";
import Layout from "../../components/layout/Layout";
import Main from "../../components/layout/Main";

export default ChildLayout;

function ChildLayout() {
  return (
    <Layout>
      <Main />
      <Main />
    </Layout>
  );
}
