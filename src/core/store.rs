#[derive(thiserror::Error, Debug)]
pub enum StoreError {
    #[error("store error: {0}")]
    Error(String),
    #[error("collection not found: {0}")]
    CollectionNotFound(String),
}

impl<E: Into<String>> From<E> for StoreError {
    fn from(error: E) -> Self {
        StoreError::Error(error.into())
    }
}

pub trait Store {
    type Transaction<'a> : Transaction where Self: 'a;

    fn new(path: &str) -> Result<Self, StoreError> where Self: Sized;
    fn create_collection(&mut self, collection_name: String) -> Result<(), StoreError>;
    fn delete_collection(&mut self, collection_name: String) -> Result<(), StoreError>;
    fn get(&self, collection_name: String, key: String) -> Result<Option<Vec<u8>>, StoreError>;
    fn get_many(&self, collection_name: String, keys: Vec<String>) -> Result<Vec<Result<Option<Vec<u8>>, StoreError>>, StoreError>;
    fn put<'a>(&self, collection_name: String, key: String, value: Vec<u8>) -> Result<Self::Transaction<'_>, StoreError>;
    fn delete<'a>(&self, collection_name: String, key: String) -> Result<Self::Transaction<'_>, StoreError>;
}

#[derive(thiserror::Error, Debug)]
pub enum TransactionError {
    #[error("transaction error: {0}")]
    Error(String),
}

impl<E: Into<String>> From<E> for TransactionError {
    fn from(error: E) -> Self {
        TransactionError::Error(error.into())
    }
}

pub trait Transaction {
    fn commit(self) -> Result<(), TransactionError>;
    fn rollback(self) -> Result<(), TransactionError>;
}
