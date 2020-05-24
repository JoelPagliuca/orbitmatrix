#!/usr/bin/env python3

import requests, subprocess, time, json, threading

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
PROC: subprocess.Popen = None
SERVER_LOGS = [[]]
FUNC_NUMBER = 0

def showResponse(res: requests.Response):
	print(f"{cols.FAIL}[-] {res.url[len(BASE):]} {res.status_code}")
	if len(res.content):
		print(f"{cols.WARNING}[*] Response body")
		try:
			print(f"{json.dumps(res.json(), indent=2)}", end=cols.ENDC)
		except:
			print(f"{str(res.content, 'utf-8')}", end=cols.ENDC)
	logs = SERVER_LOGS[FUNC_NUMBER-1]
	print()
	if len(logs):
		print(f"{cols.WARNING}[*] Server logs")
		print("".join(logs))
	print()

def do(name: str, res: requests.Response, broke: bool):
	global FUNC_NUMBER
	FUNC_NUMBER += 1
	print(f"[*] {name}", end="")
	rid = res.headers.get("X-Request-ID", "(no ID)")
	if broke(res):
		print(f" {cols.FAIL}FAILED{cols.ENDC} {rid}")
		showResponse(res)
	else:
		print(f" {cols.OKGREEN}PASS{cols.ENDC}")

def serverLogs():
	for line in PROC.stdout:
		if len(SERVER_LOGS) >= FUNC_NUMBER:
			SERVER_LOGS.append([])
		SERVER_LOGS[FUNC_NUMBER].append(str(line, "utf-8"))

def startServer():
	print(f"{cols.BOLD}[*] starting the api{cols.ENDC}")
	return subprocess.Popen(
		"exec ./caliban",
		stderr=subprocess.STDOUT,
		stdout=subprocess.PIPE,
		shell=True,
		bufsize=1,
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
	do("Do the wrong method",
		sess.delete(f"{BASE}/health"),
		lambda res: res.status_code != 405
	)
	do("Do something unauthed",
		sess.get(f"{BASE}/item"),
		lambda res: res.status_code != 401
	)
	res = sess.post(f"{BASE}/user/register", data='{"name": "Joel"}')
	token = res.json()["Token"]
	print(f"[+] Got token {token}")
	sess.headers.update({"Authorization": f"Bearer {token}"})
	do("GET /user/me",
		sess.get(f"{BASE}/user/me"),
		lambda res: res.json()["ID"] == 0 or res.json()["Name"] == ""
	)
	do("POST /item/add",
		sess.post(f"{BASE}/item/add", data='{"name":"item1", "description":"item 1"}'),
		lambda res: res.json()["ID"] == "" or res.json()["Name"] == ""
	)
	do("GET /item",
		sess.get(f"{BASE}/item"),
		lambda res: len(res.json()) == 0
	)
	do("POST /group/add",
		sess.post(f"{BASE}/group/add", data='{"name":"group1", "description":"group 1"}'),
		lambda res: res.json()["ID"] == "" or res.json()["Name"] == ""
	)
	do("GET /group",
		sess.get(f"{BASE}/group"),
		lambda res: len(res.json()) == 0
	)
	do("OPTIONS /health",
		sess.options(f"{BASE}/health"),
		lambda res: res.status_code != 204 or res.headers["Allow"] != "GET"
	)
	sess.headers.update({"Authorization": "Bearer abadtoken"})
	do("Do something with a bad token",
		sess.get(f"{BASE}/item"),
		lambda res: res.status_code != 401
	)

def main():
	global PROC
	caliban = startServer()
	PROC = caliban
	time.sleep(0.1)
	thread = threading.Thread(target=serverLogs)
	thread.start()
	try:
		tests()
	except Exception as e:
		print(f"{cols.FAIL}\n[-] Error running tests: {e}{cols.ENDC}")
	killServer(caliban)

if __name__ == "__main__":
	main()
