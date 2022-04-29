export function getToken(){
	return JSON.parse(localStorage.getItem('user')).token || ''	
}

export function setToken(token) {
	localStorage.setItem("user", {"token": token});
}
