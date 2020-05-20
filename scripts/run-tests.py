#!/usr/bin/env python3

import requests, subprocess, time, json

class cols:
	HEADER = '\033[95m'
	OKBLUE = '\033[94m'
	OKGREEN = '\033[92m'
	WARNING = '\033[93m'
	FAIL = '\033[91m'
	ENDC = '\033[0m'
	BOLD = '\033[1m'
	UNDERLINE = '\033[4m'

BASE = "http://127.0.0.1:7080"

def showResponse(res: requests.Response):
	print(f"{cols.FAIL}[-] {res.url[len(BASE):]} {res.status_code}", end=cols.ENDC)
	if len(res.content):
		try:
			print(f": {json.dumps(res.json(), indent=2)}", end="")
		except:
			print(f": {str(res.content, 'utf-8')}", end="")
	print()

def do(name: str, res: requests.Response, broke: bool):
	print(f"[*] {name}", end="")
	if broke(res):
		print(f" {cols.FAIL}FAILED{cols.ENDC}")
		showResponse(res)
	else:
		print(f" {cols.OKGREEN}PASS{cols.ENDC}")

def startServer():
	print(f"{cols.BOLD}[*] starting the api{cols.ENDC}")
	return subprocess.Popen(
		"exec ./caliban", 
		stderr=subprocess.DEVNULL, 
		stdout=subprocess.DEVNULL, 
		shell=True
	)

def killServer(s: subprocess.Popen):
	print(f"{cols.BOLD}[*] killing the api{cols.ENDC}")
	s.kill()

def tests():
	sess = requests.Session()
	do("GET /health",
		sess.get(f"{BASE}/health"),
		lambda res: res.status_code != 200
	)
	do("GET /swagger.txt",
		sess.get(f"{BASE}/swagger.txt"),
		lambda res: len(res.content) < 204
	)
	do("Do something unauthed",
		sess.get(f"{BASE}/item"),
		lambda res: res.status_code != 401
	)
	do("Do the wrong method",
		sess.delete(f"{BASE}/health"),
		lambda res: res.status_code != 405
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
	do("OPTIONS /health",
		sess.options(f"{BASE}/health"),
		lambda res: res.status_code != 204 or res.headers["Allow"] != "GET"
	)

def main():
	caliban = startServer()
	try:
		tests()
	except Exception as e:
		print(f"{cols.FAIL}[-] Error running tests: {e}{cols.ENDC}")
	killServer(caliban)

if __name__ == "__main__":
	main()
