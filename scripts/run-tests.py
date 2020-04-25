#!/usr/bin/env python3

import requests, subprocess, time, json

BASE = "http://127.0.0.1:7080"

def showResponse(res: requests.Response):
	print(f"{res.url[len(BASE):]} {res.status_code}", end="")
	if len(res.content):
		print(f": {json.dumps(res.json(), indent=2)}", end="")
	print()

def startServer():
	return subprocess.Popen("exec ./caliban", stdout=subprocess.DEVNULL, shell=True)

def killServer(s: subprocess.Popen):
	s.kill()

def tests():
	showResponse(requests.get(f"{BASE}/health"))
	res = requests.post(f"{BASE}/user/register", data='{"name": "Joel"}')
	showResponse(res)
	showResponse(requests.get(f"{BASE}/user/me"))
	showResponse(requests.post(f"{BASE}/item/add", data='{"name":"item1", "description":"item 1"}'))
	showResponse(requests.get(f"{BASE}/item"))

def main():
	caliban = startServer()
	try:
		tests()
	except Exception as e:
		print(f"Error running tests: {e}")
	killServer(caliban)

if __name__ == "__main__":
	main()
