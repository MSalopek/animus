export default function Main({ children }) {
  // use "flex-grow" to always expand vertically
  // this way footer wil stay stuck at the bottom
  return <main className="flex-grow">{children}</main>;
}
