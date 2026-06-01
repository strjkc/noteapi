import requests


class HttpClient:
    def __init__(self, base_url):
        self.base_url = base_url

    def get(self, path, token=None):
        try:
            url = f"{self.base_url}{path}"
            headers = {}
            if token:
                headers["Authorization"] = f"Bearer {token}"
            resp = requests.get(url, headers=headers if headers else None)
            return resp.status_code, resp.json()
        except requests.JSONDecodeError as e:
            print(f"HttpClient: Error decoding body - {e}")
            return resp.status_code, None
        except Exception as e:
            print(f"HttpClient: An Exception occured: {e}")
            print("Terminating run")
            exit()

    def post_json(self, path, data, token=None):
        try:
            url = f"{self.base_url}{path}"
            headers = {}
            if token:
                headers["Authorization"] = f"Bearer {token}"
            resp = requests.post(url, json=data, headers=headers if headers else None)
            return resp.status_code, resp.json()
        except requests.JSONDecodeError as e:
            print(f"HttpClient: Error decoding body - {e}")
            return resp.status_code, None
        except Exception as e:
            print(f"HttpClient: An Exception occured: {e}")
            print("Terminating run")
            exit()

    def post_file(self, path, data, token=None):
        try:
            url = f"{self.base_url}{path}"
            headers = {}
            if token:
                headers["Authorization"] = f"Bearer {token}"
            resp = requests.post(
                url, files={"file": data}, headers=headers if headers else None
            )
            return resp.status_code, resp.json()
        except requests.JSONDecodeError as e:
            print(f"HttpClient: Error decoding body - {e}")
            return resp.status_code, None
        except Exception as e:
            print(f"HttpClient: An Exception occured: {e}")
            print("Terminating run")
            exit()

    def delete(self, path, token=None):
        try:
            url = f"{self.base_url}{path}"
            headers = {}
            if token:
                headers["Authorization"] = f"Bearer {token}"
            resp = requests.delete(url, headers=headers if headers else None)
            return resp.status_code, resp.json()
        except requests.JSONDecodeError as e:
            print(f"HttpClient: Error decoding body - {e}")
            return resp.status_code, None
        except Exception as e:
            print(f"HttpClient: An Exception occured: {e}")
            print("Terminating run")
            exit()

    def get(self, path, token=None):
        try:
            url = f"{self.base_url}{path}"
            headers = {}
            if token:
                headers["Authorization"] = f"Bearer {token}"
            resp = requests.get(url, headers=headers if headers else None)
            body = None
            if resp.headers["Content-Type"] == "application/json":
                body = resp.json()
            elif resp.headers["Content-Type"] == "text/html":
                body = resp.text
            return resp.status_code, body
        except requests.JSONDecodeError as e:
            print(f"HttpClient: Error decoding body - {e}")
            return resp.status_code, None
        except Exception as e:
            print(f"HttpClient: An Exception occured: {e}")
            print("Terminating run")
            exit()
