import axios from 'axios';
// import { ApiError } from "next/dist/server/api-utils";
// import { getToken, setToken } from './util';

export const API_URL = process.env.NEXT_API_URL;

const defaultLimit = 25;
const defaultOffset = 0;

const API = axios.create({
  baseURL: API_URL,
  headers: { 'Content-type': 'application/json' },
});

// API.interceptors.request.use(
//   (request) => {
//     const token = getToken();
//     const auth = token ? `Bearer ${token}` : '';
//     request.headers.common['Authorization'] = auth;
//     return request;
//   },
//   (error) => {
//     Promise.reject(error);
//   }
// );

export async function GetUserStorage(token, limit, offset) {
  const o = offset ? offset : defaultOffset;
  const l = limit ? limit : defaultLimit;

  return await API.get(`/auth/user/storage?limit=${l}&offset=${o}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
}

export async function UploadFile(file) {
  const formData = new FormData();
  formData.append('file', file, file.name);
  return API({
    method: 'post',
    url: `${API_URL}/auth/user/storage/add-file`,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data: formData,
  });
}

export async function CreateKey() {
  const formData = new FormData();
  formData.append('file', file, file.name);
  return API({
    method: 'post',
    url: `${API_URL}/auth/user/keys`,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data: formData,
  });
}

export async function GetKeys() {
  return API({
    method: 'get',
    url: `${API_URL}/auth/user/keys`,
  });
}

export async function DeleteKey(id) {
  return API({
    method: 'delete',
    url: `${API_URL}/auth/user/keys/id/${id}`,
  });
}

export async function DisableKey(id) {
  return API({
    method: 'patch',
    url: `${API_URL}/auth/user/keys/id/${id}`,
    data: {
      disabled: true,
    },
  });
}

export async function EnableKey(id) {
  return API({
    method: 'patch',
    url: `${API_URL}/auth/user/keys/id/${id}`,
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

// export async function Logout() {
//   localStorage.removeItem('user');
// }
