use std::collections::HashMap;

use bincode::{deserialize, serialize};

use super::{
    collection::Collection,
    index::Index,
    model::*,
    store::{Store, StoreError, Transaction, TransactionError},
};

#[derive(Debug, thiserror::Error)]
enum HeisenbergError {
    #[error("collection {0} already exists")]
    CollectionExists(String),

    #[error("collection {0} not found")]
    CollectionNotFound(String),

    #[error("store error: {0}")]
    StoreError(#[from] StoreError),

    #[error("transaction error: {0}")]
    TransactionError(#[from] TransactionError),
}

pub struct Heisenenberg<S: Store<S>>
where
    S: Transaction<S>,
{
    store: S,
    collections: HashMap<String, Collection<Box<dyn Index>>>,
}

impl<S: Store<S>> Heisenenberg<S>
where
    S: Transaction<S>,
{
    pub fn new(path: &str) -> Result<Heisenenberg<S>, HeisenbergError> {
        let store = <S as Store<S>>::new(path).map_err(|e| HeisenbergError::StoreError(e))?;

        Ok(Heisenenberg {
            store,
            collections: HashMap::new(),
        })
    }

    fn get_collection(
        &mut self,
        collection_name: &str,
    ) -> Result<&mut Collection<Box<dyn Index>>, HeisenbergError> {
        let collection = self.collections.get_mut(collection_name).ok_or(
            HeisenbergError::CollectionNotFound(collection_name.to_string()),
        )?;

        Ok(collection)
    }

    pub fn new_collection<I>(&mut self, name: String, index: I) -> Result<(), HeisenbergError>
    where
        I: Index,
    {
        if self.collections.contains_key(&name) {
            return Err(HeisenbergError::CollectionExists(name));
        }

        self.store.new_collection(name)?;

        let collection: Collection<Box<dyn Index>> =
            Collection::new(name.to_string(), Box::new(index));

        self.collections.insert(name.to_string(), collection);

        Ok(())
    }

    pub fn delete_collection(&mut self, name: String) -> Result<(), HeisenbergError> {
        if !self.collections.contains_key(&name) {
            return Err(HeisenbergError::CollectionNotFound(name));
        }

        self.store.delete_collection(name)?;

        self.collections.remove(&name);

        Ok(())
    }

    pub fn get(
        &self,
        collection_name: String,
        key: String,
    ) -> Result<Option<Value>, HeisenbergError> {
        if !self.collections.contains_key(&collection_name) {
            return Err(HeisenbergError::CollectionNotFound(
                collection_name.to_string(),
            ));
        }

        self.store
            .get(collection_name, key)
            .map_err(Into::into)
            .map(|bytes| bytes.and_then(|bytes| deserialize(&bytes).ok()))
    }

    pub fn put(
        &mut self,
        collection_name: String,
        key: String,
        value: Value,
    ) -> Result<(), HeisenbergError> {
        let collection = self.get_collection(&collection_name)?;
        
        let bytes = serialize(&value).unwrap();

        let mut tx = self.store.put(collection_name, key, bytes)?;

        collection
            .index
            .insert(key, value.vector)
            .or_else(|_| tx.rollback())?;

        tx.commit()
            .map_err(|e| HeisenbergError::TransactionError(e))
    }

    pub fn delete(&mut self, collection_name: String, key: String) -> Result<(), HeisenbergError> {
        let collection = self.get_collection(&collection_name)?;

        let mut tx = self.store.delete(collection_name, key)?;

        collection.index.remove(key).or_else(|_| tx.rollback())?;

        tx.commit()
            .map_err(|e| HeisenbergError::TransactionError(e))
    }

    // pub fn search(
    //     &self,
    //     collection_name: &str,
    //     vector: Vec<f32>,
    // ) -> Result<Vec<Value>, HeisenbergError> {
    //     let collection = self.get_collection(collection_name)?;
    // }
}

mod tests {
    use serde_json::json;

    use crate::core::{
        indices::hora_hnsw_index::{HoraHNSWIndex, HoraHNSWIndexOptions},
        stores::rocksdb_store::RocksDBStore,
    };

    #[test]
    fn test_heisenenberg() {
        use super::*;
        let h: Heisenenberg<RocksDBStore> = Heisenenberg::new("/tmp/heisenenberg.db").unwrap();

        let index = HoraHNSWIndex::new(3);

        h.new_collection("c", index)?;

        h.put(
            "c",
            "k",
            Value {
                vector: vec![1.0, 3.2, 7.1],
                meta: json!("{}"),
            },
        )?;

        let v = h.get("c", "k")?;

        assert_eq!(
            v,
            Some(Value {
                vector: vec![1.0, 3.2, 7.1],
                meta: json!("{}"),
            })
        );
    }
}
