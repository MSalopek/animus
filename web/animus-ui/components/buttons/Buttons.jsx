export { BtnContainer, ModalBtn, RoundActionBtn, RoundLinkBtn };

function BtnContainer({ children }) {
  return <div className="py-2 inline-flex gap-2">{children}</div>;
}

function ModalBtn({ Icon, title, action }) {
  return (
    <button
      onClick={action}
      className="flex items-center justify-center bg-blue-600 text-sm font-medium w-40 h-10 rounded text-blue-50 hover:bg-blue-700"
    >
      {Icon && <Icon className="w-5 h-5" />}
      <span className="px-2">{title}</span>
    </button>
  );
}

function RoundActionBtn({ Icon, onClick, padding }) {
  const p = padding ? `p-${padding}` : 'p-2';
  return (
    <button
      className={`${p} rounded-full border bg-gray-50 border-gray-200"`}
      onClick={() => onClick()}
    >
      {Icon && <Icon className="w-5 h-5" />}
    </button>
  );
}

function RoundLinkBtn({ Icon, href, padding }) {
  const p = padding ? `p-${padding}` : 'p-2';
  return (
    <a
      target="_blank"
      rel="noopener noreferrer"
      href={href}
      className={`${p} rounded-full border bg-gray-50 border-gray-200"`}
    >
      {Icon && <Icon className="w-5 h-5" />}
    </a>
  );
}
