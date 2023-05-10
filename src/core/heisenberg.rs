use std::collections::HashMap;

use super::{
    collection::Collection,
    index::Index,
    model::{Entry, Value},
    store::{Store, Transaction},
};


#[derive(Debug, thiserror::Error)]
pub enum HeisenbergError {
    #[error("heisenberg error: {0}")]
    Error(String),

    #[error("collection not found: {0}")]
    CollectionNotFound(String),

    #[error("collection already exists: {0}")]
    CollectionExists(String),
}

impl<E: Into<String>> From<E> for HeisenbergError {
    fn from(error: E) -> Self {
        HeisenbergError::Error(error.into())
    }
}
pub struct Heisenberg<S> {
    store: S,
    collections: HashMap<String, Collection>,
}

impl<S: Store> Heisenberg<S> {
    pub fn new(path: &str) -> Result<Heisenberg<S>, HeisenbergError> {
        let store = S::new(path)
            .map_err(|e|HeisenbergError::Error(e.to_string()))?;

        Ok(Heisenberg {
            store,
            collections: HashMap::new(),
        })
    }

    pub fn new_collection(&mut self, name: &str, index: Box<dyn Index>) -> Result<(), HeisenbergError> {
        self.collections
            .get_mut(name)
            .map_or(Ok(()), |_| Err(HeisenbergError::CollectionExists(name.to_string())))?;

        let collection = Collection::new(name, index);
        
        self.collections.insert(name.to_string(), collection);

        Ok(())
    }

    pub fn delete_collection(&mut self, name: &str) -> Result<(), HeisenbergError> {
        self.collections
            .get_mut(name)
            .ok_or(HeisenbergError::CollectionNotFound(name.to_string()))?;
    
        // TODO: properly flush index

        self.collections.remove(name);

        Ok(())
    }

    pub fn get(&self, collection_name: &str, key: &str) -> Result<Option<Value>, HeisenbergError> {
        self.store
            .get(collection_name.to_string(), key.to_string())
            .map_err(|e| HeisenbergError::Error(e.to_string()))
            .and_then(|o| match o {
                Some(v) => bincode::deserialize(&v)
                    .map_err(|e| HeisenbergError::Error(e.to_string()))
                    .map(Some), 
                None => Ok(None),
            })
    }

    pub fn get_many(&self, collection_name: &str, keys: Vec<String>) -> Result<Vec<Result<Option<Entry>, HeisenbergError>>, HeisenbergError> {
         self.store
            .get_many(
                collection_name.to_string(), 
                keys.clone(),
            )
            .map_err(|e| HeisenbergError::Error(e.to_string()))
            .map(|rs| 
                rs
                   .into_iter()
                   .enumerate()
                   .map(|(i, r)| {
                        match r {
                            Ok(Some(bytes)) => match bincode::deserialize::<Value>(&bytes) {
                                Ok(value) => Ok(Some(Entry { collection: collection_name.to_string(), key: keys[i].to_string(), value })),
                                Err(e) => Err(HeisenbergError::Error(e.to_string())),
                            },
                            Ok(None) => Ok(None),
                            Err(e) => Err(HeisenbergError::Error(e.to_string())),
                        }
                    })
                   .collect()
            )
    }

    pub fn put(&mut self, collection_name: &str, key: &str, value: Value) -> Result<(), HeisenbergError> {
        let collection = self.collections
            .get_mut(collection_name)
            .ok_or(HeisenbergError::CollectionNotFound(collection_name.to_string()))?;

        let raw = bincode::serialize(&value)
            .map_err(|e| HeisenbergError::Error(e.to_string()))?;

        let tx = self.store
            .put(collection_name.to_string(), key.to_string(), raw)
            .map_err(|e| HeisenbergError::Error(e.to_string()))?;

        
        if let Err(e) = collection.index.insert(value.vector, key.to_string()) {
            tx.rollback();
            return Err(HeisenbergError::Error(e.to_string()));
        }

        return tx
            .commit()
            .map_err(|e| HeisenbergError::Error(e.to_string()))
    }

    pub fn delete(&mut self, collection_name: &str, key: &str) -> Result<(), HeisenbergError> {
        let collection = self.collections
            .get_mut(collection_name)
            .ok_or(HeisenbergError::CollectionNotFound(collection_name.to_string()))?;

        let tx = self.store
            .delete(collection_name.to_string(), key.to_string())
            .map_err(|e| HeisenbergError::Error(e.to_string()))?;

        if let Err(e) = collection.index.remove(key.to_string()) {
            tx.rollback();
            return Err(HeisenbergError::Error(e.to_string()));
        }

        return tx
            .commit()
            .map_err(|e| HeisenbergError::Error(e.to_string()))
    }

    pub fn search(&self, collection_name: &str, query: Vec<f32>, k: usize) -> Result<Vec<Result<Option<Entry>, HeisenbergError>>, HeisenbergError> {
        self.collections
            .get(collection_name)
            .ok_or(HeisenbergError::CollectionNotFound(collection_name.to_string()))
            .and_then(|c| {
                let indices = c.index
                    .search(query, k)
                    .map_err(|e| HeisenbergError::Error(e.to_string()))
                    .map(|i| i.0)?;

                self.get_many(collection_name, indices)
            })
    }
}

#[cfg(tests)]
mod test {
    use crate::core::stores::rocksdb_store::RocksDBStore;

    use super::*;

    #[test]
    fn test_heisenenberg() {
        // let h: Heisenberg<RocksDBStore> = Heisenberg::new("/tmp/test.db").unwrap();
        // h.get();
    }
}
