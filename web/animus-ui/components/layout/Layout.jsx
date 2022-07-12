import Head from "next/head";
import Footer from "./Footer";
import Sidebar from "./Sidebar";

export default function Layout({ children }) {
  return (
    <div className="flex bg-gray-50 dark:bg-gray-900">
      <Head>
        <title>Animus Store</title>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="description" content="Decentralized storage made easy. Manage files, seamlessly access IPFS and regain control over your private data." />
      </Head>
      <Sidebar />
      <div className="flex flex-col flex-1">
        {children}
        <Footer />
      </div>
    </div>
  );
}
