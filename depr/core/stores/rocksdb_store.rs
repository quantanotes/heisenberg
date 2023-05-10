use crate::core::{
    model::*,
    store::{Store, StoreError, Transaction, TransactionError},
};
use rocksdb::{Options, TransactionDB};

pub struct RocksDBStore {
    db: TransactionDB,
}

impl<'a> Store<RocksDBTransaction<'a>> for RocksDBStore {
    fn new(path: &str) -> Result<Self, StoreError> {
        TransactionDB::open_default(path)
            .map(|db| Self { db })
            .map_err(Into::into)
    }

    fn open(path: &str) -> Result<Self, StoreError> {
        todo!()
    }

    fn new_collection(&mut self, name: String) -> Result<(), StoreError> {
        self.db
            .create_cf(name, &Options::default())
            .map_err(Into::into)
    }

    fn delete_collection(&mut self, name: String) -> Result<(), StoreError> {
        self.db.drop_cf(&name).map_err(Into::into)
    }

    fn get(&self, collection: String, key: String) -> Result<Option<Vec<u8>>, StoreError> {
        self.db.get_cf(&collection, key).map_err(Into::into)
    }

    fn put(
        &mut self,
        collection: String,
        key: String,
        value: Vec<u8>,
    ) -> Result<RocksDBTransaction, StoreError> {
        let tx = self.db.transaction();

        tx.put_cf(&collection, key, value)?;

        tx.map(|tx| RocksDBTransaction::new(tx)).map_err(Into::into)
    }

    fn delete(
        &mut self,
        collection: String,
        key: String,
    ) -> Result<RocksDBTransaction, StoreError> {
        self.db
            .transaction()
            .delete_cf(&collection, key)
            .map(|tx| RocksDBTransaction::new(tx))
            .map_err(Into::into)
    }
}

pub struct RocksDBTransaction<'a> {
    tx: rocksdb::Transaction<'a, TransactionDB>,
}

impl<'a> Transaction<rocksdb::Transaction<'a, TransactionDB>> for RocksDBTransaction<'a> {
    fn new(tx: rocksdb::Transaction<'a, TransactionDB>) -> Self {
        Self { tx }
    }

    fn commit(self: &mut RocksDBTransaction<'a>) -> Result<(), TransactionError> {
        self.tx.commit().map_err(Into::into)
    }

    fn rollback(&self) -> Result<(), TransactionError> {
        self.tx.rollback().map_err(Into::into)
    }
}
