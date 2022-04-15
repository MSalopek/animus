import Footer from "./Footer";
import Sidebar from "./Sidebar";

export default function Layout({ children }) {
  return (
    // TODO: just use a grid for layout :)
    <div className="flex h-screen bg-gray-50 dark:bg-gray-900">
      <Sidebar />
      <div className="flex flex-col flex-1">
        <main>{children}</main>

        <Footer />
      </div>
    </div>
  );
}
