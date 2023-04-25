import unittest
from heisenberg import Heisenberg

class TestHeisenberg(unittest.TestCase):
    def setUp(self) -> None:
        self.h: Heisenberg = Heisenberg(api_key="no api key system yet")
        
    def test_heisenberg(self):
        self.h.new_collection("col", 3, 100, "cosine", 50, 100)
        self.h.put("a", [1, 2, 3], {}, "col")
        res = self.h.get("a", "col")
        print(res)
    