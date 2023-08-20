import logging
import requests

class Heisenberg:
    def __init__(self, host, port, api_key):
        self.base_url = 'http://' + host + port
        self.api_key = api_key
    
    def _request(self, method, endpoint, data=None):
        headers = {
            'Content-Type': 'application/json',
            'X-API-Key': self.api_key,
        }
        response = requests.request(
            method,
            f'{self.base_url}/{endpoint}',
            headers=headers,
            json=data,
            timeout=10
        )
        if response.status_code == 401:
            logging.error("Invalid API key")
            return None
        elif response.status_code != 200:
            logging.error(f"Error {response.status_code}: {response.text}")
            return None
        return response.json()
    
    def new_bucket(self, name, dim, space):
        data = {
            'name': name,
            'dim': dim,
            'space': space
        }
        self._request('POST', 'newbucket', data)
    
    def delete_bucket(self, name):
        data = {
            'name': name
        }
        self._request('POST', 'deletebucket', data)
    
    def get(self, bucket, key):
        data = {
            'bucket': bucket,
            'key': key
        }
        return self._request('POST', 'get', data)

    def put(self, bucket, key, vector, meta=None):
        data = {
            'bucket': bucket,
            'key': key,
            'vector': vector,
            'meta': meta or {}
        }
        self._request('POST', 'put', data)
    
    def delete(self, bucket, key):
        data = {
            'bucket': bucket,
            'key': key
        }
        self._request('POST', 'delete', data)
    
    def search(self, bucket, query, k):
        data = {
            'bucket': bucket,
            'query': query,
            'k': k
        }
        return self._request('POST', 'search', data)