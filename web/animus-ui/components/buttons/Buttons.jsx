export { BtnContainer, ActionBtn };

function BtnContainer({ children }) {
  return <div className="py-4 inline-flex gap-2">{children}</div>;
}

function ActionBtn({ Icon, title, onClick }) {
  console.log("### ICON", Icon, title);
  return (
    <button className="flex items-center justify-center bg-blue-600 text-sm font-medium w-40 h-10 rounded text-blue-50 hover:bg-blue-700">
      {Icon && <Icon className="w-5 h-5"/>}
      <span className="px-2">{title}</span>
    </button>
  );
}
