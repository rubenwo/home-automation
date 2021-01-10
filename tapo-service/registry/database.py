import json

import redis


class Database:
    def __init__(self):
        self.client = redis.Redis(host='redis.default.svc.cluster.local', port=6379, db=0)

    def insert(self, key, value):
        print("Database->insert()")
        self.client.set(key, json.dumps(value).encode("utf-8"))

    def retrieve(self, key):
        print("Database->retrieve()")
        try:
            value = self.client.get(key)
            data = json.loads(value.decode("utf-8"))
            return data
        except Exception as e:
            raise Exception(e)
