import axios from "axios";
// import { ApiError } from "next/dist/server/api-utils";
import { getToken, setToken } from "./util";

const API = axios.create({
	baseURL: process.env.API_URL,
	headers: {"Content-type": "application/json"},
})

API.interceptors.request.use(
	(request) => {
	  const token = getToken;
	  const auth = token ? `Bearer ${token}` : '';
	  request.headers.common['Authorization'] = auth;
	  return request;
	},
	(error) => Promise.reject(error),
  );

export async function GetUserStorage(params) {
	const res = await API({
		method: 'get',
		url: `${API.baseURL}/auth/storage/user`,
		params: params,
	})
	return res.data
}

export async function GetLsCid(cid) {
	const res = await API({
		method: 'get',
		url: `${API.baseURL}/auth/storage/ls/${cid}`,
	})
	return res.data
}

export async function UploadFile(file) {
	const formData = new FormData();
	formData.append("file", file);
	const res = API({
		method: 'post',
		url: `${API_URL}/auth/storage/add`,
		headers: {
			"Content-Type": "multipart/form-data",
		},
		data: formData,

	}).catch(err => {
		console.log("# ERR #", err)
		return
	})
	return res.data
}

// public endpoints

export async function Login(email, password) {
	const res = await API.post(`${API.baseURL}/api/login`, {email, password});
	if (res.status === 200) {
		setToken(res.data.token);
	}
}

export async function Register(params) {
	const res = await API.post(`${API.baseURL}/api/register`, params);
	return res.data;
}

export async function Logout() {
	localStorage.removeItem("user");
}
