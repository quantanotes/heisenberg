import requests

class Heisenberg:
    def __init__(self, api_key: str) -> None:
        self.url = 'http://localhost:420'
        self.api_key = api_key
        
    def new_collection(self, name: str, dim: int, size: int, space: str, m: int, ef: int) -> None:
        data = {
            'name': name,
            'dim': dim,
            'size': size,
            'space': space,
            'm': m,
            'ef': ef,
        }
        headers = {'Content-Type': 'application/json', 'api_key': self.api_key}
        r = requests.post(f'{self.url}/newcollection/', json=data, headers=headers)
        r.raise_for_status()
    
    def delete_collection(self, name: str) -> None:
        data = {
            'name': name,
        }
        headers = {'Content-Type': 'application/json', 'api_key': self.api_key}
        r = requests.post(f'{self.url}/delcollection/', json=data, headers=headers)
        r.raise_for_status()
        
    def get(self, key: str, collection: str) -> dict:
        data = {
            'key': key,
            'collection': collection,
        }
        headers = {'Content-Type': 'application/json', 'api_key': self.api_key}
        r = requests.post(f'{self.url}/get/', json=data, headers=headers)
        r.raise_for_status()
        return r.json()['value']
    
    def put(self, key: str, vec: list[float], meta: dict, collection: str) -> None:
        data = {
            'key': key,
            'vec': vec,
            'meta': meta,
            'collection': collection,
        }
        headers = {'Content-Type': 'application/json', 'api_key': self.api_key}
        r = requests.post(f'{self.url}/put/', json=data, headers=headers)
        r.raise_for_status()
    
    def delete(self, key: str, collection: str) -> None:
        data = {
            'key': key,
            'collection': collection,
        }
        headers = {'Content-Type': 'application/json', 'api_key': self.api_key}
        r = requests.post(f'{self.url}/del/', json=data, headers=headers)
        r.raise_for_status()
        
    def search(self, vec: list[float], k: int, collection: str) -> dict:
        data = {
            'vec': vec,
            'k': k,
            'collection': collection,
        }
        headers = {'Content-Type': 'application/json', 'api_key': self.api_key}
        r = requests.post(f'{self.url}/search/', json=data, headers=headers)
        r.raise_for_status()
        return r.json()['value']