"use strict";

const api = async (method, path, body) => {
  const res = await fetch(path, {
    method: method,
    body: body ?? JSON.stringify(body)
  });
  return await res.json();
};

document.addEventListener("DOMContentLoaded", () => {
  document.getElementById("login").addEventListener("click", async () => {
    const res = await api("POST", "/api/user/login", {
      user_id: document.getElementById("user_id").value,
      password: document.getElementById("password").value
    });
    console.log(res);
  });
  document.getElementById("logout").addEventListener("click", async () => {
    const res = await api("POST", "/api/user/logout");
    console.log(res);
  });
  document.getElementById("me").addEventListener("click", async () => {
    const res = await api("GET", "/api/user/me");
    console.log(res);
  });
});
