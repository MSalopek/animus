export function isEmailValid(email) {
  return /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/.test(
    email
  );
}

export function cleanupString(name) {
  return name.charAt(0)?.toUpperCase() + name.slice(1).replaceAll("-", " ");
}

export function slugify(string) {
  return string
    .toString()
    .trim()
    .toLowerCase()
    .replace(/\s+/g, "-")
    .replace(/[^\w-]+/g, "")
    .replace(/--+/g, "-")
    .replace(/^-+/, "")
    .replace(/-+$/, "");
}


export function unslugify(string) {
  return string
    .toString()
    .trim()
    .replace("__", " ")
    .replace("--", " ")
    .replace("-", " ")
    .replace("_", " ");
}

export function classNames(...classes) {
  return classes.filter(Boolean).join(" ");
}

export function truncElipsis(str, len) {
  return `${str.substring(0, len)} ... `;
}
