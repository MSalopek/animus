import { useState } from "react"
import FileUploadModal from "../../components/modals/FileUploadModal";

export default Modal;

function Modal() {
  const [isFileModalOpen, setIsFileModalOpen] = useState(false);
  return (
    <main>
      <button onClick={() => setIsFileModalOpen(!isFileModalOpen)}>BUTTON</button>
      <FileUploadModal
        isOpen={isFileModalOpen}
        setIsOpen={setIsFileModalOpen}
      />
    </main>
  );
}
