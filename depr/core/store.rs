#[derive(Display, Debug, thiserror::Error)]
pub enum StoreError {
    #[error("store error: {0}")]
    Error(String),
}

impl<E: Into<String>> From<E> for StoreError {
    fn from(error: E) -> Self {
        StoreError::Error(error.into())
    }
}

pub trait Store<T: Transaction<T> + 'static> {
    fn new(path: &str) -> Result<Self, StoreError>
    where
        Self: Sized;
    fn open(path: &str) -> Result<Self, StoreError>
    where
        Self: Sized;
    fn new_collection(&mut self, name: String) -> Result<(), StoreError>;
    fn delete_collection(&mut self, name: String) -> Result<(), StoreError>;
    fn get(&self, collection: String, key: String) -> Result<Option<Vec<u8>>, StoreError>;
    fn put(&self, collection: String, key: String, value: Vec<u8>) -> Result<T, StoreError>;
    fn delete(&self, collection: String, key: String) -> Result<T, StoreError>;
}


#[derive(Debug, thiserror::Error)]
pub enum TransactionError {
    #[error("{0}")]
    Error(String),
}

impl From<String> for TransactionError {
    fn from(error: String) -> Self {
        TransactionError::Error(error)
    }
}

pub trait Transaction<T> {
    fn new(tx: T) -> Self;
    fn commit(&mut self) -> Result<(), TransactionError>;
    fn rollback(&mut self) -> Result<(), TransactionError>;
}
