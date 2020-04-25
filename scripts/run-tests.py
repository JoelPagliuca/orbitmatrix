#!/usr/bin/env python3

import requests, subprocess, time, json

BASE = "http://127.0.0.1:7080"

def showResponse(res: requests.Response):
	print(f"[-] {res.url[len(BASE):]} {res.status_code}", end="")
	if len(res.content):
		print(f": {json.dumps(res.json(), indent=2)}", end="")
	print()

def do(name: str, res: requests.Response, broke: bool):
	print(f"[*] {name}", end="")
	if broke(res):
		print(" FAILED")
		showResponse(res)
	else:
		print(" PASS")

def startServer():
	print(f"[*] starting the server")
	return subprocess.Popen(
		"exec ./caliban", 
		stderr=subprocess.DEVNULL, 
		stdout=subprocess.DEVNULL, 
		shell=True
	)

def killServer(s: subprocess.Popen):
	print(f"[*] killing the api")
	s.kill()

def tests():
	sess = requests.Session()
	do("GET /health",
		sess.get(f"{BASE}/health"),
		lambda res: res.status_code != 204
	)
	res = sess.post(f"{BASE}/user/register", data='{"name": "Joel"}')
	token = res.json()["Token"]
	print(f"[+] Got token {token}")
	sess.headers.update({"Authorization": f"Bearer {token}"})
	do("GET /user/me",
		sess.get(f"{BASE}/user/me"),
		lambda res: res.json()["ID"] == "" or res.json()["Name"] == ""
	)
	do("POST /item/add",
		sess.post(f"{BASE}/item/add", data='{"name":"item1", "description":"item 1"}'),
		lambda res: res.json()["ID"] == "" or res.json()["Name"] == ""
	)
	do("GET /item",
		sess.get(f"{BASE}/item"),
		lambda res: len(res.json()) == 0
	)

def main():
	caliban = startServer()
	try:
		tests()
	except Exception as e:
		print(f"[-] Error running tests: {e}")
	killServer(caliban)

if __name__ == "__main__":
	main()
