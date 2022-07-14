import axios from 'axios';
import { format } from 'prettier';
// import { ApiError } from "next/dist/server/api-utils";

export const API_URL = process.env.NEXT_API_URL;

const defaultLimit = 10;
const defaultOffset = 0;

const API = axios.create({
  baseURL: API_URL,
  headers: { 'Content-type': 'application/json' },
});

export async function GetUserStorage(token, limit, offset) {
  const o = offset ? offset : defaultOffset;
  const l = limit ? limit : defaultLimit;

  return await API.get(`/auth/user/storage?limit=${l}&offset=${o}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
}

export async function UploadFile(token, file) {
  const formData = new FormData();
  formData.append('file', file, file.name);
  return API({
    method: 'post',
    url: 'auth/user/storage/add-file',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'multipart/form-data',
    },
    data: formData,
  });
}

export async function UploadDirectory(token, files, dirname) {
  const formData = new FormData();
  formData.append('name', dirname);
  files.forEach((f) => {
    // NOTE: sending f.path can be used to preserve dir structure on BE
    formData.append('files', f, f.path);
  });

  return API({
    method: 'post',
    url: 'auth/user/storage/add-dir',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'multipart/form-data',
    },
    data: formData,
  });
}

export async function DeleteStorage(token, id) {
  return API({
    method: 'delete',
    url: `/auth/user/storage/id/${id}`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
}

export async function CreateKey(token) {
  return await API({
    method: 'post',
    url: '/auth/user/keys',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
}

export async function GetKeys(token) {
  return API.get('/auth/user/keys', {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
}

export async function DeleteKey(token, id) {
  return API({
    method: 'delete',
    url: `/auth/user/keys/id/${id}`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
}

export async function DisableKey(token, id) {
  return await API({
    method: 'patch',
    url: `/auth/user/keys/id/${id}`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    data: {
      disabled: true,
    },
  });
}

export async function EnableKey(token, id) {
  return await API({
    method: 'patch',
    url: `/auth/user/keys/id/${id}`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    data: {
      disabled: false,
    },
  });
}

// public endpoints

export async function LoginCredentials(email, password) {
  return await API.post('/login', { email, password });
}

export async function Register(params) {
  const res = await axios.post(`${API_URL}/register`, params);
  return res.data;
}

export async function ActivateUser(email, token) {
  return await axios.post(`${API_URL}/activate/email/${email}?token=${token}`);
}

// export async function Logout() {
//   localStorage.removeItem('user');
// }
