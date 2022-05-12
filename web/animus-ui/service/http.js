import axios from 'axios';
import { format } from 'prettier';
// import { ApiError } from "next/dist/server/api-utils";
import { getToken, setToken } from './util';


const API_URL = process.env.NEXT_PUBLIC_API_URL;

const API = axios.create({
  baseURL: API_URL,
  headers: { 'Content-type': 'application/json' },
});


API.interceptors.request.use(
  (request) => {
    const token = getToken();
    const auth = token ? `Bearer ${token}` : '';
    request.headers.common['Authorization'] = auth;
    return request;
  },
  (error) => {
    Promise.reject(error);
  }
);

export async function GetUserStorage(params) {
  const res = await API({
    method: 'get',
    url: `${API_URL}/auth/storage/user`,
    params: params,
  });
  return res.data;
}

export async function UploadFile(file) {
  const formData = new FormData();
  formData.append('file', file, file.name);
  return API({
    method: 'post',
    url:`${API_URL}/auth/storage/add-file`,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
    data: formData,
  });
}

// public endpoints

export async function Login(email, password) {
  const res = await API.post(`${API.baseURL}/api/login`, { email, password });
  if (res.status === 200) {
    setToken(res.data.token);
  }
}

export async function Register(params) {
  const res = await API.post(`${API.baseURL}/api/register`, params);
  return res.data;
}

export async function Logout() {
  localStorage.removeItem('user');
}
