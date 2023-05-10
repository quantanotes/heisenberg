use std::path::PathBuf;

use crate::core::store::*;

pub struct RocksDBStore {
    store: rocksdb::TransactionDB,
}

impl RocksDBStore {
    fn get_cf(&self, name: String) -> Result<&rocksdb::ColumnFamily, StoreError> {
        self.store
            .cf_handle(&name)
            .ok_or(StoreError::CollectionNotFound(name))
    }
}

impl Store for RocksDBStore {
    type Transaction<'a> = RocksDBTransaction<'a>;

    fn new(path: &str) -> Result<RocksDBStore, StoreError> {
        let mut opts = rocksdb::Options::default();
        opts.create_missing_column_families(true);

        // Create a new store if column families don't exist
        if !PathBuf::from(path).join("CURRENT").exists() {
            let store = rocksdb::TransactionDB::open_default(path)
                .map_err(Into::<rocksdb::Error>::into)?;

            return Ok(RocksDBStore{ store })
        }

        let cf_names = rocksdb::DB::list_cf(&opts, path)
            .map_err(Into::<rocksdb::Error>::into)?;

        let cfs = cf_names
            .iter()
            .map(|n| rocksdb::ColumnFamilyDescriptor::new(n, opts.clone()));

        let tx_db_opts = rocksdb::TransactionDBOptions::default();

        let store = rocksdb::TransactionDB::open_cf_descriptors(&opts, &tx_db_opts, path, cfs)
            .map_err(Into::<rocksdb::Error>::into)?;

        Ok(RocksDBStore { store })
    }
    
    fn create_collection(&mut self, collection_name: String) -> Result<(), StoreError> {
        self.store
            .create_cf(collection_name, &rocksdb::Options::default())
            .map_err(Into::into)
    }

    fn delete_collection(&mut self, collection_name: String) -> Result<(), StoreError> {
        self.store
            .drop_cf(&collection_name)
            .map_err(Into::into)
    }

    fn get(&self, collection_name: String, key: String) -> Result<Option<Vec<u8>>, StoreError> {
        let cf = self.get_cf(collection_name)?;

        self.store
            .get_cf(cf, key)
            .map_err(Into::into)
    }

    fn get_many(&self, collection_name: String, keys: Vec<String>) -> Result<Vec<Result<Option<Vec<u8>>, StoreError>>, StoreError> {
        let cf = self.get_cf(collection_name)?;

        Ok(self.store
            .multi_get_cf::<&str, _, _>(
            keys
                    .iter()
                    .map(|k| (cf, k.as_str()))
                    .collect::<Vec<_>>()
            )
            .into_iter()
            .map(|r| r.map_err(Into::into))
            .collect())
    }

    fn put<'a>(&self, collection_name: String, key: String, value: Vec<u8>) -> Result<Self::Transaction<'_>, StoreError> {
        let cf = self.get_cf(collection_name)?;

        let tx = self.store.transaction();

        tx.put_cf(cf, key, value)
            .map_err(Into::<rocksdb::Error>::into)?;

        Ok(Self::Transaction::new(tx))
    }

    fn delete<'a>(&self, collection_name: String, key: String) -> Result<Self::Transaction<'_>, StoreError> {
        let cf = self.get_cf(collection_name)?;

        let tx = self.store.transaction();

        tx.delete_cf(cf, key)
            .map_err(Into::<rocksdb::Error>::into)?;

        Ok(Self::Transaction::new(tx))
    }
}

pub struct RocksDBTransaction<'a> {
    tx: rocksdb::Transaction<'a, rocksdb::TransactionDB>,
}

impl<'a> RocksDBTransaction<'a> {
    fn new(tx: rocksdb::Transaction<'a, rocksdb::TransactionDB>) -> RocksDBTransaction<'a> {
        RocksDBTransaction { tx }
    }
}

impl Transaction for RocksDBTransaction<'_> {
    fn commit(self) -> Result<(), TransactionError> {
        self.tx
            .commit()
            .map_err(Into::into)
    }

    fn rollback(self) -> Result<(), TransactionError> {
        self.tx
            .rollback()
            .map_err(Into::into)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn rocksdb_store_test() {
        let mut store = RocksDBStore::new("/tmp/test.db").unwrap();

        store.create_collection("c".to_string())
            .unwrap();

        store.put("c".to_string(), "k".to_string(), "v".as_bytes().to_vec())
            .unwrap()
            .commit()
            .unwrap();

        let value = store.get("c".to_string(), "k".to_string())
            .unwrap()
            .unwrap();

        let str = String::from_utf8(value).unwrap();

        println!("{}", str);

        assert_eq!("v", str);
    }
}