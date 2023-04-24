import unittest
from heisenberg import Heisenberg

class TestHeisenberg(unittest.TestCase):
    def setUp(self) -> None:
        self.h = Heisenberg(api_key="no api key system yet")
        
    def test_heisenberg(self):
        pass
    