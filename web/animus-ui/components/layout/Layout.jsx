import Footer from "./Footer";
import Main from "./Main";
import Sidebar from "./Sidebar";

export default function Layout({ children }) {
  return (
    <div className="flex bg-gray-50 dark:bg-gray-900">
      <Sidebar />
      <div className="flex flex-col flex-1">
        {children}
        <Footer/>
      </div>
    </div>
  );
}
