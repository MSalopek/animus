export default function Tooltip({ message, children }) {
  return (
    <div className="relative flex flex-col items-center group">
      {children}
      <div className="absolute bottom-2 flex flex-col items-center opacity-0 mb-6 group-hover:flex group-hover:opacity-70 transition ease-in-out delay-2000" >
        <span className="relative z-20 p-2 w-28 text-xs text-center leading-none text-white whitespace-no-wrap bg-gray-600 shadow-lg rounded-md">
          {message}
        </span>
        <div className="w-3 h-3 -mt-2 rotate-45 bg-gray-600"></div>
      </div>
    </div>
  );
};
