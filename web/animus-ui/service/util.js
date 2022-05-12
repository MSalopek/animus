export function getToken(){
	const u = JSON.parse(localStorage.getItem('user'))
	if (u) {
		return u.token;
	}
	return "";
}

export function setToken(token) {
	localStorage.setItem("user", {"token": token});
}
