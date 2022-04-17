import Head from "next/head";
import Footer from "./Footer";
import Sidebar from "./Sidebar";

export default function Layout({ children }) {
  return (
    <div className="flex bg-gray-50 dark:bg-gray-900">
      <Head>
        <title>Animus UI</title>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="description" content="The engine that keeps your files forever" />
      </Head>
      <Sidebar />
      <div className="flex flex-col flex-1">
        {children}
        <Footer />
      </div>
    </div>
  );
}
